package session

// Shape is one of the six layouts every printed line takes, as spec 0048's
// line format names them.
type Shape int

// The six shapes: a comment ("# " and its text); a narrative actor line
// ("console › " and its text); a labeled actor line ("console › ", the
// label padded to 10 columns, one space, the value); a result line ("✓ ",
// the padded label, one space, the value); a status line (the padded label,
// one space, the value); and a component row ("✓ ", the path padded to 19
// columns, three spaces, the stack and role padded to 15 columns, one
// space, the share).
const (
	Comment Shape = iota
	Narrative
	Labeled
	Result
	Status
	ComponentRow
)
