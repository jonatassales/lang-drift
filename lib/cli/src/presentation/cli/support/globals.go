package support

import (
	"context"
	"time"

	"github.com/jonatassales/tao/lib/cli/src/application/session"
)

// Globals holds the persistent flags every command shares. The root command
// binds them; commands read them only when they run, after Cobra has parsed
// the command line.
type Globals struct {
	// Project is --project: the project root, named explicitly.
	Project string

	// Fallback is where project discovery starts when --project is absent;
	// empty means the working directory.
	Fallback string

	// JSON is --json: compact results and structured diagnostics.
	JSON bool

	// Timeout is --timeout: the total deadline of one finite operation.
	Timeout time.Duration

	// Terminal says the output is a terminal, so sessions are toned and title
	// the window.
	Terminal bool

	// Observer, when set, sees every session event the output prints.
	Observer session.Sink
}

// Bounded applies --timeout to one finite operation.
func (globals *Globals) Bounded(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, globals.Timeout)
}
