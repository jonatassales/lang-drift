package judgment

// EpisodeVerdict is the judge's verdict on one retrieved episode: its id,
// kept or dropped, and the reason, as an EPISODES row shows it.
type EpisodeVerdict struct {
	ID      string `json:"id"`
	Verdict string `json:"verdict"`
	Reason  string `json:"reason"`
}
