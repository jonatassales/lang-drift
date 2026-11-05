package main

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jonatassales/tao/lib/cli/src/application/plan"
	"github.com/jonatassales/tao/lib/cli/src/domain/finops"
	"github.com/jonatassales/tao/lib/cli/src/domain/run"
	judgeinfra "github.com/jonatassales/tao/lib/cli/src/infrastructure/judge"
)

// sampleAnswers are the judge's answers that reproduce the sample's
// analysis: the verdicts, roles and reasons of its sources, its facts, the
// criteria that apply, the figure the ticket states, and the suggestion.
func sampleAnswers(t *testing.T) map[string]string {
	t.Helper()

	var sample run.Plan
	decodeStrict(t, mustRead(t, samplePlan), &sample)

	type source struct {
		Path       string  `json:"path"`
		Verdict    string  `json:"verdict"`
		Reason     string  `json:"reason"`
		Confidence float64 `json:"confidence"`
		Role       string  `json:"role"`
		Summary    string  `json:"summary"`
	}
	sources := make([]source, 0, len(sample.Sources))
	for _, entry := range sample.Sources {
		sources = append(sources, source{Path: entry.Path, Verdict: *entry.Verdict, Reason: *entry.Reason, Confidence: *entry.Confidence, Role: *entry.InferredRole, Summary: "summary of " + filepath.Base(entry.Path)})
	}

	answer := map[string]any{
		"sources": sources,
		"facts":   sample.Facts,
		"criteria": []map[string]string{
			{"id": "FIN-02", "applies": "estimate before a paid call"},
			{"id": "FIN-03", "applies": "a budget before a paid run"},
			{"id": "FIN-04", "applies": "size cpu and memory to the load"},
		},
		"figures": []map[string]any{{"topic": "resource", "fact": "about 2,000 millicores and 3,072 MiB", "source": ".agents/workflow/tickets/0070-console-microruntime.ticket.md#L5", "cpuMillis": 2000, "memoryMiB": 3072}},
	}
	sourcesText, err := json.Marshal(answer)
	if err != nil {
		t.Fatal(err)
	}

	suggestion := map[string]any{
		"model":        "claude-sonnet-5-5",
		"effort":       "medium",
		"effortReason": "one component, no broker wiring: high effort is not warranted",
		"modelReason":  "Opus 5.5 at high effort reserves more than the per-run budget",
		"reasons":      sample.Suggestion.Reasons,
	}
	suggestionText, err := json.Marshal(suggestion)
	if err != nil {
		t.Fatal(err)
	}

	return map[string]string{plan.SourcesQuestion: string(sourcesText), plan.SuggestionQuestion: string(suggestionText)}
}

