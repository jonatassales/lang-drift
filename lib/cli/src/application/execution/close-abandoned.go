package execution

import "github.com/jonatassales/tao/lib/cli/src/domain/finops"

// closeAbandoned writes an abandoned lap's settlement on its ledger line:
// what it kept and released, the overspend, and the closing line. The
// variance stays unknown when a call's spend is.
func closeAbandoned(line *finops.LapLine, lap Lap, abandonment finops.Abandonment) {
	line.Cents.Observed = abandonment.Observed
	line.Cents.Recovered = abandonment.Released
	line.Cents.Overspend = abandonment.Overspend
	line.Closed = &abandonment.Closing

	if !lap.Unmeasured {
		line.Cents.Variance = number(abandonment.Observed - lap.EstimateCents)
	}
}
