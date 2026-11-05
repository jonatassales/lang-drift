package stages_test

import (
	"testing"

	"github.com/jonatassales/tao/lib/core/contracts/go/stages"
)

// Each stage's sample is the report behind the simulation's stage line, and
// passes the check: "console › grilling   3 open questions settled",
// "console › to-spec    0011-SPEC revised · 3 acceptance criteria",
// "console › to-tickets 3 tickets · blocking edges recorded",
// "console › implement  tdd · 5 components · 14 files · boundary only · stops at the reserved part",
// "console › code-review 2 findings · fixed".
func TestCheckReportPassesTheSimulationsReports(t *testing.T) {
	t.Parallel()

	planning, err := stages.CheckReport(stages.Planning, readSample(t, "planning.report.sample.json"))
	settled := err == nil && planning.Planning.OpenQuestionsSettled == 3
	if !settled {
		t.Errorf("planning = %+v, %v", planning.Planning, err)
	}

	spec, err := stages.CheckReport(stages.Spec, readSample(t, "spec.report.sample.json"))
	revised := err == nil && spec.Spec.Revised && spec.Spec.AcceptanceCriteria == 3 && spec.Spec.ShortID() == "0011-SPEC"
	if !revised {
		t.Errorf("spec = %+v, %v", spec.Spec, err)
	}

	tickets, err := stages.CheckReport(stages.Tickets, readSample(t, "tickets.report.sample.json"))
	recorded := err == nil && tickets.Tickets.Tickets == 3 && tickets.Tickets.BlockingEdgesRecorded
	if !recorded {
		t.Errorf("tickets = %+v, %v", tickets.Tickets, err)
	}

	implement, err := stages.CheckReport(stages.Implement, readSample(t, "implement.report.sample.json"))
	counted := err == nil && implement.Implement.Method == "tdd" && implement.Implement.Components == 5 && implement.Implement.Files == 14
	stopped := counted && implement.Implement.BoundaryOnly && implement.Implement.StoppedAtReservedPart
	if !stopped {
		t.Errorf("implement = %+v, %v", implement.Implement, err)
	}

	review, err := stages.CheckReport(stages.Review, readSample(t, "review.report.sample.json"))
	fixed := err == nil && review.Review.Findings == 2 && review.Review.Fixed
	if !fixed {
		t.Errorf("review = %+v, %v", review.Review, err)
	}
}
