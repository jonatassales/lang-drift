package setup

// Pending says which inputs are not recorded yet, so only those are asked.
type Pending struct {
	// Period is the month a budget given now covers, YYYY-MM.
	Period string
	Budget bool
	Gates  bool
}
