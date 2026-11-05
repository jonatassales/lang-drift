package configuration

import "strings"

// IsImageReference reports whether text can name one already-built image.
func IsImageReference(text string) bool {
	return IsPlainText(text) && !strings.ContainsAny(text, " \t\r\n")
}
