package cluster_test

import (
	"errors"
	"testing"

	"github.com/exphere/ex/src/domain/cluster"
	"github.com/exphere/ex/src/domain/project"
)

const kubeconfig = "/workspace/.exphere/state/cluster.kubeconfig"

func testProject(t *testing.T, text string) project.Project {
	t.Helper()

	id, err := project.ParseID(text)
	if err != nil {
		t.Fatal(err)
	}

	return project.Project{ID: id, Root: "/workspace"}
}

func owner(t *testing.T) project.Project {
	t.Helper()

	return testProject(t, "0123456789abcdef0123456789abcdef")
}

// readyRecord is a cluster that EX created, identified, and inventoried.
func readyRecord(t *testing.T) cluster.Record {
	t.Helper()

	record := cluster.Begin(owner(t), kubeconfig)
	if err := record.Identify(map[string]string{"kube-system": "uid-kube-system", "default": "uid-default"}); err != nil {
		t.Fatal(err)
	}

	if err := record.Complete(bootstrapInventory()); err != nil {
		t.Fatal(err)
	}

	return record
}

func TestBegin(t *testing.T) {
	t.Parallel()

	record := cluster.Begin(owner(t), kubeconfig)

	if record.ProjectID != owner(t).ID {
		t.Errorf("ProjectID = %v, want the owner's ID", record.ProjectID)
	}

	if record.Name != "ex-0123456789ab" || record.Context != "kind-ex-0123456789ab" {
		t.Errorf("Name, Context = %q, %q; want ex-0123456789ab, kind-ex-0123456789ab", record.Name, record.Context)
	}

	if record.Kubeconfig != kubeconfig || record.Status != cluster.StatusCreating || record.UID != "" {
		t.Errorf("Begin() = %+v, want a creating record with the private kubeconfig and no UID", record)
	}
}

func TestPlanCreation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		prior     *cluster.Record
		wantReuse bool
		wantErr   error
	}{
		{name: "no record", prior: nil},
		{name: "deleted cluster", prior: &cluster.Record{Status: cluster.StatusDeleted}},
		{name: "ready cluster is reused", prior: &cluster.Record{Status: cluster.StatusReady}, wantReuse: true},
		{name: "interrupted creation", prior: &cluster.Record{Status: cluster.StatusCreating}, wantErr: cluster.ErrUnresolved},
		{name: "failed creation", prior: &cluster.Record{Status: cluster.StatusFailed}, wantErr: cluster.ErrUnresolved},
		{name: "interrupted operation", prior: &cluster.Record{Status: cluster.StatusInterrupted}, wantErr: cluster.ErrUnresolved},
		{name: "unfinished deletion", prior: &cluster.Record{Status: cluster.StatusDeleting}, wantErr: cluster.ErrUnresolved},
		{name: "unknown status", prior: &cluster.Record{Status: "edited"}, wantErr: cluster.ErrUnresolved},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			reuse, err := cluster.PlanCreation(tt.prior)

			if !errors.Is(err, tt.wantErr) || (tt.wantErr == nil && err != nil) {
				t.Fatalf("PlanCreation() error = %v, want %v", err, tt.wantErr)
			}

			if reuse != tt.wantReuse {
				t.Fatalf("PlanCreation() reuse = %t, want %t", reuse, tt.wantReuse)
			}
		})
	}
}

func TestRecord_BelongsTo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		edit    func(*cluster.Record)
		owner   string
		wantErr bool
	}{
		{name: "record of the owner", edit: func(*cluster.Record) {}},
		{name: "record of another project", edit: func(*cluster.Record) {}, owner: "fedcba9876543210fedcba9876543210", wantErr: true},
		{name: "renamed cluster", edit: func(record *cluster.Record) { record.Name = "ex-other" }, wantErr: true},
		{name: "context of another cluster", edit: func(record *cluster.Record) { record.Context = "kind-other" }, wantErr: true},
		{
			name: "creation that never learned its context or identity",
			edit: func(record *cluster.Record) { record.Context, record.UID = "", "" },
		},
		{
			name:    "identified cluster without a context",
			edit:    func(record *cluster.Record) { record.Context = "" },
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			record := readyRecord(t)
			tt.edit(&record)

			candidate := owner(t)
			if tt.owner != "" {
				candidate = testProject(t, tt.owner)
			}

			err := record.BelongsTo(candidate)

			if tt.wantErr != errors.Is(err, cluster.ErrNotOwned) || (!tt.wantErr && err != nil) {
				t.Fatalf("BelongsTo() error = %v, want refusal %t", err, tt.wantErr)
			}
		})
	}
}

func TestRecord_RequireReady(t *testing.T) {
	t.Parallel()

	for _, status := range []cluster.Status{
		cluster.StatusCreating, cluster.StatusFailed, cluster.StatusInterrupted,
		cluster.StatusDeleting, cluster.StatusDeleted,
	} {
		record := cluster.Record{Status: status}
		if err := record.RequireReady(); !errors.Is(err, cluster.ErrNotReady) {
			t.Errorf("RequireReady() with status %s = %v, want ErrNotReady", status, err)
		}
	}

	ready := cluster.Record{Status: cluster.StatusReady}
	if err := ready.RequireReady(); err != nil {
		t.Errorf("RequireReady() with a ready cluster = %v", err)
	}
}

