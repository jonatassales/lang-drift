package stages

import "github.com/jonatassales/tao/lib/core/contracts/go/configuration"

// checkSpec accepts a spec report whose spec is a clean project-relative
// path and whose criteria are counted.
func checkSpec(report SpecReport) error {
	if err := configuration.CheckRelativePath(report.Spec); err != nil {
		return FieldError{
			Field: "spec",
			Rule:  "must be the spec's clean project-relative path: " + err.Error(),
		}
	}

	return checkCount("acceptanceCriteria", report.AcceptanceCriteria)
}
