package stages

// checkTickets accepts a tickets report with a count of tickets.
func checkTickets(report TicketsReport) error {
	return checkCount("tickets", report.Tickets)
}
