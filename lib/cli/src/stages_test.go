package main

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/jonatassales/tao/lib/cli/src/application/execution"
	"github.com/jonatassales/tao/lib/cli/src/domain/finops"
	"github.com/jonatassales/tao/lib/cli/src/domain/harness"
	"github.com/jonatassales/tao/lib/cli/src/domain/run"
)

// websiteUsage is what each skill used in the website's session: six calls
// that come to 3.7M input and 240K output tokens, $9.80 at Sonnet's prices.
var websiteUsage = map[string][2]int64{
	"grilling":    {300_000, 20_000},
	"to-spec":     {400_000, 30_000},
	"to-tickets":  {300_000, 20_000},
	"implement":   {1_800_000, 110_000},
	"tdd":         {500_000, 40_000},
	"code-review": {400_000, 20_000},
}

// websiteSummaries are the website's stage lines, one per skill.
var websiteSummaries = map[string]string{
	"grilling":    "3 open questions settled",
	"to-spec":     "0011-SPEC revised · 3 acceptance criteria",
	"to-tickets":  "3 tickets · blocking edges recorded",
	"implement":   "tdd · 5 components · 14 files · 32 tests green",
	"tdd":         "32 tests green · nothing left to refactor",
	"code-review": "2 findings · fixed",
}

// scriptedDriver plays the stages without a coding agent: each skill's
// summary and usage are scripted, and every call is kept.
type scriptedDriver struct {
	usage      map[string][2]int64
	unmeasured map[string]bool
	exits      map[string]int
	cutOff     map[string]bool
	hookExit   int
	calls      []execution.StageCall
	hooks      []execution.HookCall
}

func websiteDriver() *scriptedDriver {
	usage := map[string][2]int64{}
	for skill, tokens := range websiteUsage {
		usage[skill] = tokens
	}

	return &scriptedDriver{usage: usage, unmeasured: map[string]bool{}, exits: map[string]int{}, cutOff: map[string]bool{}}
}

func (driver *scriptedDriver) Run(_ context.Context, _ harness.Driver, call execution.StageCall) (execution.StageResult, error) {
	driver.calls = append(driver.calls, call)

	result := execution.StageResult{ExitCode: driver.exits[call.Skill], Summary: websiteSummaries[call.Skill], LogPath: ".langdrift/agent/state/runs/" + call.Runtime + "/stages/" + call.Lap + "-" + call.Stage + "-" + call.Skill + ".log", ElapsedMs: 1200}
	if driver.cutOff[call.Skill] {
		return execution.StageResult{ExitCode: -1}, errors.New("claude was interrupted; what it used is unknown")
	}
	if driver.unmeasured[call.Skill] {
		return result, nil
	}

	tokens := driver.usage[call.Skill]
	result.InputTokens, result.OutputTokens = &tokens[0], &tokens[1]

	return result, nil
}

func (driver *scriptedDriver) Hook(_ context.Context, call execution.HookCall) (int, error) {
	driver.hooks = append(driver.hooks, call)

	return driver.hookExit, nil
}

// driverFile writes the driver input: the claude CLI in print mode, its JSON
// result read for the usage and the summary.
func driverFile(t *testing.T, root, unmeasured string) string {
	t.Helper()

	path := filepath.Join(root, "driver-"+unmeasured+".json")
	body := `{"executable":"claude","args":["-p","{{prompt}}","--output-format","json"],"passEnv":["ANTHROPIC_API_KEY"],"usage":{"inputTokens":"usage.input_tokens","outputTokens":"usage.output_tokens"},"summary":"result","unmeasured":"` + unmeasured + `"}`
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	return path
}

func runtimeFolder(p project) string {
	return filepath.Join(p.root, ".langdrift", "microruntimes", "runtimes", "console-microruntime")
}

func runRecord(t *testing.T, p project) run.Run {
	t.Helper()

	var record run.Run
	decodeStrict(t, mustRead(t, filepath.Join(runtimeFolder(p), "run.json")), &record)

	return record
}

func ledgerOf(t *testing.T, p project) finops.Ledger {
	t.Helper()

	var ledger finops.Ledger
	decodeStrict(t, mustRead(t, filepath.Join(p.root, ".langdrift", "finops", "ledger", "console-microruntime.ledger.json")), &ledger)

	return ledger
}

func stageState(record run.Run, id string) run.StageState {
	for _, state := range record.Stages {
		if state.ID == id {
			return state
		}
	}

	return run.StageState{}
}

