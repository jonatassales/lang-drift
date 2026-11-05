package judgment

import (
	"slices"
	"strings"

	"github.com/jonatassales/tao/lib/core/contracts/go/configuration"
)

// checkSuggestion accepts a suggestion that names a model, one of the
// efforts, and a reason for each.
func checkSuggestion(answer SuggestionAnswer) error {
	switch {
	case answer.Model == "":
		return FieldError{
			Field: "model",
			Rule:  "must name a model",
		}
	case !slices.Contains(configuration.FenceEfforts, answer.Effort):
		return FieldError{
			Field: "effort",
			Rule:  "must be one of " + strings.Join(configuration.FenceEfforts, ", "),
		}
	case answer.EffortReason == "":
		return FieldError{
			Field: "effortReason",
			Rule:  "must weigh the effort",
		}
	case answer.ModelReason == "":
		return FieldError{
			Field: "modelReason",
			Rule:  "must weigh the model",
		}
	}

	return nil
}
