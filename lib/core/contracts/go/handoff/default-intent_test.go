package handoff_test

import "github.com/jonatassales/tao/lib/core/contracts/go/handoff"

// defaultIntent is the simulation's default intent: Claude Sonnet 5.5 at
// medium effort, Planning and Review checked, authorized by the team.
func defaultIntent() handoff.Intent {
	return handoff.Intent{
		Ticket:    "0031-TICKET-book-a-demo-runtime",
		RuntimeID: "book-a-demo-runtime",
		Revision:  1,
		Model:     "claude-sonnet-5-5",
		Effort:    "medium",
		Human: handoff.Human{
			Planning: true,
			Review:   true,
		},
		Status: handoff.IntentAuthorized,
	}
}
