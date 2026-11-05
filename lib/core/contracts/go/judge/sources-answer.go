package judge

// SourcesAnswer is the judgment over the sources a plan read: a verdict per
// source, the facts they state, the guardrail criteria that apply and the
// figures given for CPU and memory.
type SourcesAnswer struct {
	Sources  []SourceVerdict `json:"sources"`
	Facts    []Fact          `json:"facts"`
	Criteria []Criterion     `json:"criteria"`
	Figures  []Figure        `json:"figures"`
}
