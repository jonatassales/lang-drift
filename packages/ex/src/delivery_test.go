package main

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jonatassales/tao/packages/contracts/configuration"
	"github.com/jonatassales/tao/packages/ex/src/application/execution"
	"github.com/jonatassales/tao/packages/ex/src/domain/finops"
	"github.com/jonatassales/tao/packages/ex/src/domain/microruntime"
	projectdomain "github.com/jonatassales/tao/packages/ex/src/domain/project"
	"github.com/jonatassales/tao/packages/ex/src/domain/run"
)

// scriptedLifecycle brings the runtime up and down without a cluster: one
// container in one Deployment, 120m steady then 410m and 690 MiB at peak.
type scriptedLifecycle struct {
	execution    string
	cleanup      execution.Cleanup
	failUp       error
	ups          int
	downs        int
	observations int
}

func websiteLifecycle() *scriptedLifecycle {
	return &scriptedLifecycle{execution: "exec-5c1e9a", cleanup: execution.Cleanup{Confirmed: true, Removed: 4, Total: 4, ClusterKept: true}}
}

func (lifecycle *scriptedLifecycle) Up(_ context.Context, _ projectdomain.Project, _ microruntime.Identifier, config configuration.Configuration, _ []byte, report execution.Reporter) (execution.Execution, error) {
	lifecycle.ups++
	if lifecycle.failUp != nil {
		return execution.Execution{}, lifecycle.failUp
	}

	report("Configured", "`runtime.md` revision 1 materialized")
	report("Built", "images built from the components' recipes")

	return execution.Execution{ID: lifecycle.execution}, nil
}

func (lifecycle *scriptedLifecycle) Observe(context.Context, projectdomain.Project, microruntime.Identifier) (execution.Observation, error) {
	lifecycle.observations++
	cpu, memory := 410, 690
	if lifecycle.observations == 1 {
		cpu = 120
	}

	return execution.Observation{Ready: true, Containers: 1, ReadyContainers: 1, Deployments: 1, CPUMillis: &cpu, MemoryMiB: &memory, Window: "30s", ObservedAt: time.Date(2026, 10, 28, 13, 46, 20, 0, time.UTC)}, nil
}

func (lifecycle *scriptedLifecycle) Down(context.Context, projectdomain.Project, microruntime.Identifier) (execution.Cleanup, error) {
	lifecycle.downs++

	return lifecycle.cleanup, nil
}

// scriptedDelivery is what the harness delivered, as the website shows it.
type scriptedDelivery struct {
	branch       execution.Branch
	pullRequest  *run.PullRequest
	reason       string
	checks       [][]run.CheckResult
	checksReason string
	checkCalls   int
}

func websiteDelivery() *scriptedDelivery {
	pending := []run.CheckResult{{Name: "build", Result: "pending"}, {Name: "check", Result: "pending"}, {Name: "test", Result: "pending"}, {Name: "guardrails", Result: "pending"}}
	passed := []run.CheckResult{{Name: "build", Result: "passed"}, {Name: "check", Result: "passed"}, {Name: "test", Result: "passed"}, {Name: "guardrails", Result: "passed"}}

	return &scriptedDelivery{
		branch:      execution.Branch{Name: "feat/0070-console-microruntime", Commit: "3f9c2e1", Upstream: true},
		pullRequest: &run.PullRequest{Number: 19, Title: "feat(console): run the console preview as a microruntime", State: "open"},
		checks:      [][]run.CheckResult{pending, passed},
	}
}

func (delivery *scriptedDelivery) Branch(context.Context, string) execution.Branch {
	return delivery.branch
}

func (delivery *scriptedDelivery) PullRequest(context.Context, string, string) (*run.PullRequest, string) {
	return delivery.pullRequest, delivery.reason
}

func (delivery *scriptedDelivery) Checks(context.Context, string, int) ([]run.CheckResult, string) {
	if delivery.checksReason != "" {
		return nil, delivery.checksReason
	}

	checks := delivery.checks[min(delivery.checkCalls, len(delivery.checks)-1)]
	delivery.checkCalls++

	return checks, ""
}

