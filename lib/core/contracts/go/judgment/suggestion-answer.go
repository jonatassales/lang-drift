package judgment

// SuggestionAnswer is the judge's answer to the suggestion question: the
// model and the effort, and why, as "weighing effort · …" and "weighing
// model · …" print it.
type SuggestionAnswer struct {
	Model        string `json:"model"`
	Effort       string `json:"effort"`
	EffortReason string `json:"effortReason"`
	ModelReason  string `json:"modelReason"`
}
