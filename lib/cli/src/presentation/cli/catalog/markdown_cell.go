package catalog

import "strings"

// markdownCell keeps text inside one Markdown table cell: a pipe would end
// the cell and a line break the row.
func markdownCell(text string) string {
	escaped := strings.ReplaceAll(text, "|", `\|`)

	return strings.ReplaceAll(escaped, "\n", " ")
}
