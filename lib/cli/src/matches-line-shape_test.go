package main

import (
	"regexp"
	"strings"
)

var (
	componentRowPattern = regexp.MustCompile(`^✓ \S+\s{3,}.+ · \S+\s+\d+m · \d+ MiB$`)
	resultLinePattern   = regexp.MustCompile(`^✓ (?:.{10}|\S{10,}) \S`)
	statusLinePattern   = regexp.MustCompile(`^(?:[a-z]{1,10}\s+|[a-z]{10,} )\S`)
)

// matchesLineShape reports whether a printed line follows one of spec 0048's
// six column shapes, with no trailing space.
func matchesLineShape(text string) bool {
	if text == "" || strings.HasSuffix(text, " ") {
		return false
	}

	if strings.HasPrefix(text, "# ") {
		return len(text) > len("# ")
	}

	if strings.HasPrefix(text, "console › ") {
		return len(text) > len("console › ")
	}

	if componentRowPattern.MatchString(text) {
		return true
	}

	if resultLinePattern.MatchString(text) {
		return true
	}

	return statusLinePattern.MatchString(text)
}
