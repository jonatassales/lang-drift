package catalog

import (
	"context"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/presentation/cli/support"

	"github.com/spf13/cobra"
)

// RuntimeInspector inspects a runtime.
type RuntimeInspector interface {
	Execute(context.Context, application.RuntimeRequest) (application.Inspection, error)
}

// NewRuntimeInspectCommand creates the runtime inspection command.
func NewRuntimeInspectCommand(inspectRuntime RuntimeInspector, globals *support.Globals) *cobra.Command {
	return &cobra.Command{
		Use:   "inspect RUNTIME",
		Short: "Inspect ownership, workload status, and resource observations",
		Args:  validateSingleArgument,
		RunE: func(command *cobra.Command, arguments []string) error {
			return runFinite(command, globals, func(ctx context.Context, project application.ProjectSelection) (application.Inspection, error) {
				return inspectRuntime.Execute(ctx, application.RuntimeRequest{Project: project, Runtime: arguments[0]})
			})
		},
	}
}
