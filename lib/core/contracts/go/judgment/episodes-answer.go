package judgment

// EpisodesAnswer is the judge's answer to the episodes question: a verdict
// on every candidate the vector search retrieved.
type EpisodesAnswer struct {
	Episodes []EpisodeVerdict `json:"episodes"`
}
