package handoff

// Evidence is the Plan card's EVIDENCE: "runtime.md draft · revision 1 ·
// proposed to the console". A value not known yet is null.
type Evidence struct {
	File       *string `json:"file"`
	State      *string `json:"state"`
	Revision   *int    `json:"revision"`
	ProposedTo *string `json:"proposedTo"`
}
