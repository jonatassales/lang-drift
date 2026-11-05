package main

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/jonatassales/tao/lib/cli/src/domain/finops"
	"github.com/jonatassales/tao/lib/cli/src/domain/run"
)

// plannedProject seeds a project and plans ticket 0070 without a judge.
func plannedProject(t *testing.T) (project, dependencies) {
	t.Helper()

	p, deps := seedProject(t)
	if status, stdout, stderr := p.exWith(deps, "agent", "plan", ".agents/workflow/tickets/0070-console-microruntime.ticket.md", "--no-judge", "--component", consoleComponent); status != 0 {
		t.Fatalf("ex agent plan = %d\n%s%s", status, stdout, stderr)
	}

	return p, deps
}

// The run's opening act reproduces the sample's first lap: Opus at high
// estimates $34.00, asks for a $40.80 reservation and pauses at FG-03 with
// nothing paid; then the suggestion reserves $11.40 and opens lap 0002.
func TestRunPausesAtTheGateThenOpensALap(t *testing.T) {
	t.Parallel()

	p, deps := plannedProject(t)
	lifecycle := websiteLifecycle()
	deps.Lifecycle = lifecycle
	folder := filepath.Join(p.root, ".langdrift", "microruntimes", "runtimes", "console-microruntime")

	status, stdout, stderr := p.exWith(deps, "agent", "run", "console-microruntime", "--model", "claude-opus-5-5", "--effort", "high")
	if status != 3 || !strings.Contains(stderr, "FG-03") || !strings.Contains(stderr, "nothing was paid") {
		t.Fatalf("a paused lap gave %d: %s\n%s", status, stderr, stdout)
	}

	for _, line := range []string{
		"ex agent run console-microruntime",
		"# run · 0070-console-microruntime · the recorded intent, authorized by owner",
		"agent › intent     claude-opus-5-5 · effort high · runtime.md revision 1",
		"agent › estimate   $34.00 · 6M input · 500K output · list prices · FIN-02",
		"agent › budget     $25.00 per run · set by owner · FIN-03",
		"agent › gate       reservation $40.80 exceeds the $25.00 budget · paused",
		"# needs a decision: raise the budget or lower the effort · nothing was paid for",
	} {
		if !strings.Contains(stdout, line+"\n") {
			t.Errorf("the session did not print %q:\n%s", line, stdout)
		}
	}

	first := string(mustRead(t, filepath.Join(folder, "laps", "0001.lap.md")))
	for _, want := range []string{
		"# Lap 0001: pit stop at the gate",
		"**Choice:** Claude Opus 5.5 at high effort, chosen by owner over the suggestion",
		"**Outcome:** paused at gate FG-03, nothing paid",
		"| Estimated | 13:46:50 | $34.00 for 6M input and 500K output tokens |",
		"| Paused | 13:46:50 | FG-03: a $40.80 reservation exceeds the $25.00 budget |",
		"No runtime revision was materialized and no execution was claimed, so nothing was observed.",
	} {
		if !strings.Contains(first, want) {
			t.Errorf("lap 0001 lacks %q:\n%s", want, first)
		}
	}

	var fence run.LapFence
	decodeStrict(t, []byte(fenceBody(t, filepath.Join(folder, "laps", "0001.lap.md"), "langdrift-lap")), &fence)
	if fence.Lap != "0001" || fence.Outcome != "paused-at-gate" || fence.Gate == nil || *fence.Gate != "FG-03" || fence.Execution != nil || fence.Observed != nil || fence.Choice.Model != "claude-opus-5-5" {
		t.Errorf("lap 0001 fence = %+v", fence)
	}

	var planned run.Plan
	decodeStrict(t, []byte(fenceBody(t, filepath.Join(folder, "runtime.md"), "langdrift-plan")), &planned)
	if planned.State != "authorized" || planned.Runtime.State != "authorized" || planned.Authorization == nil || planned.Authorization.By != "owner" || planned.Suggestion.AuthorizedBy != "owner" {
		t.Errorf("running did not authorize the plan: %+v", planned.Authorization)
	}

	var ledger finops.Ledger
	decodeStrict(t, mustRead(t, filepath.Join(p.root, ".langdrift", "finops", "ledger", "console-microruntime.ledger.json")), &ledger)
	if len(ledger.Laps) != 1 {
		t.Fatalf("the ledger has %d laps", len(ledger.Laps))
	}
	line := ledger.Laps[0]
	if line.Cents.Estimate != 3400 || line.Cents.RequestedReservation != 4080 || line.Cents.Reserved != 0 || !line.NothingPaid || line.Cents.Billed == nil || *line.Cents.Billed != 0 || line.Cents.Variance != nil || len(line.Gates) != 3 || line.Gates[2].ID != "FG-03" || line.Gates[2].Result != "paused" {
		t.Errorf("lap 0001 line = %+v", line)
	}
	if ledger.Totals.Cents.Reserved != 0 || ledger.Totals.PeriodLeftAfterCents != 50000 {
		t.Errorf("totals = %+v", ledger.Totals)
	}

	var record run.Run
	decodeStrict(t, mustRead(t, filepath.Join(folder, "run.json")), &record)
	if record.State != "paused" || len(record.Laps) != 1 || record.Laps[0].Gate == nil || *record.Laps[0].Gate != "FG-03" || !record.Laps[0].NothingPaid || record.Laps[0].Choice.Overrode == nil {
		t.Errorf("run record = %+v", record)
	}

	status, stdout, stderr = p.exWith(deps, "agent", "run", "console-microruntime")
	if status != 0 {
		t.Fatalf("the second lap gave %d: %s\n%s", status, stderr, stdout)
	}

	for _, line := range []string{
		"agent › intent     claude-sonnet-5-5 · effort medium · runtime.md revision 1",
		"agent › estimate   $9.50 · 3.5M input · 250K output · list prices · FIN-02",
		"agent › gate       reserved $11.40 of $25.00 · opened",
	} {
		if !strings.Contains(stdout, line+"\n") {
			t.Errorf("the session did not print %q:\n%s", line, stdout)
		}
	}
	// The runtime comes up before the first stage, which waits for the owner.
	ready, waiting := strings.Index(stdout, "runtime    console-microruntime"), strings.Index(stdout, "agent › waiting    planning")
	if ready < 0 || waiting < 0 || ready > waiting || lifecycle.ups != 1 {
		t.Errorf("the runtime was not up before the first wait (%d ups):\n%s", lifecycle.ups, stdout)
	}

	decodeStrict(t, mustRead(t, filepath.Join(folder, "run.json")), &record)
	if record.State != "waiting" || len(record.Laps) != 2 || record.Laps[1].Outcome != "open" || record.Laps[1].ReservedCents != 1140 || record.Laps[1].Choice.Overrode != nil || record.Window.Started == nil {
		t.Errorf("run record after the second lap = %+v", record)
	}

	decodeStrict(t, mustRead(t, filepath.Join(p.root, ".langdrift", "finops", "ledger", "console-microruntime.ledger.json")), &ledger)
	if len(ledger.Laps) != 2 || ledger.Laps[1].Cents.Reserved != 1140 || ledger.Laps[1].Cents.Billed != nil || ledger.Totals.Cents.Reserved != 1140 || ledger.Totals.Cents.Billed != nil {
		t.Errorf("ledger after the second lap = %+v", ledger)
	}

	second := string(mustRead(t, filepath.Join(folder, "laps", "0002.lap.md")))
	if !strings.Contains(second, "**Choice:** Claude Sonnet 5.5 at medium effort, the suggestion") || !strings.Contains(second, "| Reserved | 13:46:50 | $11.40 for an estimate of $9.50, FG-03 passed |") {
		t.Errorf("lap 0002:\n%s", second)
	}

	// The status shows the runtime up, by its execution, and the stage waiting.
	if status, stdout, stderr := p.exWith(deps, "agent", "status", "console-microruntime"); status != 0 || !strings.Contains(stdout, `"number": "0002"`) || !strings.Contains(stdout, `"state": "waiting"`) || !strings.Contains(stdout, `"execution": "exec-`) || stageState(runRecord(t, p), "planning").State != "waiting" || lifecycle.downs != 0 {
		t.Errorf("ex agent status = %d\n%s%s", status, stdout, stderr)
	}

	// With a lap open, a run plays on with it: the first stage waits for the owner.
	status, stdout, stderr = p.exWith(deps, "agent", "run", "console-microruntime")
	if status != 0 || !strings.Contains(stdout, "agent › lap        0002 open · $11.40 reserved · resumed\n") || !strings.Contains(stdout, "agent › waiting    planning · the owner confirms\n") {
		t.Errorf("a run over the open lap gave %d: %s\n%s", status, stderr, stdout)
	}

	if status, _, stderr := p.exWith(deps, "agent", "run", "console-microruntime", "--model", "claude-opus-5-5"); status != 3 || !strings.Contains(stderr, "is open") {
		t.Errorf("a new choice over the open lap gave %d: %s", status, stderr)
	}
}

