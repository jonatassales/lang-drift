package judge

// EpisodeReason is the judgment over one retrieved run.
type EpisodeReason struct {
	Run    string `json:"run"`
	Reason string `json:"reason"`
}
