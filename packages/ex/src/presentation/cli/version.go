package cli

import (
	"fmt"

	"github.com/exphere/ex/src/presentation/cli/support"

	"github.com/spf13/cobra"
)

// registerVersion keeps the existing subcommand, long flag, and -V alias
// on one renderer. It needs no project selection or application operation.
func (cliApplication *application) registerVersion(root *cobra.Command) {
	version := &cobra.Command{
		Use: "version", Short: "Print the CLI version",
		Args: support.ExactArguments(0),
		RunE: func(command *cobra.Command, _ []string) error {
			if cliApplication.globals.JSON {
				return support.Render(command.OutOrStdout(), map[string]string{"version": cliApplication.options.Version}, true)
			}
			_, err := fmt.Fprintf(command.OutOrStdout(), "ex %s\n", cliApplication.options.Version)
			return err
		},
	}
	root.AddCommand(version)
	var versionFlag bool
	root.Flags().BoolVarP(&versionFlag, "version", "V", false, "print the CLI version")
	root.RunE = func(command *cobra.Command, _ []string) error {
		if versionFlag {
			return version.RunE(command, nil)
		}
		return command.Help()
	}
}
