package handoff_test

import (
	"encoding/json"
	"testing"

	"github.com/jonatassales/tao/lib/core/contracts/go/handoff"
)

// A section not revealed yet is left out of the record, not written as an
// empty or zero section.
func TestCardDataLeavesAnUnrevealedSectionOut(t *testing.T) {
	t.Parallel()

	cards := handoff.CardData{
		Ticket: &handoff.TicketCard{},
	}

	encoded, err := json.Marshal(cards)
	if err != nil || string(encoded) != `{"ticket":{}}` {
		t.Errorf("encoded = %s, %v", encoded, err)
	}
}
