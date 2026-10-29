package catalog

import (
	"context"

	"github.com/jonatassales/tao/packages/ex/src/application"
	"github.com/jonatassales/tao/packages/ex/src/application/execution"
	"github.com/jonatassales/tao/packages/ex/src/presentation/cli/session"
	"github.com/jonatassales/tao/packages/ex/src/presentation/cli/support"

	"github.com/spf13/cobra"
)

// Executor runs a planned runtime and reports where it stands.
type Executor interface {
	Start(context.Context, execution.StartRequest) (execution.StartResult, error)
	Status(context.Context, execution.StatusRequest) (execution.StatusResult, error)
}

// NewAgentRunCommand creates the run command: the second act of the session.
func NewAgentRunCommand(executor Executor, globals *support.Globals) *cobra.Command {
	var request execution.StartRequest

	command := &cobra.Command{
		Use:   "run RUNTIME",
		Short: "Run a planned runtime: estimate, reserve against the budget, pause at a gate or open a lap",
		Args:  validateSingleArgument,
		RunE: func(command *cobra.Command, arguments []string) error {
			request.Runtime = arguments[0]
			request.Session = session.For(command.OutOrStdout(), globals.JSON)

			return runSession(command, globals, func(ctx context.Context, project application.ProjectSelection) (execution.StartResult, error) {
				request.Project = project

				return executor.Start(ctx, request)
			})
		},
	}

	flags := command.Flags()
	flags.StringVar(&request.Model, "model", "", "the model to run with, over the suggestion")
	flags.StringVar(&request.Effort, "effort", "", "the effort to run with, over the suggestion")
	flags.BoolVar(&request.AcknowledgeOverspend, "acknowledge-overspend", false, "acknowledge an earlier lap's overspend, as the owner")
	byFlag(command, &request.By)

	return command
}

// NewAgentStatusCommand creates the status command.
func NewAgentStatusCommand(executor Executor, globals *support.Globals) *cobra.Command {
	return &cobra.Command{
		Use:   "status RUNTIME",
		Short: "Show where a runtime's run stands",
		Args:  validateSingleArgument,
		RunE: func(command *cobra.Command, arguments []string) error {
			return runFinite(command, globals, func(ctx context.Context, project application.ProjectSelection) (execution.StatusResult, error) {
				return executor.Status(ctx, execution.StatusRequest{Project: project, Runtime: arguments[0]})
			})
		},
	}
}
