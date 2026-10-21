package catalog

import (
	"context"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/presentation/cli/support"

	"github.com/spf13/cobra"
)

// ProjectInitializer records a project's identity.
type ProjectInitializer interface {
	Execute(context.Context, application.InitializeProjectRequest) (application.ProjectResult, error)
}

// NewProjectInitCommand creates the project initialization command.
func NewProjectInitCommand(initializeProject ProjectInitializer, globals *support.Globals) *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize project metadata without starting infrastructure",
		Args:  validateNoArguments,
		RunE: func(command *cobra.Command, _ []string) error {
			return runFinite(command, globals, func(ctx context.Context, project application.ProjectSelection) (application.ProjectResult, error) {
				return initializeProject.Execute(ctx, application.InitializeProjectRequest{Directory: project.Directory})
			})
		},
	}
}
