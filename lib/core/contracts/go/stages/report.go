package stages

// Report is one checked report: its stage, and the typed report of that
// stage.
type Report struct {
	Stage     Stage
	Planning  *PlanningReport
	Spec      *SpecReport
	Tickets   *TicketsReport
	Implement *ImplementReport
	Review    *ReviewReport
}
