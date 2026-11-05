package plan

import "github.com/jonatassales/tao/lib/cli/src/domain/run"

// judgeRecord names who judged a plan: the console's judge and its version,
// the planner model and the prompts by version and by hash; a plan no
// question reached names none.
func judgeRecord(judge *judge) *run.Judge {
	if !judge.paid() {
		return nil
	}

	return &run.Judge{
		Version:       judge.version,
		Model:         judge.planner,
		PromptVersion: PromptVersion,
		PromptHash:    promptHash(),
	}
}
