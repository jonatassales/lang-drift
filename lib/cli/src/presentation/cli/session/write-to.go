package session

import "io"

// writeTo prints text on a session's output. A sink has no error to
// return, so a failed write is the output's to keep: the CLI prints through
// a checked writer, which keeps the first failure and reports it when the
// command ends.
func writeTo(output io.Writer, text string) {
	io.WriteString(output, text)
}
