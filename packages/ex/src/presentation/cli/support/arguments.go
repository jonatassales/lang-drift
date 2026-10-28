package support

import (
	"fmt"

	"github.com/spf13/cobra"
)

// ExactArguments preserves the CLI's argument-count diagnostic and classifies
// mismatches as invalid usage, including zero-argument commands.
func ExactArguments(count int) cobra.PositionalArgs {
	return func(command *cobra.Command, arguments []string) error {
		if len(arguments) != count {
			return UsageError{Cause: fmt.Errorf("%s requires %d argument(s)", command.CommandPath(), count)}
		}

		return nil
	}
}
