package configuration

import (
	"errors"
	"fmt"
	"regexp"
)

// ErrInvalidIdentifier marks a name that is not a lowercase DNS label.
var ErrInvalidIdentifier = errors.New("invalid identifier")

const maximumIdentifierLength = 63

var dnsLabel = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]*[a-z0-9])?$`)

// CheckIdentifier accepts a runtime, unit, namespace, hook or agent name: a
// DNS label, so it is also safe as a Kubernetes label value.
func CheckIdentifier(text string) error {
	notDNSLabel := len(text) > maximumIdentifierLength ||
		!dnsLabel.MatchString(text)

	if notDNSLabel {
		return fmt.Errorf(
			"%w: %q must be 1-63 lowercase DNS-label characters",
			ErrInvalidIdentifier, text,
		)
	}

	return nil
}
