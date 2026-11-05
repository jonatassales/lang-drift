package session_test

import (
	"bytes"
	"testing"

	app "github.com/jonatassales/tao/lib/cli/src/application/session"
	"github.com/jonatassales/tao/lib/cli/src/presentation/cli/session"
)

// On a terminal, each line carries its kind's tone: a comment dimmed, a
// success line green, a warning yellow, an output line the terminal's own
// color.
func TestHumanTonesEachKindOnATerminal(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	sink := session.For(&out, false, true)

	for _, event := range []app.Event{
		{Session: app.PlanSession, ID: app.PlanScope},
		{Session: app.PlanSession, ID: app.Suggest, Data: app.Data{"model": "claude-sonnet-5-5", "effort": "medium"}},
		{Session: app.PlanSession, ID: app.Limits, Data: app.Data{"requested": "none requested"}},
		{Session: app.PlanSession, ID: app.Embed},
	} {
		sink.Emit(event)
	}

	want := "\x1b[2m# plan only · nothing is built, started or paid for in this step\x1b[0m\n" +
		"\x1b[32m✓ suggests   claude-sonnet-5-5 · effort medium\x1b[0m\n" +
		"\x1b[33mlimits     gpu and disk are not runtime.md fields · none requested\x1b[0m\n" +
		"console › embedding the ticket intent · 1 query vector\n"
	if out.String() != want {
		t.Errorf("printed %q, want %q", out.String(), want)
	}
}

// Off a terminal, the bytes are the plain text, with no escape sequence.
func TestHumanPrintsPlainTextOffATerminal(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	sink := session.For(&out, false, false)

	for _, event := range []app.Event{
		{Session: app.PlanSession, ID: app.PlanScope},
		{Session: app.PlanSession, ID: app.Suggest, Data: app.Data{"model": "claude-sonnet-5-5", "effort": "medium"}},
	} {
		sink.Emit(event)
	}

	want := "# plan only · nothing is built, started or paid for in this step\n✓ suggests   claude-sonnet-5-5 · effort medium\n"
	if out.String() != want {
		t.Errorf("printed %q, want %q", out.String(), want)
	}
}

// When the session opens its project, a terminal's title names the tool and
// the project's folder, "dft — project"; off a terminal nothing is written.
func TestHumanTitlesTheTerminalWithTheProject(t *testing.T) {
	t.Parallel()

	var onTerminal, offTerminal bytes.Buffer

	app.Announce(session.For(&onTerminal, false, true), "project")
	app.Announce(session.For(&offTerminal, false, false), "project")

	if onTerminal.String() != "\x1b]0;dft — project\x07" || offTerminal.Len() != 0 {
		t.Errorf("on a terminal %q, off one %q", onTerminal.String(), offTerminal.String())
	}
}
