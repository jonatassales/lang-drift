package catalog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jonatassales/tao/lib/cli/src/application"
	"github.com/jonatassales/tao/lib/cli/src/application/inputs"
	"github.com/jonatassales/tao/lib/cli/src/domain/finops"
	"github.com/jonatassales/tao/lib/cli/src/domain/harness"
	"github.com/jonatassales/tao/lib/cli/src/presentation/cli/support"

	"github.com/spf13/cobra"
)

// Inputs records the guardrail inputs: the harness and the money rules.
type Inputs interface {
	SetEnvelope(context.Context, inputs.EnvelopeRequest) (inputs.Result, error)
	SetWorkflow(context.Context, inputs.WorkflowRequest) (inputs.Result, error)
	SetModel(context.Context, inputs.ModelRequest) (inputs.Result, error)
	SetDriver(context.Context, inputs.DriverRequest) (inputs.Result, error)
	SetBudget(context.Context, inputs.BudgetRequest) (inputs.Result, error)
	SetGates(context.Context, inputs.GatesRequest) (inputs.Result, error)
	SetPrices(context.Context, inputs.PricesRequest) (inputs.Result, error)
	SetForecasts(context.Context, inputs.ForecastsRequest) (inputs.Result, error)
}

// byFlag names who sets an input; every set command carries it.
func byFlag(command *cobra.Command, by *string) {
	command.Flags().StringVar(by, "by", inputs.DefaultBy, "who sets the input")
}

// readInputFile decodes one input body from a JSON file the owner wrote;
// unknown fields are refused, because a typo would otherwise set nothing.
func readInputFile[Body any](path string) (Body, error) {
	var body Body

	data, err := os.ReadFile(path)
	if err != nil {
		return body, support.UsageError{Cause: fmt.Errorf("read the input file: %w", err)}
	}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&body); err != nil {
		return body, support.UsageError{Cause: fmt.Errorf("the input file %s: %w", path, err)}
	}

	return body, nil
}

// parseAllocations reads ticket=cents pairs.
func parseAllocations(pairs []string) ([]finops.Allocation, error) {
	allocations := make([]finops.Allocation, 0, len(pairs))

	for _, pair := range pairs {
		ticket, amount, found := strings.Cut(pair, "=")
		if !found || ticket == "" {
			return nil, support.UsageError{Cause: fmt.Errorf("--allocate takes ticket=cents, not %q", pair)}
		}

		cents, err := strconv.Atoi(amount)
		if err != nil {
			return nil, support.UsageError{Cause: fmt.Errorf("--allocate %s: %q is not a number of cents", ticket, amount)}
		}

		allocations = append(allocations, finops.Allocation{Ticket: ticket, Cents: cents})
	}

	return allocations, nil
}

// NewEnvelopeSetCommand sets the hardware envelope.
func NewEnvelopeSetCommand(setter Inputs, globals *support.Globals) *cobra.Command {
	var request inputs.EnvelopeRequest

	command := &cobra.Command{
		Use:   "set",
		Short: "Set the hardware envelope: the per-execution ceiling and the default",
		Args:  validateNoArguments,
		RunE: func(command *cobra.Command, _ []string) error {
			return runFinite(command, globals, func(ctx context.Context, project application.ProjectSelection) (inputs.Result, error) {
				request.Project = project

				return setter.SetEnvelope(ctx, request)
			})
		},
	}

	flags := command.Flags()
	flags.IntVar(&request.PerExecution.CPUMillis, "per-execution-cpu-millis", 0, "the most CPU one execution may take, in millicores")
	flags.IntVar(&request.PerExecution.MemoryMiB, "per-execution-memory-mib", 0, "the most memory one execution may take, in MiB")
	flags.IntVar(&request.Default.CPUMillis, "default-cpu-millis", 0, "the CPU a new runtime starts from, in millicores")
	flags.IntVar(&request.Default.MemoryMiB, "default-memory-mib", 0, "the memory a new runtime starts from, in MiB")
	flags.StringArrayVar(&request.NotRuntimeFields, "not-runtime-field", nil, "a resource runtime.md has no field for, such as gpu")
	flags.StringVar(&request.Owner, "owner", "team", "who owns the envelope")
	flags.StringVar(&request.Status, "status", "draft", "the envelope's status")
	byFlag(command, &request.By)

	return command
}

