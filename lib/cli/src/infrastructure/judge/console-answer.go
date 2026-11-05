package judge

import (
	"net/http"
	"strconv"

	"github.com/jonatassales/tao/lib/cli/src/application/plan"
)

// consoleAnswer reads what the judge route answered one question: a judgment
// or the model's refusal, priced, from a 200; the secret refused from a 401;
// a question the route would not take from a 400; the judge unavailable from
// a 503; a request the route's interceptors refuse, which the CLI and the
// console should never disagree on, from a 405, 413 or 415. Any other server
// error may come after the model ran, so its cost is unknown. Usage and cost
// are kept whenever the route reports them, so a refused call is accounted.
func consoleAnswer(question plan.Question, status int, payload []byte) (plan.Answer, plan.AskOutcome) {
	switch status {
	case http.StatusOK:
		return judgedAnswer(question.Kind, payload)
	case http.StatusUnauthorized:
		return plan.Answer{}, outcome(question.Kind, plan.AskUnauthorized, "the console refused RUNTIME_SECRET")
	case http.StatusBadRequest:
		return plan.Answer{}, refusedQuestion(question, payload)
	case http.StatusServiceUnavailable:
		return plan.Answer{}, unavailableJudge(question.Kind, payload)
	case http.StatusMethodNotAllowed, http.StatusRequestEntityTooLarge, http.StatusUnsupportedMediaType:
		return plan.Answer{}, outcome(question.Kind, plan.AskMalformed, "the console refused the request with HTTP "+strconv.Itoa(status))
	}

	unmeasured := plan.Answer{Unmeasured: status >= http.StatusInternalServerError}

	return unmeasured, outcome(question.Kind, plan.AskUnavailable, "the console answered HTTP "+strconv.Itoa(status))
}
