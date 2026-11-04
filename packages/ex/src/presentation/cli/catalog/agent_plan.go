package catalog

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/jonatassales/tao/packages/ex/src/application"
	"github.com/jonatassales/tao/packages/ex/src/application/plan"
	"github.com/jonatassales/tao/packages/ex/src/domain/collection"
	runrecord "github.com/jonatassales/tao/packages/ex/src/domain/run"
	"github.com/jonatassales/tao/packages/ex/src/presentation/cli/session"
	"github.com/jonatassales/tao/packages/ex/src/presentation/cli/support"

	"github.com/spf13/cobra"
)

// Planner plans a ticket as a runtime.
type Planner interface {
	Plan(context.Context, plan.Request) (plan.Result, error)
}

// NewAgentPlanCommand creates the plan command: the first act of the session.
func NewAgentPlanCommand(planner Planner, globals *support.Globals) *cobra.Command {
	var request plan.Request
	var components []string
	var maxBytes int64

	command := &cobra.Command{
		Use:   "plan TICKET",
		Short: "Read a ticket and the files it links, and write the runtime file that plans it",
		Args:  validateSingleArgument,
		RunE: func(command *cobra.Command, arguments []string) error {
			parsed, err := parseComponents(components)
			if err != nil {
				return err
			}

			request.Ticket = arguments[0]
			request.Components = parsed
			request.Limits.MaxBytes = maxBytes
			request.Session = session.For(command.OutOrStdout(), globals.JSON)

			return runSession(command, globals, func(ctx context.Context, project application.ProjectSelection) (plan.Result, error) {
				request.Project = project

				return planner.Plan(ctx, request)
			})
		},
	}

	flags := command.Flags()
	flags.StringVar(&request.Runtime, "runtime", "", "the runtime id (defaults to the ticket's slug)")
	flags.StringVar(&request.Name, "name", "", "the runtime's name (defaults to the ticket's title)")
	flags.StringArrayVar(&components, "component", nil, "a component: name=,path=[,dockerfile=][,port=][,readiness=tcp|http:/path][,env-file=][,pass-env=A;B][,kind=][,tech=id:Label;...]")
	flags.IntVar(&request.CPUMillis, "cpu-millis", 0, "override the default envelope's CPU, with --memory-mib")
	flags.IntVar(&request.MemoryMiB, "memory-mib", 0, "override the default envelope's memory, with --cpu-millis")
	flags.IntVar(&request.Limits.MaxDepth, "max-depth", collection.DefaultLimits().MaxDepth, "how many links deep to follow")
	flags.IntVar(&request.Limits.MaxFiles, "max-files", collection.DefaultLimits().MaxFiles, "how many files to read at most")
	flags.Int64Var(&maxBytes, "max-bytes", collection.DefaultLimits().MaxBytes, "how many bytes to read at most")
	flags.BoolVar(&request.NoJudge, "no-judge", false, "plan without a model: verdicts, facts and the suggestion stay the inputs' own")
	byFlag(command, &request.By)

	return command
}

// parseComponents reads the --component values.
func parseComponents(values []string) ([]plan.ComponentInput, error) {
	components := make([]plan.ComponentInput, 0, len(values))

	for _, value := range values {
		component, err := parseComponent(value)
		if err != nil {
			return nil, err
		}

		components = append(components, component)
	}

	return components, nil
}

func parseComponent(value string) (plan.ComponentInput, error) {
	var component plan.ComponentInput

	for _, pair := range strings.Split(value, ",") {
		key, field, found := strings.Cut(pair, "=")
		if !found {
			return component, support.UsageError{Cause: fmt.Errorf("--component takes key=value pairs, not %q", pair)}
		}

		switch key {
		case "name":
			component.Name = field
		case "path":
			component.Path = field
		case "dockerfile":
			component.Dockerfile = field
		case "kind":
			component.Kind = field
		case "readiness":
			readiness, err := parseReadiness(field)
			if err != nil {
				return component, err
			}
			component.Readiness = readiness
		case "env-file":
			component.EnvFile = field
		case "pass-env":
			component.PassEnv = strings.FieldsFunc(field, isListSeparator)
		case "port":
			port, err := strconv.Atoi(field)
			if err != nil {
				return component, support.UsageError{Cause: fmt.Errorf("--component port %q is not a number", field)}
			}
			component.Port = port
		case "tech":
			component.Technologies = parseTechnologies(field)
		default:
			return component, support.UsageError{Cause: fmt.Errorf("--component has no key %q", key)}
		}
	}

	if component.Name == "" || component.Path == "" {
		return component, support.UsageError{Cause: fmt.Errorf("--component %q needs name= and path=", value)}
	}

	return component, nil
}

func parseTechnologies(field string) []runrecord.Technology {
	var technologies []runrecord.Technology
	for _, entry := range strings.Split(field, ";") {
		id, label, _ := strings.Cut(entry, ":")
		if id == "" {
			continue
		}
		technologies = append(technologies, runrecord.Technology{ID: id, Label: orLabel(label, id)})
	}

	return technologies
}

func orLabel(label, id string) string {
	if label == "" {
		return id
	}

	return label
}

// runSession runs an operation that prints its session as it goes: the
// result is rendered only for scripts, since the lines say it all.
func runSession[Result any](command *cobra.Command, globals *support.Globals, operation projectOperation[Result]) error {
	ctx, cancel := globals.Bounded(command.Context())
	defer cancel()

	selection, err := globals.Selection()
	if err != nil {
		return err
	}

	result, err := operation(ctx, selection)
	if err != nil {
		return err
	}

	if !globals.JSON {
		return nil
	}

	return support.Render(command.OutOrStdout(), result, true)
}
