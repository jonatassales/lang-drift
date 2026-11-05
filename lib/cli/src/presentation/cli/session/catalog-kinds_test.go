package session_test

import (
	"testing"

	app "github.com/jonatassales/tao/lib/cli/src/application/session"
	"github.com/jonatassales/tao/lib/cli/src/presentation/cli/session"
)

// Each line's kind is the website's, by its id: the comments, the success
// lines and the two warnings, limits and the paused gate. Every other line
// the catalog renders is output.
func TestCatalogGivesEachLineTheWebsitesKind(t *testing.T) {
	t.Parallel()

	for _, scenario := range []struct {
		session string
		id      string
		kind    app.Kind
	}{
		{app.PlanSession, app.PlanScope, app.Comment},
		{app.PlanSession, app.Analysis, app.Comment},
		{app.PlanSession, app.Planning, app.Comment},
		{app.PlanSession, app.Handoff, app.Comment},
		{app.PlanSession, app.Reads, app.Success},
		{app.PlanSession, app.Suggest, app.Success},
		{app.PlanSession, app.Component, app.Success},
		{app.PlanSession, app.Size, app.Success},
		{app.PlanSession, app.Draft, app.Success},
		{app.PlanSession, app.Limits, app.Warning},
		{app.PlanSession, app.ReadTicket, app.Output},
		{app.PlanSession, app.Guardrail, app.Output},
		{app.PlanSession, app.Price, app.Output},
		{app.RunSession, app.RunScope, app.Comment},
		{app.RunSession, app.Paused, app.Comment},
		{app.RunSession, app.NextReview, app.Comment},
		{app.RunSession, app.NextMerge, app.Comment},
		{app.RunSession, app.GateOpened, app.Success},
		{app.RunSession, app.Confirmed, app.Success},
		{app.RunSession, app.Completed, app.Success},
		{app.RunSession, app.Pushed, app.Success},
		{app.RunSession, app.PullRequest, app.Success},
		{app.RunSession, app.ChecksPassed, app.Success},
		{app.RunSession, app.ReadyForReview, app.Success},
		{app.RunSession, app.ReadyToMerge, app.Success},
		{app.RunSession, app.GatePaused, app.Warning},
		{app.RunSession, app.Intent, app.Output},
		{app.RunSession, app.Waiting, app.Output},
		{app.RunSession, app.Stage, app.Output},
		{app.RunSession, app.Checks, app.Output},
		{app.RunSession, app.RuntimeReady, app.Output},
		{app.RunSession, app.RuntimeStopped, app.Output},
	} {
		event := app.Event{
			Session: scenario.session,
			ID:      scenario.id,
		}
		if got := session.KindOf(event); got != scenario.kind {
			t.Errorf("%s.%s is %s, want %s", scenario.session, scenario.id, got, scenario.kind)
		}
	}
}

// The typed command is the person's shell, never a line dft prints: the
// catalog has no words for the website's command lines.
func TestCatalogHasNoWordsForTheTypedCommand(t *testing.T) {
	t.Parallel()

	if session.Covers(app.PlanSession, "plan") || session.Covers(app.RunSession, "command") {
		t.Error("the catalog renders the typed command")
	}
}
