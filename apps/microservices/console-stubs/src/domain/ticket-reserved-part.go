package domain

// TicketReservedPart is the part of a ticket a person writes by hand under
// Safe Seat: the paths it covers, its goal, its category and who chose it,
// the judge or the lead. Illustrative stub data.
type TicketReservedPart struct {
	Paths    []string `json:"paths"`
	Goal     string   `json:"goal"`
	Category string   `json:"category"`
	ChosenBy string   `json:"chosenBy"`
}
