package cli

// These historical catalog fixtures preserve request/timeout/error assertions
// for the native adapter reused by the agent. They do not expose retired
// commands in the production executable.
import (
	"context"
	"errors"
	"fmt"
	"github.com/jonatassales/tao/packages/ex/src/presentation/cli/catalog"
	"github.com/jonatassales/tao/packages/ex/src/presentation/cli/support"
	"github.com/spf13/cobra"
)

func executeCatalogs(
	applicationContext context.Context,
	arguments []string,
	options support.Options,
	useCases UseCases,
) int {
	options = options.WithDefaults()

	outputWriter := support.NewCheckedWriter(options.Output)
	options.Output = outputWriter

	cliApplication := &application{
		options:  options,
		useCases: useCases,
		globals:  support.Globals{Fallback: options.Directory},
	}
	rootCommand := cliApplication.command()
	for _, command := range rootCommand.Commands() {
		if command.Name() == "runtime" {
			rootCommand.RemoveCommand(command)
		}
	}
	runtimeRoot := &cobra.Command{Use: "runtime", Args: validateGroupArguments, RunE: requireSubcommand}
	globals := &cliApplication.globals
	runtimeRoot.AddCommand(catalog.NewRuntimeInitCommand(useCases.RuntimeInit, globals), catalog.NewRuntimeValidateCommand(useCases.RuntimeValidate, globals), catalog.NewRuntimeUpCommand(useCases.RuntimeUp, globals), catalog.NewRuntimeDownCommand(useCases.RuntimeDown, globals), catalog.NewRuntimeInspectCommand(useCases.RuntimeInspect, globals), catalog.NewRuntimeLogsCommand(useCases.RuntimeLogs, globals), catalog.NewRuntimeForwardCommand(useCases.RuntimeForward, globals), catalog.NewRuntimeHookCommand(useCases.RuntimeHook, globals))
	rootCommand.AddCommand(runtimeRoot)

	rootCommand.SetArgs(arguments)
	rootCommand.SetIn(options.Input)
	rootCommand.SetOut(outputWriter)
	rootCommand.SetErr(options.Diagnostics)

	_, executionError := rootCommand.ExecuteContextC(applicationContext)

	if outputError := outputWriter.Err(); outputError != nil {
		executionError = errors.Join(
			executionError,
			fmt.Errorf("write output: %w", outputError),
		)
	}

	if executionError == nil {
		return 0
	}

	exitStatus := exitCode(applicationContext, executionError)

	reportError(
		options.Diagnostics,
		executionError,
		exitStatus,
		cliApplication.globals.JSON,
	)

	return exitStatus
}
