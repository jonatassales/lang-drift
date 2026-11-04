package manifest

import "strings"

// cutWord splits text at its first space: the first word, and the rest
// with its surrounding space trimmed, empty when there is none.
func cutWord(text string) (string, string) {
	word, rest, found := strings.Cut(text, " ")
	if !found {
		return text, ""
	}

	return word, strings.TrimSpace(rest)
}
