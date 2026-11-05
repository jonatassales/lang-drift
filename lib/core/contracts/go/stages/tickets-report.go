package stages

// TicketsReport is the to-tickets skill's report: "3 tickets · blocking
// edges recorded".
type TicketsReport struct {
	Tickets               int  `json:"tickets"`
	BlockingEdgesRecorded bool `json:"blockingEdgesRecorded"`
}
