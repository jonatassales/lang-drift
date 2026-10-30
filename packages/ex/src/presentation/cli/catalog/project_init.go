package catalog

import (
	"context"

	"github.com/jonatassales/tao/packages/ex/src/application"
	"github.com/jonatassales/tao/packages/ex/src/application/inputs"
	"github.com/jonatassales/tao/packages/ex/src/application/setup"
	"github.com/jonatassales/tao/packages/ex/src/presentation/cli/support"

	"github.com/spf13/cobra"
)

// ProjectSetup records a project's identity and the inputs given on the way
// in.
type ProjectSetup interface {
	Execute(context.Context, setup.Request) (setup.Result, error)
}

// NewProjectInitCommand creates the project initialization command. A person
// at the terminal is asked for the inputs the flags leave unanswered; --ask
// decides that either way, and --json never asks.
func NewProjectInitCommand(setUpProject ProjectSetup, globals *support.Globals) *cobra.Command {
	var values initFlags
	var ask bool
	var by string

	command := &cobra.Command{
		Use:   "init",
		Short: "Initialize project metadata and record the first inputs, without starting infrastructure",
		Args:  validateNoArguments,
		RunE: func(command *cobra.Command, _ []string) error {
			answers, err := initAnswers(command.Flags(), values)
			if err != nil {
				return err
			}

			request := setup.Request{
				By:      by,
				Answers: answers,
			}

			asks := support.Interactive(command.InOrStdin())
			if command.Flags().Changed("ask") {
				asks = ask
			}

			runCommand := runFinite[setup.Result]

			if asks && !globals.JSON {
				request.Guide = initGuide{prompter: support.NewPrompter(command.InOrStdin(), command.ErrOrStderr())}
				runCommand = runInteractive[setup.Result]
			}

			return runCommand(command, globals, func(ctx context.Context, project application.ProjectSelection) (setup.Result, error) {
				request.Directory = project.Directory

				return setUpProject.Execute(ctx, request)
			})
		},
	}

	flags := command.Flags()
	flags.BoolVar(&ask, "ask", false, "ask for the inputs the flags leave unanswered; the default when stdin is a terminal")
	flags.StringVar(&values.budget, "budget", "", "this month's team budget, in dollars")
	flags.StringVar(&values.perRun, "per-run", "", "the most one run may reserve, in dollars")
	flags.IntVar(&values.reserveMargin, "reserve-margin", 0, "the margin a reservation adds to an estimate, in percent")
	flags.BoolVar(&values.defaultGates, "default-gates", false, "apply the reference gates, FG-01 to FG-06")
	flags.StringVar(&values.artifacts, "artifacts", "", "the folder of the team's tickets and specs, relative to the project")
	flags.StringVar(&by, "by", inputs.DefaultBy, "who sets the inputs")

	return command
}
