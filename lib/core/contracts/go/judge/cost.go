package judge

// Cost is what the platform charged for a judgment and where the figure
// comes from.
type Cost struct {
	Cents        int64  `json:"cents"`
	Provenance   string `json:"provenance"`
	Illustrative bool   `json:"illustrative"`
}
