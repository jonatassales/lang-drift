package catalog

import (
	"context"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/presentation/cli/support"

	"github.com/spf13/cobra"
)

// ClusterImageLoader loads a local image into the project's cluster.
type ClusterImageLoader interface {
	Execute(context.Context, application.LoadClusterImageRequest) (application.ImageResult, error)
}

// NewClusterLoadImageCommand creates the image loading command.
func NewClusterLoadImageCommand(loadImage ClusterImageLoader, globals *support.Globals) *cobra.Command {
	return &cobra.Command{
		Use:   "load-image IMAGE",
		Short: "Load an already-built local image into the cluster; never build it",
		Args:  validateSingleArgument,
		RunE: func(command *cobra.Command, arguments []string) error {
			return runFinite(command, globals, func(ctx context.Context, project application.ProjectSelection) (application.ImageResult, error) {
				return loadImage.Execute(ctx, application.LoadClusterImageRequest{Project: project, Image: arguments[0]})
			})
		},
	}
}
