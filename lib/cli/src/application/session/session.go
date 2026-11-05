// Package session is the event stream of a plan or a run: every line the
// terminal prints is one event the use case emits, with the figures behind
// it. The words belong to the presentation; the ids are the website's.
package session

// Kind is the kind of line an event prints as.
type Kind string

// The kinds of line.
const (
	Command Kind = "command"
	Comment Kind = "comment"
	Output  Kind = "output"
	Success Kind = "success"
	Warning Kind = "warning"
	Error   Kind = "error"
)

// The sessions.
const (
	PlanSession = "plan"
	RunSession  = "run"
)

// Data carries the figures a line shows, by name.
type Data map[string]any

// Event is one line of a session, before it is rendered.
type Event struct {
	Session string
	ID      string
	Kind    Kind
	Data    Data
}

// Sink receives the events of a session as they happen.
type Sink interface {
	Emit(Event)
}

// Discard is a sink that drops every event.
type Discard struct{}

// Emit drops the event.
func (Discard) Emit(Event) {}

// The plan session's line ids, as the website names them. Lines the website
// prints once per document or component carry one generic id here, with
// the document or component in the data.
const (
	PlanCommand      = "plan"
	PlanScope        = "scope"
	ReadTicket       = "readTicket"
	Links            = "links"
	Analysis         = "analysis"
	ReadDocument     = "readDocument"
	Reads            = "reads"
	Guardrail        = "guardrail"
	JudgeGuardrail   = "judgeGuardrail"
	Bounds           = "bounds"
	Discrepancy      = "discrepancy"
	Budget           = "budget"
	Memory           = "memory"
	Embed            = "embed"
	VectorSearch     = "vectorSearch"
	Candidates       = "candidates"
	JudgeRelevance   = "judgeRelevance"
	Drop             = "drop"
	Episodes         = "episodes"
	ReasonResources  = "reasonResources"
	WeighEffort      = "weighEffort"
	WeighModel       = "weighModel"
	Suggest          = "suggest"
	Planning         = "planning"
	Resolve          = "resolve"
	Inspect          = "inspect"
	Component        = "component"
	CrossCheck       = "crossCheck"
	Size             = "size"
	Limits           = "limits"
	Price            = "price"
	Validate         = "validate"
	Draft            = "draft"
	Handoff          = "handoff"
	JudgeUnavailable = "judgeUnavailable"
)

// The run session's line ids.
const (
	RunCommand         = "command"
	RunScope           = "scope"
	Intent             = "intent"
	Estimate           = "estimate"
	RunBudget          = "budget"
	GateOpened         = "gateOpened"
	GatePaused         = "gatePaused"
	Paused             = "paused"
	Waiting            = "waiting"
	Confirmed          = "confirmed"
	Stage              = "stage"
	Hook               = "hook"
	RuntimeReady       = "runtimeReady"
	Completed          = "completed"
	Pushed             = "pushed"
	NotPushed          = "notPushed"
	PullRequest        = "pullRequest"
	PullRequestUnknown = "pullRequestUnknown"
	Checks             = "checks"
	ChecksPassed       = "checksPassed"
	ChecksUnknown      = "checksUnknown"
	RuntimeStopped     = "runtimeStopped"
	ReadyForReview     = "readyForReview"
	ReadyToMerge       = "readyToMerge"
	NextReview         = "nextReview"
	NextMerge          = "nextMerge"
)

// PlanBlocked is the plan's own line when a check fails: the website has no
// such line, since its session never fails.
const PlanBlocked = "blocked"

// JudgeStopped is the plan's own line when the judge stops mid-plan: a gate
// paused its calls, a call failed, or an answer was refused.
const JudgeStopped = "judgeStopped"

// GateBlocked is the run's own line when FG-01, FG-02 or FG-05 blocks a lap
// before anything is paid; the website's session never meets one.
const GateBlocked = "gateBlocked"

// GateExhausted is the run's own line when FG-04 pauses the paid calls of
// an open lap: the reservation is reached, or a call reported no usage.
const GateExhausted = "gateExhausted"

// LapResumed is the run's own line when dft agent run finds a lap open and
// plays on from where it stood.
const LapResumed = "lapResumed"

// RuntimeUnknown is the run's own line when the runtime's cleanup was not
// confirmed: the outcome is unknown and the reservation is held.
const RuntimeUnknown = "runtimeUnknown"

// ChecksFailed is the run's own line when the pull request's checks ended
// with a failure.
const ChecksFailed = "checksFailed"
