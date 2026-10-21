package cli

import (
	"fmt"

	"github.com/exphere/ex/src/presentation/cli/support"

	"github.com/spf13/cobra"
)

// registerHelp refuses unknown command paths instead of printing unrelated
// root help and reporting success. Help remains independent of project state.
func registerHelp(root *cobra.Command) {
	root.SetHelpCommand(&cobra.Command{
		Use: "help [command]", Short: "Help about any command",
		RunE: func(command *cobra.Command, arguments []string) error {
			selected, remainder, err := root.Find(arguments)
			if err != nil || len(remainder) > 0 {
				return support.UsageError{Cause: fmt.Errorf("unknown help command")}
			}
			return selected.Help()
		},
	})
}
