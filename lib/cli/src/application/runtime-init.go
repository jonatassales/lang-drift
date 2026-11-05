package application

import (
	"context"
	"fmt"
	"path"

	"github.com/jonatassales/tao/lib/cli/src/domain/manifest"
	"github.com/jonatassales/tao/lib/cli/src/domain/microruntime"
	"github.com/jonatassales/tao/lib/cli/src/domain/project"
	"github.com/jonatassales/tao/lib/core/contracts/go/configuration"
)

// InitializeRuntime publishes a new runtime configuration. Without a guide it
// writes an intentionally incomplete descriptor for the consumer to finish.
// With one, it asks the guide for explicit choices once the project and the
// runtime ID are known to be valid, then generates complete native files. It
// never executes anything.
type InitializeRuntime struct {
	Projects   ProjectOpener
	Scaffolder RuntimeScaffolder
	Files      ComponentFiles
}

// ComponentFiles reads a file a component names, such as its Dockerfile, by
// its project-relative path.
type ComponentFiles interface {
	File(owner project.Project, relativePath string) ([]byte, error)
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
// named after its component directory and builds from the Dockerfile in it.
// Readiness is healthcheck, tcp or http; ReadinessPath serves http only.
type WorkloadInput struct {
	Component     string
	Kind          string
	Dockerfile    string
	CPURequest    string
	CPULimit      string
	MemoryRequest string
	MemoryLimit   string
	Port          int
	Replicas      int
	Readiness     string
	ReadinessPath string
}

// The readiness kinds guided setup offers: the Dockerfile's HEALTHCHECK, or
// an explicit kind.
const (
	ReadinessHealthcheck = "healthcheck"
	ReadinessTCP         = configuration.ReadinessTCP
	ReadinessHTTP        = configuration.ReadinessHTTP
)

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

	blueprint, err := input.blueprint(guidedSource{
		owner:   owner,
		runtime: runtime,
		files:   useCase.Files,
	})
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

// guidedSource is what a workload's choices are read against: the project its
// Dockerfiles live in, and the runtime its images are tagged for.
type guidedSource struct {
	owner   project.Project
	runtime microruntime.Identifier
	files   ComponentFiles
}

func (input BlueprintInput) blueprint(source guidedSource) (microruntime.Blueprint, error) {
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
		workload, err := workloadInput.workload(source)
		if err != nil {
			return microruntime.Blueprint{}, err
		}

		blueprint.Workloads = append(blueprint.Workloads, workload)
	}

	return blueprint, nil
}

func (input WorkloadInput) workload(source guidedSource) (microruntime.Workload, error) {
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

	workload := microruntime.Workload{
		Name:      name,
		Component: component,
		Image:     configuration.ImageTag(source.runtime.String(), name.String(), 1),
		Kind:      microruntime.WorkloadKind(input.Kind),
		Resources: resources,
	}

	if workload.Kind != microruntime.DeploymentWorkload {
		return workload, nil
	}

	probe, err := input.probe(source, component)
	if err != nil {
		return microruntime.Workload{}, fmt.Errorf("workload %s: %w", name, err)
	}

	workload.Service = microruntime.ServiceSpec{
		Port:     input.Port,
		Replicas: input.Replicas,
		Probe:    probe,
	}

	return workload, nil
}

// probe is the readiness the choices name: the HEALTHCHECK of the
// component's Dockerfile, or an explicit tcp or http kind.
func (input WorkloadInput) probe(source guidedSource, component microruntime.RelativePath) (microruntime.Probe, error) {
	switch input.Readiness {
	case ReadinessTCP:
		return microruntime.ExplicitProbe(configuration.Readiness{Kind: configuration.ReadinessTCP}), nil

	case ReadinessHTTP:
		return microruntime.ExplicitProbe(configuration.Readiness{
			Kind: configuration.ReadinessHTTP,
			Path: input.ReadinessPath,
		}), nil

	case ReadinessHealthcheck:
		dockerfile := path.Join(component.String(), input.Dockerfile)

		content, err := source.files.File(source.owner, dockerfile)
		if err != nil {
			return microruntime.Probe{}, fmt.Errorf("read %s: %w", dockerfile, err)
		}

		probe, err := manifest.Healthcheck(content)
		if err != nil {
			return microruntime.Probe{}, fmt.Errorf("%s: %w", dockerfile, err)
		}

		return probe, nil
	}

	return microruntime.Probe{}, fmt.Errorf("readiness %q is not %s, %s or %s", input.Readiness, ReadinessHealthcheck, ReadinessTCP, ReadinessHTTP)
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
