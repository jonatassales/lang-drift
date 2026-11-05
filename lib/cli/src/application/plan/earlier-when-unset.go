package plan

import "github.com/jonatassales/tao/lib/cli/src/domain/finops"

// earlierWhenUnset is the provenance and the illustrative mark of a summed
// planning line: the latest plan's, or the earlier plans' when the latest
// paid nothing and so carries neither.
func earlierWhenUnset(latest, earlier finops.PlanningLine) (*string, *bool) {
	if latest.Provenance != nil {
		return latest.Provenance, latest.Illustrative
	}

	return earlier.Provenance, earlier.Illustrative
}
