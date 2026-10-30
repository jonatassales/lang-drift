package setup

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jonatassales/tao/packages/ex/src/application"
	"github.com/jonatassales/tao/packages/ex/src/application/inputs"
	"github.com/jonatassales/tao/packages/ex/src/domain/finops"
	"github.com/jonatassales/tao/packages/ex/src/domain/project"
)

// Initializer records the project identity.
type Initializer interface {
	Execute(context.Context, application.InitializeProjectRequest) (application.ProjectResult, error)
}

// Recorder records an input as its own command does.
type Recorder interface {
	SetBudget(context.Context, inputs.BudgetRequest) (inputs.Result, error)
	SetGates(context.Context, inputs.GatesRequest) (inputs.Result, error)
}

// ArtifactsFolder keeps the folder of the team's tickets and specs.
type ArtifactsFolder interface {
	Folder(owner project.Project) (string, error)
	Check(owner project.Project, relative string) (string, error)
	Record(owner project.Project, relative string) (string, error)
}

// Request names the project directory, who sets the inputs, the flags'
// answers and, when a person is there, the guide that asks for the rest.
type Request struct {
	Directory string
	By        string
	Answers   Answers
	Guide     Guide
}

// Result is the project with what the setup recorded, what it kept because
// it was already recorded, and the commands still to run before a first run.
type Result struct {
	application.ProjectResult
	Recorded []string `json:"recorded,omitempty"`
	Kept     []string `json:"kept,omitempty"`
	Next     []string `json:"next,omitempty"`
}

// The names the result gives the inputs that are not money rules.
const (
	inputArtifacts = "artifacts"
)

// The inputs a first run still needs after the setup, with their commands.
var nextInputs = []struct {
	kind    string
	command string
}{
	{
		kind:    inputs.KindEnvelope,
		command: "ex harness envelope set",
	},
	{
		kind:    inputs.KindChoice,
		command: "ex harness model set",
	},
	{
		kind:    inputs.KindDriver,
		command: "ex harness driver set",
	},
}

// SetUpProject initializes the project, then records the answers it was
// given or asked for.
type SetUpProject struct {
	Initialize Initializer
	Projects   application.ProjectOpener
	Inputs     Recorder
	Store      inputs.Store
	Artifacts  ArtifactsFolder
	Clock      application.Clock
}

// Execute records the identity, asks for what is pending, and records each
// answer that is not recorded yet.
func (useCase SetUpProject) Execute(ctx context.Context, request Request) (Result, error) {
	initialized, err := useCase.Initialize.Execute(ctx, application.InitializeProjectRequest{Directory: request.Directory})
	if err != nil {
		return Result{}, err
	}

	selection := application.ProjectSelection{
		Directory: initialized.Root,
		Explicit:  true,
	}

	owner, err := useCase.Projects.Open(selection)
	if err != nil {
		return Result{}, err
	}

	pending, err := useCase.pending(owner)
	if err != nil {
		return Result{}, err
	}

	answers := request.Answers

	if request.Guide != nil {
		answers, err = request.Guide.Ask(ctx, pending, answers)
		if err != nil {
			return Result{}, err
		}
	}

	result := Result{ProjectResult: initialized}

	if err := useCase.checkArtifacts(owner, pending, answers); err != nil {
		return Result{}, err
	}

	if err := useCase.recordFinOps(ctx, selection, request.By, pending, answers, &result); err != nil {
		return Result{}, err
	}

	if err := useCase.recordArtifacts(owner, pending, answers, &result); err != nil {
		return Result{}, err
	}

	nothingAskedOrAnswered := request.Guide == nil && answers == (Answers{})
	if nothingAskedOrAnswered {
		return result, nil
	}

	result.Next, err = useCase.next(owner)

	return result, err
}

// pending reads which inputs are not recorded yet: a budget for this month,
// and the gates.
func (useCase SetUpProject) pending(owner project.Project) (Pending, error) {
	pending := Pending{Period: useCase.Clock.Now().UTC().Format("2006-01")}

	budget, err := useCase.Store.Load(owner, inputs.KindBudget)

	switch {
	case errors.Is(err, application.ErrNotFound):
		pending.Budget = true

	case err != nil:
		return Pending{}, err

	default:
		var recorded finops.Budget
		if err := json.Unmarshal(budget.Value, &recorded); err != nil {
			return Pending{}, err
		}

		pending.Budget = recorded.Period != pending.Period
	}

	_, err = useCase.Store.Load(owner, inputs.KindGates)

	switch {
	case errors.Is(err, application.ErrNotFound):
		pending.Gates = true

	case err != nil:
		return Pending{}, err
	}

	folder, err := useCase.Artifacts.Folder(owner)
	if err != nil {
		return Pending{}, err
	}

	pending.Artifacts = folder == ""

	return pending, nil
}

// checkArtifacts refuses a pending artifacts folder the project does not
// have, before any input is recorded.
func (useCase SetUpProject) checkArtifacts(owner project.Project, pending Pending, answers Answers) error {
	if answers.Artifacts == "" || !pending.Artifacts {
		return nil
	}

	_, err := useCase.Artifacts.Check(owner, answers.Artifacts)

	return err
}

// recordArtifacts stores the artifacts folder, once the budget and the gates
// the service validates are recorded.
func (useCase SetUpProject) recordArtifacts(owner project.Project, pending Pending, answers Answers, result *Result) error {
	if answers.Artifacts == "" {
		return nil
	}

	if !pending.Artifacts {
		result.Kept = append(result.Kept, inputArtifacts)

		return nil
	}

	if _, err := useCase.Artifacts.Record(owner, answers.Artifacts); err != nil {
		return err
	}

	result.Recorded = append(result.Recorded, inputArtifacts)

	return nil
}

// recordFinOps sets the budget and the gates when they are pending, and keeps
// an answer whose input is already recorded.
func (useCase SetUpProject) recordFinOps(ctx context.Context, selection application.ProjectSelection, by string, pending Pending, answers Answers, result *Result) error {
	if answers.Budget != nil && !pending.Budget {
		result.Kept = append(result.Kept, inputs.KindBudget)
	}

	if answers.Budget != nil && pending.Budget {
		_, err := useCase.Inputs.SetBudget(ctx, inputs.BudgetRequest{
			Project: selection,
			By:      by,
			Period:  pending.Period,
			Bounded: *answers.Budget,
		})
		if err != nil {
			return err
		}

		result.Recorded = append(result.Recorded, inputs.KindBudget)
	}

	if answers.DefaultGates && !pending.Gates {
		result.Kept = append(result.Kept, inputs.KindGates)
	}

	if answers.DefaultGates && pending.Gates {
		_, err := useCase.Inputs.SetGates(ctx, inputs.GatesRequest{
			Project: selection,
			By:      by,
			Gates:   finops.DefaultGates(),
		})
		if err != nil {
			return err
		}

		result.Recorded = append(result.Recorded, inputs.KindGates)
	}

	return nil
}

// next lists the commands of the inputs a first run needs that nobody set.
func (useCase SetUpProject) next(owner project.Project) ([]string, error) {
	var commands []string

	for _, input := range nextInputs {
		_, err := useCase.Store.Load(owner, input.kind)

		switch {
		case errors.Is(err, application.ErrNotFound):
			commands = append(commands, input.command)

		case err != nil:
			return nil, err
		}
	}

	return commands, nil
}
