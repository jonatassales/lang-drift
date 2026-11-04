package configuration

import "strings"

// isFileName accepts one plain file name: no folder, no parent, no blank.
func isFileName(text string) bool {
	return IsPlainText(text) &&
		text != "." &&
		text != ".." &&
		!strings.ContainsAny(text, `/\`)
}
