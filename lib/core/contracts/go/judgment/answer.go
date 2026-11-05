package judgment

// Answer is one parsed answer: its kind, and the typed answer of that kind.
type Answer struct {
	Kind       Kind
	Criteria   *CriteriaAnswer
	Episodes   *EpisodesAnswer
	Suggestion *SuggestionAnswer
}
