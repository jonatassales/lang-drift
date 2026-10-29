package support

import (
	"io"
	"os"
)

// Interactive reports whether the input is a terminal a person types on,
// rather than a pipe or a file.
func Interactive(input io.Reader) bool {
	described, describable := input.(interface{ Stat() (os.FileInfo, error) })
	if !describable {
		return false
	}

	info, err := described.Stat()
	if err != nil {
		return false
	}

	return info.Mode()&os.ModeCharDevice != 0
}
