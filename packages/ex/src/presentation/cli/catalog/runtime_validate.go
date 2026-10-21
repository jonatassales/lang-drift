package catalog

import (
	"context"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/presentation/cli/support"

	"github.com/spf13/cobra"
)

// RuntimeValidator validates a runtime's configuration.
type RuntimeValidator interface {
	Execute(context.Context, application.RuntimeRequest) (application.ValidationResult, error)
}

// NewRuntimeValidateCommand creates the runtime validation command.
func NewRuntimeValidateCommand(validateRuntime RuntimeValidator, globals *support.Globals) *cobra.Command {
	return &cobra.Command{
		Use:   "validate RUNTIME",
		Short: "Validate the descriptor and native files without provisioning",
		Args:  validateSingleArgument,
		RunE: func(command *cobra.Command, arguments []string) error {
			return runFinite(command, globals, func(ctx context.Context, project application.ProjectSelection) (application.ValidationResult, error) {
				return validateRuntime.Execute(ctx, application.RuntimeRequest{Project: project, Runtime: arguments[0]})
			})
		},
	}
}
