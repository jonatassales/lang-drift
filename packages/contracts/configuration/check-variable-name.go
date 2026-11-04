package configuration

import (
	"errors"
	"fmt"
	"regexp"
)

// ErrInvalidVariable marks an environment variable a component may not set.
var ErrInvalidVariable = errors.New("invalid environment variable")

var variableName = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)

// CheckVariableName accepts an uppercase environment variable name, the
// only shape an env file or a pass-through list may name.
func CheckVariableName(text string) error {
	if len(text) > 128 || !variableName.MatchString(text) {
		return fmt.Errorf("%w: %q must be an uppercase variable name", ErrInvalidVariable, text)
	}

	return nil
}
