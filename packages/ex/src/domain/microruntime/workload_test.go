package microruntime_test

import (
	"errors"
	"testing"

	"github.com/jonatassales/tao/packages/ex/src/domain/microruntime"
)

func deploymentWorkload(t testing.TB) microruntime.Workload {
	t.Helper()

	return microruntime.Workload{
		Name:      identifier(t, "web"),
		Component: relativePath(t, "services/web"),
		Image:     "web:dev",
		Kind:      microruntime.DeploymentWorkload,
		Resources: microruntime.Resources{
			CPURequest:    quantity(t, "100m"),
			CPULimit:      quantity(t, "500m"),
			MemoryRequest: quantity(t, "128Mi"),
			MemoryLimit:   quantity(t, "256Mi"),
		},
		Service: microruntime.ServiceSpec{Port: 8080, Replicas: 2, ReadinessPath: "/healthz"},
	}
}

func jobWorkload(t testing.TB) microruntime.Workload {
	t.Helper()

	workload := deploymentWorkload(t)
	workload.Name = identifier(t, "batch")
	workload.Component = relativePath(t, "jobs/batch")
	workload.Image = "batch:dev"
	workload.Kind = microruntime.JobWorkload
	workload.Service = microruntime.ServiceSpec{}

	return workload
}

func TestWorkload_Validate(t *testing.T) {
	t.Parallel()

	otherComponent := relativePath(t, "services/api")
	insideConfiguration := relativePath(t, ".langdrift/microruntimes/runtimes/web")

	tests := []struct {
		name    string
		job     bool
		edit    func(*microruntime.Workload)
		wantErr bool
	}{
		{name: "deployment", edit: func(*microruntime.Workload) {}},
		{name: "job", job: true, edit: func(*microruntime.Workload) {}},
		{name: "missing name", edit: func(w *microruntime.Workload) { w.Name = microruntime.Identifier{} }, wantErr: true},
		{name: "missing component", edit: func(w *microruntime.Workload) { w.Component = microruntime.RelativePath{} }, wantErr: true},
		{
			name:    "name differs from the component directory",
			edit:    func(w *microruntime.Workload) { w.Component = otherComponent },
			wantErr: false,
		},
		{
			name:    "source inside runtime configuration",
			edit:    func(w *microruntime.Workload) { w.Component = insideConfiguration },
			wantErr: true,
		},
		{name: "missing image", edit: func(w *microruntime.Workload) { w.Image = "" }, wantErr: true},
		{name: "image with whitespace", edit: func(w *microruntime.Workload) { w.Image = "web: dev" }, wantErr: true},
		{name: "image with a newline", edit: func(w *microruntime.Workload) { w.Image = "web:dev\n--privileged" }, wantErr: true},
		{
			name:    "missing memory limit",
			edit:    func(w *microruntime.Workload) { w.Resources.MemoryLimit = microruntime.Quantity{} },
			wantErr: true,
		},
		{name: "port zero", edit: func(w *microruntime.Workload) { w.Service.Port = 0 }, wantErr: true},
		{name: "port above range", edit: func(w *microruntime.Workload) { w.Service.Port = 65536 }, wantErr: true},
		{name: "no replicas", edit: func(w *microruntime.Workload) { w.Service.Replicas = 0 }, wantErr: true},
		{name: "readiness path without slash", edit: func(w *microruntime.Workload) { w.Service.ReadinessPath = "healthz" }, wantErr: true},
		{
			name:    "readiness path with a host",
			edit:    func(w *microruntime.Workload) { w.Service.ReadinessPath = "http://example.com/healthz" },
			wantErr: true,
		},
		{name: "empty readiness path", edit: func(w *microruntime.Workload) { w.Service.ReadinessPath = "" }, wantErr: true},
		{
			name:    "job with a service port",
			job:     true,
			edit:    func(w *microruntime.Workload) { w.Service.Port = 8080 },
			wantErr: true,
		},
		{name: "unsupported kind", edit: func(w *microruntime.Workload) { w.Kind = "StatefulSet" }, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			workload := deploymentWorkload(t)
			if tt.job {
				workload = jobWorkload(t)
			}

			tt.edit(&workload)
			err := workload.Validate()

			if tt.wantErr != errors.Is(err, microruntime.ErrInvalidBlueprint) || (!tt.wantErr && err != nil) {
				t.Fatalf("Validate() error = %v, want refusal %t", err, tt.wantErr)
			}
		})
	}
}

func TestWorkloadKind_IsSupported(t *testing.T) {
	t.Parallel()

	for kind, want := range map[microruntime.WorkloadKind]bool{
		microruntime.DeploymentWorkload: true,
		microruntime.JobWorkload:        true,
		"deployment":                    false,
		"CronJob":                       false,
		"":                              false,
	} {
		if got := kind.IsSupported(); got != want {
			t.Errorf("%q.IsSupported() = %t, want %t", kind, got, want)
		}
	}
}

func TestIsImageReference(t *testing.T) {
	t.Parallel()

	for text, want := range map[string]bool{
		"web:dev":                        true,
		"registry.local:5000/team/web:1": true,
		"web@sha256:0123":                true,
		"":                               false,
		"   ":                            false,
		"web dev":                        false,
		"web\tdev":                       false,
		"web:dev\n":                      false,
	} {
		if got := microruntime.IsImageReference(text); got != want {
			t.Errorf("IsImageReference(%q) = %t, want %t", text, got, want)
		}
	}
}
