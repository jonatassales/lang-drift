package protocol

// UnresolvedExecution identifies work requiring reconciliation. It never
// authorizes a new start or replaces the console claim check.
type UnresolvedExecution struct {
	ExecutionID   string `json:"executionId"`
	IntentVersion int    `json:"intentVersion"`
	Revision      int    `json:"revision"`
	State         string `json:"state"`
}
