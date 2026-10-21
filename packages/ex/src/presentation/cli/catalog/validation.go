package catalog

import (
	"github.com/exphere/ex/src/presentation/cli/support"

	"github.com/spf13/cobra"
)

// validateNoArguments classifies unexpected positional arguments as CLI usage errors.
func validateNoArguments(command *cobra.Command, arguments []string) error {
	return support.ExactArguments(0)(command, arguments)
}

// validateSingleArgument requires exactly one positional argument.
func validateSingleArgument(command *cobra.Command, arguments []string) error {
	return support.ExactArguments(1)(command, arguments)
}
