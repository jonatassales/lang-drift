package support

import (
	"context"
	"time"
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
}

// Bounded applies --timeout to one finite operation.
func (globals *Globals) Bounded(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, globals.Timeout)
}
