package application_test

import (
	"bytes"
	"context"
	"errors"
	"slices"
	"testing"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/domain/cluster"
	"github.com/exphere/ex/src/domain/microruntime"
)

func (w *world) readLogs() application.ReadRuntimeLogs {
	return application.ReadRuntimeLogs{
		Projects: w.projects, Locks: w.locks, Configuration: w.configuration, Clusters: w.clusters,
		Observer: w.observer, Executions: w.ledger, Objects: w.kubernetes, Logs: w.logs,
	}
}

func (w *world) logsRequest(unit string, follow bool, output *bytes.Buffer) application.ReadRuntimeLogsRequest {
	return application.ReadRuntimeLogsRequest{
		RuntimeUnitRequest: application.RuntimeUnitRequest{RuntimeRequest: w.runtimeRequest(), Unit: unit},
		Follow:             follow,
		Output:             output,
		Diagnostics:        &bytes.Buffer{},
	}
}

func TestReadRuntimeLogs_StreamsTheRecordedWorkload(t *testing.T) {
	t.Parallel()

	tests := []struct {
		unit   string
		follow bool
		want   string
	}{
		{unit: "web", want: namespaceName + " Deployment/web follow=false"},
		{unit: "batch", follow: true, want: namespaceName + " Job/batch follow=true"},
	}

	for _, tt := range tests {
		t.Run(tt.unit, func(t *testing.T) {
			t.Parallel()

			w := newWorld(t)
			w.up(t)

			var output bytes.Buffer
			if err := w.readLogs().Execute(context.Background(), w.logsRequest(tt.unit, tt.follow, &output)); err != nil {
				t.Fatalf("Execute() error = %v", err)
			}

			if !slices.Equal(w.logs.requests, []string{tt.want}) || output.String() != "log line\n" {
				t.Fatalf("requests %v, output %q; want %q and the streamed output", w.logs.requests, output.String(), tt.want)
			}

			if w.locks.locked != 2 || w.locks.unlocked != 2 {
				t.Errorf("the project lock was not held for the read")
			}
		})
	}
}

// TestReadRuntimeLogs_RefusesUnverifiedTargets also covers target selection
// shared with forwarding.
func TestReadRuntimeLogs_RefusesUnverifiedTargets(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		unit    string
		arrange func(*testing.T, *world)
		wantErr error
	}{
		{name: "invalid unit", unit: "Web", arrange: func(t *testing.T, w *world) { w.up(t) }, wantErr: microruntime.ErrInvalidIdentifier},
		{name: "unit without a workload", unit: "absent", arrange: func(t *testing.T, w *world) { w.up(t) }, wantErr: microruntime.ErrInconsistentExecution},
		{name: "runtime never started", unit: "web", arrange: func(*testing.T, *world) {}, wantErr: application.ErrNotFound},
		{
			name:    "broken references",
			unit:    "web",
			arrange: func(t *testing.T, w *world) { w.up(t); w.configuration.referencesErr = errInjected },
			wantErr: errInjected,
		},
		{
			name:    "descriptor for another runtime",
			unit:    "web",
			arrange: func(t *testing.T, w *world) { w.up(t); w.configuration.descriptor.ID = mustIdentifier(t, "other") },
			wantErr: microruntime.ErrInvalidDescriptor,
		},
		{
			name:    "cluster not ready",
			unit:    "web",
			arrange: func(t *testing.T, w *world) { w.up(t); w.clusters.record.Status = cluster.StatusDeleting },
			wantErr: cluster.ErrNotReady,
		},
		{
			name: "configuration targets another namespace",
			unit: "web",
			arrange: func(t *testing.T, w *world) {
				w.up(t)
				w.configuration.descriptor.Target.Namespace = mustIdentifier(t, "elsewhere")
			},
			wantErr: application.ErrTargetMismatch,
		},
		{
			name:    "workload removed",
			unit:    "web",
			arrange: func(t *testing.T, w *world) { w.up(t); w.kubernetes.objects = nil },
			wantErr: microruntime.ErrOwnershipChanged,
		},
		{
			name: "workload replaced",
			unit: "web",
			arrange: func(t *testing.T, w *world) {
				w.up(t)
				w.kubernetes.edit("Deployment", "web", func(object microruntime.Object) { object.Metadata()["uid"] = "uid-replacement" })
			},
			wantErr: microruntime.ErrOwnershipChanged,
		},
		{name: "cluster unreachable", unit: "web", arrange: func(t *testing.T, w *world) { w.up(t); w.kubernetes.getErr = errInjected }, wantErr: errInjected},
		{name: "lock held elsewhere", unit: "web", arrange: func(t *testing.T, w *world) { w.up(t); w.locks.err = errInjected }, wantErr: errInjected},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := newWorld(t)
			tt.arrange(t, w)

			err := w.readLogs().Execute(context.Background(), w.logsRequest(tt.unit, false, &bytes.Buffer{}))
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Execute() error = %v, want %v", err, tt.wantErr)
			}

			if len(w.logs.requests) != 0 {
				t.Error("logs were read from an unverified target")
			}

			w.requireLockBalanced(t)
		})
	}
}

func TestReadRuntimeLogs_ReportsTheStreamFailure(t *testing.T) {
	t.Parallel()

	w := newWorld(t)
	w.up(t)
	w.logs.err = errInjected

	if err := w.readLogs().Execute(context.Background(), w.logsRequest("web", false, &bytes.Buffer{})); !errors.Is(err, errInjected) {
		t.Fatalf("Execute() error = %v, want the stream failure", err)
	}
}