// The run's third act: after the stages, the runtime comes up from its
// fence, the implementation is complete, delivery is observed from the
// repository and the forge, the runtime goes down, the unused reservation
// returns to the budget, and the run remembers what it measured.
func TestRunBringsTheRuntimeUpObservesDeliveryAndRemembers(t *testing.T) {
	t.Parallel()

	p, deps := plannedProject(t)
	lifecycle, delivery := websiteLifecycle(), websiteDelivery()
	deps.Driver, deps.Lifecycle, deps.Delivery, deps.Poll = websiteDriver(), lifecycle, delivery, time.Millisecond

	status, stdout, stderr := p.exInput(deps, "y\n", "agent", "run", "console-microruntime", "--confirm")
	if status != 0 {
		t.Fatalf("ex agent run = %d: %s\n%s", status, stderr, stdout)
	}
	expectLines(t, stdout,
		"agent › code-review 2 findings · fixed",
		"runtime    console-microruntime · revision 1 · 1 Deployments ready · 500m · 1024 MiB",
		"✓ completed  implementation · $9.80 observed of $11.40 reserved · usage record",
		"✓ pushed     feat/0070-console-microruntime · 3f9c2e1",
		"✓ pr #19     feat(console): run the console preview as a microruntime",
		"checks     build · check · test · guardrails · running",
		"✓ checks     4 of 4 passed",
		"runtime    stopped · 1 Deployments removed · nothing left",
		"✓ ready      #19 for review · $1.60 released to the envelope",
		"# next: the owner reviews #19 and decides whether to merge",
	)
	if lifecycle.ups != 1 || lifecycle.downs != 1 || lifecycle.observations != 2 {
		t.Errorf("the lifecycle saw %d ups, %d downs, %d observations", lifecycle.ups, lifecycle.downs, lifecycle.observations)
	}

	record := runRecord(t, p)
	if record.State != "finished" || record.Outcome == nil || *record.Outcome != "readyForReview" || record.Window.Finished == nil {
		t.Errorf("run record = state %s, outcome %v, window %+v", record.State, record.Outcome, record.Window)
	}
	if record.Delivery.Branch == nil || *record.Delivery.Branch != "feat/0070-console-microruntime" || record.Delivery.Commit == nil || *record.Delivery.Commit != "3f9c2e1" || record.Delivery.PullRequest == nil || record.Delivery.PullRequest.Number != 19 || len(record.Delivery.Checks) != 4 || record.Delivery.Checks[3].Result != "passed" {
		t.Errorf("delivery = %+v", record.Delivery)
	}
	if record.ProposedTicketStatus == nil || record.ProposedTicketStatus.Status != "in-review" || record.ProposedTicketStatus.Decided {
		t.Errorf("proposed ticket status = %+v", record.ProposedTicketStatus)
	}
	if stageState(record, "review").State != "waiting" || stageState(record, "merge").State != "pending" {
		t.Errorf("the stages after the run = %+v", record.Stages)
	}
	if lap := record.Laps[0]; lap.Outcome != "completed" || lap.Execution == nil || *lap.Execution != "exec-5c1e9a" || lap.ObservedCents != 980 {
		t.Errorf("lap summary = %+v", lap)
	}

	ledger := ledgerOf(t, p)
	line := ledger.Laps[0]
	if line.Cents.Observed != 980 || line.Cents.Recovered != 160 || line.Cents.Overspend != 0 || line.Cents.Variance == nil || *line.Cents.Variance != 30 || line.Cents.Billed != nil {
		t.Errorf("lap 0001 line = %+v", line.Cents)
	}
	if gate := line.Gates[len(line.Gates)-1]; gate.ID != "FG-06" || gate.Result != "keep-provisional" {
		t.Errorf("the last gate = %+v, want FG-06 keeping the billing provisional", gate)
	}
	if ledger.Totals.Cents.Reserved != 1140 || ledger.Totals.Cents.Observed != 980 || ledger.Totals.Cents.Recovered != 160 || ledger.Totals.PeriodLeftAfterCents != 49020 {
		t.Errorf("totals = %+v", ledger.Totals)
	}

	runtimeMarkdown := mustRead(t, filepath.Join(runtimeFolder(p), "runtime.md"))
	lap := string(mustRead(t, filepath.Join(runtimeFolder(p), "laps", "0001.lap.md")))
	for _, want := range []string{
		"# Lap 0001: console-microruntime ran",
		"**Outcome:** completed, $9.80 observed of $11.40 reserved, $1.60 released",
		"**Runtime:** `console-microruntime`, revision 1, `runtime.md` sha256 " + sha256Hex(runtimeMarkdown),
		"**Evidence:** `evidence/exec-5c1e9a/` in this runtime",
		"| Configured | 13:46:50 | `runtime.md` revision 1 materialized |",
		"| Claimed | 13:46:50 | execution `exec-5c1e9a` |",
		"| Ready | 13:46:50 | 1 of 1 containers ready on `/icon.svg` |",
		"| Delivered | 13:46:50 | feat/0070-console-microruntime at 3f9c2e1, pull request #19 |",
		"| Cleaned up | 13:46:50 | 4 of 4 owned resources removed; cluster kept |",
		"| Released | 13:46:50 | $1.60 of the reservation returned to the budget |",
		"| Resource | Proposed | Executed | Observed | Measured |",
		"| CPU | 500m | ",
		"| 410m peak, 120m steady | 30s windows, oldest sample 13:46:20 |",
		"| Memory | 1024 MiB | ",
		"| 690 MiB peak | 30s windows, oldest sample 13:46:20 |",
		"| Model spend | | $11.40 reserved | $9.80, billing pending | Ledger, lap 0001 |",
	} {
		if !strings.Contains(lap, want) {
			t.Errorf("lap 0001 lacks %q:\n%s", want, lap)
		}
	}

	var fence run.LapFence
	decodeStrict(t, []byte(fenceBody(t, filepath.Join(runtimeFolder(p), "laps", "0001.lap.md"), "langdrift-lap")), &fence)
	if fence.Outcome != "completed" || fence.Execution == nil || *fence.Execution != "exec-5c1e9a" || fence.Runtime == nil || fence.Runtime.Revision != 1 || fence.Runtime.MarkdownSHA256 != sha256Hex(runtimeMarkdown) {
		t.Errorf("lap fence = %+v", fence)
	}
	if fence.Observed == nil || fence.Observed.CPUMillisPeak == nil || *fence.Observed.CPUMillisPeak != 410 || *fence.Observed.MemoryMiBPeak != 690 || fence.Observed.TokensThousands == nil || *fence.Observed.TokensThousands != (finops.TokensThousands{Input: 3700, Output: 240}) {
		t.Errorf("lap fence observed = %+v", fence.Observed)
	}

	memory := string(mustRead(t, filepath.Join(runtimeFolder(p), "run.md")))
	for _, want := range []string{
		"# Run console-microruntime: ",
		"**Ticket:** 0070-console-microruntime.ticket.md",
		"**Outcome:** completed in one lap, ready for review",
		"**Proposed ticket status:** in-review, pending the owner",
		"Lap 0001 followed the suggestion, Claude Sonnet 5.5 at medium effort and finished with $9.80 observed of $11.40 reserved.",
		"## Worth remembering",
		"- Claude Sonnet 5.5 at medium effort used 3.7M input and 240K output tokens; the next estimate for this ticket can start from these instead of the effort table.",
		"- `/icon.svg` passed readiness for `.`.",
	} {
		if !strings.Contains(memory, want) {
			t.Errorf("run.md lacks %q:\n%s", want, memory)
		}
	}
	var episode run.RunFence
	decodeStrict(t, []byte(fenceBody(t, filepath.Join(runtimeFolder(p), "run.md"), "langdrift-run")), &episode)
	if episode.Run != "console-microruntime" || episode.Ticket == nil || *episode.Ticket != "0070-console-microruntime" || episode.Outcome == nil || *episode.Outcome != "completed" || episode.Laps != 1 || !episode.MeasuredResources || !episode.MeasuredTokens {
		t.Errorf("memory fence = %+v", episode)
	}
	if strings.Join(episode.ComponentPaths, ",") != "." || strings.Join(episode.Stacks, ",") != "nextjs,go" || episode.Observed == nil || *episode.Observed.MemoryMiBPeak != 690 {
		t.Errorf("memory fence components = %v, stacks %v, observed %+v", episode.ComponentPaths, episode.Stacks, episode.Observed)
	}
	for _, fact := range []string{"memory-peak-690-of-", "claude-sonnet-5-5-medium-used-3.7m-input-240k-output", "readiness-icon-svg", "checks-4-of-4-passed"} {
		if !strings.Contains(strings.Join(episode.Facts, "\n"), fact) {
			t.Errorf("the facts %v lack %q", episode.Facts, fact)
		}
	}

	if status, stdout, stderr := p.exWith(deps, "agent", "reindex"); status != 0 || !strings.Contains(stdout, `"episodes": 1`) {
		t.Errorf("ex agent reindex = %d: %s\n%s", status, stderr, stdout)
	}
	if status, _, stderr := p.exWith(deps, "agent", "run", "console-microruntime"); status != 3 || !strings.Contains(stderr, "finished") {
		t.Errorf("a run over the finished run gave %d: %s", status, stderr)
	}

	// The next plan of this project cites the run as an episode.
	// The episode's text opens with its title and its ticket; the ticket's own
	// text opens with the title alone, so the key below matches the episode only.
	_, title, _ := strings.Cut(strings.SplitN(memory, "\n", 2)[0], ": ")
	deps.Embedder = scriptedEmbedder{similarities: map[string]float64{title + "\n0070-console-microruntime\n": 0.81}}
	status, stdout, stderr = p.exWith(deps, "agent", "plan", ".agents/workflow/tickets/0070-console-microruntime.ticket.md", "--no-judge", "--component", consoleComponent, "--runtime", "console-again")
	if status != 0 {
		t.Fatalf("the next plan gave %d: %s\n%s", status, stderr, stdout)
	}
	expectLines(t, stdout,
		"agent › 1 candidates · similarity 0.81 to 0.81",
		"agent › dropped 0 below 0.55 · kept 1",
		"agent › 1 episodes cited · 1 with measured resources",
		"agent › reasoning · 1 cited episodes measured cpu and memory · sizing follows their peaks with headroom",
	)
}

