package judge

// Kind is what a question asks for.
type Kind string

// The kinds of question the planner asks, one to one with its own names.
const (
	KindSources      Kind = "sources"
	KindEpisodes     Kind = "episodes"
	KindSuggestion   Kind = "suggestion"
	KindReservedPart Kind = "reserved-part"
)