// With a judge that answers as the sample's agent did, the plan prints the
// judged lines, records the verdicts, the facts, the criteria and the
// consistent figure, and pays the planning line.
func TestPlanJudgesTheSourcesAndSuggests(t *testing.T) {
	t.Parallel()

	p, deps := seedProject(t)
	model := &scriptedJudge{answers: sampleAnswers(t)}
	deps.Judge = model

	status, stdout, stderr := p.dftWith(deps, "agent", "plan", ".agents/workflow/tickets/0070-console-microruntime.ticket.md", "--component", consoleComponent)
	if status != 0 {
		t.Fatalf("dft agent plan = %d\n%s%s", status, stdout, stderr)
	}

	for _, line := range []string{
		"agent › tool call: read 0010-exphere-and-the-langdrift-folder.spec.md · summary of 0010-exphere-and-the-langdrift-folder.spec.md",
		"agent › judging · FIN-02 estimate before a paid call · FIN-03 a budget before a paid run · FIN-04 size cpu and memory to the load",
		"agent › weighing effort · one component, no broker wiring: high effort is not warranted",
		"agent › weighing model · Opus 5.5 at high effort reserves more than the per-run budget",
		"✓ suggests   claude-sonnet-5-5 · effort medium",
	} {
		if !strings.Contains(stdout, line+"\n") {
			t.Errorf("the session did not print %q:\n%s", line, stdout)
		}
	}
	if strings.Contains(stdout, "discrepancy") {
		t.Errorf("a figure was reported as differing:\n%s", stdout)
	}

	judged := strings.Index(stdout, "agent › judging · ")
	reasoning := strings.Index(stdout, "agent › reasoning · ")
	if judged < 0 || reasoning < judged {
		t.Errorf("the reasoning line printed before the judged answer, or not at all:\n%s", stdout)
	}

	if len(model.asked) != 2 || model.asked[0].Kind != plan.SourcesQuestion || model.asked[1].Kind != plan.SuggestionQuestion || model.asked[0].Model != "claude-sonnet-5-5" || model.asked[0].PromptVersion != "judge-1" {
		t.Errorf("the judge was asked %+v", model.asked)
	}
	for _, connection := range model.connections {
		if !connection.Recorded || connection.URL != "https://console.test" || connection.Secret != testRuntimeSecret {
			t.Errorf("the judge was asked through %v", connection)
		}
	}
	if !strings.Contains(model.asked[0].Prompt, "<<<.agents/workflow/tickets/0070-console-microruntime.ticket.md") {
		t.Error("the sources question did not frame the ticket as material")
	}

	runtimePath := filepath.Join(p.root, ".langdrift", "microruntimes", "runtimes", "console-microruntime", "runtime.md")
	var written run.Plan
	decodeStrict(t, []byte(fenceBody(t, runtimePath, "langdrift-plan")), &written)

	var sample run.Plan
	decodeStrict(t, mustRead(t, samplePlan), &sample)

	for index, source := range written.Sources {
		want := sample.Sources[index]
		if source.Verdict == nil || *source.Verdict != *want.Verdict || *source.Reason != *want.Reason || *source.Confidence != *want.Confidence || *source.InferredRole != *want.InferredRole {
			t.Errorf("source %d judged as %+v, want %+v", index, source, want)
		}
	}

	if len(written.Facts) != len(sample.Facts) || written.Facts[0] != sample.Facts[0] {
		t.Errorf("facts = %+v", written.Facts)
	}
	if len(written.Discrepancies) != 1 || written.Discrepancies[0] != sample.Discrepancies[0] {
		t.Errorf("discrepancies = %+v, want %+v", written.Discrepancies, sample.Discrepancies)
	}
	if written.Collection.Agent == nil || written.Collection.Agent.Model != "claude-sonnet-5-5" || written.Collection.Agent.PromptVersion != "judge-1" || !written.Collection.Agent.Paid || written.Collection.Agent.LedgerLine != "planning" {
		t.Errorf("collection agent = %+v", written.Collection.Agent)
	}
	if written.Suggestion.By != "judge" || len(written.Suggestion.Reasons) != 2 || *written.Suggestion.ModelName != "Claude Sonnet 5.5" {
		t.Errorf("suggestion = %+v", written.Suggestion)
	}
	judgedBy := run.Judge{
		Version:       "console-judge-1",
		Model:         "claude-sonnet-5-5",
		PromptVersion: "judge-1",
	}
	if written.Judge == nil {
		t.Fatal("the plan names no judge")
	}

	hashed := strings.HasPrefix(written.Judge.PromptHash, "sha256:") && len(written.Judge.PromptHash) == len("sha256:")+64
	named := *written.Judge
	named.PromptHash = ""
	if named != judgedBy || !hashed {
		t.Errorf("judge = %+v", written.Judge)
	}

	var ledger finops.Ledger
	decodeStrict(t, mustRead(t, filepath.Join(p.root, ".langdrift", "finops", "ledger", "console-microruntime.ledger.json")), &ledger)
	planning := ledger.Planning
	if planning.Model == nil || *planning.Model != "claude-sonnet-5-5" || planning.Cents.Estimate != 12 || planning.Cents.Reserved != 15 || planning.Cents.Observed != 10 || planning.Cents.Recovered != 5 || planning.Cents.Billed != nil {
		t.Errorf("planning line = %+v, want 12 estimated, 15 reserved, the console's 10 observed and 5 recovered", planning)
	}
	if planning.Provenance == nil || *planning.Provenance != "console" || planning.Illustrative == nil || !*planning.Illustrative {
		t.Errorf("planning line provenance = %v, illustrative %v", planning.Provenance, planning.Illustrative)
	}
	if planning.Usage == nil || planning.Usage.InputTokens != 30000 || planning.Usage.OutputTokens != 3000 {
		t.Errorf("planning usage = %+v", planning.Usage)
	}
	if ledger.Totals.Cents.Observed != 10 || ledger.Totals.PeriodLeftAfterCents != 49990 {
		t.Errorf("totals = %+v", ledger.Totals)
	}
}

