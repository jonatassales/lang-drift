package cli

import (
	"context"
	"errors"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/jonatassales/tao/lib/cli/src/application/execution"
)

// programName opens the causes that name a command; a sentence that starts
// with it keeps it as typed.
const programName = "dft"

// Stop ends a command by spec 0048's stop rule. A command that ended on its
// last line, and the over-budget pause, which the simulation shows, exit 0
// and write nothing. Every other ending is a stop: it exits 1 and writes
// one plain sentence naming the cause, with no prefix and no hint. An
// interrupt is a stop.
func Stop(err error, stderr io.Writer) int {
	if err == nil {
		return 0
	}

	interrupted := errors.Is(err, context.Canceled)
	shownPause := !interrupted && errors.Is(err, execution.ErrPausedAtGate)
	if shownPause {
		return 0
	}

	cause := strings.Join(strings.Split(err.Error(), "\n"), "; ")
	if interrupted {
		cause = "the command was interrupted"
	}

	sentence := cause
	if !strings.HasPrefix(cause, programName+" ") {
		first, size := utf8.DecodeRuneInString(cause)
		sentence = string(unicode.ToUpper(first)) + cause[size:]
	}

	if !strings.HasSuffix(sentence, ".") {
		sentence += "."
	}

	if _, failed := io.WriteString(stderr, sentence+"\n"); failed != nil {
		// Standard error itself failed: the status is all that is left to
		// say the command did not finish.
		return 1
	}

	return 1
}
