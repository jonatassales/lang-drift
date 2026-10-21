package microruntime_test

import (
	"errors"
	"testing"

	"github.com/exphere/ex/src/domain/microruntime"
)

func validBlueprint(t testing.TB) microruntime.Blueprint {
	t.Helper()

	return microruntime.Blueprint{
		Owner:     "team-payments",
		Target:    microruntime.KubernetesTarget{Context: "kind-ex-0123456789ab", Namespace: identifier(t, "demo-local")},
		Workloads: []microruntime.Workload{deploymentWorkload(t), jobWorkload(t)},
		Quotas: microruntime.ResourceQuotas{
			microruntime.QuotaRequestsCPU: quantity(t, "1"),
			microruntime.QuotaPods:        quantity(t, "10"),
		},
	}
}

func TestBlueprint_Validate(t *testing.T) {
	t.Parallel()

	one := quantity(t, "1")

	tests := []struct {
		name    string
		edit    func(*microruntime.Blueprint)
		wantErr bool
	}{
		{name: "valid", edit: func(*microruntime.Blueprint) {}},
		{name: "blank owner", edit: func(b *microruntime.Blueprint) { b.Owner = " " }, wantErr: true},
		{name: "missing context", edit: func(b *microruntime.Blueprint) { b.Target.Context = "" }, wantErr: true},
		{name: "context with a control character", edit: func(b *microruntime.Blueprint) { b.Target.Context = "kind\x1b" }, wantErr: true},
		{name: "missing namespace", edit: func(b *microruntime.Blueprint) { b.Target.Namespace = microruntime.Identifier{} }, wantErr: true},
		{name: "no workloads", edit: func(b *microruntime.Blueprint) { b.Workloads = nil }, wantErr: true},
		{name: "no quotas", edit: func(b *microruntime.Blueprint) { b.Quotas = nil }, wantErr: true},
		{
			name:    "unsupported quota",
			edit:    func(b *microruntime.Blueprint) { b.Quotas["count/secrets"] = one },
			wantErr: true,
		},
		{
			name:    "unset quota quantity",
			edit:    func(b *microruntime.Blueprint) { b.Quotas[microruntime.QuotaLimitsCPU] = microruntime.Quantity{} },
			wantErr: true,
		},
		{
			name:    "invalid workload",
			edit:    func(b *microruntime.Blueprint) { b.Workloads[0].Image = "" },
			wantErr: true,
		},
		{
			name:    "duplicate workload",
			edit:    func(b *microruntime.Blueprint) { b.Workloads = append(b.Workloads, b.Workloads[0]) },
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			blueprint := validBlueprint(t)
			tt.edit(&blueprint)
			err := blueprint.Validate()

			if tt.wantErr != errors.Is(err, microruntime.ErrInvalidBlueprint) || (!tt.wantErr && err != nil) {
				t.Fatalf("Validate() error = %v, want refusal %t", err, tt.wantErr)
			}
		})
	}
}

func TestBlueprint_Descriptor(t *testing.T) {
	t.Parallel()

	blueprint := validBlueprint(t)
	descriptor := blueprint.Descriptor(identifier(t, "demo"))

	if err := descriptor.Validate(); err != nil {
		t.Fatalf("the derived descriptor is invalid: %v", err)
	}

	if descriptor.Version != microruntime.SupportedVersion || descriptor.Environment != microruntime.LocalEnvironment ||
		descriptor.Owner != blueprint.Owner || descriptor.Target != blueprint.Target ||
		descriptor.TerraformDirectory.String() != "terraform" {
		t.Errorf("Descriptor() = %+v, want the blueprint's identity, target, and the default Terraform directory", descriptor)
	}

	for _, workload := range blueprint.Workloads {
		unit, declared := descriptor.Units[workload.Name]
		if !declared {
			t.Fatalf("unit %s is missing", workload.Name)
		}

		wantManifests := "units/" + workload.Name.String() + "/kubernetes"
		if unit.Component != workload.Component || unit.Manifests.String() != wantManifests {
			t.Errorf("unit %s = %+v, want component %s and manifests %s", workload.Name, unit, workload.Component, wantManifests)
		}
	}
}
