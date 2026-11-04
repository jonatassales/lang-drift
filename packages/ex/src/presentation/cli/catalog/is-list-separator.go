package catalog

// isListSeparator splits a list inside one --component value, whose pairs
// are already separated by commas.
func isListSeparator(character rune) bool {
	return character == ';'
}
