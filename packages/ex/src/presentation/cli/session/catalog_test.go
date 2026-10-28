package session_test

import (
	"testing"

	app "github.com/jonatassales/tao/packages/ex/src/application/session"
	"github.com/jonatassales/tao/packages/ex/src/presentation/cli/session"
)

// The website's plan lines, in order, and its run messages. A line the
// catalog cannot render would be a line the CLI forgets.
var (
	websitePlanLines = []string{
		"plan", "scope", "readTicket", "links", "analysis", "readPrd", "readDesignDoc", "readSpec", "reads",
		"guardrail", "judgeGuardrail", "bounds", "budget", "memory", "embed", "vectorSearch", "candidates",
		"judgeRelevance", "drop", "episodes", "reasonResources", "weighEffort", "weighModel", "suggest",
		"planning", "resolve", "inspectGo", "go", "inspectNode", "node", "inspectNats", "nats", "inspectRedis",
		"redis", "inspectNext", "next", "crossCheck", "size", "limits", "price", "validate", "draft", "handoff",
	}
	websiteRunMessages = []string{
		"command", "scope", "intent", "estimate", "budget", "gateOpened", "gatePaused", "paused", "waiting",
		"confirmed", "grilling", "toSpec", "toTickets", "hook", "implement", "codeReview", "runtimeReady",
		"completed", "pushed", "pullRequest", "checks", "checksPassed", "runtimeStopped", "readyForReview",
		"readyToMerge", "nextReview", "nextMerge",
	}
)

func TestCatalogCoversEveryWebsiteLine(t *testing.T) {
	t.Parallel()

	if len(websitePlanLines) != 43 {
		t.Fatalf("the website prints 43 plan lines, the list has %d", len(websitePlanLines))
	}

	for _, id := range websitePlanLines {
		if !session.Covers(app.PlanSession, id) {
			t.Errorf("plan line %q has no words", id)
		}
	}

	for _, id := range websiteRunMessages {
		if !session.Covers(app.RunSession, id) {
			t.Errorf("run message %q has no words", id)
		}
	}
}

func TestCatalogRendersTheWebsiteWordsWithTrueFigures(t *testing.T) {
	t.Parallel()

	tests := []struct {
		event app.Event
		want  string
	}{
		{app.Event{Session: app.PlanSession, ID: app.PlanCommand, Data: app.Data{"ticket": "0070-console-microruntime"}}, "ex agent plan 0070-console-microruntime"},
		{app.Event{Session: app.PlanSession, ID: app.Links, Data: app.Data{"status": "needs-triage", "links": []string{"0010-SPEC", "guardrail 0001-finops"}}}, "agent › needs-triage · links 0010-SPEC · guardrail 0001-finops"},
		{app.Event{Session: app.PlanSession, ID: app.Reads, Data: app.Data{"count": 4, "refused": 0}}, "✓ read       4 documents · nothing outside the ticket's links"},
		{app.Event{Session: app.PlanSession, ID: app.Budget, Data: app.Data{"perRunCents": 2500, "by": "owner"}}, "budget     $25.00 per run for model calls · set by owner in the CLI"},
		{app.Event{Session: app.PlanSession, ID: app.Drop, Data: app.Data{"dropped": 3, "threshold": "0.55", "kept": 3}}, "agent › dropped 3 below 0.55 · kept 3"},
		{app.Event{Session: app.PlanSession, ID: app.Suggest, Data: app.Data{"model": "claude-sonnet-5-5", "effort": "medium"}}, "✓ suggests   claude-sonnet-5-5 · effort medium"},
		{app.Event{Session: app.PlanSession, ID: app.Component, Data: app.Data{"path": "services/signup-api", "technology": "Go", "kind": "api", "cpuMillis": 250, "memoryMiB": 256}}, "✓ services/signup-api   Go · api        250m · 256 MiB"},
		{app.Event{Session: app.PlanSession, ID: app.Size, Data: app.Data{"cpuMillis": 1450, "boundedCpuMillis": 2000, "memoryMiB": 2304, "boundedMemoryMiB": 3072, "deployments": 5}}, "✓ size       1450m of 2000m · 2304 of 3072 MiB · 5 one-replica Deployments · proposed shares"},
		{app.Event{Session: app.PlanSession, ID: app.Draft, Data: app.Data{"revision": 1}}, "✓ plan       runtime.md draft, revision 1 · proposed · authorize with ex agent run"},
		{app.Event{Session: app.RunSession, ID: app.Estimate, Data: app.Data{"estimateCents": 950, "inputThousands": 3500, "outputThousands": 250}}, "agent › estimate   $9.50 · 3.5M input · 250K output · list prices · FIN-02"},
		{app.Event{Session: app.RunSession, ID: app.GateOpened, Data: app.Data{"reservedCents": 1140, "perRunCents": 2500}}, "agent › gate       reserved $11.40 of $25.00 · opened"},
		{app.Event{Session: app.RunSession, ID: app.GatePaused, Data: app.Data{"reservedCents": 4080, "perRunCents": 2500}}, "agent › gate       reservation $40.80 exceeds the $25.00 budget · paused"},
		{app.Event{Session: app.RunSession, ID: app.Stage, Data: app.Data{"skill": "grilling", "summary": "3 open questions settled"}}, "agent › grilling   3 open questions settled"},
		{app.Event{Session: app.RunSession, ID: app.Stage, Data: app.Data{"skill": "code-review", "summary": "2 findings · fixed"}}, "agent › code-review 2 findings · fixed"},
		{app.Event{Session: app.RunSession, ID: app.Completed, Data: app.Data{"observedCents": 980, "reservedCents": 1140}}, "✓ completed  implementation · $9.80 observed of $11.40 reserved · usage record"},
		{app.Event{Session: app.RunSession, ID: app.ReadyForReview, Data: app.Data{"number": 19, "releasedCents": 160}}, "✓ ready      #19 for review · $1.60 released to the envelope"},
		{app.Event{Session: app.RunSession, ID: app.NextMerge}, "# next: the merge stays the owner's"},
	}

	for _, tt := range tests {
		if got := session.Text(tt.event); got != tt.want {
			t.Errorf("%s.%s rendered\n%q\nwant\n%q", tt.event.Session, tt.event.ID, got, tt.want)
		}
	}
}

func TestCatalogGivesEachLineTheWebsiteKind(t *testing.T) {
	t.Parallel()

	if kind := session.KindOf(app.Event{Session: app.RunSession, ID: app.GatePaused}); kind != app.Warning {
		t.Errorf("gatePaused is %s, want a warning", kind)
	}

	if kind := session.KindOf(app.Event{Session: app.PlanSession, ID: app.Suggest}); kind != app.Success {
		t.Errorf("suggest is %s, want a success", kind)
	}
}
