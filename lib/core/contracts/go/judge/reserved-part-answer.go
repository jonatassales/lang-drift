package judge

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
