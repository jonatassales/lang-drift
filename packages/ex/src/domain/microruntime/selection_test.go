package microruntime_test

import (
	"errors"
	"testing"

	"github.com/exphere/ex/src/domain/microruntime"
)

func TestExecution_SelectResource(t *testing.T) {
	t.Parallel()

	web, batch, idle := identifier(t, "web"), identifier(t, "batch"), identifier(t, "idle")
	execution := &microruntime.Execution{Resources: []microruntime.Resource{
		{Kind: "Deployment", Name: "web", Unit: web},
		{Kind: "Service", Name: "web", Unit: web},
		{Kind: "Job", Name: "batch", Unit: batch},
		{Kind: "Deployment", Name: "first", Unit: idle},
		{Kind: "Deployment", Name: "second", Unit: idle},
	}}

	tests := []struct {
		name     string
		unit     microruntime.Identifier
		workload bool
		want     string
	}{
		{name: "Deployment of a unit", unit: web, workload: true, want: "Deployment/web"},
		{name: "Service of a unit", unit: web, workload: false, want: "Service/web"},
		{name: "Job of a unit", unit: batch, workload: true, want: "Job/batch"},
		{name: "unit without a Service", unit: batch, workload: false},
		{name: "unit with two workloads", unit: idle, workload: true},
		{name: "unknown unit", unit: identifier(t, "unknown"), workload: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := execution.SelectResource(tt.unit, tt.workload)

			if tt.want == "" {
				if !errors.Is(err, microruntime.ErrInconsistentExecution) {
					t.Fatalf("SelectResource() error = %v, want ErrInconsistentExecution", err)
				}

				return
			}

			if err != nil || got.Kind+"/"+got.Name != tt.want {
				t.Fatalf("SelectResource() = %s/%s, %v; want %s", got.Kind, got.Name, err, tt.want)
			}
		})
	}
}

func TestObject_ForwardPort(t *testing.T) {
	t.Parallel()

	service := func(ports ...any) microruntime.Object {
		object := liveObject("v1", "Service", "web", "uid-web", nil)
		object["spec"] = map[string]any{"ports": ports}

		return object
	}
	port := func(number any) map[string]any { return map[string]any{"port": number} }

	tests := []struct {
		name     string
		service  microruntime.Object
		explicit int
		want     int
	}{
		{name: "explicit port", service: service(port(80.0), port(443.0)), explicit: 443, want: 443},
		{name: "the Service's only port", service: service(port(8080.0)), want: 8080},
		{name: "explicit port below range", service: service(port(8080.0)), explicit: -1},
		{name: "explicit port above range", service: service(port(8080.0)), explicit: 65536},
		{name: "ambiguous Service", service: service(port(80.0), port(443.0))},
		{name: "Service without ports", service: service()},
		{name: "fractional port", service: service(port(80.5))},
		{name: "port above range", service: service(port(70000.0))},
		{name: "missing port number", service: service(map[string]any{"name": "http"})},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.service.ForwardPort(tt.explicit)

			if tt.want == 0 {
				if !errors.Is(err, microruntime.ErrUnsupportedManifest) {
					t.Fatalf("ForwardPort(%d) = %d, %v; want ErrUnsupportedManifest", tt.explicit, got, err)
				}

				return
			}

			if err != nil || got != tt.want {
				t.Fatalf("ForwardPort(%d) = %d, %v; want %d", tt.explicit, got, err, tt.want)
			}
		})
	}
}
