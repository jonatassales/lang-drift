package application_test

import (
	"context"
	"errors"
	"slices"
	"testing"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/domain/cluster"
	"github.com/exphere/ex/src/domain/microruntime"
)

func (w *world) down(recoverPartial bool) (application.ExecutionOutcome, error) {
	return w.stopRuntime().Execute(context.Background(), application.StopRuntimeRequest{
		RuntimeRequest: w.runtimeRequest(),
		Recover:        recoverPartial,
	})
}

func TestStopRuntime_RemovesOnlyWhatTheRuntimeOwns(t *testing.T) {
	t.Parallel()

	w := newWorld(t)
	up := w.up(t)
	token := w.ledger.last(t, "demo").Token

	outcome, err := w.down(false)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if outcome.Status != "down" || !outcome.Settled || outcome.Operation != "down" || outcome.RunID == up.RunID ||
		outcome.NamespaceUID != "" || len(outcome.Resources) != 0 {
		t.Fatalf("outcome = %+v, want a settled down without obligations and with its own run ID", outcome)
	}

	if outcome.Units["web"] != "completed" || outcome.Units["batch"] != "completed" {
		t.Errorf("units = %v, want every unit completed", outcome.Units)
	}

	wantDeletes := []string{"Deployment/web", "Service/web", "Job/batch", "Namespace/" + namespaceName}
	if !slices.Equal(w.kubernetes.deletes, wantDeletes) {
		t.Errorf("deletes = %v, want the recorded objects in order, then the namespace", w.kubernetes.deletes)
	}

	if len(w.preparation.destroyed) != 1 || w.preparation.destroyed[0].Token != token {
		t.Errorf("preparation destroyed %+v, want one destroy with the execution's token", w.preparation.destroyed)
	}

	if inventory, _ := w.kubernetes.NamespaceInventory(context.Background(), application.Target{}); len(inventory) != 0 {
		t.Errorf("objects remain after cleanup: %v", inventory)
	}

	if err := w.ledger.last(t, "demo").RequireClosed(); err != nil {
		t.Errorf("the ledger still records obligations: %v", err)
	}

	if len(w.evidence.executions) != 2 || w.evidence.executions[1].RunID != outcome.RunID {
		t.Errorf("evidence = %+v, want the up and the down", w.evidence.executions)
	}

	w.requireLockBalanced(t)
}

func TestStopRuntime_ReturnsAStoppedRuntimeUnchanged(t *testing.T) {
	t.Parallel()

	w := newWorld(t)
	w.up(t)

	first, err := w.down(false)
	if err != nil {
		t.Fatal(err)
	}

	saves, deletes, evidence := len(w.ledger.saved), len(w.kubernetes.deletes), len(w.evidence.executions)

	second, err := w.down(false)
	if err != nil || second.RunID != first.RunID || second.Status != "down" {
		t.Fatalf("second Execute() = %+v, %v; want the recorded down", second, err)
	}

	if len(w.ledger.saved) != saves || len(w.kubernetes.deletes) != deletes || len(w.evidence.executions) != evidence {
		t.Fatal("a repeated down changed the ledger, the cluster, or the evidence")
	}
}