// A planning reservation above the per-run cap pauses before the first
// call: nothing is paid, nothing is written.
func TestPlanPausesAtTheGateBeforeTheJudgeRuns(t *testing.T) {
	t.Parallel()

	p, deps := seedProject(t)
	model := &scriptedJudge{answers: sampleAnswers(t)}
	deps.Judge = model

	if status, _, stderr := p.dftWith(deps, "finops", "budget", "set", "--period", "2026-10", "--period-cents", "50000", "--per-run-cents", "10", "--allocate", "0070-console-microruntime=10"); status != 0 {
		t.Fatal(stderr)
	}

	status, stdout, stderr := p.dftWith(deps, "agent", "plan", ".agents/workflow/tickets/0070-console-microruntime.ticket.md", "--component", consoleComponent)
	if status != 3 || !strings.Contains(stderr, "FG-03 paused") || !strings.Contains(stderr, "nothing was paid") {
		t.Fatalf("a paused planner gave %d: %s\n%s", status, stderr, stdout)
	}

	if len(model.asked) != 0 {
		t.Errorf("the judge was asked %d questions before the gate", len(model.asked))
	}

	if _, err := os.Stat(filepath.Join(p.root, ".langdrift", "microruntimes", "runtimes", "console-microruntime", "runtime.md")); err == nil {
		t.Error("a paused plan wrote the runtime file")
	}
}

// Without a recorded connection there is no judge to ask, so there is no
// plan: the person is told to connect, and nothing is asked or written.
func TestPlanIsRefusedWithoutAConnection(t *testing.T) {
	t.Parallel()

	p, deps := seedProject(t)
	model := deps.Judge.(*scriptedJudge)

	disconnectProject(t, p.root)

	status, stdout, stderr := p.dftWith(deps, "agent", "plan", ".agents/workflow/tickets/0070-console-microruntime.ticket.md", "--component", consoleComponent)
	if status != 4 || !strings.Contains(stderr, "connect") {
		t.Fatalf("an unconnected plan gave %d: %s\n%s", status, stderr, stdout)
	}

	if len(model.asked) != 0 {
		t.Errorf("the judge was asked %d questions without a connection", len(model.asked))
	}

	assertNoPlanWritten(t, p.root)
}

// Each way the judge route can fail refuses the plan with the route's own
// reason: the secret refused, the judge unavailable, no answer in time, or
// the model's refusal. No plan is written and the secret never reaches the
// session. A refusal before the model ran pays nothing; a call cut off holds
// its reservation, since what it cost is unknown; the model's refusal is on
// the ledger at what the console charged.
func TestPlanIsRefusedWhenTheConsoleDoesNotJudge(t *testing.T) {
	t.Parallel()

	for _, scenario := range []struct {
		name    string
		console func(*fakeConsole)
		reason  string
		spend   string
	}{
		{"the secret refused", func(console *fakeConsole) { console.status = http.StatusUnauthorized }, "the console refused RUNTIME_SECRET", "nothing"},
		{"the judge unavailable", func(console *fakeConsole) {
			console.status, console.body = http.StatusServiceUnavailable, `{"code":"UNAVAILABLE","reason":"platform-key-missing"}`
		}, "the console's judge is unavailable · platform-key-missing", "nothing"},
		{"no answer in time", func(console *fakeConsole) { console.delay = 2 * time.Second }, "the console did not answer in time", "held"},
		{"the model's refusal", func(console *fakeConsole) { console.refusal = "the model stopped with refusal, not end_turn" }, "the model stopped with refusal, not end_turn", "charged"},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			t.Parallel()

			console := newFakeConsole(t, minimalAnswers())
			scenario.console(console)

			p, deps := seedProject(t)
			connectProject(t, p.root, console.server.URL)
			deps.Judge = judgeinfra.Console{Timeout: 200 * time.Millisecond}

			status, stdout, stderr := p.dftWith(deps, "agent", "plan", ".agents/workflow/tickets/0070-console-microruntime.ticket.md", "--component", consoleComponent)
			if status == 0 || !strings.Contains(stderr, scenario.reason) {
				t.Fatalf("the plan gave %d, want a refusal naming %q: %s\n%s", status, scenario.reason, stderr, stdout)
			}

			if strings.Contains(stdout+stderr, testRuntimeSecret) {
				t.Errorf("the session carries the secret:\n%s%s", stdout, stderr)
			}

			if strings.Contains(stdout, "agent › reasoning") {
				t.Errorf("a refused plan printed a reasoning line:\n%s", stdout)
			}

			asked, _ := console.questions()
			if len(asked) != 1 {
				t.Errorf("the console was asked %d times; a failed call is never retried", len(asked))
			}

			assertNoPlanWritten(t, p.root)
			assertPlanningSpend(t, p, scenario.spend)
		})
	}
}

