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

func (w *world) deleteCluster() application.DeleteCluster {
	return application.DeleteCluster{
		Projects: w.projects, Locks: w.locks, Clusters: w.clusters, Credentials: w.credentials,
		Executions: w.ledger, Provisioner: w.provisioner, Observer: w.observer,
	}
}

func TestDeleteCluster_DeletesAClusterHoldingOnlyItsBaseline(t *testing.T) {
	t.Parallel()

	w := newWorld(t)
	w.up(t)

	if _, err := w.down(false); err != nil {
		t.Fatal(err)
	}

	result, err := w.deleteCluster().Execute(context.Background(), w.clusterRequest())
	if err != nil || result.Status != "deleted" {
		t.Fatalf("Execute() = %+v, %v; want a deleted cluster", result, err)
	}

	if !slices.Equal(w.provisioner.deleted, []string{w.owner.ClusterName()}) || w.credentials.discarded != 1 {
		t.Errorf("deleted %v and discarded %d kubeconfigs, want one of each", w.provisioner.deleted, w.credentials.discarded)
	}

	if statuses := w.clusters.statuses(); !slices.Equal(statuses, []cluster.Status{cluster.StatusDeleting, cluster.StatusDeleted}) {
		t.Errorf("saved statuses = %v, want intent before deletion, then the deletion", statuses)
	}

	w.requireLockBalanced(t)
}

func TestDeleteCluster_ForgetsAClusterThatIsAlreadyGone(t *testing.T) {
	t.Parallel()

	w := newWorld(t)
	w.clusters.record.Status = cluster.StatusFailed
	w.provisioner.clusters = nil

	result, err := w.deleteCluster().Execute(context.Background(), w.clusterRequest())
	if err != nil || result.Status != "deleted" {
		t.Fatalf("Execute() = %+v, %v; want the record reconciled as deleted", result, err)
	}

	if len(w.provisioner.deleted) != 0 || w.credentials.discarded != 1 {
		t.Fatalf("deleted %v and discarded %d kubeconfigs; want no deletion and the credentials discarded", w.provisioner.deleted, w.credentials.discarded)
	}
}

func TestDeleteCluster_ReturnsADeletedRecordUnchanged(t *testing.T) {
	t.Parallel()

	w := newWorld(t)
	w.clusters.record.MarkDeleted()

	if result, err := w.deleteCluster().Execute(context.Background(), w.clusterRequest()); err != nil || result.Status != "deleted" {
		t.Fatalf("Execute() = %+v, %v", result, err)
	}

	if len(w.clusters.saved) != 0 || len(w.provisioner.deleted) != 0 {
		t.Fatal("a deleted record was changed or deleted again")
	}
}

func TestDeleteCluster_RefusesWithoutDeleting(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		arrange func(*testing.T, *world)
		wantErr error
	}{
		{name: "no cluster", arrange: func(_ *testing.T, w *world) { w.clusters.record = nil }, wantErr: application.ErrNotFound},
		{name: "runtime still running", arrange: func(t *testing.T, w *world) { w.up(t) }, wantErr: microruntime.ErrExistingObligations},
		{
			name: "inconsistent runtime record",
			arrange: func(t *testing.T, w *world) {
				w.up(t)
				if _, err := w.down(false); err != nil {
					t.Fatal(err)
				}
				w.ledger.last(t, "demo").Token = ""
			},
			wantErr: microruntime.ErrInconsistentExecution,
		},
		{name: "runtime records unreadable", arrange: func(_ *testing.T, w *world) { w.ledger.listErr = errInjected }, wantErr: errInjected},
		{
			name:    "replaced cluster",
			arrange: func(_ *testing.T, w *world) { w.observer.namespaces["kube-system"] = "uid-replacement" },
			wantErr: cluster.ErrUntrustedIdentity,
		},
		{
			name: "creation that never recorded the cluster's identity",
			arrange: func(_ *testing.T, w *world) {
				w.clusters.record.UID, w.clusters.record.Status = "", cluster.StatusInterrupted
			},
			wantErr: cluster.ErrUntrustedIdentity,
		},
		{
			name:    "foreign namespace",
			arrange: func(_ *testing.T, w *world) { w.observer.namespaces["consumer"] = "uid-consumer" },
			wantErr: cluster.ErrUnsafeDeletion,
		},
		{
			name: "foreign object",
			arrange: func(_ *testing.T, w *world) {
				w.observer.inventory.Objects = map[string]cluster.ObjectIdentity{
					"true/pods/kube-system/coredns":    {UID: "uid-coredns"},
					"true/configmaps/default/consumer": {UID: "uid-consumer"},
				}
			},
			wantErr: cluster.ErrUnsafeDeletion,
		},
		{name: "inventory unavailable", arrange: func(_ *testing.T, w *world) { w.observer.inventoryErr = errInjected }, wantErr: errInjected},
		{name: "cluster list unavailable", arrange: func(_ *testing.T, w *world) { w.provisioner.listErr = errInjected }, wantErr: errInjected},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := newWorld(t)
			tt.arrange(t, w)
			saves := len(w.clusters.saved)

			if _, err := w.deleteCluster().Execute(context.Background(), w.clusterRequest()); !errors.Is(err, tt.wantErr) {
				t.Fatalf("Execute() error = %v, want %v", err, tt.wantErr)
			}

			if len(w.provisioner.deleted) != 0 || len(w.clusters.saved) != saves || w.credentials.discarded != 0 {
				t.Error("a refused deletion deleted, recorded, or discarded something")
			}

			w.requireLockBalanced(t)
		})
	}
}

func TestDeleteCluster_RecordsAFailedDeletion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		arrange func(*world)
		wantErr error
	}{
		{name: "kind fails", arrange: func(w *world) { w.provisioner.deleteErr = errInjected }, wantErr: errInjected},
		{name: "cluster survives deletion", arrange: func(w *world) { w.provisioner.keepAfterDelete = true }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := newWorld(t)
			tt.arrange(w)

			_, err := w.deleteCluster().Execute(context.Background(), w.clusterRequest())
			if err == nil || (tt.wantErr != nil && !errors.Is(err, tt.wantErr)) {
				t.Fatalf("Execute() error = %v, want %v", err, tt.wantErr)
			}

			if w.clusters.record.Status != cluster.StatusFailed || w.clusters.record.Error == "" || w.credentials.discarded != 0 {
				t.Fatalf("record = %s (%q), discarded %d; want the failure recorded and the credentials kept",
					w.clusters.record.Status, w.clusters.record.Error, w.credentials.discarded)
			}
		})
	}
}

func TestDeleteCluster_ReportsCredentialsThatCannotBeDiscarded(t *testing.T) {
	t.Parallel()

	w := newWorld(t)
	w.credentials.discardErr = errInjected

	if _, err := w.deleteCluster().Execute(context.Background(), w.clusterRequest()); !errors.Is(err, errInjected) {
		t.Fatalf("Execute() error = %v, want the discard failure", err)
	}

	if w.clusters.record.Status != cluster.StatusDeleted {
		t.Fatalf("record = %s, want the deletion recorded before discarding credentials", w.clusters.record.Status)
	}
}
