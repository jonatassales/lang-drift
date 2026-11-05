package session_test

import (
	"testing"

	"github.com/jonatassales/tao/lib/cli/src/presentation/cli/session"
)

// Each shape lays out as spec 0048's line format says, counting characters:
// a label shorter than 10 columns is padded, a label of 10 or more is
// followed by exactly one space, a result's value starts at the same column
// whatever its label, the component rows align in three columns, and no
// line ends with a space.
func TestLayoutLaysEachShapeInTheSimulationsColumns(t *testing.T) {
	t.Parallel()

	for _, scenario := range []struct {
		name  string
		shape session.Shape
		cells session.Cells
		want  string
	}{
		{"a comment", session.Comment, session.Cells{Text: "plan only · nothing is built, started or paid for in this step"}, "# plan only · nothing is built, started or paid for in this step"},
		{"a narrative actor line", session.Narrative, session.Cells{Text: "tool call: read 0031-TICKET-book-a-demo-runtime"}, "console › tool call: read 0031-TICKET-book-a-demo-runtime"},
		{"a short label", session.Labeled, session.Cells{Label: "gate", Value: "reserved $11.40 of $25.00 · opened"}, "console › gate       reserved $11.40 of $25.00 · opened"},
		{"a label of 10", session.Labeled, session.Cells{Label: "to-tickets", Value: "3 tickets · blocking edges recorded"}, "console › to-tickets 3 tickets · blocking edges recorded"},
		{"a label of 11", session.Labeled, session.Cells{Label: "code-review", Value: "2 findings · fixed"}, "console › code-review 2 findings · fixed"},
		{"a result", session.Result, session.Cells{Label: "read", Value: "4 documents · nothing outside the ticket's links"}, "✓ read       4 documents · nothing outside the ticket's links"},
		{"a result with a numbered label", session.Result, session.Cells{Label: "pr #19", Value: "feat(demo): run book a demo as a microruntime"}, "✓ pr #19     feat(demo): run book a demo as a microruntime"},
		{"a status line", session.Status, session.Cells{Label: "runtime", Value: "stopped · 5 Deployments removed · nothing left"}, "runtime    stopped · 5 Deployments removed · nothing left"},
		{"a label of 11 in a status line", session.Status, session.Cells{Label: "discrepancy", Value: "the spec says 4000m · the input holds"}, "discrepancy the spec says 4000m · the input holds"},
		{"a component row", session.ComponentRow, session.Cells{Path: "services/signup-api", StackRole: "Go · api", Share: "250m · 256 MiB"}, "✓ services/signup-api   Go · api        250m · 256 MiB"},
		{"a shorter component row", session.ComponentRow, session.Cells{Path: "apps/signup", StackRole: "Next.js · web", Share: "500m · 1024 MiB"}, "✓ apps/signup           Next.js · web   500m · 1024 MiB"},
		{"no value", session.Status, session.Cells{Label: "limits"}, "limits"},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			t.Parallel()

			if got := session.Layout(scenario.shape, scenario.cells); got != scenario.want {
				t.Errorf("Layout = %q, want %q", got, scenario.want)
			}
		})
	}
}
