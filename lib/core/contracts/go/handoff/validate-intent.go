package handoff

import (
	"slices"
	"strings"

	"github.com/jonatassales/tao/lib/core/contracts/go/configuration"
)

// Validate returns every problem of a recorded intent, in field order: a
// missing ticket, runtime id or model, a revision below 1, an effort the
// Workflow's EFFORT choice does not offer, Planning without a person, since
// the simulation's Planning cannot be unchecked, and a status other than
// proposed or authorized. An intent with no problem is one a run may apply.
func (intent Intent) Validate() []Problem {
	var problems []Problem

	if intent.Ticket == "" {
		problems = append(problems, problemAt("ticket", "is missing"))
	}

	if intent.RuntimeID == "" {
		problems = append(problems, problemAt("runtimeId", "is missing"))
	}

	if intent.Revision < 1 {
		problems = append(problems, problemAt("revision", "must be at least 1"))
	}

	if intent.Model == "" {
		problems = append(problems, problemAt("model", "is missing"))
	}

	if !slices.Contains(configuration.FenceEfforts, intent.Effort) {
		problems = append(problems, problemAt("effort", "must be one of "+strings.Join(configuration.FenceEfforts, ", ")))
	}

	if !intent.Human.Planning {
		problems = append(problems, problemAt("human.planning", "must have a person in the loop"))
	}

	recordedStatus := intent.Status == IntentProposed || intent.Status == IntentAuthorized
	if !recordedStatus {
		problems = append(problems, problemAt("status", "must be proposed or authorized"))
	}

	return problems
}
