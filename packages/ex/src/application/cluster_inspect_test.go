package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/domain/cluster"
)

func (w *world) inspectCluster() application.InspectCluster {
	return application.InspectCluster{Projects: w.projects, Clusters: w.clusters, Observer: w.observer}
}

func TestInspectCluster_ReportsTheVerifiedRecord(t *testing.T) {
	t.Parallel()

	w := newWorld(t)

	result, err := w.inspectCluster().Execute(context.Background(), w.clusterRequest())
	if err != nil || result.Status != "ready" || result.UID != kubeSystemUID || result.ProjectID != w.owner.ID.String() {
		t.Fatalf("Execute() = %+v, %v; want the verified ready record", result, err)
	}
}

func TestInspectCluster_ReportsADeletedRecordWithoutContactingTheCluster(t *testing.T) {
	t.Parallel()

	w := newWorld(t)
	w.clusters.record.MarkDeleted()
	w.observer.namespacesErr = errInjected

	if result, err := w.inspectCluster().Execute(context.Background(), w.clusterRequest()); err != nil || result.Status != "deleted" {
		t.Fatalf("Execute() = %+v, %v; want the deleted record", result, err)
	}
}

func TestInspectCluster_RefusesAnUntrustedRecord(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		arrange func(*world)
		wantErr error
	}{
		{name: "no record", arrange: func(w *world) { w.clusters.record = nil }, wantErr: application.ErrNotFound},
		{name: "unknown project", arrange: func(w *world) { w.projects.err = errInjected }, wantErr: errInjected},
		{name: "replaced cluster", arrange: func(w *world) { w.observer.namespaces["kube-system"] = "uid-new" }, wantErr: cluster.ErrUntrustedIdentity},
		{name: "unreachable cluster", arrange: func(w *world) { w.observer.namespacesErr = errInjected }, wantErr: errInjected},
		{name: "record of another project", arrange: func(w *world) { w.clusters.record.Context = "kind-other" }, wantErr: cluster.ErrNotOwned},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := newWorld(t)
			tt.arrange(w)

			if _, err := w.inspectCluster().Execute(context.Background(), w.clusterRequest()); !errors.Is(err, tt.wantErr) {
				t.Fatalf("Execute() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
