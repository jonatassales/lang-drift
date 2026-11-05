package judge

// SourceVerdict says whether one source is input for planning and why.
type SourceVerdict struct {
	Path       string  `json:"path"`
	Verdict    string  `json:"verdict"`
	Reason     string  `json:"reason"`
	Confidence float64 `json:"confidence"`
	Role       string  `json:"role"`
	Summary    string  `json:"summary"`
}
