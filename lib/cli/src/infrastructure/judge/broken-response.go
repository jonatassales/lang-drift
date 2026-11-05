package judge

import contract "github.com/jonatassales/tao/lib/core/contracts/go/judge"

// brokenResponse names the rule of the contract's response a decoded answer
// breaks, or nothing: the judge's version is named, the cost is the
// console's, and no count or amount is negative.
func brokenResponse(response contract.Response) string {
	negative := response.Usage.InputTokens < 0 || response.Usage.OutputTokens < 0 || response.Cost.Cents < 0

	switch {
	case response.JudgeVersion == "":
		return "judgeVersion is empty"
	case response.Cost.Provenance != "console":
		return "cost.provenance is " + response.Cost.Provenance + ", not console"
	case negative:
		return "a usage count or the cost is negative"
	}

	return ""
}
