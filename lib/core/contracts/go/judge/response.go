package judge

// Response is the console's answer to one question. Text is the model's
// answer, already validated by the console against the kind's schema; the
// planner validates it again with ParseAnswer before using it.
type Response struct {
	Text         string `json:"text"`
	Usage        Usage  `json:"usage"`
	Refused      bool   `json:"refused"`
	Reason       string `json:"reason"`
	Cost         Cost   `json:"cost"`
	JudgeVersion string `json:"judgeVersion"`
}
