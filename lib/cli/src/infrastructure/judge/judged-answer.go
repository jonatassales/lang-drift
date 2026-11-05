package judge

import (
	"bytes"
	"encoding/json"

	"github.com/jonatassales/tao/lib/cli/src/application/plan"
	contract "github.com/jonatassales/tao/lib/core/contracts/go/judge"
)

// judgedAnswer reads a 200 strictly as the contract's response. A body
// outside it leaves what the call cost unknown.
func judgedAnswer(kind string, payload []byte) (plan.Answer, plan.AskOutcome) {
	decoder := json.NewDecoder(bytes.NewReader(payload))
	decoder.DisallowUnknownFields()

	var response contract.Response
	if err := decoder.Decode(&response); err != nil {
		return plan.Answer{Unmeasured: true}, outcome(kind, plan.AskMalformed, "the console's answer is outside the contract: "+err.Error())
	}

	if broken := brokenResponse(response); broken != "" {
		return plan.Answer{Unmeasured: true}, outcome(kind, plan.AskMalformed, "the console's answer is outside the contract: "+broken)
	}

	answer := plan.Answer{
		Text:         response.Text,
		InputTokens:  response.Usage.InputTokens,
		OutputTokens: response.Usage.OutputTokens,
		Cents:        int(response.Cost.Cents),
		Provenance:   response.Cost.Provenance,
		Illustrative: response.Cost.Illustrative,
		JudgeVersion: response.JudgeVersion,
	}

	if response.Refused {
		return answer, outcome(kind, plan.AskRefused, response.Reason)
	}

	return answer, outcome(kind, plan.AskAnswered, "")
}
