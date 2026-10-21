// Package catalog holds EX's commands, one per file. Each command declares the
// narrow use-case interface it consumes, maps its arguments and flags to an
// application request, and renders the result. It never orchestrates work or
// reaches infrastructure itself.
package catalog

import (
	"context"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/presentation/cli/support"

	"github.com/spf13/cobra"
)

// projectOperation runs one use case for the selected project.
type projectOperation[Result any] func(context.Context, application.ProjectSelection) (Result, error)

// runFinite is the shared shape of a finite command: it maps the project
// flags to a selection, bounds the operation with --timeout, and renders the
// result.
func runFinite[Result any](command *cobra.Command, globals *support.Globals, operation projectOperation[Result]) error {
	ctx, cancel := globals.Bounded(command.Context())
	defer cancel()

	return run(ctx, command, globals, operation)
}

// runInteractive is runFinite for an operation that waits for the user, so
// --timeout does not bound it. An interrupt still cancels it.
func runInteractive[Result any](command *cobra.Command, globals *support.Globals, operation projectOperation[Result]) error {
	return run(command.Context(), command, globals, operation)
}

func run[Result any](
	ctx context.Context,
	command *cobra.Command,
	globals *support.Globals,
	operation projectOperation[Result],
) error {
	selection, err := globals.Selection()
	if err != nil {
		return err
	}

	result, err := operation(ctx, selection)
	if err != nil {
		return err
	}

	return support.Render(command.OutOrStdout(), result, globals.JSON)
}
