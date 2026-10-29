package catalog

import (
	"github.com/jonatassales/tao/packages/ex/src/presentation/cli/support"

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

// validateTwoArguments requires exactly two positional arguments.
func validateTwoArguments(command *cobra.Command, arguments []string) error {
	return support.ExactArguments(2)(command, arguments)
}
