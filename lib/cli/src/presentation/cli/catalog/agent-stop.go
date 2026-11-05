package catalog

import (
	"context"

	"github.com/jonatassales/tao/lib/cli/src/application"
	"github.com/jonatassales/tao/lib/cli/src/application/execution"
	"github.com/jonatassales/tao/lib/cli/src/presentation/cli/support"

	"github.com/spf13/cobra"
)

// NewAgentStopCommand creates the stop command: a run that waits at a stage
// is let go, its runtime taken down and its reservation released, and the
// runtime can be planned again.
func NewAgentStopCommand(executor Executor, globals *support.Globals) *cobra.Command {
	var by string

	command := &cobra.Command{
		Use:   "stop RUNTIME",
		Short: "Abandon a run that waits at a stage: stop its runtime, release the reservation, keep the evidence",
		Args:  validateSingleArgument,
		RunE: func(command *cobra.Command, arguments []string) error {
			request := execution.StopRequest{
				Runtime: arguments[0],
				By:      by,
				Session: sessionOf(command, globals),
			}

			return runSession(command, globals, func(ctx context.Context, project application.ProjectSelection) (execution.StopResult, error) {
				request.Project = project

				return executor.Stop(ctx, request)
			})
		},
	}
	byFlag(command, &by)

	return command
}
