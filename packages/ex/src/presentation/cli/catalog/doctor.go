package catalog

import (
	"context"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/presentation/cli/support"

	"github.com/spf13/cobra"
)

// Diagnoser checks the development environment.
type Diagnoser interface {
	Execute(context.Context, application.DiagnoseRequest) (application.Diagnosis, error)
}

// NewDoctorCommand creates the environment diagnostics command.
func NewDoctorCommand(runDiagnostics Diagnoser, globals *support.Globals) *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Check local prerequisites without installing anything",
		Args:  validateNoArguments,
		RunE: func(command *cobra.Command, _ []string) error {
			return runFinite(command, globals, func(ctx context.Context, project application.ProjectSelection) (application.Diagnosis, error) {
				return runDiagnostics.Execute(ctx, application.DiagnoseRequest{Project: project})
			})
		},
	}
}
