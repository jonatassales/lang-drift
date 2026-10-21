package application_test

import (
	"context"
	"errors"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/domain/cluster"
)

func (w *world) createCluster() application.CreateCluster {
	return application.CreateCluster{
		Projects: w.projects, Locks: w.locks, Clusters: w.clusters, Credentials: w.credentials,
		Provisioner: w.provisioner, Observer: w.observer,
	}
}

// withoutCluster is a project whose cluster was never created.
func withoutCluster(t *testing.T) *world {
	t.Helper()

	w := newWorld(t)
	w.clusters.record = nil
	w.credentials.exists = false
	w.provisioner.clusters = nil

	return w
}

func TestCreateCluster_CreatesAndRecordsTheCluster(t *testing.T) {
	t.Parallel()

	w := withoutCluster(t)

	result, err := w.createCluster().Execute(context.Background(), w.clusterRequest())
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	name := w.owner.ClusterName()
	if result.Status != "ready" || result.Name != name || result.Context != "kind-"+name ||
		result.UID != kubeSystemUID || result.Kubeconfig != kubeconfigPath || result.InventoryVersion != 1 {
		t.Fatalf("result = %+v, want a ready, identified cluster with the private kubeconfig", result)
	}

	if len(result.Bootstrap) != 1 || len(result.ResourceTypes) != 1 || len(result.NamespaceUIDs) != 2 {
		t.Errorf("result baseline = %+v, want the observed namespaces and inventory", result)
	}

	statuses := w.clusters.statuses()
	if !slices.Equal(statuses, []cluster.Status{cluster.StatusCreating, cluster.StatusCreating, cluster.StatusReady}) {
		t.Errorf("saved statuses = %v, want intent, identity, then readiness", statuses)
	}

	if w.clusters.saved[0].UID != "" || w.clusters.saved[1].UID != kubeSystemUID {
		t.Errorf("saved UIDs = %q, %q; want the identity recorded before the inventory", w.clusters.saved[0].UID, w.clusters.saved[1].UID)
	}

	if !slices.Equal(w.provisioner.created, []string{name}) {
		t.Errorf("created = %v, want %s", w.provisioner.created, name)
	}

	w.requireLockBalanced(t)
}

func TestCreateCluster_ReusesAVerifiedReadyCluster(t *testing.T) {
	t.Parallel()

	w := newWorld(t)

	result, err := w.createCluster().Execute(context.Background(), w.clusterRequest())
	if err != nil || result.Status != "ready" {
		t.Fatalf("Execute() = %+v, %v; want the existing ready cluster", result, err)
	}

	if len(w.provisioner.created) != 0 || len(w.clusters.saved) != 0 {
		t.Fatal("reusing a ready cluster created or recorded something")
	}
}

func TestCreateCluster_RefusesBeforeCreating(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		arrange func(*world)
		wantErr error
	}{
		{name: "unknown project", arrange: func(w *world) { w.projects.err = application.NotFound("no project") }, wantErr: application.ErrNotFound},
		{name: "operation in progress", arrange: func(w *world) { w.locks.err = errInjected }, wantErr: errInjected},
		{name: "unreadable record", arrange: func(w *world) { w.clusters.loadErr = errInjected }, wantErr: errInjected},
		{
			name:    "record of another project",
			arrange: func(w *world) { w.clusters.record.Name = "ex-somebody" },
			wantErr: cluster.ErrNotOwned,
		},
		{
			name:    "ready cluster was replaced",
			arrange: func(w *world) { w.observer.namespaces["kube-system"] = "uid-replacement" },
			wantErr: cluster.ErrUntrustedIdentity,
		},
		{
			name:    "unresolved earlier creation",
			arrange: func(w *world) { w.clusters.record.Status = cluster.StatusInterrupted },
			wantErr: cluster.ErrUnresolved,
		},
		{
			name: "cluster with the project's name already exists",
			arrange: func(w *world) {
				w.clusters.record = nil
				w.credentials.exists = false
			},
			wantErr: application.ErrAdoptionRefused,
		},
		{
			name: "unrecorded kubeconfig",
			arrange: func(w *world) {
				w.clusters.record = nil
				w.provisioner.clusters = nil
			},
			wantErr: application.ErrAdoptionRefused,
		},
		{
			name: "cluster list unavailable",
			arrange: func(w *world) {
				w.clusters.record = nil
				w.provisioner.listErr = errInjected
			},
			wantErr: errInjected,
		},
		{
			name: "intent cannot be recorded",
			arrange: func(w *world) {
				w.clusters.record = nil
				w.credentials.exists = false
				w.provisioner.clusters = nil
				w.clusters.saveErr = errInjected
			},
			wantErr: errInjected,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := newWorld(t)
			tt.arrange(w)

			if _, err := w.createCluster().Execute(context.Background(), w.clusterRequest()); !errors.Is(err, tt.wantErr) {
				t.Fatalf("Execute() error = %v, want %v", err, tt.wantErr)
			}

			if len(w.provisioner.created) != 0 {
				t.Error("a refused creation created a cluster")
			}

			w.requireLockBalanced(t)
		})
	}
}

