package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/domain/microruntime"
	"github.com/exphere/ex/src/domain/project"
)

type fakeScaffolder struct {
	err        error
	scaffolded []microruntime.Descriptor
	generated  []microruntime.Blueprint
}

func (scaffolder *fakeScaffolder) Scaffold(_ project.Project, descriptor microruntime.Descriptor) (string, error) {
	scaffolder.scaffolded = append(scaffolder.scaffolded, descriptor)

	return projectRoot + "/.exphere/runtimes/" + descriptor.ID.String(), scaffolder.err
}

func (scaffolder *fakeScaffolder) Generate(_ project.Project, runtime microruntime.Identifier, blueprint microruntime.Blueprint) (string, error) {
	scaffolder.generated = append(scaffolder.generated, blueprint)

	return projectRoot + "/.exphere/runtimes/" + runtime.String(), scaffolder.err
}

type fakeGuide struct {
	input application.BlueprintInput
	err   error
}

func (guide fakeGuide) Blueprint(context.Context) (application.BlueprintInput, error) {
	return guide.input, guide.err
}

func guidedInput() application.BlueprintInput {
	return application.BlueprintInput{
		Owner:     "team-payments",
		Context:   "kind-ex-0123456789ab",
		Namespace: namespaceName,
		Workloads: []application.WorkloadInput{
			{
				Component: "services/web", Kind: "Deployment", Image: "web:dev",
				CPURequest: "100m", CPULimit: "500m", MemoryRequest: "128Mi", MemoryLimit: "256Mi",
				Port: 8080, Replicas: 1, ReadinessPath: "/healthz",
			},
			{
				Component: "jobs/batch", Kind: "Job", Image: "batch:dev",
				CPURequest: "100m", CPULimit: "200m", MemoryRequest: "64Mi", MemoryLimit: "128Mi",
			},
		},
		Quotas: map[string]string{"requests.cpu": "1", "limits.memory": "1Gi", "pods": "10"},
	}
}

func TestInitializeRuntime_ScaffoldsAnIncompleteDescriptor(t *testing.T) {
	t.Parallel()

	w := newWorld(t)
	scaffolder := &fakeScaffolder{}

	result, err := application.InitializeRuntime{Projects: w.projects, Scaffolder: scaffolder}.Execute(
		context.Background(), application.InitializeRuntimeRequest{Project: w.runtimeRequest().Project, Runtime: "demo"})
	if err != nil {
		t.Fatal(err)
	}

	if result.Guided || result.Runtime != "demo" || result.Path != projectRoot+"/.exphere/runtimes/demo" {
		t.Errorf("result = %+v, want the unguided scaffold location", result)
	}

	if len(scaffolder.scaffolded) != 1 || scaffolder.scaffolded[0].ID.String() != "demo" || len(scaffolder.generated) != 0 {
		t.Errorf("scaffolded %v and generated %v, want one scaffold", scaffolder.scaffolded, scaffolder.generated)
	}
}

func TestInitializeRuntime_GeneratesFromGuidedChoices(t *testing.T) {
	t.Parallel()

	w := newWorld(t)
	scaffolder := &fakeScaffolder{}

	result, err := application.InitializeRuntime{Projects: w.projects, Scaffolder: scaffolder}.Execute(
		context.Background(), application.InitializeRuntimeRequest{Project: w.runtimeRequest().Project, Runtime: "demo", Guide: fakeGuide{input: guidedInput()}})
	if err != nil {
		t.Fatal(err)
	}

	if !result.Guided || len(scaffolder.generated) != 1 {
		t.Fatalf("result = %+v with %d generations, want one guided generation", result, len(scaffolder.generated))
	}

	blueprint := scaffolder.generated[0]
	if blueprint.Owner != "team-payments" || blueprint.Target.Namespace.String() != namespaceName || len(blueprint.Quotas) != 3 {
		t.Errorf("blueprint = %+v, want the parsed guided choices", blueprint)
	}

	web, batch := blueprint.Workloads[0], blueprint.Workloads[1]
	if web.Name.String() != "web" || web.Kind != microruntime.DeploymentWorkload || web.Service.Port != 8080 ||
		web.Resources.MemoryLimit.String() != "256Mi" || batch.Name.String() != "batch" || batch.Kind != microruntime.JobWorkload {
		t.Errorf("workloads = %+v, want names derived from component directories", blueprint.Workloads)
	}
}