func expectLines(t *testing.T, stdout string, lines ...string) {
	t.Helper()

	for _, line := range lines {
		if !strings.Contains(stdout, line+"\n") {
			t.Errorf("the session did not print %q:\n%s", line, stdout)
		}
	}
}

// The run's second act: the first stage waits for the owner, the owner
// confirms, and the next run plays every stage through the coding agent,
// each call under FG-04, to the website's completed line.
func TestRunWaitsForTheOwnerThenPlaysTheStages(t *testing.T) {
	t.Parallel()

	p, deps := plannedProject(t)
	driver := websiteDriver()
	deps.Driver, deps.Lifecycle, deps.Delivery, deps.Poll = driver, websiteLifecycle(), websiteDelivery(), time.Millisecond

	status, stdout, stderr := p.dftWith(deps, "agent", "run", "console-microruntime")
	if status != 0 {
		t.Fatalf("dft agent run = %d: %s\n%s", status, stderr, stdout)
	}
	expectLines(t, stdout, "console › gate       reserved $11.40 of $25.00 · opened", "console › waiting    planning · the owner confirms")
	if len(driver.calls) != 0 {
		t.Errorf("the run called the driver before the owner confirmed: %+v", driver.calls)
	}

	record := runRecord(t, p)
	if planning := stageState(record, "planning"); record.State != "waiting" || planning.State != "waiting" || !planning.Human {
		t.Errorf("run record while waiting = %+v", record)
	}

	status, stdout, _ = p.dftWith(deps, "agent", "run", "console-microruntime")
	if status != 0 {
		t.Fatalf("a run over the waiting lap gave %d:\n%s", status, stdout)
	}
	expectLines(t, stdout, "console › lap        0001 open · $11.40 reserved · resumed", "console › waiting    planning · the owner confirms")

	status, stdout, stderr = p.dftWith(deps, "agent", "confirm", "console-microruntime", "planning")
	if status != 0 || !strings.Contains(stdout, "✓ confirmed  planning · by owner\n") {
		t.Fatalf("dft agent confirm = %d: %s\n%s", status, stderr, stdout)
	}
	if planning := stageState(runRecord(t, p), "planning"); planning.State != "active" || planning.ConfirmedBy == nil || *planning.ConfirmedBy != "owner" {
		t.Errorf("planning after the confirmation = %+v", planning)
	}
	if status, _, stderr := p.dftWith(deps, "agent", "confirm", "console-microruntime", "planning"); status != 3 || !strings.Contains(stderr, "not waiting") {
		t.Errorf("a second confirmation gave %d: %s", status, stderr)
	}
	if status, _, stderr := p.dftWith(deps, "agent", "confirm", "console-microruntime", "deploy"); status != 2 {
		t.Errorf("confirming an unknown stage gave %d: %s", status, stderr)
	}

	status, stdout, stderr = p.dftWith(deps, "agent", "run", "console-microruntime")
	if status != 0 {
		t.Fatalf("the resumed run gave %d: %s\n%s", status, stderr, stdout)
	}
	expectLines(t, stdout,
		"console › lap        0001 open · $11.40 reserved · resumed",
		"console › grilling   3 open questions settled",
		"console › to-spec    0011-SPEC revised · 3 acceptance criteria",
		"console › to-tickets 3 tickets · blocking edges recorded",
		"console › implement  tdd · 5 components · 14 files · 32 tests green",
		"console › tdd        32 tests green · nothing left to refactor",
		"console › code-review 2 findings · fixed",
		"✓ completed  implementation · $9.80 observed of $11.40 reserved · usage record",
	)
	if strings.Contains(stdout, "waiting") || strings.Contains(stdout, "confirmed") {
		t.Errorf("the resumed run asked the owner again:\n%s", stdout)
	}

	var calls []string
	for _, call := range driver.calls {
		calls = append(calls, call.Stage+"/"+call.Skill)
	}
	if want := []string{"planning/grilling", "spec/to-spec", "tickets/to-tickets", "implement/implement", "implement/tdd", "implement/code-review"}; !slices.Equal(calls, want) {
		t.Errorf("the driver ran %v, want %v", calls, want)
	}
	if first := driver.calls[0]; first.ProjectRoot != p.root || first.Runtime != "console-microruntime" || first.Lap != "0001" || first.Ticket != ".agents/workflow/tickets/0070-console-microruntime.ticket.md" {
		t.Errorf("the first call = %+v", first)
	}

	record = runRecord(t, p)
	if record.State != "finished" || len(record.Laps) != 1 || record.Laps[0].Outcome != "completed" || record.Laps[0].ObservedCents != 980 {
		t.Errorf("run record after the stages = %+v", record)
	}
	for _, id := range []string{"planning", "spec", "tickets", "implement"} {
		if state := stageState(record, id); state.State != "done" || state.At == nil || state.Summary == nil {
			t.Errorf("stage %s = %+v, want done with its summary", id, state)
		}
	}
	if stageState(record, "review").State != "waiting" || stageState(record, "merge").State != "pending" {
		t.Errorf("the stages the run never plays = %+v, want review waiting and merge pending", record.Stages)
	}
	implement := stageState(record, "implement")
	if *implement.Summary != "2 findings · fixed" || implement.Usage == nil || implement.Usage.InputTokens != 2_700_000 || implement.Usage.OutputTokens != 170_000 || implement.Usage.Cents != 710 {
		t.Errorf("implement = %+v, usage %+v", implement, implement.Usage)
	}

	ledger := ledgerOf(t, p)
	line := ledger.Laps[0]
	if line.Cents.Observed != 980 || line.Cents.Reserved != 1140 || line.Cents.Billed != nil || line.Cents.Recovered != 160 || line.Cents.Variance == nil || *line.Cents.Variance != 30 || line.NothingPaid {
		t.Errorf("lap 0001 line = %+v", line)
	}
	if line.TokensThousands.Observed == nil || *line.TokensThousands.Observed != (finops.TokensThousands{Input: 3700, Output: 240}) {
		t.Errorf("lap 0001 observed tokens = %+v", line.TokensThousands.Observed)
	}
	if gate := line.Gates[len(line.Gates)-2]; gate.ID != "FG-04" || gate.Result != "not-reached" {
		t.Errorf("the gate before the last = %+v, want FG-04 not reached", gate)
	}
	if ledger.Totals.Cents.Observed != 990 || ledger.Totals.PeriodLeftAfterCents != 49010 || ledger.Totals.Cents.Billed != nil {
		t.Errorf("totals = %+v", ledger.Totals)
	}

	lap := string(mustRead(t, filepath.Join(runtimeFolder(p), "laps", "0001.lap.md")))
	for _, want := range []string{
		"**Outcome:** completed, $9.80 observed of $11.40 reserved, $1.60 released",
		"| Confirmed | 13:46:50 | planning · by owner |",
		"| Stage | 13:46:50 | planning · grilling · 3 open questions settled |",
		"| Stage | 13:46:50 | implement · code-review · 2 findings · fixed |",
		"| Completed | 13:46:50 | $9.80 observed of $11.40 reserved |",
	} {
		if !strings.Contains(lap, want) {
			t.Errorf("lap 0001 lacks %q:\n%s", want, lap)
		}
	}

	var evidence execution.StageEvidence
	decodeStrict(t, mustRead(t, filepath.Join(runtimeFolder(p), "evidence", "laps", "0001", "stages", "planning-grilling.json")), &evidence)
	if evidence.Stage != "planning" || evidence.Skill != "grilling" || evidence.InputTokens == nil || *evidence.InputTokens != 300_000 || evidence.Summary != "3 open questions settled" || evidence.Log != ".langdrift/agent/state/runs/console-microruntime/stages/0001-planning-grilling.log" {
		t.Errorf("evidence = %+v", evidence)
	}
	if entries, err := os.ReadDir(filepath.Join(runtimeFolder(p), "evidence", "laps", "0001", "stages")); err != nil || len(entries) != 6 {
		t.Errorf("stage evidence = %d files, %v; want one per call", len(entries), err)
	}

	if status, _, stderr := p.dftWith(deps, "agent", "run", "console-microruntime"); status != 3 || !strings.Contains(stderr, "finished") || len(driver.calls) != 6 {
		t.Errorf("a run after the run finished gave %d with %d calls: %s", status, len(driver.calls), stderr)
	}

	if status, stdout, _ := p.dftWith(deps, "agent", "status", "console-microruntime"); status != 0 || !strings.Contains(stdout, `"observedCents": 980`) {
		t.Errorf("dft agent status = %d:\n%s", status, stdout)
	}
}

