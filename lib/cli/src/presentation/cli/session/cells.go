package session

// Cells are the words a line's shape lays out: Text for a comment or a
// narrative actor line, Label and Value for a labeled actor, result or
// status line, and Path, StackRole and Share for a component row.
type Cells struct {
	Text      string
	Label     string
	Value     string
	Path      string
	StackRole string
	Share     string
}
