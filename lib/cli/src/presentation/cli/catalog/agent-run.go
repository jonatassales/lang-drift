package catalog

import (
	"context"
	"time"

	"github.com/jonatassales/tao/lib/cli/src/application"
	"github.com/jonatassales/tao/lib/cli/src/application/execution"
	"github.com/jonatassales/tao/lib/cli/src/presentation/cli/support"

	"github.com/spf13/cobra"
)

// Executor runs a planned runtime, confirms its stages and reports where it
// stands.
type Executor interface {
	Start(context.Context, execution.StartRequest) (execution.StartResult, error)
	Confirm(context.Context, execution.ConfirmRequest) (execution.ConfirmResult, error)
	Stop(context.Context, execution.StopRequest) (execution.StopResult, error)
	Status(context.Context, execution.StatusRequest) (execution.StatusResult, error)
}

// NewAgentRunCommand creates the run command: the second act of the session.
func NewAgentRunCommand(executor Executor, globals *support.Globals) *cobra.Command {
	var request execution.StartRequest
	var confirm bool

	command := &cobra.Command{
		Use:   "run RUNTIME",
		Short: "Run a planned runtime: reserve against the budget, then play the stages through the coding agent",
		Args:  validateSingleArgument,
		RunE: func(command *cobra.Command, arguments []string) error {
			request.Runtime = arguments[0]
			request.Session = sessionOf(command, globals)
			if confirmsOnTheTerminal(command, confirm) {
				prompter := support.NewPrompter(command.InOrStdin(), command.ErrOrStderr())
				request.Confirm = func(ctx context.Context, stage string) (bool, error) {
					return prompter.Confirm(ctx, "Confirm "+stage+"? [y/N]")
				}
			}

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
	flags.BoolVar(&confirm, "confirm", false, "answer each human stage's confirmation on the terminal; the default when stdin is one")
	flags.DurationVar(&request.ChecksWait, "checks-wait", 10*time.Minute, "how long to follow the pull request's checks before leaving them unknown")
	byFlag(command, &request.By)

	return command
}

// confirmsOnTheTerminal says whether the run asks for the confirmations
// itself: when told to, or by default when a person is at the terminal.
func confirmsOnTheTerminal(command *cobra.Command, confirm bool) bool {
	if command.Flags().Changed("confirm") {
		return confirm
	}

	return support.Interactive(command.InOrStdin())
}

// NewAgentConfirmCommand creates the confirm command: the owner's yes to a
// stage that waits, for the next run to play on from.
func NewAgentConfirmCommand(executor Executor, globals *support.Globals) *cobra.Command {
	var by string

	command := &cobra.Command{
		Use:   "confirm RUNTIME STAGE",
		Short: "Confirm a waiting stage as the owner; the next run plays on from it",
		Args:  validateTwoArguments,
		RunE: func(command *cobra.Command, arguments []string) error {
			request := execution.ConfirmRequest{
				Runtime: arguments[0],
				Stage:   arguments[1],
				By:      by,
				Session: sessionOf(command, globals),
			}

			return runSession(command, globals, func(ctx context.Context, project application.ProjectSelection) (execution.ConfirmResult, error) {
				request.Project = project

				return executor.Confirm(ctx, request)
			})
		},
	}
	byFlag(command, &by)

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
