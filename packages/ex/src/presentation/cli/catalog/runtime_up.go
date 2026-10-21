package catalog

import (
	"context"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/presentation/cli/support"

	"github.com/spf13/cobra"
)

// RuntimeStarter starts a runtime.
type RuntimeStarter interface {
	Execute(context.Context, application.RuntimeRequest) (application.ExecutionOutcome, error)
}

// NewRuntimeUpCommand creates the runtime startup command.
func NewRuntimeUpCommand(startRuntime RuntimeStarter, globals *support.Globals) *cobra.Command {
	return &cobra.Command{
		Use:   "up RUNTIME",
		Short: "Prepare the runtime's namespace and run every configured unit",
		Args:  validateSingleArgument,
		RunE: func(command *cobra.Command, arguments []string) error {
			return runFinite(command, globals, func(ctx context.Context, project application.ProjectSelection) (application.ExecutionOutcome, error) {
				return startRuntime.Execute(ctx, application.RuntimeRequest{Project: project, Runtime: arguments[0]})
			})
		},
	}
}
