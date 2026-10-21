package catalog

import (
	"context"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/presentation/cli/support"

	"github.com/spf13/cobra"
)

// ClusterInspector inspects the project's cluster.
type ClusterInspector interface {
	Execute(context.Context, application.ClusterRequest) (application.ClusterResult, error)
}

// NewClusterInspectCommand creates the cluster inspection command.
func NewClusterInspectCommand(inspectCluster ClusterInspector, globals *support.Globals) *cobra.Command {
	return &cobra.Command{
		Use:   "inspect",
		Short: "Verify and report the recorded cluster identity",
		Args:  validateNoArguments,
		RunE: func(command *cobra.Command, _ []string) error {
			return runFinite(command, globals, func(ctx context.Context, project application.ProjectSelection) (application.ClusterResult, error) {
				return inspectCluster.Execute(ctx, application.ClusterRequest{Project: project})
			})
		},
	}
}
