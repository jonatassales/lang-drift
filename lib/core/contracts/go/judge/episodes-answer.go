package judge

// EpisodesAnswer says, per earlier run retrieved, why it is comparable.
type EpisodesAnswer struct {
	Episodes []EpisodeReason `json:"episodes"`
}
