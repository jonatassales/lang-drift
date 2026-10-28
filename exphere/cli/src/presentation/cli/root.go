package cli

import (
	"fmt"
	"time"

	"github.com/jonatassales/tao/exphere/cli/src/presentation/cli/catalog"
	"github.com/jonatassales/tao/exphere/cli/src/presentation/cli/support"

	"github.com/spf13/cobra"
)

// defaultTimeout is the --timeout default for one finite operation.
const defaultTimeout = 2 * time.Minute

type application struct {
	options  support.Options
	useCases UseCases
	globals  support.Globals
}

// command creates the EX root command.
func (cliApplication *application) command() *cobra.Command {
	rootCommand := &cobra.Command{
		Use:           "ex",
		Short:         "LangDrift CLI",
		Long:          "LangDrift CLI — project setup and connected console agent",
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          validateRootArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			return command.Help()
		},
	}

	rootCommand.SetFlagErrorFunc(
		func(_ *cobra.Command, flagError error) error {
			return support.UsageError{Cause: flagError}
		},
	)

	cliApplication.registerGlobalFlags(rootCommand)
	cliApplication.registerVersion(rootCommand)
	cliApplication.registerCommands(rootCommand)
	registerHelp(rootCommand)

	return rootCommand
}

// registerGlobalFlags configures flags shared by CLI commands.
func (cliApplication *application) registerGlobalFlags(
	rootCommand *cobra.Command,
) {
	persistentFlags := rootCommand.PersistentFlags()

	persistentFlags.StringVar(
		&cliApplication.globals.Project,
		"project",
		"",
		"project root (defaults to discovery from the working directory)",
	)

	persistentFlags.BoolVar(
		&cliApplication.globals.JSON,
		"json",
		false,
		"emit compact results and structured diagnostics",
	)

	persistentFlags.DurationVar(
		&cliApplication.globals.Timeout,
		"timeout",
		defaultTimeout,
		"total deadline for a finite operation",
	)

	rootCommand.PersistentPreRunE = func(
		_ *cobra.Command,
		_ []string,
	) error {
		if cliApplication.globals.Timeout <= 0 {
			return support.UsageError{
				Cause: fmt.Errorf("--timeout must be positive"),
			}
		}

		return nil
	}
}

// registerCommands connects the top-level commands and command groups.
func (cliApplication *application) registerCommands(rootCommand *cobra.Command) {
	rootCommand.AddCommand(
		catalog.NewProjectInitCommand(cliApplication.useCases.ProjectInit, &cliApplication.globals),
		catalog.NewDoctorCommand(cliApplication.useCases.Doctor, &cliApplication.globals),
		cliApplication.runtimeCommand(),
		cliApplication.clusterCommand(),
		cliApplication.agentCommand(),
	)
}

// requireSubcommand provides a consistent error for command groups.
func requireSubcommand(command *cobra.Command, _ []string) error {
	return support.UsageError{
		Cause: fmt.Errorf(
			"%s requires a subcommand; run '%s --help'",
			command.CommandPath(),
			command.CommandPath(),
		),
	}
}

// validateGroupArguments rejects arguments that are not a valid subcommand.
func validateGroupArguments(command *cobra.Command, arguments []string) error {
	if validationError := cobra.NoArgs(command, arguments); validationError != nil {
		return support.UsageError{Cause: validationError}
	}

	return nil
}

// validateRootArgs rejects unexpected root arguments.
func validateRootArgs(
	_ *cobra.Command,
	arguments []string,
) error {
	if len(arguments) == 0 {
		return nil
	}

	return support.UsageError{
		Cause: fmt.Errorf(
			"unknown command %q; run 'ex --help'",
			arguments[0],
		),
	}
}
