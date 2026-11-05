// Package judge is the contract between the planner of the ex CLI and the
// judge the platform serves at <console-url>/api/judge: the question the
// planner asks, the response the console returns and the typed answer of
// each kind. The JSON Schemas under schema/ define the same shapes for the
// console; both sides test against them.
package judge

// Kind is what a question asks for.
type Kind string

// The kinds of question the planner asks, one to one with its own names.
const (
	KindSources      Kind = "sources"
	KindEpisodes     Kind = "episodes"
	KindSuggestion   Kind = "suggestion"
	KindReservedPart Kind = "reserved-part"
)

// Request is one question, as the planner sends it to the console.
type Request struct {
	Kind          Kind   `json:"kind"`
	Model         string `json:"model"`
	System        string `json:"system"`
	Prompt        string `json:"prompt"`
	MaxTokens     int    `json:"maxTokens"`
	PromptVersion string `json:"promptVersion"`
}

// Usage is what a judgment consumed, in tokens.
type Usage struct {
	InputTokens  int64 `json:"inputTokens"`
	OutputTokens int64 `json:"outputTokens"`
}

// Cost is what the platform charged for a judgment and where the figure
// comes from.
type Cost struct {
	Cents        int64  `json:"cents"`
	Provenance   string `json:"provenance"`
	Illustrative bool   `json:"illustrative"`
}

// Response is the console's answer to one question. Text is the model's
// answer, already validated by the console against the kind's schema; the
// planner validates it again with ParseAnswer before using it.
type Response struct {
	Text         string `json:"text"`
	Usage        Usage  `json:"usage"`
	Refused      bool   `json:"refused"`
	Reason       string `json:"reason"`
	Cost         Cost   `json:"cost"`
	JudgeVersion string `json:"judgeVersion"`
}

// Answer is a parsed answer: its kind and the one field of that kind.
type Answer struct {
	Kind         Kind
	Sources      *SourcesAnswer
	Episodes     *EpisodesAnswer
	Suggestion   *SuggestionAnswer
	ReservedPart *ReservedPartAnswer
}

// SourcesAnswer is the judgment over the sources a plan read: a verdict per
// source, the facts they state, the guardrail criteria that apply and the
// figures given for CPU and memory.
type SourcesAnswer struct {
	Sources  []SourceVerdict `json:"sources"`
	Facts    []Fact          `json:"facts"`
	Criteria []Criterion     `json:"criteria"`
	Figures  []Figure        `json:"figures"`
}

// SourceVerdict says whether one source is input for planning and why.
type SourceVerdict struct {
	Path       string  `json:"path"`
	Verdict    string  `json:"verdict"`
	Reason     string  `json:"reason"`
	Confidence float64 `json:"confidence"`
	Role       string  `json:"role"`
	Summary    string  `json:"summary"`
}

// Fact is one statement the material makes, with its citation.
type Fact struct {
	Topic  string `json:"topic"`
	Text   string `json:"text"`
	Source string `json:"source"`
}

// Criterion is one guardrail criterion that applies to the run.
type Criterion struct {
	ID      string `json:"id"`
	Applies string `json:"applies"`
}

// Figure is one figure the material gives for CPU and memory.
type Figure struct {
	Topic     string `json:"topic"`
	Fact      string `json:"fact"`
	Source    string `json:"source"`
	CPUMillis int    `json:"cpuMillis"`
	MemoryMiB int    `json:"memoryMiB"`
}

// EpisodesAnswer says, per earlier run retrieved, why it is comparable.
type EpisodesAnswer struct {
	Episodes []EpisodeReason `json:"episodes"`
}

// EpisodeReason is the judgment over one retrieved run.
type EpisodeReason struct {
	Run    string `json:"run"`
	Reason string `json:"reason"`
}

// SuggestionAnswer is the suggested model and effort for the ticket's
// workflow, with the reasons.
type SuggestionAnswer struct {
	Model        string   `json:"model"`
	Effort       string   `json:"effort"`
	EffortReason string   `json:"effortReason"`
	ModelReason  string   `json:"modelReason"`
	Reasons      []string `json:"reasons"`
}

// ReservedPartAnswer is the part of the ticket the judge proposes a person
// writes by hand under Safe Seat: a problem with its own goal and acceptance
// criteria, located by project-relative paths. Empty paths mean no part
// qualifies, and the reason says why; ExcludedKind names the kind the policy
// never reserves when the best candidate was one of those.
type ReservedPartAnswer struct {
	Paths             []string `json:"paths"`
	Goal              string   `json:"goal"`
	Acceptance        []string `json:"acceptance"`
	Category          *string  `json:"category"`
	PlannedLines      int      `json:"plannedLines"`
	TotalPlannedLines int      `json:"totalPlannedLines"`
	ExcludedKind      *string  `json:"excludedKind"`
	Reason            string   `json:"reason"`
	Confidence        float64  `json:"confidence"`
}
