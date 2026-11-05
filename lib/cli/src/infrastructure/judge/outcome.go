package judge

import "github.com/jonatassales/tao/lib/cli/src/application/plan"

// outcome is a question's outcome with its cause.
func outcome(kind string, status plan.AskStatus, reason string) plan.AskOutcome {
	return plan.AskOutcome{
		Kind:   kind,
		Status: status,
		Reason: reason,
	}
}
