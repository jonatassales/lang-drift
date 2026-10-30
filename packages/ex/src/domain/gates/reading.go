package gates

// StatusUnknown marks a reading CI cannot make. Unknown is never zero.
const StatusUnknown = "unknown"

// Reading is one value the report shows, or the reason it cannot.
type Reading struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Reason string `json:"reason"`
}
