package main

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/jonatassales/tao/lib/cli/src/application/session"
)

// On a terminal, the plan titles the window "dft — <project folder>",
// prints its scope dimmed as its first line and never echoes the typed
// command. Its acts open at their moments: the analysis right after the
// links, the plan right after the suggestion, the hand-off last. The
// observer the dependencies hold sees every line, by id, in the order
// printed, each with its tone.
func TestThePlanOnATerminalIsTitledTonedAndObserved(t *testing.T) {
	t.Parallel()

	seeded, deps := seedProject(t)
	recorder := &sessionRecorder{}
	deps.Observer = recorder

	status, stdout, stderr := seeded.dftOnATerminal(deps, "agent", "plan", consoleTicket, "--component", consoleComponent)
	if status != 0 {
		t.Fatalf("dft agent plan = %d: %s\n%s", status, stderr, stdout)
	}

	title := "\x1b]0;dft — " + filepath.Base(seeded.root) + "\x07"
	scope := "\x1b[2m# plan only · nothing is built, started or paid for in this step\x1b[0m\n"
	opens := strings.HasPrefix(stdout, title+scope)
	if !opens || strings.Contains(stdout, "dft agent plan") {
		t.Fatalf("the plan opened with %q", stdout[:min(160, len(stdout))])
	}

	events := recorder.events
	analysisAfterLinks := eventPosition(events, session.Analysis) == eventPosition(events, session.Links)+1
	planAfterSuggestion := eventPosition(events, session.Planning) == eventPosition(events, session.Suggest)+1
	handoffLast := eventPosition(events, session.Handoff) == len(events)-1
	if !analysisAfterLinks || !planAfterSuggestion || !handoffLast {
		t.Errorf("the acts open out of place: %+v", events)
	}

	var expected strings.Builder
	for _, event := range events {
		expected.WriteString(tonedLine(event))
	}

	printed := strings.TrimPrefix(stdout, title)
	if printed != expected.String() || len(events) == 0 {
		t.Errorf("the observer saw %d events that do not match the printed lines:\n%s", len(events), printed)
	}
}
