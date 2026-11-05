package run

// Judge is who judged a plan: the console's judge and its version, the
// planner model it asked, and the prompts the CLI sent, by version and by
// hash.
type Judge struct {
	Version       string `json:"version"`
	Model         string `json:"model"`
	PromptVersion string `json:"promptVersion"`
	PromptHash    string `json:"promptHash"`
}
