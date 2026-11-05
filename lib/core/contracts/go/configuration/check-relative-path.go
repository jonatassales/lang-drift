package configuration

import (
	"errors"
	"fmt"
	"path"
	"strings"
)

// ErrInvalidPath marks a path that is not a clean relative path.
var ErrInvalidPath = errors.New("invalid relative path")

// CheckRelativePath accepts a clean, slash-separated path that cannot leave
// its boundary. Whether it exists or crosses a symlink is not a lexical rule.
func CheckRelativePath(text string) error {
	notCleanRelative := !IsPlainText(text) ||
		text == "." ||
		strings.Contains(text, `\`) ||
		path.IsAbs(text) ||
		path.Clean(text) != text
	if notCleanRelative {
		return fmt.Errorf("%w: %q must be a clean relative path", ErrInvalidPath, text)
	}

	for _, segment := range strings.Split(text, "/") {
		if segment == ".." {
			return fmt.Errorf("%w: %q escapes its boundary", ErrInvalidPath, text)
		}
	}

	return nil
}
