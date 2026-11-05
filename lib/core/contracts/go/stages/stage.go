package stages

// Stage names the stage a report belongs to.
type Stage string

// The five stages whose skill writes a report.
const (
	Planning  Stage = "planning"
	Spec      Stage = "spec"
	Tickets   Stage = "tickets"
	Implement Stage = "implement"
	Review    Stage = "review"
)
