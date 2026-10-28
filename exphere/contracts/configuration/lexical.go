package configuration

import (
	"errors"
	"fmt"
	"path"
	"regexp"
	"strings"
	"unicode"
)

// The lexical rules every Microruntimes value obeys, wherever it is checked:
// in a runtime.md fence, in the CLI's domain, or in the server that authors
// the fence. They need no filesystem, cluster or network.

var (
	// ErrInvalidIdentifier marks a name that is not a lowercase DNS label.
	ErrInvalidIdentifier = errors.New("invalid identifier")

	// ErrInvalidPath marks a path that is not a clean relative path.
	ErrInvalidPath = errors.New("invalid relative path")
)

const maximumIdentifierLength = 63

var dnsLabel = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]*[a-z0-9])?$`)

// CheckIdentifier accepts a runtime, unit, namespace, hook or agent name: a
// DNS label, so it is also safe as a Kubernetes label value.
func CheckIdentifier(text string) error {
	if len(text) > maximumIdentifierLength || !dnsLabel.MatchString(text) {
		return fmt.Errorf(
			"%w: %q must be 1-63 lowercase DNS-label characters",
			ErrInvalidIdentifier, text,
		)
	}

	return nil
}

// CheckRelativePath accepts a clean, slash-separated path that cannot leave
// its boundary. Whether it exists or crosses a symlink is not a lexical rule.
func CheckRelativePath(text string) error {
	if !IsPlainText(text) ||
		text == "." ||
		strings.Contains(text, `\`) ||
		path.IsAbs(text) ||
		path.Clean(text) != text {
		return fmt.Errorf("%w: %q must be a clean relative path", ErrInvalidPath, text)
	}

	for _, segment := range strings.Split(text, "/") {
		if segment == ".." {
			return fmt.Errorf("%w: %q escapes its boundary", ErrInvalidPath, text)
		}
	}

	return nil
}

// CheckComponentPath accepts the repository root, ".", as a component source,
// and otherwise the same paths as CheckRelativePath.
func CheckComponentPath(text string) error {
	if text == "." {
		return nil
	}

	return CheckRelativePath(text)
}

// IsPlainText reports whether text is non-blank and free of control characters.
func IsPlainText(text string) bool {
	return strings.TrimSpace(text) != "" &&
		strings.IndexFunc(text, unicode.IsControl) == -1
}

// ValidPort reports whether a TCP port is usable.
func ValidPort(port int) bool {
	return port >= 1 && port <= 65535
}

// IsPlatformPath reports whether a path lies inside the platform's own
// .langdrift directory, which component sources never may.
func IsPlatformPath(text string) bool {
	return text == ".langdrift" || strings.HasPrefix(text, ".langdrift/")
}

// IsImageReference reports whether text can name one already-built image.
func IsImageReference(text string) bool {
	return IsPlainText(text) && !strings.ContainsAny(text, " \t\r\n")
}
