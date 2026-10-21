package application

import (
	"context"
	"fmt"

	"github.com/exphere/ex/src/domain/microruntime"
	"github.com/exphere/ex/src/domain/project"
)

// InitializeRuntime publishes a new runtime configuration. Without a guide it
// writes an intentionally incomplete descriptor for the consumer to finish.
// With one, it asks the guide for explicit choices once the project and the
// runtime ID are known to be valid, then generates complete native files. It
// never executes anything.
type InitializeRuntime struct {
	Projects   ProjectOpener
	Scaffolder RuntimeScaffolder
}

// BlueprintGuide supplies the choices a runtime is generated from, for
// example by asking the user. It only collects them; the use case parses and
// validates them.
type BlueprintGuide interface {
	Blueprint(ctx context.Context) (BlueprintInput, error)
}

// InitializeRuntimeRequest names the runtime and, optionally, the guide that
// supplies the choices to generate it from.
type InitializeRuntimeRequest struct {
	Project ProjectSelection
	Runtime string
	Guide   BlueprintGuide
}

// BlueprintInput carries guided choices exactly as entered. The use case
// parses them into domain values, so every caller gets the same validation.
type BlueprintInput struct {
	Owner     string
	Context   string
	Namespace string
	Workloads []WorkloadInput
	Quotas    map[string]string
}

// WorkloadInput carries one workload's choices as entered. The workload is
// named after its component directory.
type WorkloadInput struct {
	Component     string
	Kind          string
	Image         string
	CPURequest    string
	CPULimit      string
	MemoryRequest string
	MemoryLimit   string
	Port          int
	Replicas      int
	ReadinessPath string
}

// RuntimeInitResult describes the published configuration.
type RuntimeInitResult struct {
	Runtime string `json:"runtime"`
	Path    string `json:"path"`
	Guided  bool   `json:"guided"`
	Status  string `json:"status"`
}

// Execute publishes the scaffold, or the configuration generated from the
// guide's choices.
func (useCase InitializeRuntime) Execute(ctx context.Context, request InitializeRuntimeRequest) (RuntimeInitResult, error) {
	runtime, err := parseRuntime(request.Runtime)
	if err != nil {
		return RuntimeInitResult{}, err
	}

	owner, err := useCase.Projects.Open(request.Project)
	if err != nil {
		return RuntimeInitResult{}, err
	}

	if request.Guide == nil {
		return useCase.scaffold(owner, runtime)
	}

	return useCase.generate(ctx, owner, runtime, request.Guide)
}

func (useCase InitializeRuntime) scaffold(owner project.Project, runtime microruntime.Identifier) (RuntimeInitResult, error) {
	location, err := useCase.Scaffolder.Scaffold(owner, microruntime.Scaffold(runtime))
	if err != nil {
		return RuntimeInitResult{}, err
	}

	return RuntimeInitResult{
		Runtime: runtime.String(),
		Path:    location,
		Status:  "incomplete; fill in the descriptor and native files, then validate",
	}, nil
}

func (useCase InitializeRuntime) generate(
	ctx context.Context,
	owner project.Project,
	runtime microruntime.Identifier,
	guide BlueprintGuide,
) (RuntimeInitResult, error) {
	input, err := guide.Blueprint(ctx)
	if err != nil {
		return RuntimeInitResult{}, err
	}

	blueprint, err := input.blueprint()
	if err != nil {
		return RuntimeInitResult{}, err
	}

	if err := blueprint.Validate(); err != nil {
		return RuntimeInitResult{}, err
	}

	location, err := useCase.Scaffolder.Generate(owner, runtime, blueprint)
	if err != nil {
		return RuntimeInitResult{}, err
	}

	return RuntimeInitResult{
		Runtime: runtime.String(),
		Path:    location,
		Guided:  true,
		Status:  "generated; infrastructure has not been started",
	}, nil
}

func (input BlueprintInput) blueprint() (microruntime.Blueprint, error) {
	namespace, err := microruntime.ParseIdentifier(input.Namespace)
	if err != nil {
		return microruntime.Blueprint{}, fmt.Errorf("namespace: %w", err)
	}

	blueprint := microruntime.Blueprint{
		Owner:  input.Owner,
		Target: microruntime.KubernetesTarget{Context: input.Context, Namespace: namespace},
		Quotas: make(microruntime.ResourceQuotas, len(input.Quotas)),
	}

	for resource, text := range input.Quotas {
		quantity, err := microruntime.ParseQuantity(text)
		if err != nil {
			return microruntime.Blueprint{}, fmt.Errorf("quota %s: %w", resource, err)
		}

		blueprint.Quotas[microruntime.QuotaResource(resource)] = quantity
	}

	for _, workloadInput := range input.Workloads {
		workload, err := workloadInput.workload()
		if err != nil {
			return microruntime.Blueprint{}, err
		}

		blueprint.Workloads = append(blueprint.Workloads, workload)
	}

	return blueprint, nil
}

func (input WorkloadInput) workload() (microruntime.Workload, error) {
	component, err := microruntime.ParseRelativePath(input.Component)
	if err != nil {
		return microruntime.Workload{}, fmt.Errorf("component: %w", err)
	}

	name, err := microruntime.ParseIdentifier(component.Base())
	if err != nil {
		return microruntime.Workload{}, fmt.Errorf("component %s directory name: %w", component, err)
	}

	resources, err := input.resources()
	if err != nil {
		return microruntime.Workload{}, fmt.Errorf("workload %s: %w", name, err)
	}

	return microruntime.Workload{
		Name:      name,
		Component: component,
		Image:     input.Image,
		Kind:      microruntime.WorkloadKind(input.Kind),
		Resources: resources,
		Service: microruntime.ServiceSpec{
			Port:          input.Port,
			Replicas:      input.Replicas,
			ReadinessPath: input.ReadinessPath,
		},
	}, nil
}

func (input WorkloadInput) resources() (microruntime.Resources, error) {
	var resources microruntime.Resources

	fields := []struct {
		label  string
		text   string
		target *microruntime.Quantity
	}{
		{"CPU request", input.CPURequest, &resources.CPURequest},
		{"CPU limit", input.CPULimit, &resources.CPULimit},
		{"memory request", input.MemoryRequest, &resources.MemoryRequest},
		{"memory limit", input.MemoryLimit, &resources.MemoryLimit},
	}

	for _, field := range fields {
		quantity, err := microruntime.ParseQuantity(field.text)
		if err != nil {
			return microruntime.Resources{}, fmt.Errorf("%s: %w", field.label, err)
		}

		*field.target = quantity
	}

	return resources, nil
}
