package judge

// Usage is what a judgment consumed, in tokens.
type Usage struct {
	InputTokens  int64 `json:"inputTokens"`
	OutputTokens int64 `json:"outputTokens"`
}
