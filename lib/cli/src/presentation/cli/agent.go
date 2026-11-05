package cli

import (
	"github.com/jonatassales/tao/lib/cli/src/presentation/cli/catalog"

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

	if cliApplication.useCases.Plan != nil {
		agentRoot.AddCommand(
			catalog.NewAgentPlanCommand(cliApplication.useCases.Plan, &cliApplication.globals),
			catalog.NewAgentReindexCommand(cliApplication.useCases.Plan, &cliApplication.globals),
		)
	}

	if cliApplication.useCases.Execution != nil {
		agentRoot.AddCommand(
			catalog.NewAgentRunCommand(cliApplication.useCases.Execution, &cliApplication.globals),
			catalog.NewAgentConfirmCommand(cliApplication.useCases.Execution, &cliApplication.globals),
			catalog.NewAgentStatusCommand(cliApplication.useCases.Execution, &cliApplication.globals),
		)
	}

	if cliApplication.useCases.ProjectAgent != nil {
		agentRoot.AddCommand(
			catalog.NewAgentConnectCommand(cliApplication.useCases.ProjectAgent, &cliApplication.globals),
			catalog.NewAgentProjectCommand(cliApplication.useCases.ProjectAgent, &cliApplication.globals),
			catalog.NewAgentSuperviseCommand(cliApplication.useCases.ProjectAgent, &cliApplication.globals),
			catalog.NewAgentMigrateCommand(cliApplication.useCases.ProjectAgent, &cliApplication.globals),
		)
	}
	return agentRoot
}
