package plan

import contract "github.com/jonatassales/tao/lib/core/contracts/go/judge"

// parsedAnswer reads an answered question against the judge contract and the
// planner's own check. An answer that does not fit turns the outcome
// malformed, so it is refused like any other cause; any other outcome passes
// through untouched.
func parsedAnswer(answer Answer, outcome AskOutcome, check func(contract.Answer) error) (contract.Answer, AskOutcome) {
	if outcome.Status != AskAnswered {
		return contract.Answer{}, outcome
	}

	parsed, err := contract.ParseAnswer(contract.Kind(outcome.Kind), answer.Text)
	if err == nil {
		err = check(parsed)
	}

	if err != nil {
		return contract.Answer{}, AskOutcome{
			Kind:   outcome.Kind,
			Status: AskMalformed,
			Reason: err.Error(),
		}
	}

	return parsed, outcome
}
