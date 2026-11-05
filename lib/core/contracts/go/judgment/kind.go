package judgment

// Kind names the question an answer answers.
type Kind string

// The three questions the plan asks its judge.
const (
	KindCriteria   Kind = "criteria"
	KindEpisodes   Kind = "episodes"
	KindSuggestion Kind = "suggestion"
)
