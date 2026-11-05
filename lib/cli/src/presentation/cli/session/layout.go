package session

import (
	"strings"
	"unicode/utf8"
)

// The columns of the line format: a label is padded to 10, a component's
// path to 19 and its stack and role to 15.
const (
	labelColumns     = 10
	pathColumns      = 19
	stackRoleColumns = 15
)

// Layout lays out one line in its shape, counting characters rather than
// bytes, so "›", "·", "✓" and "—" count one column each. A label shorter
// than its columns is padded with spaces; a label as long or longer is
// followed by exactly one space. No line ends with a space.
func Layout(shape Shape, cells Cells) string {
	padded := func(text string, columns int) string {
		short := columns - utf8.RuneCountInString(text)
		if short <= 0 {
			return text
		}

		return text + strings.Repeat(" ", short)
	}

	labeled := func() string {
		return padded(cells.Label, labelColumns) + " " + cells.Value
	}

	var line string

	switch shape {
	case Comment:
		line = "# " + cells.Text
	case Narrative:
		line = Actor + cells.Text
	case Labeled:
		line = Actor + labeled()
	case Result:
		line = "✓ " + labeled()
	case Status:
		line = labeled()
	case ComponentRow:
		line = "✓ " + padded(cells.Path, pathColumns) + "   " + padded(cells.StackRole, stackRoleColumns) + " " + cells.Share
	default:
		line = cells.Text
	}

	return strings.TrimRight(line, " ")
}
