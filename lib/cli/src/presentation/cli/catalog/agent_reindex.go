package catalog

import (
	"context"

	"github.com/jonatassales/tao/lib/cli/src/application"
	"github.com/jonatassales/tao/lib/cli/src/application/plan"
	"github.com/jonatassales/tao/lib/cli/src/presentation/cli/support"

	"github.com/spf13/cobra"
)

// Reindexer rebuilds the project's episode index.
type Reindexer interface {
	Reindex(context.Context, plan.ReindexRequest) (plan.ReindexResult, error)
}

// NewAgentReindexCommand creates the reindex command.
func NewAgentReindexCommand(reindexer Reindexer, globals *support.Globals) *cobra.Command {
	return &cobra.Command{
		Use:   "reindex",
		Short: "Rebuild the episode index from every runtime's memory",
		Args:  validateNoArguments,
		RunE: func(command *cobra.Command, _ []string) error {
			return runFinite(command, globals, func(ctx context.Context, project application.ProjectSelection) (plan.ReindexResult, error) {
				return reindexer.Reindex(ctx, plan.ReindexRequest{Project: project})
			})
		},
	}
}

// PlanOperations are what the plan commands consume together.
type PlanOperations interface {
	Planner
	Reindexer
}
