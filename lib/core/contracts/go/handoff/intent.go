package handoff

// Intent is the recorded intent (intent.json) that dft agent run applies:
// the ticket, the runtime and its revision, the model and effort, the human
// checkpoints and whether the team authorized it, for "console › intent
// claude-sonnet-5-5 · effort medium · runtime.md revision 1".
type Intent struct {
	Ticket    string       `json:"ticket"`
	RuntimeID string       `json:"runtimeId"`
	Revision  int          `json:"revision"`
	Model     string       `json:"model"`
	Effort    string       `json:"effort"`
	Human     Human        `json:"human"`
	Status    IntentStatus `json:"status"`
}
