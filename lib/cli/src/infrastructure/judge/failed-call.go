package judge

import (
	"context"
	"errors"
	"net/http"

	"github.com/jonatassales/tao/lib/cli/src/application/plan"
)

// failedCall names a call that got no answer. Past its bound it timed out
// after the console had it, so the model may have run and its cost is
// unknown; otherwise the console could not be reached, or refused to be, as
// the transport says.
func failedCall(kind string, bounded context.Context, err error) (plan.Answer, plan.AskOutcome) {
	if errors.Is(bounded.Err(), context.DeadlineExceeded) {
		return plan.Answer{Unmeasured: true}, outcome(kind, plan.AskTimedOut, "the console did not answer in time")
	}

	var oversized *http.MaxBytesError
	if errors.As(err, &oversized) {
		return plan.Answer{Unmeasured: true}, outcome(kind, plan.AskMalformed, "the console's answer exceeds 4 MiB")
	}

	return plan.Answer{}, outcome(kind, plan.AskUnavailable, "the console could not be reached · "+err.Error())
}