func TestStopRuntime_RefusesBeforeRecording(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		runtime string
		arrange func(*testing.T, *world)
		recover bool
		wantErr error
	}{
		{name: "invalid runtime ID", runtime: "demo_", arrange: func(*testing.T, *world) {}, wantErr: microruntime.ErrInvalidIdentifier},
		{name: "runtime never started", arrange: func(*testing.T, *world) {}, wantErr: application.ErrNotFound},
		{name: "operation already in progress", arrange: func(t *testing.T, w *world) { w.up(t); w.locks.err = errInjected }, wantErr: errInjected},
		{
			name:    "previous operation never settled",
			arrange: func(t *testing.T, w *world) { w.up(t); w.ledger.last(t, "demo").Settled = false },
			recover: true,
			wantErr: microruntime.ErrCleanupBlocked,
		},
		{
			name:    "failed up without explicit recovery",
			arrange: func(t *testing.T, w *world) { w.up(t); w.ledger.last(t, "demo").Status = microruntime.StatusFailed },
			wantErr: microruntime.ErrCleanupBlocked,
		},
		{
			name:    "inconsistent record",
			arrange: func(t *testing.T, w *world) { w.up(t); w.ledger.last(t, "demo").Resources[0].Kind = "Secret" },
			wantErr: microruntime.ErrInconsistentExecution,
		},
		{
			name:    "configuration changed since the up",
			arrange: func(t *testing.T, w *world) { w.up(t); w.configuration.hash = "hash-edited" },
			wantErr: application.ErrConfigurationChanged,
		},
		{
			name: "configuration removed",
			arrange: func(t *testing.T, w *world) {
				w.up(t)
				w.configuration.loadErr = application.NotFound("runtime demo is not initialized")
			},
			wantErr: application.ErrNotFound,
		},
		{
			name:    "replaced cluster",
			arrange: func(t *testing.T, w *world) { w.up(t); w.observer.namespaces["kube-system"] = "uid-replacement" },
			wantErr: cluster.ErrUntrustedIdentity,
		},
		{
			name:    "execution on another cluster",
			arrange: func(t *testing.T, w *world) { w.up(t); w.ledger.last(t, "demo").ClusterUID = "uid-other-cluster" },
			wantErr: application.ErrTargetMismatch,
		},
		{
			name: "configuration now targets another namespace",
			arrange: func(t *testing.T, w *world) {
				w.up(t)
				w.configuration.descriptor.Target.Namespace = mustIdentifier(t, "elsewhere")
			},
			wantErr: application.ErrTargetMismatch,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := newWorld(t)
			tt.arrange(t, w)

			saves, deletes, destroyed := len(w.ledger.saved), len(w.kubernetes.deletes), len(w.preparation.destroyed)

			request := application.StopRuntimeRequest{RuntimeRequest: w.runtimeRequest(), Recover: tt.recover}
			if tt.runtime != "" {
				request.Runtime = tt.runtime
			}

			_, err := w.stopRuntime().Execute(context.Background(), request)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Execute() error = %v, want %v", err, tt.wantErr)
			}

			if len(w.ledger.saved) != saves || len(w.kubernetes.deletes) != deletes || len(w.preparation.destroyed) != destroyed {
				t.Error("a refused down recorded or deleted something")
			}

			w.requireLockBalanced(t)
		})
	}
}

func TestStopRuntime_SettlesAFailedCleanup(t *testing.T) {
	t.Parallel()

	replace := func(kind, name string) func(*world) {
		return func(w *world) {
			w.kubernetes.edit(kind, name, func(object microruntime.Object) { object.Metadata()["uid"] = "uid-replacement" })
		}
	}

	tests := []struct {
		name        string
		arrange     func(*world)
		wantErr     error
		wantDeletes []string
	}{
		{
			name: "foreign object in the namespace",
			arrange: func(w *world) {
				w.kubernetes.put(microruntime.Object{"apiVersion": "v1", "kind": "ConfigMap", "metadata": map[string]any{"name": "consumer", "uid": "uid-consumer"}})
			},
			wantErr: microruntime.ErrCleanupBlocked,
		},
		{name: "recorded object replaced", arrange: replace("Deployment", "web"), wantErr: microruntime.ErrOwnershipChanged},
		{name: "namespace replaced", arrange: replace("Namespace", namespaceName), wantErr: microruntime.ErrOwnershipChanged},
		{
			name:    "inventory unavailable",
			arrange: func(w *world) { w.kubernetes.inventoryErr = errInjected },
			wantErr: errInjected,
		},
		{
			name:        "deletion request fails",
			arrange:     func(w *world) { w.kubernetes.deleteErr = errInjected },
			wantErr:     errInjected,
			wantDeletes: []string{"Deployment/web"},
		},
		{
			name:        "preparation cannot be destroyed",
			arrange:     func(w *world) { w.preparation.destroyErr = errInjected },
			wantErr:     errInjected,
			wantDeletes: []string{"Deployment/web", "Service/web", "Job/batch", "Namespace/" + namespaceName},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := newWorld(t)
			w.up(t)
			tt.arrange(w)

			_, err := w.down(false)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Execute() error = %v, want %v", err, tt.wantErr)
			}

			record := w.ledger.last(t, "demo")
			if record.Status != microruntime.StatusFailed || !record.Settled || record.Operation != microruntime.OperationDown {
				t.Fatalf("record = %s, settled %t, operation %s; want a settled failed down", record.Status, record.Settled, record.Operation)
			}

			if len(record.Resources) != 3 || record.NamespaceUID == "" {
				t.Errorf("record = %+v, want the obligations kept for a later recovery", record)
			}

			if !slices.Equal(w.kubernetes.deletes, tt.wantDeletes) {
				t.Errorf("deletes = %v, want %v", w.kubernetes.deletes, tt.wantDeletes)
			}

			if len(w.evidence.executions) != 2 || w.evidence.executions[1].Status != "failed" {
				t.Errorf("evidence = %+v, want the failed down recorded", w.evidence.executions)
			}

			w.requireLockBalanced(t)
		})
	}
}

