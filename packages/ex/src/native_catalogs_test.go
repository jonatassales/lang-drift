package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/exphere/ex/src/presentation/cli"
	"github.com/exphere/ex/src/presentation/cli/catalog"
	"github.com/exphere/ex/src/presentation/cli/support"
	"github.com/spf13/cobra"
)

// executeNativeCatalogs is a test-only adapter for retained native lifecycle
// contracts. Public console-agent behavior is covered by separate tests.
func executeNativeCatalogs(ctx context.Context, args []string, options support.Options, cases cli.UseCases) int {
	if len(args) == 0 || args[0] != "runtime" {
		return cli.Execute(ctx, args, options, cases)
	}
	options = options.WithDefaults()
	globals := support.Globals{Fallback: options.Directory, Timeout: 2 * time.Minute}
	root := &cobra.Command{Use: "runtime", SilenceUsage: true, SilenceErrors: true}
	root.PersistentFlags().StringVar(&globals.Project, "project", "", "project root")
	root.PersistentFlags().BoolVar(&globals.JSON, "json", false, "JSON results")
	root.PersistentFlags().DurationVar(&globals.Timeout, "timeout", 2*time.Minute, "finite operation deadline")
	root.AddCommand(catalog.NewRuntimeInitCommand(cases.RuntimeInit, &globals), catalog.NewRuntimeValidateCommand(cases.RuntimeValidate, &globals), catalog.NewRuntimeUpCommand(cases.RuntimeUp, &globals), catalog.NewRuntimeDownCommand(cases.RuntimeDown, &globals), catalog.NewRuntimeInspectCommand(cases.RuntimeInspect, &globals), catalog.NewRuntimeLogsCommand(cases.RuntimeLogs, &globals), catalog.NewRuntimeForwardCommand(cases.RuntimeForward, &globals), catalog.NewRuntimeHookCommand(cases.RuntimeHook, &globals))
	root.SetArgs(args[1:])
	root.SetIn(options.Input)
	root.SetOut(options.Output)
	root.SetErr(options.Diagnostics)
	if err := root.ExecuteContext(ctx); err != nil {
		status := 1
		var usage support.UsageError
		if ctx.Err() != nil || errors.Is(err, context.Canceled) {
			status = 130
		} else if errors.As(err, &usage) {
			status = 2
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
