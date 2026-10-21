package catalog

import (
	"context"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/presentation/cli/support"

	"github.com/spf13/cobra"
)

// RuntimeStopper stops a runtime.
type RuntimeStopper interface {
	Execute(context.Context, application.StopRuntimeRequest) (application.ExecutionOutcome, error)
}

// NewRuntimeDownCommand creates the runtime shutdown command.
func NewRuntimeDownCommand(stopRuntime RuntimeStopper, globals *support.Globals) *cobra.Command {
	var recoverPartial bool

	command := &cobra.Command{
		Use:   "down RUNTIME",
		Short: "Remove owned resources; keep configuration, evidence, and the cluster",
		Args:  validateSingleArgument,
		RunE: func(command *cobra.Command, arguments []string) error {
			return runFinite(command, globals, func(ctx context.Context, project application.ProjectSelection) (application.ExecutionOutcome, error) {
				return stopRuntime.Execute(ctx, application.StopRuntimeRequest{
					RuntimeRequest: application.RuntimeRequest{Project: project, Runtime: arguments[0]},
					Recover:        recoverPartial,
				})
			})
		},
	}

	command.Flags().BoolVar(&recoverPartial, "recover", false, "clean up after an inspected failed or interrupted operation")

	return command
}
