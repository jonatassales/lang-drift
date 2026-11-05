package main

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jonatassales/tao/lib/cli/src/domain/finops"
	"github.com/jonatassales/tao/lib/cli/src/domain/harness"
	"github.com/jonatassales/tao/lib/cli/src/domain/microruntime"
	projectdomain "github.com/jonatassales/tao/lib/cli/src/domain/project"
	"github.com/jonatassales/tao/lib/cli/src/domain/run"
	"github.com/jonatassales/tao/lib/cli/src/infrastructure/records"
)

// A plan that fails after the judge was paid still leaves the money record:
// the ledger's planning line and the private entry the gates count.
func TestAFailedPlanKeepsTheJudgesSpend(t *testing.T) {
	t.Parallel()

	p, deps := seedProject(t)
	deps.Judge = &scriptedJudge{answers: sampleAnswers(t)}

	status, stdout, stderr := p.dftWith(deps, "agent", "plan", consoleTicket, "--component", "name=missing,path=services/missing,port=80,readiness=tcp")
	if status == 0 || !strings.Contains(stderr, "services/missing") {
		t.Fatalf("a plan over a missing component gave %d: %s\n%s", status, stderr, stdout)
	}

	ledger := ledgerOf(t, p)
	if ledger.Planning.Model == nil || ledger.Planning.Cents.Observed != 10 || ledger.Planning.Cents.Reserved != 15 {
		t.Errorf("the planning line after a failed plan = %+v", ledger.Planning)
	}
	if spent := mustRead(t, filepath.Join(p.root, ".langdrift", "agent", "state", "runs", "console-microruntime", "planning.json")); !strings.Contains(string(spent), `"observed": 10`) {
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

// A call the driver cut off leaves what it used unknown: the lap fails with
// its reservation held instead of counting the call as zero, and the next
// run opens a new lap.
func TestAnInterruptedCallLeavesTheSpendUnknown(t *testing.T) {
	t.Parallel()

	p, deps := plannedProject(t)
	driver := websiteDriver()
	driver.cutOff["to-spec"] = true
	lifecycle := websiteLifecycle()
	deps.Driver, deps.Lifecycle, deps.Delivery, deps.Poll = driver, lifecycle, websiteDelivery(), time.Millisecond

	if status, stdout, stderr := p.dftInput(deps, "y\n", "agent", "run", "console-microruntime", "--confirm"); status != 1 || !strings.Contains(stderr, "interrupted") {
		t.Fatalf("an interrupted call gave %d: %s\n%s", status, stderr, stdout)
	}
	if record := runRecord(t, p); record.State != "failed" || record.Laps[0].Outcome != "failed" || lifecycle.downs != 1 {
		t.Errorf("after an interrupted call: %+v, %d downs", record, lifecycle.downs)
	}
	lap := string(mustRead(t, filepath.Join(runtimeFolder(p), "laps", "0001.lap.md")))
	if !strings.Contains(lap, "| Held | ") || strings.Contains(lap, "| Released | ") {
		t.Errorf("the interrupted lap did not hold its reservation:\n%s", lap)
	}
	if line := ledgerOf(t, p).Laps[0]; line.TokensThousands.Observed != nil || line.Cents.Variance != nil {
		t.Errorf("an interrupted lap's line = %+v", line)
	}

	driver.cutOff = map[string]bool{}
	status, stdout, stderr := p.dftInput(deps, "y\n", "agent", "run", "console-microruntime", "--confirm")
	if status != 0 {
		t.Fatalf("the next run after an interrupted one gave %d: %s\n%s", status, stderr, stdout)
	}
	if runRecord(t, p).Laps[1].Number != "0002" || ledgerOf(t, p).Laps[0].TokensThousands.Observed != nil {
		t.Errorf("the next run did not open lap 0002 over the held lap 0001:\n%s", stdout)
	}
}

// A runtime that started but was not ready fails the lap before any stage:
// it is taken down, the reservation returns, and the next run opens a new
// lap that brings it up again.
func TestARuntimeThatIsNotReadyFailsTheLap(t *testing.T) {
	t.Parallel()

	p, deps := plannedProject(t)
	lifecycle := websiteLifecycle()
	lifecycle.notReady = 1
	driver := websiteDriver()
	deps.Driver, deps.Lifecycle, deps.Delivery, deps.Poll = driver, lifecycle, websiteDelivery(), time.Millisecond

	if status, stdout, stderr := p.dftInput(deps, "y\n", "agent", "run", "console-microruntime", "--confirm"); status != 1 || !strings.Contains(stderr, "not ready") {
		t.Fatalf("a runtime that was not ready gave %d: %s\n%s", status, stderr, stdout)
	}
	if record := runRecord(t, p); record.State != "failed" || record.Laps[0].Outcome != "failed" || len(driver.calls) != 0 || lifecycle.downs != 1 {
		t.Errorf("after a runtime that was not ready: %+v, %d calls, %d downs", record, len(driver.calls), lifecycle.downs)
	}
	if lap := string(mustRead(t, filepath.Join(runtimeFolder(p), "laps", "0001.lap.md"))); !strings.Contains(lap, "| Released | ") {
		t.Errorf("the failed lap released nothing:\n%s", lap)
	}

	status, stdout, stderr := p.dftInput(deps, "y\n", "agent", "run", "console-microruntime", "--confirm")
	if status != 0 {
		t.Fatalf("the next run gave %d: %s\n%s", status, stderr, stdout)
	}
	expectLines(t, stdout, "runtime    console-microruntime · revision 1 · 1 Deployments ready · 500m · 1024 MiB", "✓ ready      #19 for review · $1.60 released to the envelope")
	if lifecycle.ups != 2 || runRecord(t, p).Laps[1].Number != "0002" {
		t.Errorf("the next run brought the runtime up %d times in total, want twice, in lap 0002", lifecycle.ups)
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
	if status, stdout, stderr := p.dftWith(deps, "harness", "workflow", "set", "--file", workflow); status != 0 {
		t.Fatalf("dft harness workflow set = %d: %s\n%s", status, stderr, stdout)
	}
	if status, stdout, stderr := p.dftWith(deps, "agent", "plan", consoleTicket, "--component", consoleComponent); status != 0 {
		t.Fatalf("dft agent plan = %d: %s\n%s", status, stderr, stdout)
	}

	delivery := websiteDelivery()
	delivery.checks[1][2].Result = "failed"
	deps.Driver, deps.Lifecycle, deps.Delivery, deps.Poll = websiteDriver(), websiteLifecycle(), delivery, time.Millisecond

	status, stdout, stderr := p.dftInput(deps, "y\n", "agent", "run", "console-microruntime", "--confirm")
	if status != 0 {
		t.Fatalf("dft agent run = %d: %s\n%s", status, stderr, stdout)
	}
	expectLines(t, stdout, "✗ checks     3 of 4 passed · test failed", "✓ ready      #19 for review · $1.60 released to the envelope")
	if record := runRecord(t, p); *record.Outcome != "readyForReview" {
		t.Errorf("a failed check ended %s", *record.Outcome)
	}

	other, otherDeps := seedProject(t)
	if status, stdout, stderr := other.dftWith(otherDeps, "harness", "workflow", "set", "--file", workflow); status != 0 {
		t.Fatalf("dft harness workflow set = %d: %s\n%s", status, stderr, stdout)
	}
	if status, stdout, stderr := other.dftWith(otherDeps, "agent", "plan", consoleTicket, "--component", consoleComponent); status != 0 {
		t.Fatalf("dft agent plan = %d: %s\n%s", status, stderr, stdout)
	}
	otherDeps.Driver, otherDeps.Lifecycle, otherDeps.Delivery, otherDeps.Poll = websiteDriver(), websiteLifecycle(), websiteDelivery(), time.Millisecond
	status, stdout, stderr = other.dftInput(otherDeps, "y\n", "agent", "run", "console-microruntime", "--confirm")
	if status != 0 {
		t.Fatalf("dft agent run = %d: %s\n%s", status, stderr, stdout)
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

	if status, _, stderr := p.dftWith(deps, "agent", "run", "console-microruntime"); status != 3 || !strings.Contains(stderr, "busy") {
		t.Errorf("a run while the lock is held gave %d: %s", status, stderr)
	}
	if status, _, stderr := p.dftWith(deps, "agent", "plan", consoleTicket, "--component", consoleComponent); status != 3 || !strings.Contains(stderr, "busy") {
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

	if status, stdout, stderr := p.dftInput(deps, "y\n", "agent", "run", "console-microruntime", "--confirm"); status != 3 || !strings.Contains(stderr, "FG-04") {
		t.Fatalf("the exhausted lap gave %d: %s\n%s", status, stderr, stdout)
	}
	if implement := stageState(runRecord(t, p), "implement"); strings.Join(implement.SkillsDone, ",") != "implement" {
		t.Errorf("implement after the pause = %+v", implement)
	}

	before := len(driver.calls)
	if status, stdout, stderr := p.dftWith(deps, "agent", "run", "console-microruntime", "--acknowledge-overspend"); status != 0 {
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