func TestRecord_Identify(t *testing.T) {
	t.Parallel()

	t.Run("records the kube-system UID and the namespace baseline", func(t *testing.T) {
		t.Parallel()

		record := cluster.Begin(owner(t), kubeconfig)
		namespaces := map[string]string{"kube-system": "uid-kube-system", "default": "uid-default"}

		if err := record.Identify(namespaces); err != nil {
			t.Fatal(err)
		}

		if record.UID != "uid-kube-system" || len(record.Baseline.Namespaces) != 2 {
			t.Fatalf("Identify() recorded UID %q and baseline %v", record.UID, record.Baseline.Namespaces)
		}
	})

	t.Run("refuses a cluster without a kube-system UID", func(t *testing.T) {
		t.Parallel()

		record := cluster.Begin(owner(t), kubeconfig)

		err := record.Identify(map[string]string{"default": "uid-default"})
		if !errors.Is(err, cluster.ErrUntrustedIdentity) {
			t.Fatalf("Identify() error = %v, want ErrUntrustedIdentity", err)
		}

		if record.UID != "" || record.Baseline.Namespaces != nil {
			t.Fatalf("Identify() changed the record on refusal: %+v", record)
		}
	})
}

func TestRecord_Complete(t *testing.T) {
	t.Parallel()

	t.Run("records the inventory and marks the cluster ready", func(t *testing.T) {
		t.Parallel()

		record := cluster.Begin(owner(t), kubeconfig)
		record.Error = "earlier failure"

		if err := record.Complete(bootstrapInventory()); err != nil {
			t.Fatal(err)
		}

		if record.Status != cluster.StatusReady || record.Error != "" || record.Baseline.InventoryVersion != 1 {
			t.Fatalf("Complete() = status %s, error %q, inventory version %d", record.Status, record.Error, record.Baseline.InventoryVersion)
		}
	})

	for name, inventory := range map[string]cluster.Inventory{
		"no objects":        {ResourceTypes: []string{"true/pods"}},
		"no resource types": {Objects: bootstrapInventory().Objects},
	} {
		t.Run("refuses an inventory with "+name, func(t *testing.T) {
			t.Parallel()

			record := cluster.Begin(owner(t), kubeconfig)

			if err := record.Complete(inventory); !errors.Is(err, cluster.ErrUnsafeDeletion) {
				t.Fatalf("Complete() error = %v, want ErrUnsafeDeletion", err)
			}

			if record.Status != cluster.StatusCreating || record.Baseline.InventoryVersion != 0 {
				t.Fatalf("Complete() changed the record on refusal: %+v", record)
			}
		})
	}
}

func TestRecord_Fail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		interrupted bool
		want        cluster.Status
	}{
		{name: "failure", want: cluster.StatusFailed},
		{name: "interruption", interrupted: true, want: cluster.StatusInterrupted},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			record := cluster.Begin(owner(t), kubeconfig)
			record.Fail(errors.New("kind exited"), tt.interrupted)

			if record.Status != tt.want || record.Error != "kind exited" {
				t.Fatalf("Fail() = status %s, error %q; want %s, kind exited", record.Status, record.Error, tt.want)
			}
		})
	}
}

func TestRecord_MarkDeletingAndDeleted(t *testing.T) {
	t.Parallel()

	record := readyRecord(t)
	record.Error = "earlier failure"

	record.MarkDeleting()
	if record.Status != cluster.StatusDeleting || record.Error != "" {
		t.Fatalf("MarkDeleting() = status %s, error %q", record.Status, record.Error)
	}

	record.Error = "interrupted deletion"
	record.MarkDeleted()
	if record.Status != cluster.StatusDeleted || record.Error != "" {
		t.Fatalf("MarkDeleted() = status %s, error %q", record.Status, record.Error)
	}
}

func TestRecord_VerifyIdentity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		uid        string
		namespaces map[string]string
		wantErr    bool
	}{
		{name: "same cluster", uid: "uid-kube-system", namespaces: map[string]string{"kube-system": "uid-kube-system"}},
		{name: "identity never recorded", uid: "", namespaces: map[string]string{"kube-system": "uid-kube-system"}, wantErr: true},
		{name: "replaced cluster", uid: "uid-kube-system", namespaces: map[string]string{"kube-system": "uid-new"}, wantErr: true},
		{name: "no observation", uid: "uid-kube-system", namespaces: nil, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			record := cluster.Record{UID: tt.uid}
			err := record.VerifyIdentity(tt.namespaces)

			if tt.wantErr != errors.Is(err, cluster.ErrUntrustedIdentity) || (!tt.wantErr && err != nil) {
				t.Fatalf("VerifyIdentity() error = %v, want refusal %t", err, tt.wantErr)
			}
		})
	}
}
