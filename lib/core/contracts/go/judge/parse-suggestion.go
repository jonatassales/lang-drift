package judge

import "strconv"

// The efforts a suggestion may name, and how many reasons it may give.
var efforts = []string{"low", "medium", "high"}

const maxSuggestionReasons = 3

// parseSuggestion reads a suggestion answer: a model, an effort from the
// set, the reasons for both and at most three further reasons.
func parseSuggestion(text []byte) (*SuggestionAnswer, error) {
	var answer SuggestionAnswer
	if err := decodeObject(text, "", []string{"model", "effort", "effortReason", "modelReason", "reasons"}, &answer); err != nil {
		return nil, err
	}

	if err := requireText("model", answer.Model); err != nil {
		return nil, err
	}
	if err := requireOneOf("effort", answer.Effort, efforts); err != nil {
		return nil, err
	}
	if err := requireText("effortReason", answer.EffortReason); err != nil {
		return nil, err
	}
	if err := requireText("modelReason", answer.ModelReason); err != nil {
		return nil, err
	}
	if len(answer.Reasons) > maxSuggestionReasons {
		return nil, refuse("reasons", "must hold at most "+strconv.Itoa(maxSuggestionReasons)+" reasons")
	}

	for index, reason := range answer.Reasons {
		if err := requireText("reasons["+strconv.Itoa(index)+"]", reason); err != nil {
			return nil, err
		}
	}

	return &answer, nil
}
