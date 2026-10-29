package main

import (
	"bytes"
	"context"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jonatassales/tao/packages/ex/src/application/plan"
	"github.com/jonatassales/tao/packages/ex/src/domain/run"
	"github.com/jonatassales/tao/packages/ex/src/presentation/cli"
	"github.com/jonatassales/tao/packages/ex/src/presentation/cli/support"
)

const consoleTicket = ".agents/workflow/tickets/0070-console-microruntime.ticket.md"

// A plan over a run that is going on is refused; a plan over a run that
// ended raises the revision and keeps the laps, the delivery and the money.
func TestReplanKeepsTheRecordsAndRefusesAnOpenLap(t *testing.T) {
	t.Parallel()

	p, deps := plannedProject(t)
	deps.Driver, deps.Lifecycle, deps.Delivery, deps.Poll = websiteDriver(), websiteLifecycle(), websiteDelivery(), time.Millisecond

	if status, stdout, stderr := p.exWith(deps, "agent", "run", "console-microruntime"); status != 0 {
		t.Fatalf("ex agent run = %d: %s\n%s", status, stderr, stdout)
	}
	if status, _, stderr := p.exWith(deps, "agent", "plan", consoleTicket, "--no-judge", "--component", consoleComponent); status != 3 || !strings.Contains(stderr, "waiting with a lap open") {
		t.Errorf("a plan over an open lap gave %d: %s", status, stderr)
	}

	if status, stdout, stderr := p.exWith(deps, "agent", "confirm", "console-microruntime", "planning"); status != 0 {
		t.Fatalf("ex agent confirm = %d: %s\n%s", status, stderr, stdout)
	}
	if status, stdout, stderr := p.exWith(deps, "agent", "run", "console-microruntime"); status != 0 {
		t.Fatalf("the run gave %d: %s\n%s", status, stderr, stdout)
	}

	status, stdout, stderr := p.exWith(deps, "agent", "plan", consoleTicket, "--no-judge", "--component", consoleComponent)
	if status != 0 {
		t.Fatalf("a plan over a finished run gave %d: %s\n%s", status, stderr, stdout)
	}
	expectLines(t, stdout, "✓ plan       runtime.md draft, revision 2 · proposed · authorize with ex agent run")

	record := runRecord(t, p)
	if record.State != "planned" || record.Outcome != nil || len(record.Laps) != 1 || record.Laps[0].Outcome != "completed" || record.Delivery.PullRequest == nil || record.Delivery.PullRequest.Number != 19 || record.Window.Started == nil || record.Window.Finished != nil {
		t.Errorf("the replanned run record = state %s, laps %+v, delivery %+v, window %+v", record.State, record.Laps, record.Delivery, record.Window)
	}
	for _, state := range record.Stages {
		if state.State != "pending" || state.ConfirmedBy != nil {
			t.Errorf("stage %s = %+v after the replan, want pending", state.ID, state)
		}
	}

	ledger := ledgerOf(t, p)
	if len(ledger.Laps) != 1 || ledger.Laps[0].Cents.Observed != 980 || ledger.Totals.Cents.Observed != 980 || ledger.Bounded.PeriodLeftCents != 49020 {
		t.Errorf("the replanned ledger = laps %+v, totals %+v, bounded %+v", ledger.Laps, ledger.Totals, ledger.Bounded)
	}

	var planned run.Plan
	decodeStrict(t, []byte(fenceBody(t, filepath.Join(runtimeFolder(p), "runtime.md"), "langdrift-plan")), &planned)
	if planned.Runtime.Revision != 2 || planned.State != "proposed" {
		t.Errorf("the replanned plan = revision %d, state %s", planned.Runtime.Revision, planned.State)
	}
}