// At a terminal the run asks for each confirmation itself.
func TestRunConfirmsOnTheTerminal(t *testing.T) {
	t.Parallel()

	p, deps := plannedProject(t)
	deps.Driver, deps.Lifecycle, deps.Delivery, deps.Poll = websiteDriver(), websiteLifecycle(), websiteDelivery(), time.Millisecond

	status, stdout, stderr := p.dftInput(deps, "n\n", "agent", "run", "console-microruntime", "--confirm")
	if status != 0 || !strings.Contains(stderr, "Confirm planning? [y/N]: ") || strings.Contains(stdout, "confirmed") {
		t.Fatalf("a declined confirmation gave %d: %s\n%s", status, stderr, stdout)
	}
	expectLines(t, stdout, "console › waiting    planning · the owner confirms")

	status, stdout, stderr = p.dftInput(deps, "y\n", "agent", "run", "console-microruntime", "--confirm", "--by", "jonatas")
	if status != 0 {
		t.Fatalf("a confirmed run gave %d: %s\n%s", status, stderr, stdout)
	}
	expectLines(t, stdout,
		"console › waiting    planning · the owner confirms",
		"✓ confirmed  planning · by jonatas",
		"console › grilling   3 open questions settled",
		"✓ completed  implementation · $9.80 observed of $11.40 reserved · usage record",
	)
}

