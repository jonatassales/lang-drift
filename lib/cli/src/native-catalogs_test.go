package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jonatassales/tao/lib/cli/src/application"
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
	if err := root.ExecuteContext(ctx); err != nil {
		// Mirrors the real CLI's documented contract until this retained
		// adapter is deleted with the retired runtime tree.
		status := 1
		var usage support.UsageError
		switch {
		case ctx.Err() != nil || errors.Is(err, context.Canceled):
			status = 130
		case errors.As(err, &usage):
			status = 2
		case errors.Is(err, application.ErrAdoptionRefused) || errors.Is(err, application.ErrTargetMismatch) || errors.Is(err, application.ErrConfigurationChanged) || errors.Is(err, application.ErrAlreadyExists):
			status = 3
		case errors.Is(err, application.ErrNotFound):
			status = 4
		}
		if globals.JSON {
			_ = json.NewEncoder(options.Diagnostics).Encode(map[string]any{"error": err.Error(), "exitCode": status})
		} else {
			_, _ = fmt.Fprintln(options.Diagnostics, "ex:", err)
		}
		return status
	}
	return 0
}
