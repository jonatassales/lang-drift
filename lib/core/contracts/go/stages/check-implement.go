package stages

import "errors"

// testDriven is the one method the implement line shows.
const testDriven = "tdd"

// checkImplement accepts an implement report whose method is tdd and whose
// components and files are counted.
func checkImplement(report ImplementReport) error {
	if report.Method != testDriven {
		return FieldError{
			Field: "method",
			Rule:  "must be tdd",
		}
	}

	return errors.Join(checkCount("components", report.Components), checkCount("files", report.Files))
}
