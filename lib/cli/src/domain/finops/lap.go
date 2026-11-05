package finops

// Lap is a lap's money as it stands when the lap closes: what it reserved,
// what its calls were observed to spend, and whether a call was cut off, so
// that what it spent is unknown.
type Lap struct {
	Reserved   int
	Observed   int
	Unmeasured bool
}
