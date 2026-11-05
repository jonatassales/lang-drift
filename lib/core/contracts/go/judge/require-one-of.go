package judge

import (
	"slices"
	"strings"
)

// requireOneOf refuses a value outside the set the schema allows.
func requireOneOf(path string, value string, allowed []string) error {
	if slices.Contains(allowed, value) {
		return nil
	}

	return refuse(path, "must be one of "+strings.Join(allowed, ", "))
}
