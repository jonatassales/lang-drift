package domain_test

import (
	"testing"
	"time"

	"github.com/jonatassales/tao/apps/microservices/console-stubs/src/domain"
	"github.com/jonatassales/tao/apps/microservices/console-stubs/src/infrastructure/seed"
)

// TestSeededTicketsCoverEverySafeSeatState checks the Safe Seat data the
// stub serves: the seeded tickets cover the five states, only an override
// carries a reason, every record is illustrative, every workflow has one
// person-driven dev-ticket stage, and a sized plan reserves a part.
func TestSeededTicketsCoverEverySafeSeatState(t *testing.T) {
	store, err := seed.New(time.Now())
	if err != nil {
		t.Fatal(err)
	}

	tickets, err := store.Snapshot().Tickets(domain.TicketProductID)
	if err != nil {
		t.Fatal(err)
	}

	seen := map[string]bool{}
	for _, ticket := range tickets {
		sized := ticket.Plan != nil && ticket.Plan.Suggested != nil
		if sized && ticket.Plan.ReservedPart == nil {
			t.Errorf("%s: a sized plan has no reserved part", ticket.ID)
		}

		if ticket.Workflow == nil {
			if ticket.SafeSeat != nil {
				t.Errorf("%s: a Safe Seat record before the workflow", ticket.ID)
			}

			continue
		}

		persons := 0
		for _, stage := range ticket.Workflow.Stages {
			if stage.DrivenBy == "person" {
				persons++
			}
		}

		if persons != 1 {
			t.Errorf("%s: %d person-driven stages, want 1", ticket.ID, persons)
		}

		safeSeat := ticket.SafeSeat
		if safeSeat == nil || !safeSeat.Illustrative {
			t.Fatalf("%s: Safe Seat record %+v, want an illustrative one", ticket.ID, safeSeat)
		}

		hasReason := safeSeat.OverrideReason != nil
		if hasReason != (safeSeat.State == "overridden") {
			t.Errorf("%s: state %s with reason %v", ticket.ID, safeSeat.State, safeSeat.OverrideReason)
		}

		seen[safeSeat.State] = true
	}

	for _, state := range []string{"met", "needs-override", "overridden", "not-applicable", "unknown"} {
		if !seen[state] {
			t.Errorf("no seeded ticket is %s", state)
		}
	}
}
