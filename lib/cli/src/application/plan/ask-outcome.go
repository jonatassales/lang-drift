package plan

// AskOutcome is how one question of the planner went: its kind, where it
// stands, and what the gate, the console or the model said about it, in
// their own words, for a refusal that names the cause.
type AskOutcome struct {
	Kind   string
	Status AskStatus
	Reason string
}
