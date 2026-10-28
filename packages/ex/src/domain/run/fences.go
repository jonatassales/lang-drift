package run

import "github.com/jonatassales/tao/packages/ex/src/domain/finops"

// RunFence is the memory's JSON fence: what the run was and what it measured,
// the episode the RAG indexes.
type RunFence struct {
	Run               string    `json:"run"`
	Sample            bool      `json:"sample,omitempty"`
	Ticket            *string   `json:"ticket"`
	Kind              string    `json:"kind"`
	Outcome           *string   `json:"outcome"`
	Laps              int       `json:"laps"`
	ComponentPaths    []string  `json:"componentPaths"`
	Stacks            []string  `json:"stacks"`
	MeasuredResources bool      `json:"measuredResources"`
	MeasuredTokens    bool      `json:"measuredTokens"`
	Observed          *Observed `json:"observed"`
	Facts             []string  `json:"facts"`
}

// Observed is what a lap measured: resource peaks and tokens.
type Observed struct {
	CPUMillisPeak   *int                    `json:"cpuMillisPeak"`
	MemoryMiBPeak   *int                    `json:"memoryMiBPeak"`
	TokensThousands *finops.TokensThousands `json:"tokensThousands"`
}

// LapFence is a lap's JSON fence: the choice, the outcome, the execution and
// the runtime it ran, and what it observed.
type LapFence struct {
	Run       string        `json:"run"`
	Lap       string        `json:"lap"`
	Sample    bool          `json:"sample,omitempty"`
	Choice    ChoiceFence   `json:"choice"`
	Outcome   string        `json:"outcome"`
	Gate      *string       `json:"gate,omitempty"`
	Execution *string       `json:"execution"`
	Runtime   *RuntimeFence `json:"runtime,omitempty"`
	Observed  *Observed     `json:"observed"`
}

// ChoiceFence is the model and effort of a lap.
type ChoiceFence struct {
	Model  string `json:"model"`
	Effort string `json:"effort"`
}

// RuntimeFence names the runtime revision a lap materialized.
type RuntimeFence struct {
	ID             string `json:"id"`
	Revision       int    `json:"revision"`
	MarkdownSHA256 string `json:"markdownSha256"`
}
