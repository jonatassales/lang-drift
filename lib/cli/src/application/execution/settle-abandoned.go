package execution

import (
	"github.com/jonatassales/tao/lib/cli/src/domain/finops"
	"github.com/jonatassales/tao/lib/cli/src/domain/harness"
)

// settleAbandoned closes the lap as abandoned, by whom and when, as
// finops.Abandon settles its money, and keeps the abandonment for the
// lap's ledger line. It returns what was released.
func (playing *playing) settleAbandoned(now harness.Timestamp) int {
	lap := &playing.lap
	lapMoney := finops.Lap{
		Reserved:   lap.Reserved,
		Observed:   lap.ObservedCents,
		Unmeasured: lap.Unmeasured,
	}

	abandonment := finops.Abandon(lapMoney, now.Time, playing.by)
	playing.abandonment = &abandonment

	settled := Step{
		Name: "Released",
		At:   now,
		Note: finops.Dollars(abandonment.Released) + " of the reservation returned to the budget",
	}
	if lap.Unmeasured {
		settled = Step{
			Name: "Held",
			At:   now,
			Note: "a call was cut off and its spend is unknown; the reservation is held",
		}
	}

	abandoned := Step{
		Name: "Abandoned",
		At:   now,
		Note: "by " + playing.by + ", while the run waited",
	}

	lap.Outcome = Abandoned
	lap.Steps = append(lap.Steps, abandoned, settled)

	return abandonment.Released
}
