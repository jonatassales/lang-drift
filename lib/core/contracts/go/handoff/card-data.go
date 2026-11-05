package handoff

// CardData is the plan's card data (plan.json): one section per card, each
// absent until it is revealed.
type CardData struct {
	Ticket   *TicketCard   `json:"ticket,omitempty"`
	Analysis *AnalysisCard `json:"analysis,omitempty"`
	Plan     *PlanCard     `json:"plan,omitempty"`
	Workflow *WorkflowCard `json:"workflow,omitempty"`
}
