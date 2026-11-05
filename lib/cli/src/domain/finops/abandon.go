package finops

import "time"

// Abandon settles a lap a person let go while it waited. The spend observed
// before the wait is kept; the reservation it did not spend is released,
// never more than was reserved, all of it when nothing was observed; spend
// above the reservation releases nothing and stays visible as overspend. A
// lap with a call cut off releases nothing and names no overspend: what it
// spent is unknown, and unknown is never zero, so the reservation is held and
// the observed figure is a floor. The closing line names the outcome, the
// person and the time, in RFC 3339 and UTC.
func Abandon(lap Lap, at time.Time, by string) Abandonment {
	abandonment := Abandonment{
		Released:  max(0, lap.Reserved-lap.Observed),
		Observed:  lap.Observed,
		Overspend: max(0, lap.Observed-lap.Reserved),
		Closing: LapClosing{
			Outcome: LapAbandoned,
			By:      by,
			At:      at.UTC().Format(time.RFC3339),
		},
	}

	if lap.Unmeasured {
		abandonment.Released = 0
		abandonment.Overspend = 0
	}

	return abandonment
}
