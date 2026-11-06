package cli_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/jonatassales/tao/lib/cli/src/application/execution"
	"github.com/jonatassales/tao/lib/cli/src/presentation/cli"
)

// A command that ends on its last line, and the over-budget pause, exit 0
// and write nothing; every other ending is a stop: one plain sentence on
// standard error, with no prefix and no hint, and exit 1, an interrupt
// included.
func TestStopTurnsAnEndingIntoItsStatusAndSentence(t *testing.T) {
	t.Parallel()

	for _, scenario := range []struct {
		name     string
		err      error
		status   int
		sentence string
	}{
		{"a finished command", nil, 0, ""},
		{"the over-budget pause", fmt.Errorf("%w: FG-03 · a $40.80 reservation exceeds the $25.00 budget; nothing was paid", execution.ErrPausedAtGate), 0, ""},
		{"failed checks", fmt.Errorf("%w: build failed", execution.ErrChecksFailed), 1, "The pull request's checks failed: build failed.\n"},
		{"a cause that ends its sentence", errors.New("no runtime is planned."), 1, "No runtime is planned.\n"},
		{"two joined causes", errors.Join(errors.New("the run failed"), errors.New("the record could not be written")), 1, "The run failed; the record could not be written.\n"},
		{"an interrupt", fmt.Errorf("plan: %w", context.Canceled), 1, "The command was interrupted.\n"},
		{"an interrupted pause", errors.Join(context.Canceled, execution.ErrPausedAtGate), 1, "The command was interrupted.\n"},
		{"a cause naming the command", errors.New("dft cluster requires a subcommand"), 1, "dft cluster requires a subcommand.\n"},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			t.Parallel()

			var stderr bytes.Buffer

			status := cli.Stop(scenario.err, &stderr)
			if status != scenario.status || stderr.String() != scenario.sentence {
				t.Errorf("Stop = %d, %q; want %d, %q", status, stderr.String(), scenario.status, scenario.sentence)
			}
		})
	}
}