// FG-03 pauses on the ticket's allocation and on what is left of the period,
// not only on the per-run cap, for a run and for the planner alike.
func TestGatePausesOnTheAllocationAndThePeriod(t *testing.T) {
	t.Parallel()

	p, deps := plannedProject(t)
	deps.Driver, deps.Lifecycle, deps.Delivery, deps.Poll = websiteDriver(), websiteLifecycle(), websiteDelivery(), time.Millisecond

	// One run spends $9.80 of the period first.
	if status, stdout, stderr := p.exInput(deps, "y\n", "agent", "run", "console-microruntime", "--confirm"); status != 0 {
		t.Fatalf("the first run gave %d: %s\n%s", status, stderr, stdout)
	}
	if status, stdout, stderr := p.exWith(deps, "agent", "plan", consoleTicket, "--no-judge", "--component", consoleComponent, "--runtime", "console-again"); status != 0 {
		t.Fatalf("the second plan gave %d: %s\n%s", status, stderr, stdout)
	}

	budget := func(periodCents, perRunCents, allocation string) {
		t.Helper()
		if status, stdout, stderr := p.exWith(deps, "finops", "budget", "set", "--period", "2026-10", "--period-cents", periodCents, "--per-run-cents", perRunCents, "--reserve-margin-percent", "20", "--allocate", "0070-console-microruntime="+allocation); status != 0 {
			t.Fatalf("ex finops budget set = %d: %s\n%s", status, stderr, stdout)
		}
	}

	budget("50000", "2500", "1000")
	status, stdout, stderr := p.exWith(deps, "agent", "run", "console-again")
	if status != 3 || !strings.Contains(stderr, "ticket allocation") {
		t.Fatalf("a reservation over the allocation gave %d: %s\n%s", status, stderr, stdout)
	}
	expectLines(t, stdout, "agent › gate       reservation $11.40 exceeds the $10.00 allocation · paused")
	if lap := string(mustRead(t, filepath.Join(p.root, ".langdrift", "microruntimes", "runtimes", "console-again", "laps", "0001.lap.md"))); !strings.Contains(lap, "| Paused | 13:46:50 | FG-03: a $11.40 reservation exceeds the $10.00 allocation |") {
		t.Errorf("lap 0001 of console-again:\n%s", lap)
	}

	budget("2000", "2000", "2000")
	status, stdout, stderr = p.exWith(deps, "agent", "run", "console-again")
	if status != 3 || !strings.Contains(stderr, "left of the period") {
		t.Fatalf("a reservation over the period gave %d: %s\n%s", status, stderr, stdout)
	}
	expectLines(t, stdout, "agent › gate       reservation $11.40 exceeds the $10.20 left of the period · paused")

	// The planner's own reservation, 15 cents, meets the same gate.
	budget("985", "985", "985")
	deps.Judge = &scriptedJudge{answers: sampleAnswers(t)}
	if status, _, stderr := p.exWith(deps, "agent", "plan", consoleTicket, "--runtime", "console-third", "--component", consoleComponent); status != 3 || !strings.Contains(stderr, "left of the period") || !strings.Contains(stderr, "nothing was paid") {
		t.Errorf("a planning reservation over the period gave %d: %s", status, stderr)
	}
}

// With --json, every line of the session is one event with the website's
// line id and the card and field it fills.
func TestSessionEmitsEventsAsJSON(t *testing.T) {
	t.Parallel()

	p, deps := seedProject(t)

	status, stdout, stderr := p.exWith(deps, "--json", "agent", "plan", consoleTicket, "--no-judge", "--component", consoleComponent)
	if status != 0 {
		t.Fatalf("ex --json agent plan = %d: %s\n%s", status, stderr, stdout)
	}

	type event struct {
		At      string `json:"at"`
		Session string `json:"session"`
		Kind    string `json:"kind"`
		ID      string `json:"id"`
		Text    string `json:"text"`
		Card    string `json:"card"`
		Field   string `json:"field"`
	}
	fields := map[string]event{}
	var ids []string
	lines := strings.Split(strings.TrimSpace(stdout), "\n")
	for _, line := range lines[:len(lines)-1] {
		var record event
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			t.Fatalf("not an event: %s (%v)", line, err)
		}
		if record.Session != "plan" || record.At == "" || record.Kind == "" || record.ID == "" || record.Text == "" {
			t.Errorf("event lacks a field: %+v", record)
		}
		ids = append(ids, record.ID)
		fields[record.ID] = record
	}

	joined := strings.Join(ids, " ")
	for _, want := range []string{"plan scope readTicket links", "bounds budget memory embed vectorSearch candidates", "validate draft handoff"} {
		if !strings.Contains(joined, want) {
			t.Errorf("the events %v lack the sequence %q", ids, want)
		}
	}
	if bounds := fields["bounds"]; bounds.Card != "plan" || bounds.Field != "bounded" || bounds.Kind != "output" {
		t.Errorf("the bounds event = %+v, want the plan card's bounded field", bounds)
	}
	if draft := fields["draft"]; draft.Card != "plan" || draft.Field != "drafted" || draft.Kind != "success" {
		t.Errorf("the draft event = %+v", draft)
	}

	var result plan.Result
	if err := json.Unmarshal([]byte(lines[len(lines)-1]), &result); err != nil || result.Runtime != "console-microruntime" || result.Revision != 1 {
		t.Errorf("the result line = %s (%v)", lines[len(lines)-1], err)
	}
}

