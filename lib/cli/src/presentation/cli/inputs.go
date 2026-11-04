package cli

import (
	"github.com/jonatassales/tao/lib/cli/src/presentation/cli/catalog"

	"github.com/spf13/cobra"
)

// harnessCommand registers the inputs of the delivery harness: the envelope,
// the workflow, the model and the driver.
func (cliApplication *application) harnessCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "harness",
		Short: "Record the delivery harness: the envelope, the workflow, the model and the driver",
		Args:  validateGroupArguments,
		RunE:  requireSubcommand,
	}

	setter, globals := cliApplication.useCases.Inputs, &cliApplication.globals

	root.AddCommand(
		inputGroup("envelope", "The hardware envelope", catalog.NewEnvelopeSetCommand(setter, globals)),
		inputGroup("workflow", "The delivery stages and where a person takes part", catalog.NewWorkflowSetCommand(setter, globals)),
		inputGroup("model", "The default model and effort", catalog.NewModelSetCommand(setter, globals)),
		inputGroup("driver", "The coding agent the stages run through", catalog.NewDriverSetCommand(setter, globals)),
	)

	return root
}

// finopsCommand registers the money inputs: the budget, the gates, the prices
// and the forecasts.
func (cliApplication *application) finopsCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "finops",
		Short: "Record the money rules: the budget, the gates, the prices and the forecasts",
		Args:  validateGroupArguments,
		RunE:  requireSubcommand,
	}

	setter, globals := cliApplication.useCases.Inputs, &cliApplication.globals

	root.AddCommand(
		inputGroup("budget", "The team budget for a period", catalog.NewBudgetSetCommand(setter, globals)),
		inputGroup("gates", "The money gates a run passes", catalog.NewGatesSetCommand(setter, globals)),
		inputGroup("prices", "A dated price list", catalog.NewPricesSetCommand(setter, globals)),
		inputGroup("forecasts", "The tokens a run is expected to use", catalog.NewForecastsSetCommand(setter, globals)),
	)

	return root
}

// inputGroup is one input with its set command.
func inputGroup(name, short string, set *cobra.Command) *cobra.Command {
	group := &cobra.Command{
		Use:   name,
		Short: short,
		Args:  validateGroupArguments,
		RunE:  requireSubcommand,
	}

	group.AddCommand(set)

	return group
}
