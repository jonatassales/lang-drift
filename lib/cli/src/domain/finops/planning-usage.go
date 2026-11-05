package finops

// PlanningUsage is what the planner's calls used, in tokens, as the console
// reported it.
type PlanningUsage struct {
	InputTokens  int64 `json:"inputTokens"`
	OutputTokens int64 `json:"outputTokens"`
}