// An answer that does not fit the question is refused once, and the plan
// keeps the deterministic shell with what the call cost on the ledger.
func TestPlanRefusesAnAnswerThatDoesNotFit(t *testing.T) {
	t.Parallel()

	p, deps := seedProject(t)
	judge := &scriptedJudge{answers: map[string]string{plan.SourcesQuestion: "I would rather talk about the weather."}}
	deps.Judge = judge

	status, stdout, stderr := p.exWith(deps, "agent", "plan", consoleTicket, "--component", consoleComponent)
	if status != 0 {
		t.Fatalf("ex agent plan = %d: %s\n%s", status, stderr, stdout)
	}
	if strings.Count(stdout, "judge      stopped · the sources answer was refused") != 1 || len(judge.asked) != 1 {
		t.Errorf("the refused answer was not reported once (%d calls):\n%s", len(judge.asked), stdout)
	}
	expectLines(t, stdout, "✓ suggests   claude-sonnet-5-5 · effort medium", "✓ plan       runtime.md draft, revision 1 · proposed · authorize with ex agent run")

	ledger := ledgerOf(t, p)
	if ledger.Planning.Model == nil || ledger.Planning.Cents.Observed == 0 || ledger.Planning.Cents.Billed != nil {
		t.Errorf("the refused call is not on the planning line: %+v", ledger.Planning)
	}
}

// Without an episode indexed, the plan says so instead of citing nothing.
func TestPlanWithoutEpisodesSaysSo(t *testing.T) {
	t.Parallel()

	p, deps := seedProject(t)

	status, stdout, stderr := p.exWith(deps, "agent", "plan", consoleTicket, "--no-judge", "--component", consoleComponent)
	if status != 0 {
		t.Fatalf("ex agent plan = %d: %s\n%s", status, stderr, stdout)
	}
	expectLines(t, stdout,
		"agent › 0 candidates · no episode indexed yet",
		"agent › dropped 0 below 0.55 · kept 0",
		"agent › 0 episodes cited · none with measured resources",
	)
}

// A component that lacks what the fence needs blocks the plan with the
// fields named; a port the manifest declares needs no flag.
func TestPlanBlocksAnIncompleteComponentAndReadsThePortFromTheManifest(t *testing.T) {
	t.Parallel()

	p, deps := seedProject(t)

	status, stdout, stderr := p.exWith(deps, "agent", "plan", consoleTicket, "--no-judge", "--component", "path=.", "--component", "path=apps/console,image=langdrift-console:local,readiness=/icon.svg,build=platform/microruntime/console/build-image.sh")
	if status != 3 || !strings.Contains(stderr, "blocked") {
		t.Fatalf("an incomplete component gave %d: %s\n%s", status, stderr, stdout)
	}
	expectLines(t, stdout, "agent › tool call: inspect apps/console · package.json · :3001")
	if !strings.Contains(stdout, "✗ plan       runtime.md draft, revision 1 · blocked · ") || !strings.Contains(stdout, "components-complete") {
		t.Errorf("the plan did not say which check blocked it:\n%s", stdout)
	}

	var planned run.Plan
	decodeStrict(t, []byte(fenceBody(t, filepath.Join(runtimeFolder(p), "runtime.md"), "langdrift-plan")), &planned)
	var completeness *run.Check
	for index := range planned.Checks {
		if planned.Checks[index].ID == "components-complete" {
			completeness = &planned.Checks[index]
		}
	}
	if completeness == nil || completeness.Result != "failed" || completeness.Detail == nil || *completeness.Detail != "missing .: image, port, readiness, build" {
		t.Errorf("components-complete = %+v", completeness)
	}
	if len(planned.Components) != 2 || planned.Components[1].Port != 3001 {
		t.Errorf("components = %+v, want the console's port from its manifest", planned.Components)
	}
}

// The connected loop answers to supervise.
func TestSuperviseIsTheConnectedLoop(t *testing.T) {
	t.Parallel()

	var stdout, stderr bytes.Buffer
	status := cli.Execute(context.Background(), []string{"agent", "supervise", "--help"}, support.Options{Version: version, Directory: t.TempDir(), Output: &stdout, Diagnostics: &stderr}, mainUseCases())
	if status != 0 || !strings.Contains(stdout.String(), "Supervise the connected runtimes") {
		t.Errorf("ex agent supervise --help = %d\n%s%s", status, stdout.String(), stderr.String())
	}
}
