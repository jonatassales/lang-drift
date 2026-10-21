package microruntime

import (
	"path"
	"strings"
)

// RelativePath is a clean, slash-separated path that cannot leave its boundary.
//
// Only the lexical rules live here. Whether the path exists, or crosses a
// symlink, depends on the filesystem and is checked by infrastructure.
// The zero RelativePath means "not set".
type RelativePath struct{ value string }

// ParseRelativePath validates text as a RelativePath.
func ParseRelativePath(text string) (RelativePath, error) {
	if !isPlainText(text) ||
		text == "." ||
		strings.Contains(text, `\`) ||
		path.IsAbs(text) ||
		path.Clean(text) != text {
		return RelativePath{}, reject(ErrInvalidPath, "%q must be a clean relative path", text)
	}

	for _, segment := range strings.Split(text, "/") {
		if segment == ".." {
			return RelativePath{}, reject(ErrInvalidPath, "%q escapes its boundary", text)
		}
	}

	return RelativePath{value: text}, nil
}

// mustRelativePath builds paths this package composes from already valid parts.
func mustRelativePath(text string) RelativePath {
	relative, err := ParseRelativePath(text)
	if err != nil {
		panic(err)
	}

	return relative
}

func (relative RelativePath) String() string { return relative.value }

// IsZero reports whether the path was never set.
func (relative RelativePath) IsZero() bool { return relative.value == "" }

// Base returns the last path segment.
func (relative RelativePath) Base() string { return path.Base(relative.value) }

// HasSuffix reports whether the path ends with suffix.
func (relative RelativePath) HasSuffix(suffix string) bool {
	return strings.HasSuffix(relative.value, suffix)
}

// IsWithin reports whether the path is directory itself or lies beneath it.
func (relative RelativePath) IsWithin(directory RelativePath) bool {
	return relative == directory ||
		strings.HasPrefix(relative.value, directory.value+"/")
}

// ParseComponentPath permits the repository root as source without relaxing
// native configuration and hook paths, which must still name a child.
func ParseComponentPath(text string) (RelativePath, error) {
	if text == "." {
		return RelativePath{value: "."}, nil
	}
	return ParseRelativePath(text)
}
