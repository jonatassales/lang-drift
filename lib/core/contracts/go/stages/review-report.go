package stages

// ReviewReport is the code-review skill's report: "code-review 2 findings · fixed".
type ReviewReport struct {
	Findings int  `json:"findings"`
	Fixed    bool `json:"fixed"`
}