func TestRunRefusesWhatItCannotRun(t *testing.T) {
	t.Parallel()

	p, deps := seedProject(t)

	if status, _, stderr := p.exWith(deps, "agent", "run", "console-microruntime"); status != 4 || !strings.Contains(stderr, "runtime.md") {
		t.Errorf("an unplanned runtime gave %d: %s", status, stderr)
	}

	if status, _, stderr := p.exWith(deps, "agent", "run", "console-microruntime", "--effort", "max"); status != 2 {
		t.Errorf("an unknown effort gave %d: %s", status, stderr)
	}

	// A model the price list does not cover blocks at FG-01, nothing paid.
	if status, stdout, stderr := p.exWith(deps, "agent", "plan", ".agents/workflow/tickets/0070-console-microruntime.ticket.md", "--no-judge", "--component", consoleComponent); status != 0 {
		t.Fatalf("ex agent plan = %d\n%s%s", status, stdout, stderr)
	}
	status, stdout, stderr := p.exWith(deps, "agent", "run", "console-microruntime", "--model", "claude-mythos-5-1")
	if status != 3 || !strings.Contains(stderr, "FG-01") || !strings.Contains(stdout, "agent › gate       blocked · FG-01") {
		t.Errorf("an unpriced model gave %d: %s\n%s", status, stderr, stdout)
	}
}
