package catalog

import (
	"context"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/presentation/cli/support"

	"github.com/spf13/cobra"
)

// AgentInitializer scaffolds an agent project.
type AgentInitializer interface {
	Execute(context.Context, application.InitializeAgentRequest) (application.AgentResult, error)
}

// NewAgentInitCommand creates the agent initialization command.
func NewAgentInitCommand(initializeAgent AgentInitializer, globals *support.Globals) *cobra.Command {
	return &cobra.Command{
		Use:   "init NAME",
		Short: "Generate an editable TypeScript LangGraph agent scaffold",
		Args:  validateSingleArgument,
		RunE: func(command *cobra.Command, arguments []string) error {
			return runFinite(command, globals, func(ctx context.Context, project application.ProjectSelection) (application.AgentResult, error) {
				return initializeAgent.Execute(ctx, application.InitializeAgentRequest{Project: project, Name: arguments[0]})
			})
		},
	}
}
