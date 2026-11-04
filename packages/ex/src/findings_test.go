package main

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jonatassales/tao/packages/ex/src/domain/finops"
	"github.com/jonatassales/tao/packages/ex/src/domain/harness"
	"github.com/jonatassales/tao/packages/ex/src/domain/microruntime"
	projectdomain "github.com/jonatassales/tao/packages/ex/src/domain/project"
	"github.com/jonatassales/tao/packages/ex/src/domain/run"
	"github.com/jonatassales/tao/packages/ex/src/infrastructure/records"
)

// A plan that fails after the judge was paid still leaves the money record:
// the ledger's planning line and the private entry the gates count.
func TestAFailedPlanKeepsTheJudgesSpend(t *testing.T) {
	t.Parallel()

	p, deps := seedProject(t)
	deps.Judge = &scriptedJudge{answers: sampleAnswers(t)}

	status, stdout, stderr := p.exWith(deps, "agent", "plan", consoleTicket, "--component", "name=missing,path=services/missing,port=80,readiness=tcp")
	if status == 0 || !strings.Contains(stderr, "services/missing") {
		t.Fatalf("a plan over a missing component gave %d: %s\n%s", status, stderr, stdout)
	}

	ledger := ledgerOf(t, p)
	if ledger.Planning.Model == nil || ledger.Planning.Cents.Observed != 9 || ledger.Planning.Cents.Reserved != 15 {
		t.Errorf("the planning line after a failed plan = %+v", ledger.Planning)
	}
	if spent := mustRead(t, filepath.Join(p.root, ".langdrift", "agent", "state", "runs", "console-microruntime", "planning.json")); !strings.Contains(string(spent), `"observed": 9`) {
		t.Errorf("planning.json = %s", spent)
	}
}

// Each plan keeps its own planning entry, so a period counts only the
// plans made in it.
func TestPlanningSpendCountsByPeriod(t *testing.T) {
	t.Parallel()

	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	owner := projectdomain.Project{Root: root}
	runtime, _ := microruntime.ParseIdentifier("console-microruntime")
	store := records.LapStore{}

	ledger := ledgerLine(9)
	if err := store.SavePlanning(owner, runtime, harness.At(time.Date(2026, 9, 30, 23, 0, 0, 0, time.UTC)), ledger); err != nil {
		t.Fatal(err)
	}
	if err := store.SavePlanning(owner, runtime, harness.At(time.Date(2026, 10, 1, 9, 0, 0, 0, time.UTC)), ledgerLine(12)); err != nil {
		t.Fatal(err)
	}

	for period, want := range map[string]int{"2026-09": 9, "2026-10": 12, "2026-11": 0} {
		if spent, err := store.Period(owner, period); err != nil || spent != want {
			t.Errorf("Period(%s) = %d, %v; want %d", period, spent, err, want)
		}
	}
}

// A call the driver cut off leaves what it used unknown: the next paid call
// pauses at FG-04 instead of counting it as zero.
func TestAnInterruptedCallLeavesTheSpendUnknown(t *testing.T) {
	t.Parallel()

	p, deps := plannedProject(t)
	driver := websiteDriver()
	driver.cutOff["to-spec"] = true
	deps.Driver, deps.Lifecycle, deps.Delivery, deps.Poll = driver, websiteLifecycle(), websiteDelivery(), time.Millisecond

	if status, stdout, stderr := p.exInput(deps, "y\n", "agent", "run", "console-microruntime", "--confirm"); status != 1 || !strings.Contains(stderr, "interrupted") {
		t.Fatalf("an interrupted call gave %d: %s\n%s", status, stderr, stdout)
	}

	driver.cutOff = map[string]bool{}
	status, stdout, stderr := p.exWith(deps, "agent", "run", "console-microruntime")
	if status != 3 || !strings.Contains(stderr, "FG-04") {
		t.Fatalf("the next call after an interrupted one gave %d: %s\n%s", status, stderr, stdout)
	}
	expectLines(t, stdout, "agent › gate       paused · FG-04 · observed spend is unknown; a call reported no usage")
	if line := ledgerOf(t, p).Laps[0]; line.TokensThousands.Observed != nil || line.Cents.Variance != nil {
		t.Errorf("an interrupted lap's line = %+v", line)
	}
}

