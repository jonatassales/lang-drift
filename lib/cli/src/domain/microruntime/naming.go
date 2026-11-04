package microruntime

import "github.com/jonatassales/tao/lib/core/contracts/go/configuration"

// Identifier names a runtime, unit, namespace, hook, or agent. It is a DNS
// label, so it is also safe as a Kubernetes label value. Only ParseIdentifier
// produces a usable value; the zero Identifier means "not set".
type Identifier struct{ value string }

// ParseIdentifier validates text as an Identifier, by the Microruntimes
// contracts' lexical rule.
func ParseIdentifier(text string) (Identifier, error) {
	if err := configuration.CheckIdentifier(text); err != nil {
		return Identifier{}, err
	}

	return Identifier{value: text}, nil
}

func (identifier Identifier) String() string { return identifier.value }

// IsZero reports whether the identifier was never set.
func (identifier Identifier) IsZero() bool { return identifier.value == "" }

// isPlainText reports whether text is non-blank and free of control characters.
func isPlainText(text string) bool { return configuration.IsPlainText(text) }

// ValidPort reports whether a TCP port is usable. Every port a consumer or
// the console declares goes through this one definition.
func ValidPort(port int) bool { return configuration.ValidPort(port) }

// IsPlatformPath reports whether a path lies inside the platform's own
// .langdrift directory, which consumer sources never may.
func IsPlatformPath(path string) bool { return configuration.IsPlatformPath(path) }
