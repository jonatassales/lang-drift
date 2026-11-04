package configuration

import (
	"errors"
	"fmt"
)

// ErrInvalidConfiguration classifies every authored-configuration rejection,
// so callers match the class with errors.Is and read the reason in the text.
var ErrInvalidConfiguration = errors.New("invalid runtime configuration")

// invalid rejects an authored configuration with its reason.
func invalid(reason string) error {
	return fmt.Errorf("%w: %s", ErrInvalidConfiguration, reason)
}