// assertPlanningSpend checks what a refused plan left on the ledger: nothing,
// a reservation held for a call whose cost is unknown, or what the console
// charged.
func assertPlanningSpend(t *testing.T, p project, spend string) {
	t.Helper()

	ledgerPath := filepath.Join(p.root, ".langdrift", "finops", "ledger", "console-microruntime.ledger.json")
	_, err := os.Stat(ledgerPath)

	if spend == "nothing" {
		if err == nil {
			t.Error("a refusal before the model ran left a planning line")
		}

		return
	}

	planning := ledgerOf(t, p).Planning
	held := planning.Cents.Recovered == 0 && strings.Contains(planning.Note, "the reservation is held")
	charged := planning.Cents.Observed == 5 && planning.Cents.Recovered == 10

	if spend == "held" && !held {
		t.Errorf("a call cut off released its reservation: %+v", planning)
	}

	if spend == "charged" && !charged {
		t.Errorf("the model's refusal is not on the planning line at the console's charge: %+v", planning)
	}
}

// The gates' statuses only come before a question. A judge that reports one
// after it gave no answer, so the plan is refused as outside the contract,
// and nothing is proposed.
func TestPlanRefusesAGateStatusFromTheJudge(t *testing.T) {
	t.Parallel()

	for _, status := range []plan.AskStatus{plan.AskCleared, plan.AskPaused} {
		t.Run(string(status), func(t *testing.T) {
			t.Parallel()

			p, deps := seedProject(t)
			deps.Judge = &scriptedJudge{outcome: plan.AskOutcome{
				Kind:   plan.SourcesQuestion,
				Status: status,
			}}

			code, stdout, stderr := p.dftWith(deps, "agent", "plan", ".agents/workflow/tickets/0070-console-microruntime.ticket.md", "--component", consoleComponent)
			refusedAsOutsideTheContract := code == 1 && strings.Contains(stderr, "outside the contract") && strings.Contains(stderr, string(status))
			if !refusedAsOutsideTheContract {
				t.Fatalf("a judge reporting %s gave %d: %s\n%s", status, code, stderr, stdout)
			}

			assertNoPlanWritten(t, p.root)
		})
	}
}

// An answer is read as the question that was asked, whatever kind the judge
// says it answered: a suggestion given to the sources question is outside the
// contract, and nothing is proposed.
func TestPlanReadsAnAnswerAsTheQuestionAsked(t *testing.T) {
	t.Parallel()

	p, deps := seedProject(t)
	deps.Judge = &scriptedJudge{
		answers: minimalAnswers(),
		outcome: plan.AskOutcome{
			Kind:   plan.SuggestionQuestion,
			Status: plan.AskAnswered,
		},
	}

	code, stdout, stderr := p.dftWith(deps, "agent", "plan", ".agents/workflow/tickets/0070-console-microruntime.ticket.md", "--component", consoleComponent)
	refusedAsTheSourcesAnswer := code == 1 && strings.Contains(stderr, "the sources answer is refused")
	if !refusedAsTheSourcesAnswer {
		t.Fatalf("a suggestion given to the sources question gave %d: %s\n%s", code, stderr, stdout)
	}

	assertNoPlanWritten(t, p.root)
}

// --no-judge is gone: no plan exists without a judge.
func TestPlanHasNoFlagToSkipTheJudge(t *testing.T) {
	t.Parallel()

	p, deps := seedProject(t)

	status, _, stderr := p.dftWith(deps, "agent", "plan", ".agents/workflow/tickets/0070-console-microruntime.ticket.md", "--no-judge", "--component", consoleComponent)
	if status != 2 || !strings.Contains(stderr, "unknown flag: --no-judge") {
		t.Fatalf("--no-judge gave %d: %s", status, stderr)
	}

	assertNoPlanWritten(t, p.root)
}

// assertNoPlanWritten fails when a refused plan left a runtime file behind.
func assertNoPlanWritten(t *testing.T, root string) {
	t.Helper()

	if _, err := os.Stat(filepath.Join(root, ".langdrift", "microruntimes", "runtimes", "console-microruntime", "runtime.md")); err == nil {
		t.Error("a refused plan wrote the runtime file")
	}
}
