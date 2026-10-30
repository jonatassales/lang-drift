package support

import (
	"os"
	"strings"
	"testing"
)

// Interactive is true only for a terminal: a pipe, a string and the null
// device, a character device nobody types on, are not one.
func TestInteractive(test *testing.T) {
	test.Parallel()

	null, err := os.Open(os.DevNull)
	if err != nil {
		test.Fatal(err)
	}
	defer null.Close()

	if Interactive(null) {
		test.Error("the null device reads as a terminal")
	}

	reader, writer, err := os.Pipe()
	if err != nil {
		test.Fatal(err)
	}
	defer reader.Close()
	defer writer.Close()

	if Interactive(reader) {
		test.Error("a pipe reads as a terminal")
	}

	if Interactive(strings.NewReader("answers\n")) {
		test.Error("a string reads as a terminal")
	}
}