// FG-04 pauses the paid calls once what the lap observed reached its
// reservation: the lap ends with its overspend on the ledger, and the next
// lap needs the owner to acknowledge it, FG-05, before it plays on.
func TestRunPausesThePaidCallsAtTheReservation(t *testing.T) {
	t.Parallel()

	p, deps := plannedProject(t)
	driver := websiteDriver()
	driver.usage["grilling"] = [2]int64{5_000_000, 400_000}
	deps.Driver, deps.Lifecycle, deps.Delivery, deps.Poll = driver, websiteLifecycle(), websiteDelivery(), time.Millisecond

	if status, stdout, stderr := p.dftWith(deps, "agent", "run", "console-microruntime"); status != 0 {
		t.Fatalf("dft agent run = %d: %s\n%s", status, stderr, stdout)
	}
	if status, stdout, stderr := p.dftWith(deps, "agent", "confirm", "console-microruntime", "planning"); status != 0 {
		t.Fatalf("dft agent confirm = %d: %s\n%s", status, stderr, stdout)
	}

	status, stdout, stderr := p.dftWith(deps, "agent", "run", "console-microruntime")
	if status != 3 || !strings.Contains(stderr, "FG-04") || len(driver.calls) != 1 {
		t.Fatalf("the exhausted lap gave %d with %d calls: %s\n%s", status, len(driver.calls), stderr, stdout)
	}
	expectLines(t, stdout, "console › grilling   3 open questions settled", "console › gate       paused · FG-04 · observed 1400 reached the reservation 1140")

	record := runRecord(t, p)
	if lap := record.Laps[0]; record.State != "paused" || lap.Outcome != "paused-at-gate" || lap.Gate == nil || *lap.Gate != "FG-04" || lap.NothingPaid || lap.ObservedCents != 1400 {
		t.Errorf("run record after FG-04 = %+v", record)
	}
	if stageState(record, "planning").State != "done" || stageState(record, "spec").State != "pending" {
		t.Errorf("stages after FG-04 = %+v", record.Stages)
	}

	ledger := ledgerOf(t, p)
	if line := ledger.Laps[0]; line.Cents.Observed != 1400 || line.Cents.Overspend != 260 || line.Cents.Recovered != 0 || line.Cents.Variance == nil || *line.Cents.Variance != 450 || line.NothingPaid || line.Cents.Billed != nil {
		t.Errorf("lap 0001 line = %+v", line)
	}
	if ledger.Totals.Cents.Overspend != 260 || ledger.Totals.PeriodLeftAfterCents != 48590 {
		t.Errorf("totals = %+v", ledger.Totals)
	}

	lap := string(mustRead(t, filepath.Join(runtimeFolder(p), "laps", "0001.lap.md")))
	for _, want := range []string{
		"# Lap 0001: pit stop at the gate",
		"**Outcome:** paused at gate FG-04, $14.00 observed of $11.40 reserved",
		"| Paused | 13:46:50 | FG-04: observed 1400 reached the reservation 1140 |",
	} {
		if !strings.Contains(lap, want) {
			t.Errorf("lap 0001 lacks %q:\n%s", want, lap)
		}
	}

	status, stdout, stderr = p.dftWith(deps, "agent", "run", "console-microruntime")
	if status != 3 || !strings.Contains(stderr, "FG-05") {
		t.Fatalf("a lap over an unacknowledged overspend gave %d: %s\n%s", status, stderr, stdout)
	}
	expectLines(t, stdout, "console › gate       blocked · FG-05 · FG-05 require-owner")

	status, stdout, stderr = p.dftWith(deps, "agent", "run", "console-microruntime", "--acknowledge-overspend")
	if status != 0 {
		t.Fatalf("the acknowledged lap gave %d: %s\n%s", status, stderr, stdout)
	}
	expectLines(t, stdout,
		"console › gate       reserved $11.40 of $25.00 · opened",
		"console › to-spec    0011-SPEC revised · 3 acceptance criteria",
		"✓ completed  implementation · $9.00 observed of $11.40 reserved · usage record",
	)
	if strings.Contains(stdout, "grilling") || strings.Contains(stdout, "waiting") {
		t.Errorf("the new lap played a stage that was done:\n%s", stdout)
	}

	record = runRecord(t, p)
	if len(record.Laps) != 3 || record.Laps[1].Outcome != "blocked" || record.Laps[2].Outcome != "completed" || record.Laps[2].ObservedCents != 900 {
		t.Errorf("laps = %+v", record.Laps)
	}
	if ledger := ledgerOf(t, p); ledger.Totals.Cents.Observed != 2310 || ledger.Totals.Cents.Reserved != 2295 || ledger.Totals.PeriodLeftAfterCents != 47690 {
		t.Errorf("totals after the laps = %+v", ledger.Totals)
	}
}

