package handoff_test

import (
	"slices"
	"testing"

	"github.com/jonatassales/tao/lib/core/contracts/go/handoff"
)

// The default intent passes, as a proposed one does; every rule it can break
// is a problem naming its field, and an intent breaking several returns them
// all: Planning unchecked, an unknown status and a revision of 0 are three.
func TestValidateIntentReturnsEveryProblem(t *testing.T) {
	t.Parallel()

	for _, scenario := range []struct {
		name   string
		change func(*handoff.Intent)
		fields []string
	}{
		{"the default intent", func(*handoff.Intent) {}, nil},
		{"a proposed intent", func(intent *handoff.Intent) { intent.Status = handoff.IntentProposed }, nil},
		{"no ticket", func(intent *handoff.Intent) { intent.Ticket = "" }, []string{"ticket"}},
		{"no runtime id", func(intent *handoff.Intent) { intent.RuntimeID = "" }, []string{"runtimeId"}},
		{"revision 0", func(intent *handoff.Intent) { intent.Revision = 0 }, []string{"revision"}},
		{"no model", func(intent *handoff.Intent) { intent.Model = "" }, []string{"model"}},
		{"another effort", func(intent *handoff.Intent) { intent.Effort = "max" }, []string{"effort"}},
		{"planning unchecked", func(intent *handoff.Intent) { intent.Human.Planning = false }, []string{"human.planning"}},
		{"an unknown status", func(intent *handoff.Intent) { intent.Status = "merged" }, []string{"status"}},
		{"three at once", func(intent *handoff.Intent) {
			intent.Human.Planning = false
			intent.Status = "merged"
			intent.Revision = 0
		}, []string{"revision", "human.planning", "status"}},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			t.Parallel()

			intent := defaultIntent()
			scenario.change(&intent)

			var fields []string
			for _, problem := range intent.Validate() {
				fields = append(fields, problem.Field)
			}

			if !slices.Equal(fields, scenario.fields) {
				t.Errorf("Validate named %v, want %v", fields, scenario.fields)
			}
		})
	}
}
