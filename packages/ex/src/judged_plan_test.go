package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/anthropics/anthropic-sdk-go/option"

	"github.com/jonatassales/tao/packages/ex/src/application/plan"
	"github.com/jonatassales/tao/packages/ex/src/domain/finops"
	"github.com/jonatassales/tao/packages/ex/src/domain/run"
	judgeinfra "github.com/jonatassales/tao/packages/ex/src/infrastructure/judge"
)

// scriptedJudge answers each kind of question with a prepared text and a
// fixed usage, or fails as told.
type scriptedJudge struct {
	answers map[string]string
	fail    error
	asked   []plan.Question
}

func (model *scriptedJudge) Ask(_ context.Context, question plan.Question) (plan.Answer, error) {
	model.asked = append(model.asked, question)
	if model.fail != nil {
		return plan.Answer{}, model.fail
	}

	return plan.Answer{Text: model.answers[question.Kind], InputTokens: 15000, OutputTokens: 1500}, nil
}

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

	status, stdout, stderr := p.exWith(deps, "agent", "plan", ".agents/workflow/tickets/0070-console-microruntime.ticket.md", "--component", consoleComponent)
	if status != 0 {
		t.Fatalf("ex agent plan = %d\n%s%s", status, stdout, stderr)
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
	if strings.Contains(stdout, "judge      unavailable") || strings.Contains(stdout, "judge      stopped") || strings.Contains(stdout, "discrepancy") {
		t.Errorf("the judge was reported unavailable or a figure as differing:\n%s", stdout)
	}

	if len(model.asked) != 2 || model.asked[0].Kind != plan.SourcesQuestion || model.asked[1].Kind != plan.SuggestionQuestion || model.asked[0].Model != "claude-sonnet-5-5" {
		t.Errorf("the judge was asked %+v", model.asked)
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
	if written.Suggestion.By != "recommendation" || len(written.Suggestion.Reasons) != 2 || *written.Suggestion.ModelName != "Claude Sonnet 5.5" {
		t.Errorf("suggestion = %+v", written.Suggestion)
	}

	var ledger finops.Ledger
	decodeStrict(t, mustRead(t, filepath.Join(p.root, ".langdrift", "finops", "ledger", "console-microruntime.ledger.json")), &ledger)
	planning := ledger.Planning
	if planning.Model == nil || *planning.Model != "claude-sonnet-5-5" || planning.Cents.Estimate != 12 || planning.Cents.Reserved != 15 || planning.Cents.Observed != 9 || planning.Cents.Recovered != 6 || planning.Cents.Billed != nil {
		t.Errorf("planning line = %+v, want the sample's 12 / 15 / 9 / 6", planning)
	}
	if ledger.Totals.Cents.Observed != 9 || ledger.Totals.PeriodLeftAfterCents != 49991 {
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

	if status, _, stderr := p.exWith(deps, "finops", "budget", "set", "--period", "2026-10", "--period-cents", "50000", "--per-run-cents", "10", "--allocate", "0070-console-microruntime=10"); status != 0 {
		t.Fatal(stderr)
	}

	status, stdout, stderr := p.exWith(deps, "agent", "plan", ".agents/workflow/tickets/0070-console-microruntime.ticket.md", "--component", consoleComponent)
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

// A model that cannot be reached leaves the deterministic shell, says so
// once, and the plan still lands.
func TestPlanContinuesWhenTheJudgeIsUnavailable(t *testing.T) {
	t.Parallel()

	p, deps := seedProject(t)
	deps.Judge = &scriptedJudge{fail: plan.ErrModelUnavailable}

	status, stdout, stderr := p.exWith(deps, "agent", "plan", ".agents/workflow/tickets/0070-console-microruntime.ticket.md", "--component", consoleComponent)
	if status != 0 {
		t.Fatalf("ex agent plan = %d\n%s%s", status, stdout, stderr)
	}

	if strings.Count(stdout, "judge      stopped · the planner's model is unavailable") != 1 {
		t.Errorf("the judge's stop was not reported once:\n%s", stdout)
	}

	runtimePath := filepath.Join(p.root, ".langdrift", "microruntimes", "runtimes", "console-microruntime", "runtime.md")
	var written run.Plan
	decodeStrict(t, []byte(fenceBody(t, runtimePath, "langdrift-plan")), &written)
	if written.Sources[1].Verdict != nil || written.Suggestion.By != "input" || written.State != "proposed" {
		t.Errorf("an unjudged plan recorded %+v, %+v", written.Sources[1], written.Suggestion)
	}
}

// A provider failure reaches the session as its status and request id: the
// provider's body and the URL that was called never reach the judge's stop
// reason or the line that prints it.
func TestPlanReportsAProviderFailureByItsStatusAndRequestID(t *testing.T) {
	t.Parallel()

	api := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.Header().Set("Request-Id", "req_011plan")
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(writer, `{"type":"error","error":{"type":"api_error","message":"provider body that echoes the ticket"}}`)
	}))
	defer api.Close()

	p, deps := seedProject(t)
	deps.Judge = judgeinfra.Anthropic{Options: []option.RequestOption{option.WithBaseURL(api.URL), option.WithAPIKey("test-key"), option.WithMaxRetries(0)}}

	status, stdout, stderr := p.exWith(deps, "agent", "plan", ".agents/workflow/tickets/0070-console-microruntime.ticket.md", "--component", consoleComponent)
	if status != 0 {
		t.Fatalf("ex agent plan = %d\n%s%s", status, stdout, stderr)
	}

	warning := "judge      stopped · the sources call failed · ask claude-sonnet-5-5 about sources: the API answered HTTP 500 · request req_011plan\n"
	if strings.Count(stdout, warning) != 1 {
		t.Errorf("the session did not print %q once:\n%s", warning, stdout)
	}

	for _, leak := range []string{"provider body", "api_error", api.URL, "test-key"} {
		if strings.Contains(stdout+stderr, leak) {
			t.Errorf("the session carries %q:\n%s%s", leak, stdout, stderr)
		}
	}
}