func TestStopRuntime_RecoversAPartialUp(t *testing.T) {
	t.Parallel()

	w := newWorld(t)
	w.kubernetes.createErr = errInjected

	if _, err := w.startRuntime().Execute(context.Background(), w.runtimeRequest()); !errors.Is(err, errInjected) {
		t.Fatalf("up error = %v, want the injected creation failure", err)
	}

	w.kubernetes.createErr = nil

	if _, err := w.down(false); !errors.Is(err, microruntime.ErrCleanupBlocked) {
		t.Fatalf("down without recovery error = %v, want ErrCleanupBlocked", err)
	}

	outcome, err := w.down(true)
	if err != nil || outcome.Status != "down" {
		t.Fatalf("down with recovery = %+v, %v; want a completed down", outcome, err)
	}

	if !slices.Equal(w.kubernetes.deletes, []string{"Namespace/" + namespaceName}) {
		t.Errorf("deletes = %v, want only the namespace, since no object was created", w.kubernetes.deletes)
	}
}

func TestStopRuntime_LearnsUIDsAnInterruptedUpNeverRecorded(t *testing.T) {
	t.Parallel()

	w := newWorld(t)
	w.up(t)

	record := w.ledger.last(t, "demo")
	record.Status = microruntime.StatusInterrupted
	record.Resources[0].UID = ""

	if _, err := w.down(true); err != nil {
		t.Fatalf("down with recovery: %v", err)
	}

	if !slices.Contains(w.kubernetes.deletes, "Deployment/web") {
		t.Fatalf("deletes = %v, want the Deployment identified by its ownership token", w.kubernetes.deletes)
	}
}

func TestStopRuntime_AcceptsAnEarlierFingerprintOfTheSameConfiguration(t *testing.T) {
	t.Parallel()

	w := newWorld(t)
	w.configuration.hash = "hash-previous"
	w.up(t)
	w.configuration.hash = "hash-current"

	if outcome, err := w.down(false); err != nil || outcome.Status != "down" {
		t.Fatalf("down = %+v, %v; want the compatible fingerprint accepted", outcome, err)
	}
}

func TestStopRuntime_CompletesWhenTheNamespaceIsAlreadyGone(t *testing.T) {
	t.Parallel()

	w := newWorld(t)
	w.up(t)
	w.kubernetes.objects = nil

	outcome, err := w.down(false)
	if err != nil || outcome.Status != "down" {
		t.Fatalf("down = %+v, %v; want a completed down", outcome, err)
	}

	if len(w.kubernetes.deletes) != 0 || len(w.preparation.destroyed) != 1 {
		t.Fatalf("deletes = %v, destroys = %d; want no deletion and the preparation destroyed", w.kubernetes.deletes, len(w.preparation.destroyed))
	}
}

func TestStopRuntime_ReverifiesTheNamespaceBeforeDeletingIt(t *testing.T) {
	t.Parallel()

	w := newWorld(t)
	w.up(t)

	// Another actor replaces the namespace once the last recorded object is gone.
	w.kubernetes.onDelete = func(kubernetes *fakeKubernetes) {
		if _, job := kubernetes.find("Job", "batch"); job == nil {
			if _, namespace := kubernetes.find("Namespace", namespaceName); namespace != nil {
				namespace.Metadata()["uid"] = "uid-replacement"
			}
		}
	}

	if _, err := w.down(false); !errors.Is(err, microruntime.ErrOwnershipChanged) {
		t.Fatalf("Execute() error = %v, want ErrOwnershipChanged", err)
	}

	if slices.Contains(w.kubernetes.deletes, "Namespace/"+namespaceName) || len(w.preparation.destroyed) != 0 {
		t.Fatalf("deletes = %v, destroys = %d; a replaced namespace must not be deleted", w.kubernetes.deletes, len(w.preparation.destroyed))
	}

	if record := w.ledger.last(t, "demo"); record.Status != microruntime.StatusFailed || !record.Settled {
		t.Fatalf("record = %s, settled %t; want a settled failure", record.Status, record.Settled)
	}
}

func TestStopRuntime_RefusesANamespaceThatSurvivesCleanup(t *testing.T) {
	t.Parallel()

	w := newWorld(t)
	w.up(t)
	w.preparation.recreate = true

	if _, err := w.down(false); !errors.Is(err, microruntime.ErrCleanupBlocked) {
		t.Fatalf("Execute() error = %v, want ErrCleanupBlocked", err)
	}

	if record := w.ledger.last(t, "demo"); record.Status != microruntime.StatusFailed || record.NamespaceUID == "" {
		t.Fatalf("record = %+v, want a failed down that keeps its obligations", record)
	}
}
