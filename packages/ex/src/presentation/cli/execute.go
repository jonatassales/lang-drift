// Package cli is EX's command-line presentation: the Cobra command tree, the
// global flags, and how failures become diagnostics and exit codes. It maps
// input to application requests through the catalog and never orchestrates
// work or reaches infrastructure itself.
package cli

import (
	"context"
	"errors"
	"fmt"

	"github.com/jonatassales/tao/packages/ex/src/presentation/cli/catalog"
	"github.com/jonatassales/tao/packages/ex/src/presentation/cli/support"
)

// UseCases are the application operations the command tree exposes. Each
// field is the narrow contract its command consumes; the composition root
// supplies the implementations.
type UseCases struct {
	ProjectInit      catalog.ProjectInitializer
	ProjectAgent     catalog.ProjectAgent
	Doctor           catalog.Diagnoser
	ClusterCreate    catalog.ClusterCreator
	ClusterDelete    catalog.ClusterDeleter
	ClusterInspect   catalog.ClusterInspector
	ClusterLoadImage catalog.ClusterImageLoader
	AgentInit        catalog.AgentInitializer
	Inputs           catalog.Inputs
}

// Execute runs the CLI and returns the appropriate exit code.
func Execute(
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
