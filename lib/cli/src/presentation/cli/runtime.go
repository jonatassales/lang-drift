package cli

import (
	"errors"
	"github.com/spf13/cobra"
)

// runtimeCommand keeps an actionable retirement message for previous users.
// Native lifecycle use cases remain available only to the connected agent.
func (cliApplication *application) runtimeCommand() *cobra.Command {
	return &cobra.Command{Use: "runtime", Short: "Runtime operations moved to the console", DisableFlagParsing: true, Args: cobra.ArbitraryArgs, RunE: func(*cobra.Command, []string) error {
		return errors.New("runtime operations moved to the console; run dft agent connect <console-url> and dft agent run; migrate JSON v2 with dft agent migrate <runtime-id>")
	}}
}
