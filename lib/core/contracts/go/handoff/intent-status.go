package handoff

// IntentStatus says whether the team has authorized a recorded intent.
type IntentStatus string

// An intent is proposed by the plan and authorized by the team.
const (
	IntentProposed   IntentStatus = "proposed"
	IntentAuthorized IntentStatus = "authorized"
)