// A call that reports no usage leaves the spend unknown, which is never
// zero: FG-04 pauses the next paid call unless the driver allows it, and the
// ledger shows no tokens and no variance for the lap.
func TestRunPausesWhenACallReportsNoUsage(t *testing.T) {
	t.Parallel()

	p, deps := plannedProject(t)
	driver := websiteDriver()
	driver.unmeasured["grilling"] = true
	deps.Driver, deps.Lifecycle, deps.Delivery, deps.Poll = driver, websiteLifecycle(), websiteDelivery(), time.Millisecond

	if status, stdout, stderr := p.dftWith(deps, "agent", "run", "console-microruntime"); status != 0 {
		t.Fatalf("dft agent run = %d: %s\n%s", status, stderr, stdout)
	}
	if status, stdout, stderr := p.dftWith(deps, "agent", "confirm", "console-microruntime", "planning"); status != 0 {
		t.Fatalf("dft agent confirm = %d: %s\n%s", status, stderr, stdout)
	}

	status, stdout, stderr := p.dftWith(deps, "agent", "run", "console-microruntime")
	if status != 3 || !strings.Contains(stderr, "FG-04") {
		t.Fatalf("an unmeasured call gave %d: %s\n%s", status, stderr, stdout)
	}
	expectLines(t, stdout,
		"console › grilling   3 open questions settled · no usage reported",
		"console › gate       paused · FG-04 · observed spend is unknown; a call reported no usage",
	)

	ledger := ledgerOf(t, p)
	if line := ledger.Laps[0]; line.TokensThousands.Observed != nil || line.Cents.Variance != nil || line.Cents.Observed != 0 || line.Cents.Recovered != 0 || line.Cents.Overspend != 0 || line.NothingPaid {
		t.Errorf("an unmeasured lap's line = %+v", line)
	}
	if lap := string(mustRead(t, filepath.Join(runtimeFolder(p), "laps", "0001.lap.md"))); !strings.Contains(lap, "A call reported no usage, so the observed spend is a floor, not a measurement.") {
		t.Errorf("lap 0001 does not say its spend is a floor:\n%s", lap)
	}

	if status, stdout, stderr := p.dftWith(deps, "harness", "driver", "set", "--file", driverFile(t, p.root, "allow")); status != 0 {
		t.Fatalf("dft harness driver set = %d: %s\n%s", status, stderr, stdout)
	}
	driver.unmeasured["to-spec"] = true

	status, stdout, stderr = p.dftWith(deps, "agent", "run", "console-microruntime")
	if status != 0 {
		t.Fatalf("the allowed lap gave %d: %s\n%s", status, stderr, stdout)
	}
	expectLines(t, stdout,
		"console › to-spec    0011-SPEC revised · 3 acceptance criteria · no usage reported",
		"console › to-tickets 3 tickets · blocking edges recorded",
		"✓ completed  implementation · $7.90 observed of $11.40 reserved · usage record",
	)
	if line := ledgerOf(t, p).Laps[1]; line.TokensThousands.Observed != nil || line.Cents.Observed != 790 {
		t.Errorf("the allowed lap's line = %+v", line)
	}
}

