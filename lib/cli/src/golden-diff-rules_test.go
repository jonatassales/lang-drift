package main

import (
	"strings"
	"testing"

	"github.com/jonatassales/tao/lib/cli/src/application/session"
)

// goldenDiff returns no message when the sessions match with an empty pending
// list, and returns one message naming the line when an unlisted line differs,
// a listed line already matches, a line is missing or extra, or a line's kind
// differs even while pending.
func TestGoldenDiffRules(t *testing.T) {
	t.Parallel()

	for _, scenario := range []struct {
		name    string
		golden  []goldenLine
		printed []printedLine
		pending map[string]string
		want    string
	}{
		{
			name: "identical sessions with empty pending",
			golden: []goldenLine{
				{
					id:   "scope",
					kind: session.Comment,
					text: "# plan only",
				},
				{
					id:   "readTicket",
					kind: session.Output,
					text: "console › tool call: read ticket",
				},
			},
			printed: []printedLine{
				{
					id:   "scope",
					kind: session.Comment,
					text: "# plan only",
				},
				{
					id:   "readTicket",
					kind: session.Output,
					text: "console › tool call: read ticket",
				},
			},
			pending: map[string]string{},
			want:    "",
		},
		{
			name: "differing unlisted line",
			golden: []goldenLine{
				{
					id:   "budget",
					kind: session.Output,
					text: "budget     $25.00 per run for model calls · set by the team in the design doc",
				},
			},
			printed: []printedLine{
				{
					id:   "budget",
					kind: session.Output,
					text: "budget     $25.00 per run · set by owner in the CLI",
				},
			},
			pending: map[string]string{},
			want:    `line 1 "budget": differs from golden`,
		},
		{
			name: "matching listed line",
			golden: []goldenLine{
				{
					id:   "suggest",
					kind: session.Success,
					text: "✓ suggests   claude-sonnet-5-5 · effort medium",
				},
			},
			printed: []printedLine{
				{
					id:   "suggest",
					kind: session.Success,
					text: "✓ suggests   claude-sonnet-5-5 · effort medium",
				},
			},
			pending: map[string]string{"suggest": "0051"},
			want:    `line 1 "suggest": already matches golden, remove from pending under spec 0051`,
		},
		{
			name: "differing listed line stays pending",
			golden: []goldenLine{
				{
					id:   "bounds",
					kind: session.Output,
					text: "bounded    cpuMillis 2000 · memoryMiB 3072 · set by the team in the design doc",
				},
			},
			printed: []printedLine{
				{
					id:   "bounds",
					kind: session.Output,
					text: "bounded    cpuMillis 2000 · memoryMiB 3072 · set by owner in the CLI",
				},
			},
			pending: map[string]string{"bounds": "0054"},
			want:    "",
		},
		{
			name: "missing line",
			golden: []goldenLine{
				{
					id:   "scope",
					kind: session.Comment,
					text: "# plan only",
				},
				{
					id:   "links",
					kind: session.Output,
					text: "console › ready-for-agent · links",
				},
				{
					id:   "analysis",
					kind: session.Comment,
					text: "# analysis",
				},
			},
			printed: []printedLine{
				{
					id:   "scope",
					kind: session.Comment,
					text: "# plan only",
				},
				{
					id:   "analysis",
					kind: session.Comment,
					text: "# analysis",
				},
			},
			pending: map[string]string{},
			want:    `missing line 2 "links"`,
		},
		{
			name: "extra line",
			golden: []goldenLine{
				{
					id:   "scope",
					kind: session.Comment,
					text: "# plan only",
				},
				{
					id:   "analysis",
					kind: session.Comment,
					text: "# analysis",
				},
			},
			printed: []printedLine{
				{
					id:   "scope",
					kind: session.Comment,
					text: "# plan only",
				},
				{
					id:   "extra",
					kind: session.Output,
					text: "console › extra line",
				},
				{
					id:   "analysis",
					kind: session.Comment,
					text: "# analysis",
				},
			},
			pending: map[string]string{},
			want:    `extra line 2 "extra"`,
		},
		{
			name: "different kind even when pending",
			golden: []goldenLine{
				{
					id:   "reads",
					kind: session.Success,
					text: "✓ read       4 documents · nothing outside the ticket's links",
				},
			},
			printed: []printedLine{
				{
					id:   "reads",
					kind: session.Output,
					text: "✓ read       5 documents · nothing outside the ticket's links",
				},
			},
			pending: map[string]string{"reads": "0053"},
			want:    `line 1 "reads": kind differs: got output, want success`,
		},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			t.Parallel()

			differences := goldenDiff(scenario.golden, scenario.printed, scenario.pending)
			if scenario.want == "" {
				if len(differences) != 0 {
					t.Fatalf("expected no differences, got %v", differences)
				}

				return
			}

			if len(differences) != 1 || !strings.Contains(differences[0], scenario.want) {
				t.Fatalf("got %v, want 1 difference containing %q", differences, scenario.want)
			}
		})
	}
}
