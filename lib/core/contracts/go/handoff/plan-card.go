package handoff

// PlanCard is the Plan card's data: the components, the sizing, the price,
// the authorization and the evidence.
type PlanCard struct {
	Components    *Components    `json:"components,omitempty"`
	Sizing        *Sizing        `json:"sizing,omitempty"`
	Price         *Price         `json:"price,omitempty"`
	Authorization *Authorization `json:"authorization,omitempty"`
	Evidence      *Evidence      `json:"evidence,omitempty"`
}