// A runtime that started but was not ready is looked at again on the next
// run, never taken for ready.
func TestARuntimeThatWasNotReadyIsObservedAgain(t *testing.T) {
	t.Parallel()

	p, deps := plannedProject(t)
	lifecycle := websiteLifecycle()
	lifecycle.notReady = 1
	deps.Driver, deps.Lifecycle, deps.Delivery, deps.Poll = websiteDriver(), lifecycle, websiteDelivery(), time.Millisecond

	if status, stdout, stderr := p.exInput(deps, "y\n", "agent", "run", "console-microruntime", "--confirm"); status != 1 || !strings.Contains(stderr, "not ready") {
		t.Fatalf("a runtime that was not ready gave %d: %s\n%s", status, stderr, stdout)
	}

	status, stdout, stderr := p.exWith(deps, "agent", "run", "console-microruntime")
	if status != 0 {
		t.Fatalf("the resumed run gave %d: %s\n%s", status, stderr, stdout)
	}
	expectLines(t, stdout, "runtime    console-microruntime · revision 1 · 1 Deployments ready · 500m · 1024 MiB", "✓ ready      #19 for review · $1.60 released to the envelope")
	if lifecycle.ups != 1 {
		t.Errorf("the runtime was brought up %d times, want once", lifecycle.ups)
	}
	if lap := string(mustRead(t, filepath.Join(runtimeFolder(p), "laps", "0001.lap.md"))); !strings.Contains(lap, "| Ready | ") {
		t.Errorf("the lap has no Ready step:\n%s", lap)
	}
}

// Without a person in review, the run is ready to merge only when delivery
// shows a pushed branch and a pull request whose checks all passed.
func TestReadyToMergeNeedsACleanDelivery(t *testing.T) {
	t.Parallel()

	p, deps := seedProject(t)
	workflow := bodyFile(t, p.root, "workflow-no-review.json", "harness/workflow.json", func(body *harness.Workflow) {
		withoutMirror(body)
		body.Defaults.HumanAt = []string{"planning", "merge"}
	})
	if status, stdout, stderr := p.exWith(deps, "harness", "workflow", "set", "--file", workflow); status != 0 {
		t.Fatalf("ex harness workflow set = %d: %s\n%s", status, stderr, stdout)
	}
	if status, stdout, stderr := p.exWith(deps, "agent", "plan", consoleTicket, "--no-judge", "--component", consoleComponent); status != 0 {
		t.Fatalf("ex agent plan = %d: %s\n%s", status, stderr, stdout)
	}

	delivery := websiteDelivery()
	delivery.checks[1][2].Result = "failed"
	deps.Driver, deps.Lifecycle, deps.Delivery, deps.Poll = websiteDriver(), websiteLifecycle(), delivery, time.Millisecond

	status, stdout, stderr := p.exInput(deps, "y\n", "agent", "run", "console-microruntime", "--confirm")
	if status != 0 {
		t.Fatalf("ex agent run = %d: %s\n%s", status, stderr, stdout)
	}
	expectLines(t, stdout, "✗ checks     3 of 4 passed · test failed", "✓ ready      #19 for review · $1.60 released to the envelope")
	if record := runRecord(t, p); *record.Outcome != "readyForReview" {
		t.Errorf("a failed check ended %s", *record.Outcome)
	}

	other, otherDeps := seedProject(t)
	if status, stdout, stderr := other.exWith(otherDeps, "harness", "workflow", "set", "--file", workflow); status != 0 {
		t.Fatalf("ex harness workflow set = %d: %s\n%s", status, stderr, stdout)
	}
	if status, stdout, stderr := other.exWith(otherDeps, "agent", "plan", consoleTicket, "--no-judge", "--component", consoleComponent); status != 0 {
		t.Fatalf("ex agent plan = %d: %s\n%s", status, stderr, stdout)
	}
	otherDeps.Driver, otherDeps.Lifecycle, otherDeps.Delivery, otherDeps.Poll = websiteDriver(), websiteLifecycle(), websiteDelivery(), time.Millisecond
	status, stdout, stderr = other.exInput(otherDeps, "y\n", "agent", "run", "console-microruntime", "--confirm")
	if status != 0 {
		t.Fatalf("ex agent run = %d: %s\n%s", status, stderr, stdout)
	}
	expectLines(t, stdout, "✓ ready      #19 to merge · $1.60 released to the envelope", "# next: the merge stays the owner's")
}

