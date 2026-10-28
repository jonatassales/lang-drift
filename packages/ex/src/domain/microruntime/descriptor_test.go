package microruntime_test

import (
	"errors"
	"slices"
	"testing"

	"github.com/jonatassales/tao/packages/ex/src/domain/microruntime"
)

// validDescriptor declares a Deployment unit and a Job unit.
func validDescriptor(t testing.TB) microruntime.Descriptor {
	t.Helper()

	return microruntime.Descriptor{
		Version:     microruntime.SupportedVersion,
		ID:          identifier(t, "demo"),
		Environment: microruntime.LocalEnvironment,
		Owner:       "team-payments",
		Target: microruntime.KubernetesTarget{
			Context:   "kind-ex-0123456789ab",
			Namespace: identifier(t, "demo-local"),
		},
		TerraformDirectory: relativePath(t, "terraform"),
		Units: map[microruntime.Identifier]microruntime.Unit{
			identifier(t, "web"):   {Component: relativePath(t, "services/web"), Manifests: relativePath(t, "units/web/kubernetes")},
			identifier(t, "batch"): {Component: relativePath(t, "jobs/batch"), Manifests: relativePath(t, "units/batch/kubernetes")},
		},
		Hooks: map[microruntime.Identifier]microruntime.RelativePath{
			identifier(t, "seed"): relativePath(t, "hooks/seed.sh"),
		},
	}
}

func TestConfigurationDirectory(t *testing.T) {
	t.Parallel()

	if got := microruntime.ConfigurationDirectory(identifier(t, "demo")).String(); got != ".langdrift/microruntimes/runtimes/demo" {
		t.Fatalf("ConfigurationDirectory() = %q, want .langdrift/microruntimes/runtimes/demo", got)
	}
}

func TestScaffold(t *testing.T) {
	t.Parallel()

	scaffold := microruntime.Scaffold(identifier(t, "demo"))

	if scaffold.Version != microruntime.SupportedVersion || scaffold.Environment != microruntime.LocalEnvironment ||
		scaffold.TerraformDirectory.String() != "terraform" || scaffold.Units == nil || len(scaffold.Units) != 0 {
		t.Fatalf("Scaffold() = %+v, want version, environment, Terraform directory, and no units", scaffold)
	}

	if err := scaffold.Validate(); !errors.Is(err, microruntime.ErrInvalidDescriptor) {
		t.Fatalf("Validate() of the incomplete scaffold = %v, want ErrInvalidDescriptor", err)
	}
}

func TestDescriptor_UnitNames(t *testing.T) {
	t.Parallel()

	var names []string
	for _, name := range validDescriptor(t).UnitNames() {
		names = append(names, name.String())
	}

	if !slices.Equal(names, []string{"batch", "web"}) {
		t.Fatalf("UnitNames() = %v, want them sorted", names)
	}
}

func TestDescriptor_Validate(t *testing.T) {
	t.Parallel()

	var (
		web             = identifier(t, "web")
		batch           = identifier(t, "batch")
		apiComponent    = relativePath(t, "services/api")
		configComponent = relativePath(t, ".langdrift/microruntimes/runtimes/demo/web")
		webManifests    = relativePath(t, "units/web/kubernetes")
		seed            = identifier(t, "seed")
		textHook        = relativePath(t, "hooks/seed.txt")
	)

	tests := []struct {
		name    string
		edit    func(*microruntime.Descriptor)
		wantErr bool
	}{
		{name: "valid", edit: func(*microruntime.Descriptor) {}},
		{name: "without hooks", edit: func(d *microruntime.Descriptor) { d.Hooks = nil }},
		{name: "older version", edit: func(d *microruntime.Descriptor) { d.Version = 1 }, wantErr: true},
		{name: "remote environment", edit: func(d *microruntime.Descriptor) { d.Environment = "production" }, wantErr: true},
		{name: "missing ID", edit: func(d *microruntime.Descriptor) { d.ID = microruntime.Identifier{} }, wantErr: true},
		{name: "missing namespace", edit: func(d *microruntime.Descriptor) { d.Target.Namespace = microruntime.Identifier{} }, wantErr: true},
		{name: "blank owner", edit: func(d *microruntime.Descriptor) { d.Owner = "\t" }, wantErr: true},
		{name: "missing context", edit: func(d *microruntime.Descriptor) { d.Target.Context = "" }, wantErr: true},
		{name: "missing Terraform directory", edit: func(d *microruntime.Descriptor) { d.TerraformDirectory = microruntime.RelativePath{} }, wantErr: true},
		{name: "no units", edit: func(d *microruntime.Descriptor) { d.Units = nil }, wantErr: true},
		{
			name: "unit without manifests",
			edit: func(d *microruntime.Descriptor) {
				d.Units[web] = microruntime.Unit{Component: d.Units[web].Component}
			},
			wantErr: true,
		},
		{
			name: "unit named apart from its component",
			edit: func(d *microruntime.Descriptor) {
				d.Units[web] = microruntime.Unit{Component: apiComponent, Manifests: webManifests}
			},
			wantErr: false,
		},
		{
			name: "unit source inside runtime configuration",
			edit: func(d *microruntime.Descriptor) {
				d.Units[web] = microruntime.Unit{Component: configComponent, Manifests: webManifests}
			},
			wantErr: true,
		},
		{
			name: "units sharing manifests",
			edit: func(d *microruntime.Descriptor) {
				d.Units[batch] = microruntime.Unit{Component: d.Units[batch].Component, Manifests: webManifests}
			},
			wantErr: true,
		},
		{
			name:    "hook that is not a shell script",
			edit:    func(d *microruntime.Descriptor) { d.Hooks[seed] = textHook },
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			descriptor := validDescriptor(t)
			tt.edit(&descriptor)
			err := descriptor.Validate()

			if tt.wantErr != errors.Is(err, microruntime.ErrInvalidDescriptor) || (!tt.wantErr && err != nil) {
				t.Fatalf("Validate() error = %v, want refusal %t", err, tt.wantErr)
			}
		})
	}
}
