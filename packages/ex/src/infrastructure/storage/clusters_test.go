package storage_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/domain/cluster"
	"github.com/exphere/ex/src/domain/project"
	"github.com/exphere/ex/src/infrastructure/storage"
)

func kubeconfigOf(owner project.Project) string {
	return filepath.Join(owner.Root, ".exphere", "state", "cluster.kubeconfig")
}

func readyCluster(t *testing.T, owner project.Project) cluster.Record {
	t.Helper()

	record := cluster.Begin(owner, kubeconfigOf(owner))
	if err := record.Identify(map[string]string{"kube-system": "uid-kube-system", "default": "uid-default"}); err != nil {
		t.Fatal(err)
	}

	if err := record.Complete(cluster.Inventory{
		Objects: map[string]cluster.ObjectIdentity{
			"true/pods/kube-system/coredns": {UID: "uid-pod", Controllers: []string{"uid-replicaset"}},
		},
		ResourceTypes: []string{"true/pods"},
	}); err != nil {
		t.Fatal(err)
	}

	record.UpdatedAt = time.Date(2026, 10, 17, 12, 0, 0, 0, time.UTC)

	return record
}

func TestClusters_SaveAndLoad(t *testing.T) {
	t.Parallel()

	owner := newProject(t)
	record := readyCluster(t, owner)

	if err := (storage.Clusters{}).Save(owner, record); err != nil {
		t.Fatal(err)
	}

	loaded, err := storage.Clusters{}.Load(owner)
	if err != nil {
		t.Fatal(err)
	}

	if loaded.ProjectID != record.ProjectID || loaded.Name != record.Name || loaded.Context != record.Context ||
		loaded.UID != record.UID || loaded.Status != record.Status || !loaded.UpdatedAt.Equal(record.UpdatedAt) ||
		loaded.Baseline.InventoryVersion != 1 || len(loaded.Baseline.Namespaces) != 2 ||
		loaded.Baseline.Inventory.Objects["true/pods/kube-system/coredns"].Controllers[0] != "uid-replicaset" {
		t.Fatalf("Load() = %+v, want the saved record", loaded)
	}

	if err := loaded.AuthorizeDeletion(record.Baseline.Namespaces, record.Baseline.Inventory); err != nil {
		t.Fatalf("a loaded record no longer authorizes deleting its own baseline: %v", err)
	}
}

func TestClusters_LoadReportsAMissingRecord(t *testing.T) {
	t.Parallel()

	if _, err := (storage.Clusters{}).Load(newProject(t)); !errors.Is(err, application.ErrNotFound) {
		t.Fatalf("Load() error = %v, want ErrNotFound", err)
	}
}

// recordAt writes a cluster record document as an earlier version would have.
func recordAt(t *testing.T, owner project.Project, document string) {
	t.Helper()

	writeFile(t, filepath.Join(owner.Root, ".exphere", "state", "cluster.json"), document)
}

func TestClusters_LoadEarlierFormats(t *testing.T) {
	t.Parallel()

	bootstrap := `"bootstrap": {"true/pods/kube-system/coredns": {"uid": "uid-pod"}}, "resourceTypes": ["true/pods"],`

	tests := []struct {
		name                 string
		fields               string
		wantInventoryVersion int
	}{
		{name: "original format with versioned inventory", fields: `"inventoryVersion": 1, ` + bootstrap, wantInventoryVersion: 1},
		{name: "original format without inventory proof", fields: bootstrap, wantInventoryVersion: 0},
		{name: "standalone record version with a complete inventory", fields: `"version": 1, ` + bootstrap, wantInventoryVersion: 1},
		{name: "standalone record version with a partial inventory", fields: `"version": 1, "resourceTypes": ["true/pods"],`, wantInventoryVersion: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			owner := newProject(t)
			recordAt(t, owner, `{`+tt.fields+`
				"projectId": "`+projectIDText+`", "name": "ex-0123456789ab", "context": "kind-ex-0123456789ab",
				"kubeconfig": "`+kubeconfigOf(owner)+`", "uid": "uid-kube-system", "status": "ready",
				"namespaceUIDs": {"kube-system": "uid-kube-system"}, "updatedAt": "2026-10-13T12:00:00Z"}`)

			loaded, err := storage.Clusters{}.Load(owner)
			if err != nil {
				t.Fatal(err)
			}

			if loaded.Baseline.InventoryVersion != tt.wantInventoryVersion || loaded.Status != cluster.StatusReady {
				t.Fatalf("Load() inventory version = %d, status %s; want %d", loaded.Baseline.InventoryVersion, loaded.Status, tt.wantInventoryVersion)
			}
		})
	}
}

func TestClusters_LoadRefusesUntrustedRecords(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		document func(owner project.Project) string
	}{
		{
			name: "unsupported record version",
			document: func(owner project.Project) string {
				return `{"version": 2, "projectId": "` + projectIDText + `", "kubeconfig": "` + kubeconfigOf(owner) + `"}`
			},
		},
		{
			name: "kubeconfig outside the project's private state",
			document: func(project.Project) string {
				return `{"projectId": "` + projectIDText + `", "kubeconfig": "/home/user/.kube/config"}`
			},
		},
		{
			name: "invalid project identity",
			document: func(owner project.Project) string {
				return `{"projectId": "unknown", "kubeconfig": "` + kubeconfigOf(owner) + `"}`
			},
		},
		{name: "malformed document", document: func(project.Project) string { return `{"projectId":` }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			owner := newProject(t)
			recordAt(t, owner, tt.document(owner))

			if _, err := (storage.Clusters{}).Load(owner); err == nil || errors.Is(err, application.ErrNotFound) {
				t.Fatalf("Load() error = %v, want the record refused", err)
			}
		})
	}
}

func TestClusters_Credentials(t *testing.T) {
	t.Parallel()

	owner := newProject(t)
	clusters := storage.Clusters{}

	location, exists, err := clusters.Credentials(owner)
	if err != nil || exists || location != kubeconfigOf(owner) {
		t.Fatalf("Credentials() = %q, %t, %v; want the private location, absent", location, exists, err)
	}

	writeFile(t, location, "apiVersion: v1\n")

	if _, exists, err := clusters.Credentials(owner); err != nil || !exists {
		t.Fatalf("Credentials() after writing = %t, %v; want present", exists, err)
	}

	for range 2 {
		if err := clusters.DiscardCredentials(owner); err != nil {
			t.Fatalf("DiscardCredentials() = %v; want it to succeed even when nothing remains", err)
		}
	}

	if _, err := os.Stat(location); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("kubeconfig still exists: %v", err)
	}
}