func TestCreateCluster_RecreatesAfterDeletion(t *testing.T) {
	t.Parallel()

	w := withoutCluster(t)
	deleted := cluster.Begin(w.owner, kubeconfigPath)
	deleted.MarkDeleted()
	w.clusters.record = &deleted
	w.credentials.exists = true

	if _, err := w.createCluster().Execute(context.Background(), w.clusterRequest()); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if w.credentials.discarded != 1 || len(w.provisioner.created) != 1 {
		t.Fatalf("discarded %d kubeconfigs and created %d clusters, want the leftover credentials replaced", w.credentials.discarded, len(w.provisioner.created))
	}
}

func TestCreateCluster_RecordsWhyCreationStopped(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		arrange    func(*world)
		cancel     bool
		wantErr    error
		wantStatus cluster.Status
		wantUID    string
	}{
		{
			name:       "kind fails",
			arrange:    func(w *world) { w.provisioner.createErr = errInjected },
			wantErr:    errInjected,
			wantStatus: cluster.StatusFailed,
		},
		{
			name:       "interrupted creation",
			arrange:    func(w *world) { w.provisioner.blockCreate = true },
			cancel:     true,
			wantErr:    context.Canceled,
			wantStatus: cluster.StatusInterrupted,
		},
		{
			name:       "cluster without a kube-system identity",
			arrange:    func(w *world) { delete(w.observer.namespaces, "kube-system") },
			wantErr:    cluster.ErrUntrustedIdentity,
			wantStatus: cluster.StatusFailed,
		},
		{
			name:       "inventory unavailable after identification",
			arrange:    func(w *world) { w.observer.inventoryErr = errInjected },
			wantErr:    errInjected,
			wantStatus: cluster.StatusFailed,
			wantUID:    kubeSystemUID,
		},
		{
			name:       "empty inventory",
			arrange:    func(w *world) { w.observer.inventory = cluster.Inventory{} },
			wantErr:    cluster.ErrUnsafeDeletion,
			wantStatus: cluster.StatusFailed,
			wantUID:    kubeSystemUID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := withoutCluster(t)
			tt.arrange(w)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if tt.cancel {
				time.AfterFunc(10*time.Millisecond, cancel)
			}

			if _, err := w.createCluster().Execute(ctx, w.clusterRequest()); !errors.Is(err, tt.wantErr) {
				t.Fatalf("Execute() error = %v, want %v", err, tt.wantErr)
			}

			record := w.clusters.record
			if record.Status != tt.wantStatus || record.Error == "" || record.UID != tt.wantUID {
				t.Fatalf("record = status %s, error %q, UID %q; want %s with the cause and UID %q",
					record.Status, record.Error, record.UID, tt.wantStatus, tt.wantUID)
			}

			w.requireLockBalanced(t)
		})
	}
}

func TestCreateCluster_KeepsTheCauseWhenTheFailureCannotBeRecorded(t *testing.T) {
	t.Parallel()

	w := withoutCluster(t)
	w.provisioner.createErr = errors.New("kind exited with status 1")
	w.clusters.saveErrAfter = 1

	_, err := w.createCluster().Execute(context.Background(), w.clusterRequest())
	if err == nil || !strings.Contains(err.Error(), "kind exited with status 1") || !strings.Contains(err.Error(), "record cluster failure") {
		t.Fatalf("Execute() error = %v, want both the cause and the recording failure", err)
	}

	if w.clusters.record.Status != cluster.StatusCreating {
		t.Fatalf("record = %s, want the recorded intent kept", w.clusters.record.Status)
	}
}
