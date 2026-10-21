package catalog

import (
	"context"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/presentation/cli/support"

	"github.com/spf13/cobra"
)

// ClusterDeleter deletes the project's cluster.
type ClusterDeleter interface {
	Execute(context.Context, application.ClusterRequest) (application.ClusterResult, error)
}

// NewClusterDeleteCommand creates the cluster deletion command.
func NewClusterDeleteCommand(deleteCluster ClusterDeleter, globals *support.Globals) *cobra.Command {
	return &cobra.Command{
		Use:   "delete",
		Short: "Delete the project's cluster once every runtime is cleaned up",
		Args:  validateNoArguments,
		RunE: func(command *cobra.Command, _ []string) error {
			return runFinite(command, globals, func(ctx context.Context, project application.ProjectSelection) (application.ClusterResult, error) {
				return deleteCluster.Execute(ctx, application.ClusterRequest{Project: project})
			})
		},
	}
}