// What the run cannot observe stays unknown: a branch without an upstream,
// a forge without its CLI, a cleanup nobody confirmed; the reservation is
// held, and the next lap may play on once the owner looked.
func TestRunLeavesWhatItCannotObserveUnknown(t *testing.T) {
	t.Parallel()

	p, deps := plannedProject(t)
	lifecycle := websiteLifecycle()
	lifecycle.cleanup = execution.Cleanup{Confirmed: false, Removed: 2, Total: 4}
	delivery := &scriptedDelivery{branch: execution.Branch{Name: "feat/local-only", Commit: "abc1234", Reason: "no upstream"}, reason: "gh is not installed"}
	deps.Driver, deps.Lifecycle, deps.Delivery, deps.Poll = websiteDriver(), lifecycle, delivery, time.Millisecond

	status, stdout, stderr := p.exInput(deps, "y\n", "agent", "run", "console-microruntime", "--confirm")
	if status != 0 {
		t.Fatalf("ex agent run = %d: %s\n%s", status, stderr, stdout)
	}
	expectLines(t, stdout,
		"pushed     unknown · no upstream",
		"pr         unknown · gh is not installed",
		"checks     unknown · no pull request",
		"runtime    unknown · the cleanup was not confirmed: 2 of 4 owned resources removed · the reservation is held",
	)
	if strings.Contains(stdout, "✓ ready") || strings.Contains(stdout, "# next") {
		t.Errorf("an unknown outcome printed a ready line:\n%s", stdout)
	}

	record := runRecord(t, p)
	if record.State != "finished" || record.Outcome == nil || *record.Outcome != "unknown" || record.Laps[0].Outcome != "unknown" || record.ProposedTicketStatus != nil {
		t.Errorf("run record = state %s, outcome %v, laps %+v", record.State, record.Outcome, record.Laps)
	}
	if record.Delivery.Branch == nil || *record.Delivery.Branch != "feat/local-only" || record.Delivery.PullRequest != nil || len(record.Delivery.Checks) != 0 {
		t.Errorf("delivery = %+v", record.Delivery)
	}

	ledger := ledgerOf(t, p)
	if line := ledger.Laps[0]; line.Cents.Reserved != 1140 || line.Cents.Observed != 980 || line.Cents.Recovered != 0 || line.Cents.Variance == nil || *line.Cents.Variance != 30 {
		t.Errorf("an unknown lap's line = %+v", line.Cents)
	}
	lap := string(mustRead(t, filepath.Join(runtimeFolder(p), "laps", "0001.lap.md")))
	if !strings.Contains(lap, "# Lap 0001: outcome unknown") || !strings.Contains(lap, "**Outcome:** unknown, the cleanup was not confirmed, $11.40 reserved and held") || !strings.Contains(lap, "| Held | 13:46:50 | the reservation is held until the outcome is known |") {
		t.Errorf("lap 0001:\n%s", lap)
	}
	var episode run.RunFence
	decodeStrict(t, []byte(fenceBody(t, filepath.Join(runtimeFolder(p), "run.md"), "langdrift-run")), &episode)
	if episode.Outcome == nil || *episode.Outcome != "unknown" {
		t.Errorf("memory fence outcome = %v", episode.Outcome)
	}

	// Once the owner looked, a new lap plays on: the stages are done, so it
	// brings the runtime up again and ends with the outcome known.
	lifecycle.cleanup = execution.Cleanup{Confirmed: true, Removed: 4, Total: 4, ClusterKept: true}
	deps.Delivery = websiteDelivery()
	status, stdout, stderr = p.exWith(deps, "agent", "run", "console-microruntime")
	if status != 0 {
		t.Fatalf("the lap after the unknown one gave %d: %s\n%s", status, stderr, stdout)
	}
	expectLines(t, stdout,
		"agent › gate       reserved $11.40 of $25.00 · opened",
		"runtime    console-microruntime · revision 1 · 1 Deployments ready · 500m · 1024 MiB",
		"✓ ready      #19 for review · $11.40 released to the envelope",
	)
	if strings.Contains(stdout, "grilling") || lifecycle.ups != 2 {
		t.Errorf("the second lap replayed the stages or skipped the runtime (%d ups):\n%s", lifecycle.ups, stdout)
	}
	if record := runRecord(t, p); len(record.Laps) != 2 || *record.Outcome != "readyForReview" {
		t.Errorf("run record after the second lap = %+v", record.Laps)
	}
}

