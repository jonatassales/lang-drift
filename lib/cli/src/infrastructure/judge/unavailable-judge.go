package judge

import "github.com/jonatassales/tao/lib/cli/src/application/plan"

// unavailableJudge reads a 503: the judge route says why its judge cannot
// answer, or the outcome names the body it could not read.
func unavailableJudge(kind string, payload []byte) plan.AskOutcome {
	refusal, err := refusalOf(payload)
	if err != nil {
		return outcome(kind, plan.AskUnavailable, "the console's judge is unavailable · "+err.Error())
	}

	return outcome(kind, plan.AskUnavailable, "the console's judge is unavailable · "+refusal.Reason)
}
