package configuration

import (
	"strings"
	"unicode"
)

// IsPlainText reports whether text is non-blank and free of control characters.
func IsPlainText(text string) bool {
	return strings.TrimSpace(text) != "" &&
		strings.IndexFunc(text, unicode.IsControl) == -1
}