// NewWorkflowSetCommand sets the delivery stages from a file.
func NewWorkflowSetCommand(setter Inputs, globals *support.Globals) *cobra.Command {
	var file, by string

	command := &cobra.Command{
		Use:   "set",
		Short: "Set the delivery stages and where a person takes part, from a JSON file",
		Args:  validateNoArguments,
		RunE: func(command *cobra.Command, _ []string) error {
			workflow, err := readInputFile[harness.Workflow](file)
			if err != nil {
				return err
			}

			return runFinite(command, globals, func(ctx context.Context, project application.ProjectSelection) (inputs.Result, error) {
				return setter.SetWorkflow(ctx, inputs.WorkflowRequest{Project: project, By: by, Workflow: workflow})
			})
		},
	}

	command.Flags().StringVar(&file, "file", "", "the workflow body: status, owner, note, stages, defaults, efforts")
	_ = command.MarkFlagRequired("file")
	byFlag(command, &by)

	return command
}

// NewModelSetCommand sets the default model and effort.
func NewModelSetCommand(setter Inputs, globals *support.Globals) *cobra.Command {
	var request inputs.ModelRequest

	command := &cobra.Command{
		Use:   "set MODEL",
		Short: "Set the default model and effort, and the planner's model",
		Args:  validateSingleArgument,
		RunE: func(command *cobra.Command, arguments []string) error {
			request.Model = arguments[0]

			return runFinite(command, globals, func(ctx context.Context, project application.ProjectSelection) (inputs.Result, error) {
				request.Project = project

				return setter.SetModel(ctx, request)
			})
		},
	}

	command.Flags().StringVar(&request.Effort, "effort", "", "low, medium or high")
	_ = command.MarkFlagRequired("effort")
	command.Flags().StringVar(&request.Planner, "planner", "", "the model the planner's own calls use (defaults to MODEL)")
	byFlag(command, &request.By)

	return command
}

// NewDriverSetCommand sets the coding agent the stages run through.
func NewDriverSetCommand(setter Inputs, globals *support.Globals) *cobra.Command {
	var file, by string

	command := &cobra.Command{
		Use:   "set",
		Short: "Set the coding agent the workflow stages run through, from a JSON file",
		Args:  validateNoArguments,
		RunE: func(command *cobra.Command, _ []string) error {
			driver, err := readInputFile[harness.Driver](file)
			if err != nil {
				return err
			}

			return runFinite(command, globals, func(ctx context.Context, project application.ProjectSelection) (inputs.Result, error) {
				return setter.SetDriver(ctx, inputs.DriverRequest{Project: project, By: by, Driver: driver})
			})
		},
	}

	command.Flags().StringVar(&file, "file", "", "the driver: executable, args, passEnv, timeoutSeconds, usage, unmeasured")
	_ = command.MarkFlagRequired("file")
	byFlag(command, &by)

	return command
}

// NewBudgetSetCommand sets the team budget for a period.
func NewBudgetSetCommand(setter Inputs, globals *support.Globals) *cobra.Command {
	var request inputs.BudgetRequest
	var pairs []string

	command := &cobra.Command{
		Use:   "set",
		Short: "Set the team budget: the period, the per-run cap, the margin and the allocations",
		Args:  validateNoArguments,
		RunE: func(command *cobra.Command, _ []string) error {
			allocations, err := parseAllocations(pairs)
			if err != nil {
				return err
			}
			request.Allocations = allocations

			return runFinite(command, globals, func(ctx context.Context, project application.ProjectSelection) (inputs.Result, error) {
				request.Project = project

				return setter.SetBudget(ctx, request)
			})
		},
	}

	flags := command.Flags()
	flags.StringVar(&request.Period, "period", "", "the month, YYYY-MM")
	_ = command.MarkFlagRequired("period")
	flags.IntVar(&request.Bounded.PeriodCents, "period-cents", 0, "the period's amount, in cents")
	flags.IntVar(&request.Bounded.PerRunCents, "per-run-cents", 0, "the most one run may reserve, in cents")
	flags.IntVar(&request.Bounded.ReserveMarginPercent, "reserve-margin-percent", 20, "the margin a reservation adds to an estimate")
	flags.StringArrayVar(&pairs, "allocate", nil, "a ticket's share, ticket=cents")
	flags.StringVar(&request.Currency, "currency", "USD", "the currency")
	flags.StringVar(&request.Owner, "owner", "team", "who owns the budget")
	flags.StringVar(&request.Status, "status", "draft", "the budget's status")
	byFlag(command, &request.By)

	return command
}

