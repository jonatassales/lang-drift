package handoff_test

import (
	"testing"

	"github.com/jonatassales/tao/lib/core/contracts/go/handoff"
)

// A record with a field its shape does not define is refused, at the top
// and inside a section, and so is a text holding more than one record.
func TestDecodeRefusesAFieldOutsideTheRecord(t *testing.T) {
	t.Parallel()

	for _, text := range []string{
		`{"ticket": "0031-TICKET-book-a-demo-runtime", "lease": "1h"}`,
		`{"human": {"planning": true, "merge": true}}`,
		`{} {}`,
	} {
		if _, err := handoff.Decode[handoff.Intent](text); err == nil {
			t.Errorf("the intent %s was accepted", text)
		}
	}

	for _, text := range []string{
		`{"plan": {"evidence": {"file": "runtime.md", "signedBy": "owner"}}}`,
		`{"analysis": {"documents": {"cost": 1}}}`,
	} {
		if _, err := handoff.Decode[handoff.CardData](text); err == nil {
			t.Errorf("the card data %s was accepted", text)
		}
	}

	if _, err := handoff.Decode[handoff.RunRecord](`{"run": null, "channel": "mcp"}`); err == nil {
		t.Error("a run record with a channel was accepted")
	}
}
