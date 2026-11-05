package judge

// SuggestionAnswer is the suggested model and effort for the ticket's
// workflow, with the reasons.
type SuggestionAnswer struct {
	Model        string   `json:"model"`
	Effort       string   `json:"effort"`
	EffortReason string   `json:"effortReason"`
	ModelReason  string   `json:"modelReason"`
	Reasons      []string `json:"reasons"`
}
