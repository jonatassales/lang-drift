package domain

// measuredSafeSeatStates are the states a finished run's Safe Seat record
// cycles through; a run still going reads unknown.
var measuredSafeSeatStates = []string{"met", "needs-override", "overridden", "not-applicable"}

// overrideReason is the reason the seeded overrides record.
const overrideReason = "The person was on call this week; the reserved part moves to the next ticket."

// generateSafeSeat gives a ticket in its workflow a Safe Seat record. A run
// still going, or paused at the gate, has nothing measured yet and reads
// unknown; only an overridden record carries a reason.
func generateSafeSeat(number int, outcome string) *TicketSafeSeat {
	safeSeat := &TicketSafeSeat{
		State:        "unknown",
		Illustrative: true,
	}

	finished := outcome == "readyForReview" || outcome == "readyToMerge"
	if !finished {
		return safeSeat
	}

	safeSeat.State = measuredSafeSeatStates[(number/16)%len(measuredSafeSeatStates)]
	if isEarlyAccess(number) {
		safeSeat.State = "met"
	}

	if safeSeat.State == "overridden" {
		safeSeat.OverrideReason = stringPointer(overrideReason)
	}

	return safeSeat
}
