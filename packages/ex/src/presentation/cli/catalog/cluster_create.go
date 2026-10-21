package catalog

import (
	"context"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/presentation/cli/support"

	"github.com/spf13/cobra"
)

// ClusterCreator creates the project's cluster.
type ClusterCreator interface {
	Execute(context.Context, application.ClusterRequest) (application.ClusterResult, error)
}

// NewClusterCreateCommand creates the cluster creation command.
func NewClusterCreateCommand(createCluster ClusterCreator, globals *support.Globals) *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "Create the project's local cluster and private kubeconfig",
		Args:  validateNoArguments,
		RunE: func(command *cobra.Command, _ []string) error {
			return runFinite(command, globals, func(ctx context.Context, project application.ProjectSelection) (application.ClusterResult, error) {
				return createCluster.Execute(ctx, application.ClusterRequest{Project: project})
			})
		},
	}
}
