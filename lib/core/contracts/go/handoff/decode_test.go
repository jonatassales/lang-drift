package handoff_test

import (
	"testing"

	"github.com/jonatassales/tao/lib/core/contracts/go/handoff"
)

// The recorded intent of the simulation's default run decodes whole: the
// ticket, the runtime, revision 1, Claude Sonnet 5.5 at medium effort, Planning
// and Review checked, authorized. It feeds "console › intent
// claude-sonnet-5-5 · effort medium · runtime.md revision 1" and the HUMAN IN
// THE LOOP ticks.
func TestDecodeReadsTheSimulationsIntent(t *testing.T) {
	t.Parallel()

	text := `{
		"ticket": "0031-TICKET-book-a-demo-runtime",
		"runtimeId": "book-a-demo-runtime",
		"revision": 1,
		"model": "claude-sonnet-5-5",
		"effort": "medium",
		"human": {"planning": true, "spec": false, "tickets": false, "implement": false, "review": true},
		"status": "authorized"
	}`

	intent, err := handoff.Decode[handoff.Intent](text)
	if err != nil {
		t.Fatal(err)
	}

	want := defaultIntent()
	if intent != want {
		t.Errorf("intent = %+v, want %+v", intent, want)
	}
}

// The plan's card data carries only what has been revealed: with the
// EVIDENCE line "runtime.md draft · revision 1 · proposed to the console"
// written and nothing else, every other section stays absent, as its card
// reads "Not set yet".
func TestDecodeReadsTheRevealedSectionsOnly(t *testing.T) {
	t.Parallel()

	text := `{"plan": {"evidence": {"file": "runtime.md", "state": "draft", "revision": 1, "proposedTo": "console"}}}`

	cards, err := handoff.Decode[handoff.CardData](text)
	if err != nil {
		t.Fatal(err)
	}

	unrevealed := cards.Ticket == nil && cards.Analysis == nil && cards.Workflow == nil
	evidenceOnly := cards.Plan != nil && cards.Plan.Evidence != nil && cards.Plan.Sizing == nil
	if !unrevealed || !evidenceOnly {
		t.Fatalf("cards = %+v", cards)
	}

	evidence := *cards.Plan.Evidence
	present := evidence.File != nil && evidence.State != nil && evidence.Revision != nil && evidence.ProposedTo != nil
	written := present && *evidence.File == "runtime.md" && *evidence.State == "draft" && *evidence.Revision == 1 && *evidence.ProposedTo == "console"
	if !written {
		t.Errorf("evidence = %+v, plan = %+v", evidence, *cards.Plan)
	}
}
