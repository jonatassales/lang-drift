package judge

// Criterion is one guardrail criterion that applies to the run.
type Criterion struct {
	ID      string `json:"id"`
	Applies string `json:"applies"`
}
