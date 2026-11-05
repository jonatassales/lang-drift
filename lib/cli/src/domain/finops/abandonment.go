package finops

// Abandonment is what an abandoned lap settles to: the part of the
// reservation it releases, the spend it keeps, the overspend it leaves
// visible, and its closing line. Settlement already names what Settle
// returns, so this is its own type.
type Abandonment struct {
	Released  int
	Observed  int
	Overspend int
	Closing   LapClosing
}
