package catalog

import (
	"context"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/presentation/cli/support"

	"github.com/spf13/cobra"
)

// RuntimeHookRunner executes one explicitly declared consumer extension.
type RuntimeHookRunner interface {
	Execute(context.Context, application.RunRuntimeHookRequest) (application.HookResult, error)
}

// NewRuntimeHookCommand runs only the named hook and renders its native result.
func NewRuntimeHookCommand(runHook RuntimeHookRunner, globals *support.Globals) *cobra.Command {
	return &cobra.Command{
		Use: "hook <id> <name>", Short: "Explicitly run a declared trusted .sh extension",
		Args: support.ExactArguments(2),
		RunE: func(command *cobra.Command, arguments []string) error {
			return runFinite(command, globals, func(ctx context.Context, project application.ProjectSelection) (application.HookResult, error) {
				return runHook.Execute(ctx, application.RunRuntimeHookRequest{
					RuntimeRequest: application.RuntimeRequest{Project: project, Runtime: arguments[0]},
					Name:           arguments[1],
				})
			})
		},
	}
}
