package catalog

import (
	"context"
	"io"
	"os"

	"github.com/jonatassales/tao/packages/ex/src/application"
	"github.com/jonatassales/tao/packages/ex/src/application/connected"
	"github.com/jonatassales/tao/packages/ex/src/presentation/cli/support"
	"github.com/spf13/cobra"
)

// ProjectAgent is the independent connected-agent use case, not an SDK API.
type ProjectAgent interface {
	Connect(context.Context, connected.ConnectRequest) (connected.Result, error)
	Run(context.Context, connected.RunRequest) error
	RunProject(context.Context, connected.ProjectRunRequest) error
	Migrate(context.Context, connected.MigrateRequest) (connected.Result, error)
}

func NewAgentConnectCommand(agent ProjectAgent, globals *support.Globals) *cobra.Command {
	return &cobra.Command{Use: "connect <console-url>", Short: "Bind this checkout using RUNTIME_SECRET from the environment", Args: cobra.ExactArgs(1), RunE: func(command *cobra.Command, args []string) error {
		return runFinite(command, globals, func(ctx context.Context, selection application.ProjectSelection) (connected.Result, error) {
			return agent.Connect(ctx, connected.ConnectRequest{Project: selection, URL: args[0], Secret: os.Getenv("RUNTIME_SECRET")})
		})
	}}
}
func NewAgentRunCommand(agent ProjectAgent, globals *support.Globals) *cobra.Command {
	return &cobra.Command{Use: "run", Short: "Synchronize connected runtimes until interrupted", Args: validateNoArguments, RunE: func(command *cobra.Command, _ []string) error {
		selection, err := globals.Selection()
		if err != nil {
			return err
		}
		return agent.Run(command.Context(), connected.RunRequest{Project: selection})
	}}
}
func NewAgentMigrateCommand(agent ProjectAgent, globals *support.Globals) *cobra.Command {
	var recover bool
	command := &cobra.Command{Use: "migrate <runtime-id>", Short: "Prepare migration of a stopped JSON v2 runtime without changing native files", Args: cobra.ExactArgs(1), RunE: func(command *cobra.Command, args []string) error {
		return runFinite(command, globals, func(ctx context.Context, selection application.ProjectSelection) (connected.Result, error) {
			return agent.Migrate(ctx, connected.MigrateRequest{Project: selection, Runtime: args[0], Recover: recover})
		})
	}}
	command.Flags().BoolVar(&recover, "recover", false, "Explicitly clean recorded legacy runtime resources before preparing migration")
	return command
}

// NewAgentProjectCommand runs the single agent belonging to this project's .langdrift.
// A supervising console owns stdin; its closure cancels work even after a crash.
func NewAgentProjectCommand(agent ProjectAgent, globals *support.Globals) *cobra.Command {
	var supervised bool
	command := &cobra.Command{Use: "project <console-url>", Short: "Connect this project and synchronize all its microruntimes", Args: cobra.ExactArgs(1), RunE: func(command *cobra.Command, args []string) error {
		selection, err := globals.Selection()
		if err != nil {
			return err
		}
		ctx, cancel := context.WithCancel(command.Context())
		defer cancel()
		if supervised {
			input := os.Stdin
			done := make(chan struct{})
			go func() { defer close(done); _, _ = io.Copy(io.Discard, input); cancel() }()
			defer func() { _ = input.Close(); <-done }()
		}
		return agent.RunProject(ctx, connected.ProjectRunRequest{Project: selection, URL: args[0], Secret: os.Getenv("PROJECT_SECRET")})
	}}
	command.Flags().BoolVar(&supervised, "supervised", false, "Stop when the supervising console closes its input pipe")
	return command
}
