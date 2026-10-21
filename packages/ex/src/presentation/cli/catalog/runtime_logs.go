package catalog

import (
	"context"
	"fmt"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/presentation/cli/support"

	"github.com/spf13/cobra"
)

// RuntimeLogReader streams one owned workload's logs.
type RuntimeLogReader interface {
	Execute(context.Context, application.ReadRuntimeLogsRequest) error
}

// NewRuntimeLogsCommand preserves native text output and explicit follow.
func NewRuntimeLogsCommand(readLogs RuntimeLogReader, globals *support.Globals) *cobra.Command {
	var unit string
	var follow bool
	command := &cobra.Command{
		Use: "logs <id>", Short: "Read the selected workload's logs", Args: validateSingleArgument,
		RunE: func(command *cobra.Command, arguments []string) error {
			if unit == "" {
				return support.UsageError{Cause: fmt.Errorf("--unit is required")}
			}
			if globals.JSON {
				return support.UsageError{Cause: fmt.Errorf("logs emit native text; --json is unsupported")}
			}
			selection, err := globals.Selection()
			if err != nil {
				return err
			}
			ctx := command.Context()
			if !follow {
				var cancel context.CancelFunc
				ctx, cancel = globals.Bounded(ctx)
				defer cancel()
			}
			return readLogs.Execute(ctx, application.ReadRuntimeLogsRequest{
				RuntimeUnitRequest: application.RuntimeUnitRequest{
					RuntimeRequest: application.RuntimeRequest{Project: selection, Runtime: arguments[0]},
					Unit:           unit,
				},
				Follow: follow, Output: command.OutOrStdout(), Diagnostics: command.ErrOrStderr(),
			})
		},
	}
	command.Flags().StringVar(&unit, "unit", "", "selected runtime unit")
	command.Flags().BoolVar(&follow, "follow", false, "stream until interrupted")

	return command
}