// A stage's committed hooks run first; a hook or a skill that fails stops
// the run with the stage and its lap failed, and the next run opens a new
// lap that plays the stage again.
func TestRunRunsTheHooksAndReplaysAFailedStage(t *testing.T) {
	t.Parallel()

	p, deps := seedProject(t)
	workflow := bodyFile(t, p.root, "workflow-hooks.json", "harness/workflow.json", func(body *harness.Workflow) {
		withoutMirror(body)
		body.Stages[3].Hooks = []string{"platform/hooks/lint.sh"}
	})
	if status, stdout, stderr := p.dftWith(deps, "harness", "workflow", "set", "--file", workflow); status != 0 {
		t.Fatalf("dft harness workflow set = %d: %s\n%s", status, stderr, stdout)
	}
	if status, stdout, stderr := p.dftWith(deps, "agent", "plan", ".agents/workflow/tickets/0070-console-microruntime.ticket.md", "--component", consoleComponent); status != 0 {
		t.Fatalf("dft agent plan = %d\n%s%s", status, stdout, stderr)
	}

	driver := websiteDriver()
	driver.exits["to-spec"] = 1
	deps.Driver, deps.Lifecycle, deps.Delivery, deps.Poll = driver, websiteLifecycle(), websiteDelivery(), time.Millisecond

	if status, stdout, stderr := p.dftInput(deps, "y\n", "agent", "run", "console-microruntime", "--confirm"); status != 1 || !strings.Contains(stderr, "to-spec of stage spec exited 1") {
		t.Fatalf("a failing skill gave %d: %s\n%s", status, stderr, stdout)
	}
	record := runRecord(t, p)
	if record.State != "failed" || stageState(record, "spec").State != "failed" || len(driver.calls) != 2 {
		t.Errorf("after the failed skill: state %s, stages %+v, %d calls", record.State, record.Stages, len(driver.calls))
	}

	driver.exits = map[string]int{}
	driver.hookExit = 1
	status, stdout, stderr := p.dftWith(deps, "agent", "run", "console-microruntime")
	if status != 1 || !strings.Contains(stderr, "hook platform/hooks/lint.sh of stage implement exited 1") {
		t.Fatalf("a failing hook gave %d: %s\n%s", status, stderr, stdout)
	}
	expectLines(t, stdout, "console › to-spec    0011-SPEC revised · 3 acceptance criteria", "console › hook       platform/hooks/lint.sh · exit 1")
	if record := runRecord(t, p); stageState(record, "implement").State != "failed" || stageState(record, "tickets").State != "done" || record.Laps[1].Number != "0002" || record.Laps[1].Outcome != "failed" {
		t.Errorf("after the failed hook: %+v, laps %+v", record.Stages, record.Laps)
	}

	driver.hookExit = 0
	status, stdout, stderr = p.dftWith(deps, "agent", "run", "console-microruntime")
	if status != 0 {
		t.Fatalf("the replayed stage gave %d: %s\n%s", status, stderr, stdout)
	}
	expectLines(t, stdout,
		"console › hook       platform/hooks/lint.sh · exit 0",
		"console › implement  tdd · 5 components · 14 files · 32 tests green",
		"console › gate       reserved $11.40 of $25.00 · opened",
		"✓ completed  implementation · $7.10 observed of $11.40 reserved · usage record",
	)
	if record := runRecord(t, p); len(record.Laps) != 3 || record.Laps[0].Outcome != "failed" || record.Laps[2].Number != "0003" {
		t.Errorf("the replay did not run in lap 0003 after two failed laps: %+v", record.Laps)
	}
	if len(driver.hooks) != 2 || driver.hooks[1].Stage != "implement" || driver.hooks[1].Script != "platform/hooks/lint.sh" || driver.hooks[1].ProjectRoot != p.root {
		t.Errorf("hooks = %+v", driver.hooks)
	}
}
