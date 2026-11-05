package stages_test

import (
	"errors"
	"testing"

	"github.com/jonatassales/tao/lib/core/contracts/go/stages"
)

// A report outside its stage's shape is refused naming the field: a missing
// field, a null one, a field the stage does not define, a negative count,
// another method, a spec path that climbs and a value of the wrong type. A
// stage with no report is refused.
func TestCheckReportRefusesAReportOutsideItsShape(t *testing.T) {
	t.Parallel()

	for _, scenario := range []struct {
		stage stages.Stage
		text  string
		field string
	}{
		{stages.Planning, `{}`, "openQuestionsSettled"},
		{stages.Review, `{"findings":null,"fixed":true}`, "findings"},
		{stages.Review, `{"findings":2,"fixed":true,"cost":1}`, "cost"},
		{stages.Tickets, `{"tickets":-1,"blockingEdgesRecorded":true}`, "tickets"},
		{stages.Implement, `{"method":"vibes","components":5,"files":14,"boundaryOnly":true,"stoppedAtReservedPart":true}`, "method"},
		{stages.Spec, `{"spec":"../0011-SPEC.spec.md","revised":true,"acceptanceCriteria":3}`, "spec"},
		{stages.Review, `{"findings":"two","fixed":true}`, "findings"},
		{stages.Review, `[]`, "report"},
	} {
		_, err := stages.CheckReport(scenario.stage, scenario.text)

		var refusal stages.FieldError
		namesTheField := errors.As(err, &refusal) && refusal.Field == scenario.field
		if !namesTheField {
			t.Errorf("%s %s gave %v, want a refusal naming %s", scenario.stage, scenario.text, err, scenario.field)
		}
	}

	if _, err := stages.CheckReport("merge", `{}`); err == nil {
		t.Error("a stage with no report was accepted")
	}
}
