package cli

import (
	"github.com/jonatassales/tao/packages/ex/src/presentation/cli/catalog"

	"github.com/spf13/cobra"
)

// agentCommand registers all agent-related CLI commands.
func (cliApplication *application) agentCommand() *cobra.Command {
	agentRoot := &cobra.Command{
		Use:   "agent",
		Short: "Connect and run the independent project agent",
		Args:  validateGroupArguments,
		RunE:  requireSubcommand,
	}

	agentRoot.AddCommand(
		catalog.NewAgentInitCommand(cliApplication.useCases.AgentInit, &cliApplication.globals),
	)

	if cliApplication.useCases.ProjectAgent != nil {
		agentRoot.AddCommand(
			catalog.NewAgentConnectCommand(cliApplication.useCases.ProjectAgent, &cliApplication.globals),
			catalog.NewAgentProjectCommand(cliApplication.useCases.ProjectAgent, &cliApplication.globals),
			catalog.NewAgentRunCommand(cliApplication.useCases.ProjectAgent, &cliApplication.globals),
			catalog.NewAgentMigrateCommand(cliApplication.useCases.ProjectAgent, &cliApplication.globals),
		)
	}
	return agentRoot
}
