package protocol

// Claim is the server's grant of one start for one intent version and revision.
type Claim struct {
	ExecutionID   string `json:"executionId"`
	IntentVersion int    `json:"intentVersion"`
	Revision      int    `json:"revision"`
}
