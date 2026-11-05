package finops

// LapClosing is how a lap closed, on its ledger line: the outcome, the
// person who closed it and the time, in RFC 3339.
type LapClosing struct {
	Outcome string `json:"outcome"`
	By      string `json:"by"`
	At      string `json:"at"`
}
