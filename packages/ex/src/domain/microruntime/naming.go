package microruntime

import (
	"regexp"
	"strings"
	"unicode"
)

const maximumIdentifierLength = 63

var dnsLabel = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]*[a-z0-9])?$`)

// Identifier names a runtime, unit, namespace, hook, or agent. It is a DNS
// label, so it is also safe as a Kubernetes label value. Only ParseIdentifier
// produces a usable value; the zero Identifier means "not set".
type Identifier struct{ value string }

// ParseIdentifier validates text as an Identifier.
func ParseIdentifier(text string) (Identifier, error) {
	if len(text) > maximumIdentifierLength || !dnsLabel.MatchString(text) {
		return Identifier{}, reject(
			ErrInvalidIdentifier,
			"%q must be 1-63 lowercase DNS-label characters",
			text,
		)
	}

	return Identifier{value: text}, nil
}

func (identifier Identifier) String() string { return identifier.value }

// IsZero reports whether the identifier was never set.
func (identifier Identifier) IsZero() bool { return identifier.value == "" }

// isPlainText reports whether text is non-blank and free of control characters.
func isPlainText(text string) bool {
	return strings.TrimSpace(text) != "" &&
		strings.IndexFunc(text, unicode.IsControl) == -1
}
