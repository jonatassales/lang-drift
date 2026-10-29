package session

import "github.com/jonatassales/tao/packages/ex/src/application/session"

// Reveal names the card and the field a line fills, as the website's cards
// reveal them when the line prints.
type Reveal struct {
	Card  string
	Field string
}

// The cards.
const (
	TicketCard   = "ticket"
	AnalysisCard = "analysis"
	PlanCard     = "plan"
	WorkflowCard = "workflow"
)

var reveals = map[string]Reveal{
	key(session.PlanSession, session.PlanCommand):     {TicketCard, "cards"},
	key(session.PlanSession, session.Links):           {TicketCard, "links"},
	key(session.PlanSession, session.Analysis):        {AnalysisCard, "act"},
	key(session.PlanSession, session.ReadDocument):    {AnalysisCard, "document"},
	key(session.PlanSession, session.Reads):           {AnalysisCard, "documentsRead"},
	key(session.PlanSession, session.Guardrail):       {AnalysisCard, "guardrail"},
	key(session.PlanSession, session.JudgeGuardrail):  {AnalysisCard, "guardrailJudged"},
	key(session.PlanSession, session.Bounds):          {PlanCard, "bounded"},
	key(session.PlanSession, session.Memory):          {AnalysisCard, "suggestionSkeleton"},
	key(session.PlanSession, session.VectorSearch):    {AnalysisCard, "episodesSkeleton"},
	key(session.PlanSession, session.Candidates):      {AnalysisCard, "candidates"},
	key(session.PlanSession, session.Drop):            {AnalysisCard, "verdicts"},
	key(session.PlanSession, session.Episodes):        {AnalysisCard, "cited"},
	key(session.PlanSession, session.ReasonResources): {AnalysisCard, "measured"},
	key(session.PlanSession, session.Suggest):         {AnalysisCard, "suggestion"},
	key(session.PlanSession, session.Planning):        {PlanCard, "act"},
	key(session.PlanSession, session.Resolve):         {PlanCard, "componentsSkeleton"},
	key(session.PlanSession, session.Component):       {PlanCard, "component"},
	key(session.PlanSession, session.CrossCheck):      {PlanCard, "crossChecked"},
	key(session.PlanSession, session.Size):            {PlanCard, "suggested"},
	key(session.PlanSession, session.Limits):          {PlanCard, "limits"},
	key(session.PlanSession, session.Price):           {PlanCard, "priceBasis"},
	key(session.PlanSession, session.Validate):        {PlanCard, "validated"},
	key(session.PlanSession, session.Draft):           {PlanCard, "drafted"},
	key(session.PlanSession, session.Handoff):         {WorkflowCard, "finished"},
	key(session.RunSession, session.RunCommand):       {WorkflowCard, "result"},
	key(session.RunSession, session.GateOpened):       {WorkflowCard, "reserved"},
	key(session.RunSession, session.GatePaused):       {WorkflowCard, "paused"},
	key(session.RunSession, session.Waiting):          {WorkflowCard, "stage"},
	key(session.RunSession, session.Confirmed):        {WorkflowCard, "stage"},
	key(session.RunSession, session.Stage):            {WorkflowCard, "stage"},
	key(session.RunSession, session.Completed):        {WorkflowCard, "observed"},
	key(session.RunSession, session.PullRequest):      {WorkflowCard, "pullRequest"},
	key(session.RunSession, session.Checks):           {WorkflowCard, "checks"},
	key(session.RunSession, session.ChecksPassed):     {WorkflowCard, "checks"},
	key(session.RunSession, session.ReadyForReview):   {WorkflowCard, "released"},
	key(session.RunSession, session.ReadyToMerge):     {WorkflowCard, "released"},
}

// RevealOf is the card and field an event fills, if any.
func RevealOf(event session.Event) (Reveal, bool) {
	reveal, known := reveals[key(event.Session, event.ID)]

	return reveal, known
}
