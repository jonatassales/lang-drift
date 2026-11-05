package stages

import "fmt"

// CheckReport reads the report a stage's skill wrote and refuses one outside
// the stage's shape: a missing field, a field the stage does not define, a
// null value, a value of the wrong type or a rule broken, naming the field.
func CheckReport(stage Stage, text string) (Report, error) {
	report := Report{Stage: stage}

	var err error

	switch stage {
	case Planning:
		report.Planning, err = checked(stage, text, checkPlanning)
	case Spec:
		report.Spec, err = checked(stage, text, checkSpec)
	case Tickets:
		report.Tickets, err = checked(stage, text, checkTickets)
	case Implement:
		report.Implement, err = checked(stage, text, checkImplement)
	case Review:
		report.Review, err = checked(stage, text, checkReview)
	default:
		return Report{}, fmt.Errorf("stages: no report for the stage %q", stage)
	}

	if err != nil {
		return Report{}, fmt.Errorf("stages: the %s report is refused: %w", stage, err)
	}

	return report, nil
}