// NewGatesSetCommand sets the money gates.
func NewGatesSetCommand(setter Inputs, globals *support.Globals) *cobra.Command {
	var file, status, by string
	var useDefault bool

	command := &cobra.Command{
		Use:   "set",
		Short: "Set the money gates: the reference six, or a JSON file",
		Args:  validateNoArguments,
		RunE: func(command *cobra.Command, _ []string) error {
			if useDefault == (file != "") {
				return support.UsageError{Cause: fmt.Errorf("choose --default or --file")}
			}

			gates := finops.DefaultGates()
			if file != "" {
				body, err := readInputFile[finops.Gates](file)
				if err != nil {
					return err
				}
				gates = body.Gates
			}

			return runFinite(command, globals, func(ctx context.Context, project application.ProjectSelection) (inputs.Result, error) {
				return setter.SetGates(ctx, inputs.GatesRequest{Project: project, By: by, Status: status, Gates: gates})
			})
		},
	}

	command.Flags().BoolVar(&useDefault, "default", false, "the reference gates, FG-01 to FG-06")
	command.Flags().StringVar(&file, "file", "", "the gates body: gates")
	command.Flags().StringVar(&status, "status", "draft", "the gates' status")
	byFlag(command, &by)

	return command
}

// NewPricesSetCommand sets a dated price list.
func NewPricesSetCommand(setter Inputs, globals *support.Globals) *cobra.Command {
	var file, slug, by string

	command := &cobra.Command{
		Use:   "set",
		Short: "Set the current price list from a JSON file, mirrored under the next number",
		Args:  validateNoArguments,
		RunE: func(command *cobra.Command, _ []string) error {
			prices, err := readInputFile[finops.Prices](file)
			if err != nil {
				return err
			}

			return runFinite(command, globals, func(ctx context.Context, project application.ProjectSelection) (inputs.Result, error) {
				return setter.SetPrices(ctx, inputs.PricesRequest{Project: project, By: by, Slug: slug, Prices: prices})
			})
		},
	}

	command.Flags().StringVar(&file, "file", "", "the price list: note, readOn, currency, unit, cache, models")
	_ = command.MarkFlagRequired("file")
	command.Flags().StringVar(&slug, "slug", "", "the mirror's name after its number")
	_ = command.MarkFlagRequired("slug")
	byFlag(command, &by)

	return command
}

// NewForecastsSetCommand sets the token forecasts.
func NewForecastsSetCommand(setter Inputs, globals *support.Globals) *cobra.Command {
	var file, by string

	command := &cobra.Command{
		Use:   "set",
		Short: "Set the tokens a run is expected to use by effort, from a JSON file",
		Args:  validateNoArguments,
		RunE: func(command *cobra.Command, _ []string) error {
			forecasts, err := readInputFile[finops.Forecasts](file)
			if err != nil {
				return err
			}

			return runFinite(command, globals, func(ctx context.Context, project application.ProjectSelection) (inputs.Result, error) {
				return setter.SetForecasts(ctx, inputs.ForecastsRequest{Project: project, By: by, Forecasts: forecasts})
			})
		},
	}

	command.Flags().StringVar(&file, "file", "", "the forecasts: status, note, unit, byEffort, planning")
	_ = command.MarkFlagRequired("file")
	byFlag(command, &by)

	return command
}
