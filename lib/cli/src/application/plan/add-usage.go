package plan

import "github.com/jonatassales/tao/lib/cli/src/domain/finops"

// addUsage is the tokens two planning lines used together; a line that
// reported none adds nothing.
func addUsage(earlier, latest *finops.PlanningUsage) *finops.PlanningUsage {
	if earlier == nil {
		return latest
	}

	if latest == nil {
		return earlier
	}

	return &finops.PlanningUsage{
		InputTokens:  earlier.InputTokens + latest.InputTokens,
		OutputTokens: earlier.OutputTokens + latest.OutputTokens,
	}
}
