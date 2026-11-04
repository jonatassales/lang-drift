package support

import (
	"io"
	"os"
)

// Interactive reports whether the input is a terminal a person types on,
// rather than a pipe, a file or the null device, which is a character device
// nobody types on.
func Interactive(input io.Reader) bool {
	described, describable := input.(interface{ Stat() (os.FileInfo, error) })
	if !describable {
		return false
	}

	info, err := described.Stat()
	if err != nil {
		return false
	}

	if info.Mode()&os.ModeCharDevice == 0 {
		return false
	}

	// Without a null device to compare with, the character device stands.
	null, err := os.Stat(os.DevNull)
	if err != nil {
		return true
	}

	return !os.SameFile(info, null)
}
