package domain

// TicketSafeSeat is a ticket's Safe Seat record: its state, the override
// reason when a lead overrode it, and the mark that it is illustrative stub
// data. It never holds the person's share.
type TicketSafeSeat struct {
	State          string  `json:"state"`
	OverrideReason *string `json:"overrideReason"`
	Illustrative   bool    `json:"illustrative"`
}