func TestInitializeRuntime_RejectsInvalidChoicesBeforeWriting(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		edit    func(*application.BlueprintInput)
		wantErr error
	}{
		{name: "invalid namespace", edit: func(input *application.BlueprintInput) { input.Namespace = "Demo" }, wantErr: microruntime.ErrInvalidIdentifier},
		{name: "invalid quota quantity", edit: func(input *application.BlueprintInput) { input.Quotas["pods"] = "ten" }, wantErr: microruntime.ErrInvalidQuantity},
		{name: "unsupported quota", edit: func(input *application.BlueprintInput) { input.Quotas["count/secrets"] = "1" }, wantErr: microruntime.ErrInvalidBlueprint},
		{
			name:    "component outside the project",
			edit:    func(input *application.BlueprintInput) { input.Workloads[0].Component = "../web" },
			wantErr: microruntime.ErrInvalidPath,
		},
		{
			name:    "component named outside the identifier rules",
			edit:    func(input *application.BlueprintInput) { input.Workloads[0].Component = "services/Web_API" },
			wantErr: microruntime.ErrInvalidIdentifier,
		},
		{
			name:    "invalid CPU limit",
			edit:    func(input *application.BlueprintInput) { input.Workloads[1].CPULimit = "lots" },
			wantErr: microruntime.ErrInvalidQuantity,
		},
		{
			name:    "Deployment without replicas",
			edit:    func(input *application.BlueprintInput) { input.Workloads[0].Replicas = 0 },
			wantErr: microruntime.ErrInvalidBlueprint,
		},
		{name: "no workloads", edit: func(input *application.BlueprintInput) { input.Workloads = nil }, wantErr: microruntime.ErrInvalidBlueprint},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := newWorld(t)
			scaffolder := &fakeScaffolder{}
			input := guidedInput()
			tt.edit(&input)

			_, err := application.InitializeRuntime{Projects: w.projects, Scaffolder: scaffolder}.Execute(
				context.Background(), application.InitializeRuntimeRequest{Project: w.runtimeRequest().Project, Runtime: "demo", Guide: fakeGuide{input: input}})
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Execute() error = %v, want %v", err, tt.wantErr)
			}

			if len(scaffolder.generated) != 0 {
				t.Error("invalid choices were generated")
			}
		})
	}
}

func TestInitializeRuntime_Refuses(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		runtime    string
		guide      application.BlueprintGuide
		projects   error
		scaffolder error
		wantErr    error
	}{
		{name: "invalid runtime ID", runtime: "demo!", wantErr: microruntime.ErrInvalidIdentifier},
		{name: "unknown project", runtime: "demo", projects: errInjected, wantErr: errInjected},
		{name: "guide canceled", runtime: "demo", guide: fakeGuide{err: context.Canceled}, wantErr: context.Canceled},
		{name: "existing destination", runtime: "demo", scaffolder: errInjected, wantErr: errInjected},
		{name: "existing guided destination", runtime: "demo", guide: fakeGuide{input: guidedInput()}, scaffolder: errInjected, wantErr: errInjected},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := newWorld(t)
			w.projects.err = tt.projects

			_, err := application.InitializeRuntime{Projects: w.projects, Scaffolder: &fakeScaffolder{err: tt.scaffolder}}.Execute(
				context.Background(), application.InitializeRuntimeRequest{Project: w.runtimeRequest().Project, Runtime: tt.runtime, Guide: tt.guide})
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Execute() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
