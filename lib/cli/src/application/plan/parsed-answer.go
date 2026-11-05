package plan

import contract "github.com/jonatassales/tao/lib/core/contracts/go/judge"

// parsedAnswer reads an answered question against the judge contract and the
// planner's own check, as the kind of question that was asked, whatever kind
// the judge says it answered. An answer that does not fit turns the outcome
// malformed, so it is refused like any other cause. The gates' statuses only
// come before a question: a judge that reports one after it gave no answer,
// and is refused the same way. Any other outcome passes through, under the
// asked kind.
func parsedAnswer(kind string, answer Answer, outcome AskOutcome, check func(contract.Answer) error) (contract.Answer, AskOutcome) {
	outcome.Kind = kind

	reportedAGateStatus := outcome.Status == AskCleared || outcome.Status == AskPaused
	if reportedAGateStatus {
		return contract.Answer{}, AskOutcome{
			Kind:   kind,
			Status: AskMalformed,
			Reason: "the judge reported " + string(outcome.Status) + " after the question, a status only the money gates give",
		}
	}

	if outcome.Status != AskAnswered {
		return contract.Answer{}, outcome
	}

	parsed, err := contract.ParseAnswer(contract.Kind(kind), answer.Text)
	if err == nil {
		err = check(parsed)
	}

	if err != nil {
		return contract.Answer{}, AskOutcome{
			Kind:   kind,
			Status: AskMalformed,
			Reason: err.Error(),
		}
	}

	return parsed, outcome
}
