package catalog

import (
	appsession "github.com/jonatassales/tao/lib/cli/src/application/session"
	"github.com/jonatassales/tao/lib/cli/src/presentation/cli/session"
	"github.com/jonatassales/tao/lib/cli/src/presentation/cli/support"

	"github.com/spf13/cobra"
)

// sessionOf is the sink a session command prints through: toned on a
// terminal, and seen by the observer when one is set.
func sessionOf(command *cobra.Command, globals *support.Globals) appsession.Sink {
	printed := session.For(command.OutOrStdout(), globals.JSON, globals.Terminal)

	return session.Observed(printed, globals.Observer)
}