// One plan, run or confirmation changes a runtime at a time.
func TestTheRunLockRefusesASecondOperation(t *testing.T) {
	t.Parallel()

	p, deps := plannedProject(t)
	deps.Driver, deps.Lifecycle, deps.Delivery, deps.Poll = websiteDriver(), websiteLifecycle(), websiteDelivery(), time.Millisecond

	runtime, _ := microruntime.ParseIdentifier("console-microruntime")
	unlock, err := records.LapStore{}.Lock(projectdomain.Project{Root: p.root}, runtime)
	if err != nil {
		t.Fatal(err)
	}
	defer unlock()

	if status, _, stderr := p.exWith(deps, "agent", "run", "console-microruntime"); status != 3 || !strings.Contains(stderr, "busy") {
		t.Errorf("a run while the lock is held gave %d: %s", status, stderr)
	}
	if status, _, stderr := p.exWith(deps, "agent", "plan", consoleTicket, "--no-judge", "--component", consoleComponent); status != 3 || !strings.Contains(stderr, "busy") {
		t.Errorf("a plan while the lock is held gave %d: %s", status, stderr)
	}
}

// A lap FG-04 paused halfway through a stage plays on from the next skill.
func TestAStagePausedHalfwayPlaysOnFromTheNextSkill(t *testing.T) {
	t.Parallel()

	p, deps := plannedProject(t)
	driver := websiteDriver()
	driver.usage["implement"] = [2]int64{5_000_000, 300_000}
	deps.Driver, deps.Lifecycle, deps.Delivery, deps.Poll = driver, websiteLifecycle(), websiteDelivery(), time.Millisecond

	if status, stdout, stderr := p.exInput(deps, "y\n", "agent", "run", "console-microruntime", "--confirm"); status != 3 || !strings.Contains(stderr, "FG-04") {
		t.Fatalf("the exhausted lap gave %d: %s\n%s", status, stderr, stdout)
	}
	if implement := stageState(runRecord(t, p), "implement"); strings.Join(implement.SkillsDone, ",") != "implement" {
		t.Errorf("implement after the pause = %+v", implement)
	}

	before := len(driver.calls)
	if status, stdout, stderr := p.exWith(deps, "agent", "run", "console-microruntime", "--acknowledge-overspend"); status != 0 {
		t.Fatalf("the next lap gave %d: %s\n%s", status, stderr, stdout)
	}
	var replayed []string
	for _, call := range driver.calls[before:] {
		replayed = append(replayed, call.Skill)
	}
	if strings.Join(replayed, ",") != "tdd,code-review" {
		t.Errorf("the next lap ran %v, want tdd and code-review only", replayed)
	}
	if implement := stageState(runRecord(t, p), "implement"); implement.State != run.StageDone {
		t.Errorf("implement after the next lap = %+v", implement)
	}
}

func ledgerLine(observed int) finops.PlanningLine {
	return finops.PlanningLine{Cents: finops.PlanningCents{Observed: observed, Reserved: observed}}
}
