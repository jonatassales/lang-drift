package main

import (
	"context"
	"errors"
	"time"

	"github.com/jonatassales/tao/lib/cli/src/presentation/cli"
	"github.com/jonatassales/tao/lib/cli/src/presentation/cli/catalog"
	"github.com/jonatassales/tao/lib/cli/src/presentation/cli/support"
	"github.com/spf13/cobra"
)

// executeNativeCatalogs is a test-only adapter for retained native lifecycle
// contracts. Public console-agent behavior is covered by separate tests.
func executeNativeCatalogs(ctx context.Context, args []string, options support.Options, cases cli.UseCases, native lifecycles) int {
	if len(args) == 0 || args[0] != "runtime" {
		return cli.Execute(ctx, args, options, cases)
	}
	options = options.WithDefaults()
	globals := support.Globals{Fallback: options.Directory, Timeout: 2 * time.Minute}
	root := &cobra.Command{Use: "runtime", SilenceUsage: true, SilenceErrors: true}
	root.PersistentFlags().StringVar(&globals.Project, "project", "", "project root")
	root.PersistentFlags().BoolVar(&globals.JSON, "json", false, "JSON results")
	root.PersistentFlags().DurationVar(&globals.Timeout, "timeout", 2*time.Minute, "finite operation deadline")
	root.AddCommand(catalog.NewRuntimeInitCommand(native.Init, &globals), catalog.NewRuntimeValidateCommand(native.Validate, &globals), catalog.NewRuntimeUpCommand(native.Up, &globals), catalog.NewRuntimeDownCommand(native.Down, &globals), catalog.NewRuntimeInspectCommand(native.Inspect, &globals), catalog.NewRuntimeLogsCommand(native.Logs, &globals), catalog.NewRuntimeForwardCommand(native.Forward, &globals), catalog.NewRuntimeHookCommand(native.Hook, &globals))
	root.SetArgs(args[1:])
	root.SetIn(options.Input)
	root.SetOut(options.Output)
	root.SetErr(options.Diagnostics)

	err := root.ExecuteContext(ctx)
	if interrupted := ctx.Err(); interrupted != nil {
		err = errors.Join(interrupted, err)
	}

	// Mirrors the real CLI's stop rule until this retained adapter is
	// deleted with the retired runtime tree.
	return cli.Stop(err, options.Diagnostics)
}
