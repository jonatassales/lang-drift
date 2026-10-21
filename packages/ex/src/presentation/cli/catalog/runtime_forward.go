package catalog

import (
	"context"
	"fmt"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/presentation/cli/support"

	"github.com/spf13/cobra"
)

// RuntimeForwarder owns a foreground tunnel to one recorded Service.
type RuntimeForwarder interface {
	Execute(context.Context, application.ForwardRuntimeRequest) error
}

// NewRuntimeForwardCommand selects ports without applying the finite-operation
// timeout. The tunnel remains in the foreground until interrupted.
func NewRuntimeForwardCommand(forward RuntimeForwarder, globals *support.Globals) *cobra.Command {
	var unit string
	var localPort, remotePort int
	command := &cobra.Command{
		Use: "forward <id>", Short: "Forward a service to localhost until interrupted", Args: validateSingleArgument,
		RunE: func(command *cobra.Command, arguments []string) error {
			if unit == "" {
				return support.UsageError{Cause: fmt.Errorf("--unit is required")}
			}
			selection, err := globals.Selection()
			if err != nil {
				return err
			}
			return forward.Execute(command.Context(), application.ForwardRuntimeRequest{
				RuntimeUnitRequest: application.RuntimeUnitRequest{
					RuntimeRequest: application.RuntimeRequest{Project: selection, Runtime: arguments[0]},
					Unit:           unit,
				},
				LocalPort: localPort, RemotePort: remotePort, Diagnostics: command.ErrOrStderr(),
			})
		},
	}
	command.Flags().StringVar(&unit, "unit", "", "selected runtime unit")
	command.Flags().IntVar(&localPort, "port", 8080, "local TCP port")
	command.Flags().IntVar(&remotePort, "remote-port", 0, "service port (required when ambiguous)")

	return command
}
