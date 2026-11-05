package judge

// Request is one question, as the planner sends it to the console.
type Request struct {
	Kind          Kind   `json:"kind"`
	Model         string `json:"model"`
	System        string `json:"system"`
	Prompt        string `json:"prompt"`
	MaxTokens     int    `json:"maxTokens"`
	PromptVersion string `json:"promptVersion"`
}
