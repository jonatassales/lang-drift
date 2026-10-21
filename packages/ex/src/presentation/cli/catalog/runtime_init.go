package catalog

import (
	"context"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/presentation/cli/support"

	"github.com/spf13/cobra"
)

// RuntimeInitializer publishes a new runtime configuration.
type RuntimeInitializer interface {
	Execute(context.Context, application.InitializeRuntimeRequest) (application.RuntimeInitResult, error)
}

// NewRuntimeInitCommand creates the runtime initialization command. With
// --guided, the use case asks for every choice once the project is known to
// exist; nothing is written until all are given, and --timeout does not bound
// the questions.
func NewRuntimeInitCommand(initializeRuntime RuntimeInitializer, globals *support.Globals) *cobra.Command {
	var guided bool

	command := &cobra.Command{
		Use:   "init RUNTIME",
		Short: "Generate editable runtime configuration without executing it",
		Args:  validateSingleArgument,
		RunE: func(command *cobra.Command, arguments []string) error {
			var guide application.BlueprintGuide

			runCommand := runFinite[application.RuntimeInitResult]

			if guided {
				guide = promptGuide{prompter: support.NewPrompter(command.InOrStdin(), command.ErrOrStderr())}
				runCommand = runInteractive[application.RuntimeInitResult]
			}

			return runCommand(command, globals, func(ctx context.Context, project application.ProjectSelection) (application.RuntimeInitResult, error) {
				return initializeRuntime.Execute(ctx, application.InitializeRuntimeRequest{
					Project: project,
					Runtime: arguments[0],
					Guide:   guide,
				})
			})
		},
	}

	command.Flags().BoolVar(&guided, "guided", false, "collect explicit component and resource choices")

	return command
}
