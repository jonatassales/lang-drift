package judgment

// CriteriaAnswer is the judge's answer to the criteria question: the
// guardrail criteria that apply to the ticket.
type CriteriaAnswer struct {
	Criteria []Criterion `json:"criteria"`
}
