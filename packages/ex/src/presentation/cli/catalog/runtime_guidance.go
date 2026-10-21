package catalog

import (
	"context"
	"fmt"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/domain/microruntime"
	"github.com/exphere/ex/src/presentation/cli/support"
)

// guidedQuotas are the namespace quotas guided setup asks for.
var guidedQuotas = []microruntime.QuotaResource{
	microruntime.QuotaRequestsCPU,
	microruntime.QuotaLimitsCPU,
	microruntime.QuotaRequestsMemory,
	microruntime.QuotaLimitsMemory,
	microruntime.QuotaPods,
}

// promptGuide supplies guided runtime choices by asking the user. It only
// collects answers; the application parses and validates them.
type promptGuide struct {
	prompter *support.Prompter
}

// textQuestion is one required question and where its answer goes.
type textQuestion struct {
	label  string
	answer *string
}

// Blueprint asks every guided question in order.
func (guide promptGuide) Blueprint(ctx context.Context) (application.BlueprintInput, error) {
	var input application.BlueprintInput

	err := guide.askTexts(ctx, []textQuestion{
		{"Owner reference", &input.Owner},
		{"Kubernetes context of the project cluster", &input.Context},
		{"Runtime namespace", &input.Namespace},
	})
	if err != nil {
		return application.BlueprintInput{}, err
	}

	components, err := guide.prompter.List(ctx, "Project-relative component paths (comma-separated)")
	if err != nil {
		return application.BlueprintInput{}, err
	}

	for _, component := range components {
		workload, err := guide.askWorkload(ctx, component)
		if err != nil {
			return application.BlueprintInput{}, err
		}

		input.Workloads = append(input.Workloads, workload)
	}

	input.Quotas = make(map[string]string, len(guidedQuotas))

	for _, resource := range guidedQuotas {
		quantity, err := guide.prompter.Text(ctx, "Namespace quota "+string(resource))
		if err != nil {
			return application.BlueprintInput{}, err
		}

		input.Quotas[string(resource)] = quantity
	}

	return input, nil
}

// askWorkload asks about one component. The kind is checked at once, because
// it decides which questions follow.
func (guide promptGuide) askWorkload(ctx context.Context, component string) (application.WorkloadInput, error) {
	if err := guide.prompter.Say("Configure " + component); err != nil {
		return application.WorkloadInput{}, err
	}

	kind, err := guide.prompter.Text(ctx, "Workload kind (Deployment or Job)")
	if err != nil {
		return application.WorkloadInput{}, err
	}

	if !microruntime.WorkloadKind(kind).IsSupported() {
		return application.WorkloadInput{}, support.UsageError{
			Cause: fmt.Errorf("workload kind must be %s or %s", microruntime.DeploymentWorkload, microruntime.JobWorkload),
		}
	}

	workload := application.WorkloadInput{Component: component, Kind: kind}

	err = guide.askTexts(ctx, []textQuestion{
		{"Prebuilt image reference", &workload.Image},
		{"CPU request (for example 100m)", &workload.CPURequest},
		{"CPU limit", &workload.CPULimit},
		{"Memory request (for example 128Mi)", &workload.MemoryRequest},
		{"Memory limit", &workload.MemoryLimit},
	})
	if err != nil {
		return application.WorkloadInput{}, err
	}

	if microruntime.WorkloadKind(kind) != microruntime.DeploymentWorkload {
		return workload, nil
	}

	return guide.askService(ctx, workload)
}

// askService asks the questions only a Deployment has.
func (guide promptGuide) askService(ctx context.Context, workload application.WorkloadInput) (application.WorkloadInput, error) {
	var err error

	if workload.Port, err = guide.prompter.Integer(ctx, "Application port"); err != nil {
		return application.WorkloadInput{}, err
	}

	if workload.Replicas, err = guide.prompter.Integer(ctx, "Replicas"); err != nil {
		return application.WorkloadInput{}, err
	}

	if workload.ReadinessPath, err = guide.prompter.Text(ctx, "HTTP readiness path"); err != nil {
		return application.WorkloadInput{}, err
	}

	return workload, nil
}

func (guide promptGuide) askTexts(ctx context.Context, questions []textQuestion) error {
	for _, question := range questions {
		answer, err := guide.prompter.Text(ctx, question.label)
		if err != nil {
			return err
		}

		*question.answer = answer
	}

	return nil
}
