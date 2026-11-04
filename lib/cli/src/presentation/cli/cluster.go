package cli

import (
	"github.com/jonatassales/tao/lib/cli/src/presentation/cli/catalog"

	"github.com/spf13/cobra"
)

// clusterCommand registers all cluster-related CLI commands.
func (cliApplication *application) clusterCommand() *cobra.Command {
	clusterRoot := &cobra.Command{
		Use:   "cluster",
		Short: "Manage the project's explicitly owned kind cluster",
		Args:  validateGroupArguments,
		RunE:  requireSubcommand,
	}

	useCases, globals := cliApplication.useCases, &cliApplication.globals

	clusterRoot.AddCommand(
		catalog.NewClusterCreateCommand(useCases.ClusterCreate, globals),
		catalog.NewClusterDeleteCommand(useCases.ClusterDelete, globals),
		catalog.NewClusterInspectCommand(useCases.ClusterInspect, globals),
		catalog.NewClusterLoadImageCommand(useCases.ClusterLoadImage, globals),
	)

	return clusterRoot
}
