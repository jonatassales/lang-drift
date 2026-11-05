package judge

// Answer is a parsed answer: its kind and the one field of that kind.
type Answer struct {
	Kind         Kind
	Sources      *SourcesAnswer
	Episodes     *EpisodesAnswer
	Suggestion   *SuggestionAnswer
	ReservedPart *ReservedPartAnswer
}
