package handoff

// AnalysisCard is the Analysis card's data: the documents read, the episodes
// cited and the suggestion.
type AnalysisCard struct {
	Documents  *Documents  `json:"documents,omitempty"`
	Episodes   *Episodes   `json:"episodes,omitempty"`
	Suggestion *Suggestion `json:"suggestion,omitempty"`
}
