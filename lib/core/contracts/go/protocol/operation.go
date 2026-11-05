package protocol

// Operation identifies who requested the current intent, and why.
type Operation struct {
	ID          string `json:"id"`
	Actor       string `json:"actor"`
	Cause       string `json:"cause"`
	RequestedAt string `json:"requestedAt"`
}