// A runtime that does not come up stops the run with the lap open; the next
// run plays on from the runtime, not from the stages.
func TestRunStopsWhenTheRuntimeDoesNotComeUp(t *testing.T) {
	t.Parallel()

	p, deps := plannedProject(t)
	lifecycle := websiteLifecycle()
	lifecycle.failUp = errors.New("build recipe docker failed for .")
	driver := websiteDriver()
	deps.Driver, deps.Lifecycle, deps.Delivery, deps.Poll = driver, lifecycle, websiteDelivery(), time.Millisecond

	status, stdout, stderr := p.exInput(deps, "y\n", "agent", "run", "console-microruntime", "--confirm")
	if status != 1 || !strings.Contains(stderr, "did not come up") || !strings.Contains(stderr, "build recipe") {
		t.Fatalf("a runtime that did not come up gave %d: %s\n%s", status, stderr, stdout)
	}
	if record := runRecord(t, p); record.State != "failed" || record.Laps[0].Outcome != "open" {
		t.Errorf("run record = %+v", record)
	}

	lifecycle.failUp = nil
	status, stdout, stderr = p.exWith(deps, "agent", "run", "console-microruntime")
	if status != 0 {
		t.Fatalf("the resumed run gave %d: %s\n%s", status, stderr, stdout)
	}
	expectLines(t, stdout, "agent › lap        0001 open · $11.40 reserved · resumed", "runtime    console-microruntime · revision 1 · 1 Deployments ready · 500m · 1024 MiB", "✓ ready      #19 for review · $1.60 released to the envelope")
	if len(driver.calls) != 6 || lifecycle.ups != 2 {
		t.Errorf("the resumed run made %d calls and %d ups", len(driver.calls), lifecycle.ups)
	}
}
