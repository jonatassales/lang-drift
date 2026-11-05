package judge

import (
	"strings"

	"github.com/jonatassales/tao/lib/cli/src/application/plan"
)

// refusedQuestion reads a 400: a model the console does not offer, with the
// ones it does; any other code means the CLI sent a question outside the
// contract.
func refusedQuestion(question plan.Question, payload []byte) plan.AskOutcome {
	refusal, err := refusalOf(payload)
	if err != nil {
		return outcome(question.Kind, plan.AskMalformed, err.Error())
	}

	if refusal.Code == "MODEL_NOT_OFFERED" {
		return outcome(question.Kind, plan.AskNotOffered, "the console does not offer "+question.Model+"; it offers "+strings.Join(refusal.Offered, ", "))
	}

	return outcome(question.Kind, plan.AskMalformed, "the console refused the question · "+refusal.Code+" "+refusal.Field)
}
