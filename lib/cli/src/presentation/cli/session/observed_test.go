package session_test

import (
	"bytes"
	"testing"

	app "github.com/jonatassales/tao/lib/cli/src/application/session"
	"github.com/jonatassales/tao/lib/cli/src/presentation/cli/session"
)

// An observer receives every event the printed session receives, in order,
// and without an observer the printed session is the sink itself.
func TestObservedTeesEveryEventToTheObserver(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	observer := &recording{}
	sink := session.Observed(session.For(&out, false, false), observer)

	for _, event := range []app.Event{
		{Session: app.PlanSession, ID: app.PlanScope},
		{Session: app.PlanSession, ID: app.Embed},
	} {
		sink.Emit(event)
	}

	seen := len(observer.events) == 2 && observer.events[0].ID == app.PlanScope && observer.events[1].ID == app.Embed
	if !seen || out.String() == "" {
		t.Errorf("observer saw %+v, printed %q", observer.events, out.String())
	}

	alone := session.For(&out, false, false)
	if session.Observed(alone, nil) != alone {
		t.Error("without an observer the sink was wrapped")
	}
}
