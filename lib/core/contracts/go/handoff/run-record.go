package handoff

// RunRecord is the run's record (run.json), keyed as spec 0065 keys it: the
// run's id and outcome, null until known, and one section per part of the
// run, absent until the run reaches it. The Workflow card's PULL REQUEST and
// SAFE SEAT read from its sections; its RESERVED, OBSERVED and RELEASED are
// not copied here, and read from the run's entry in the money ledger,
// finops/ledger/<run>.ledger.json (spec 0056), which the run's id names.
type RunRecord struct {
	Run      *string        `json:"run"`
	Intent   *AppliedChoice `json:"intent,omitempty"`
	Stages   []StageState   `json:"stages,omitempty"`
	Effort   *EffortRecord  `json:"effort,omitempty"`
	SafeSeat *SafeSeat      `json:"safeSeat,omitempty"`
	Runtime  *RuntimeState  `json:"runtime,omitempty"`
	Delivery *Delivery      `json:"delivery,omitempty"`
	Outcome  *RunOutcome    `json:"outcome"`
}
