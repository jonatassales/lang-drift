package handoff

// StageState is one stage of the run, the dev ticket among them: its state
// (spec 0057, and spec 0061 for the dev ticket) and its report (spec 0059).
// Its other fields join with those specs.
type StageState struct {
	Report *StageReport `json:"report,omitempty"`
}
