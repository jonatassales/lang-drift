package cli

import (
	"github.com/jonatassales/tao/lib/cli/src/presentation/cli/catalog"

	"github.com/spf13/cobra"
)

// gatesCommand registers CI Gates: the report a pull request's workflow runs.
func (cliApplication *application) gatesCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "gates",
		Short: "CI Gates: the report a pull request's workflow runs",
		Args:  validateGroupArguments,
		RunE:  requireSubcommand,
	}

	root.AddCommand(catalog.NewGatesReportCommand(cliApplication.useCases.Gates, &cliApplication.globals))

	return root
}
