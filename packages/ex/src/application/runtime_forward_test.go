package application_test

import (
	"bytes"
	"context"
	"errors"
	"slices"
	"testing"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/domain/microruntime"
)

func (w *world) forwardRuntime() application.ForwardRuntime {
	return application.ForwardRuntime{
		Projects: w.projects, Locks: w.locks, Configuration: w.configuration, Clusters: w.clusters,
		Observer: w.observer, Executions: w.ledger, Objects: w.kubernetes, Forwarder: w.forwarder,
	}
}

func (w *world) forwardRequest(unit string, local, remote int) application.ForwardRuntimeRequest {
	return application.ForwardRuntimeRequest{
		RuntimeUnitRequest: application.RuntimeUnitRequest{RuntimeRequest: w.runtimeRequest(), Unit: unit},
		LocalPort:          local,
		RemotePort:         remote,
		Diagnostics:        &bytes.Buffer{},
	}
}

func TestForwardRuntime_ForwardsTheRecordedService(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		remote int
		want   string
	}{
		{name: "the Service's only port", want: namespaceName + " Service/web 9090:8080"},
		{name: "an explicit port", remote: 8081, want: namespaceName + " Service/web 9090:8081"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := newWorld(t)
			w.up(t)

			if err := w.forwardRuntime().Execute(context.Background(), w.forwardRequest("web", 9090, tt.remote)); err != nil {
				t.Fatalf("Execute() error = %v", err)
			}

			if !slices.Equal(w.forwarder.requests, []string{tt.want}) {
				t.Fatalf("forwarded %v, want %q", w.forwarder.requests, tt.want)
			}

			w.requireLockBalanced(t)
		})
	}
}

func TestForwardRuntime_Refuses(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		unit    string
		local   int
		remote  int
		arrange func(*world)
		wantErr error
	}{
		{name: "local port zero", unit: "web", local: 0, arrange: func(*world) {}},
		{name: "local port above range", unit: "web", local: 65536, arrange: func(*world) {}},
		{name: "negative remote port", unit: "web", local: 9090, remote: -1, arrange: func(*world) {}},
		{name: "invalid unit", unit: "-web", local: 9090, arrange: func(*world) {}, wantErr: microruntime.ErrInvalidIdentifier},
		{name: "unit without a Service", unit: "batch", local: 9090, arrange: func(*world) {}, wantErr: microruntime.ErrInconsistentExecution},
		{
			name:  "ambiguous Service",
			unit:  "web",
			local: 9090,
			arrange: func(w *world) {
				w.kubernetes.edit("Service", "web", func(object microruntime.Object) {
					object.Spec()["ports"] = []any{map[string]any{"port": 80.0}, map[string]any{"port": 443.0}}
				})
			},
			wantErr: microruntime.ErrUnsupportedManifest,
		},
		{name: "unknown project", unit: "web", local: 9090, arrange: func(w *world) { w.projects.err = errInjected }, wantErr: errInjected},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := newWorld(t)
			w.up(t)
			tt.arrange(w)

			err := w.forwardRuntime().Execute(context.Background(), w.forwardRequest(tt.unit, tt.local, tt.remote))
			if err == nil || (tt.wantErr != nil && !errors.Is(err, tt.wantErr)) {
				t.Fatalf("Execute() error = %v, want %v", err, tt.wantErr)
			}

			if len(w.forwarder.requests) != 0 {
				t.Error("a refused forward started a tunnel")
			}

			w.requireLockBalanced(t)
		})
	}
}

func TestForwardRuntime_RejectsAnInvalidRuntimeID(t *testing.T) {
	t.Parallel()

	w := newWorld(t)
	request := w.forwardRequest("web", 9090, 0)
	request.Runtime = "Demo"

	if err := w.forwardRuntime().Execute(context.Background(), request); !errors.Is(err, microruntime.ErrInvalidIdentifier) {
		t.Fatalf("Execute() error = %v, want ErrInvalidIdentifier", err)
	}
}
