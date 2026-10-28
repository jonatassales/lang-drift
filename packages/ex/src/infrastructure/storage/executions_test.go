package storage_test

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/jonatassales/tao/exphere/cli/src/application"
	"github.com/jonatassales/tao/exphere/cli/src/domain/microruntime"
	"github.com/jonatassales/tao/exphere/cli/src/domain/project"
	"github.com/jonatassales/tao/exphere/cli/src/infrastructure/storage"
)

func runtimeID(t *testing.T, text string) microruntime.Identifier {
	t.Helper()

	identifier, err := microruntime.ParseIdentifier(text)
	if err != nil {
		t.Fatal(err)
	}

	return identifier
}

// runningExecution is a completed up of one Deployment and its Service.
func runningExecution(t *testing.T, owner project.Project, runtime string) *microruntime.Execution {
	t.Helper()

	web := runtimeID(t, "web")
	container := map[string]any{"name": "web", "image": "web:dev", "imagePullPolicy": "Never"}
	manifests := []microruntime.Manifest{
		{Unit: web, Object: microruntime.Object{"apiVersion": "apps/v1", "kind": "Deployment", "metadata": map[string]any{"name": "web"},
			"spec": map[string]any{"template": map[string]any{"spec": map[string]any{"containers": []any{container}}}}}},
		{Unit: web, Object: microruntime.Object{"apiVersion": "v1", "kind": "Service", "metadata": map[string]any{"name": "web"}}},
	}

	execution := microruntime.StartExecution(microruntime.Start{
		ProjectID: owner.ID, Runtime: runtimeID(t, runtime), Namespace: runtimeID(t, runtime+"-local"),
		ClusterUID: "uid-kube-system", Context: "kind-ex-0123456789ab", ConfigHash: "hash", Token: strings.Repeat("a", 32), RunID: "run-1",
		Tools:     map[string]microruntime.ToolVersion{"kubectl": {Version: "v1.32.0", Observed: true}, "terraform": {}},
		Manifests: manifests,
	})

	owned := func(uid string) microruntime.Object {
		return microruntime.Object{"metadata": map[string]any{"uid": uid, "labels": map[string]any{microruntime.LabelOwnershipToken: execution.Token}}}
	}

	if err := execution.RecordNamespace(owned("uid-namespace")); err != nil {
		t.Fatal(err)
	}

	if err := execution.BeginApplying(); err != nil {
		t.Fatal(err)
	}

	for index, manifest := range manifests {
		created := owned("uid-" + manifest.Object.Kind())
		created["apiVersion"], created["kind"] = manifest.Object.APIVersion(), manifest.Object.Kind()
		created.Metadata()["name"], created.Metadata()["namespace"] = manifest.Object.Name(), runtime+"-local"

		if err := execution.RecordCreated(index, created); err != nil {
			t.Fatal(err)
		}
	}

	execution.MarkUnit(web, microruntime.UnitCompleted)
	if err := execution.CompleteUp(); err != nil {
		t.Fatal(err)
	}

	execution.Settle(nil, false)
	execution.Images[0].Pod, execution.Images[0].PodUID, execution.Images[0].ImageID = "web-1", "uid-pod", "sha256:feed"
	execution.UpdatedAt = time.Date(2026, 10, 17, 12, 0, 0, 0, time.UTC)

	return execution
}

func TestExecutions_SaveAndLoad(t *testing.T) {
	t.Parallel()

	owner := newProject(t)
	saved := runningExecution(t, owner, "demo")

	if err := (storage.Executions{}).Save(owner, saved); err != nil {
		t.Fatal(err)
	}

	loaded, err := storage.Executions{}.Load(owner, runtimeID(t, "demo"))
	if err != nil {
		t.Fatal(err)
	}

	if !loaded.UpdatedAt.Equal(saved.UpdatedAt) {
		t.Errorf("UpdatedAt = %v, want %v", loaded.UpdatedAt, saved.UpdatedAt)
	}

	loaded.UpdatedAt, saved.UpdatedAt = time.Time{}, time.Time{}
	if !reflect.DeepEqual(loaded, saved) {
		t.Fatalf("Load() = %+v\nwant %+v", loaded, saved)
	}

	info, err := os.Stat(filepath.Join(storage.LayoutOf(owner).RuntimeState(runtimeID(t, "demo")), "state.json"))
	if err != nil || info.Mode().Perm() != 0o600 {
		t.Fatalf("ledger mode = %v, %v; want a private file", info.Mode().Perm(), err)
	}
}

func TestExecutions_LoadReportsAMissingRecord(t *testing.T) {
	t.Parallel()

	if _, err := (storage.Executions{}).Load(newProject(t), runtimeID(t, "demo")); !errors.Is(err, application.ErrNotFound) {
		t.Fatalf("Load() error = %v, want ErrNotFound", err)
	}
}

// savedDocument saves a running execution and returns its document for editing.
func savedDocument(t *testing.T, owner project.Project) (map[string]any, string) {
	t.Helper()

	if err := (storage.Executions{}).Save(owner, runningExecution(t, owner, "demo")); err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(storage.LayoutOf(owner).RuntimeState(runtimeID(t, "demo")), "state.json")

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	var document map[string]any
	if err := json.Unmarshal(data, &document); err != nil {
		t.Fatal(err)
	}

	return document, path
}

func rewrite(t *testing.T, path string, document map[string]any) {
	t.Helper()

	data, err := json.Marshal(document)
	if err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestExecutions_LoadEarlierFormats(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		edit         func(map[string]any)
		wantObserved bool
	}{
		{name: "historical inventory left empty as an array", edit: func(d map[string]any) { d["inventory"] = []any{} }, wantObserved: true},
		{name: "historical inventory left empty as an object", edit: func(d map[string]any) { d["inventory"] = map[string]any{} }, wantObserved: true},
		{name: "historical inventory null", edit: func(d map[string]any) { d["inventory"] = nil }, wantObserved: true},
		{
			name: "standalone tool observation without a status",
			edit: func(d map[string]any) {
				d["toolVersions"] = map[string]any{"kubectl": map[string]any{"version": "v1.32.0", "observed": true}}
			},
			wantObserved: true,
		},
		{
			name: "standalone tool observation agreeing with its status",
			edit: func(d map[string]any) {
				d["toolVersions"] = map[string]any{"kubectl": map[string]any{"status": "unavailable", "observed": false}}
			},
		},
		{name: "image without a recorded status", edit: func(d map[string]any) { delete(d["images"].([]any)[0].(map[string]any), "imageStatus") }, wantObserved: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			owner := newProject(t)
			document, path := savedDocument(t, owner)
			tt.edit(document)
			rewrite(t, path, document)

			loaded, err := storage.Executions{}.Load(owner, runtimeID(t, "demo"))
			if err != nil {
				t.Fatal(err)
			}

			if loaded.ToolVersions["kubectl"].Observed != tt.wantObserved {
				t.Fatalf("kubectl observed = %t, want %t", loaded.ToolVersions["kubectl"].Observed, tt.wantObserved)
			}
		})
	}
}

func TestExecutions_LoadRefusesUntrustedRecords(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		edit         func(map[string]any)
		inconsistent bool
	}{
		{name: "unsupported version", edit: func(d map[string]any) { d["version"] = 2 }},
		{name: "invalid project identity", edit: func(d map[string]any) { d["projectId"] = "nope" }},
		{
			name:         "historical inventory with obligations",
			edit:         func(d map[string]any) { d["inventory"] = []any{map[string]any{"kind": "Deployment"}} },
			inconsistent: true,
		},
		{name: "unreadable historical inventory", edit: func(d map[string]any) { d["inventory"] = "objects" }, inconsistent: true},
		{name: "invalid runtime identifier", edit: func(d map[string]any) { d["id"] = "Demo" }, inconsistent: true},
		{name: "invalid namespace", edit: func(d map[string]any) { d["namespace"] = "" }, inconsistent: true},
		{
			name:         "resource without a valid unit",
			edit:         func(d map[string]any) { d["resources"].([]any)[0].(map[string]any)["unit"] = "Web" },
			inconsistent: true,
		},
		{name: "invalid unit name", edit: func(d map[string]any) { d["units"] = map[string]any{"-web": "completed"} }, inconsistent: true},
		{
			name:         "image status contradicting the observation",
			edit:         func(d map[string]any) { d["images"].([]any)[0].(map[string]any)["imageStatus"] = "not-observed" },
			inconsistent: true,
		},
		{
			name: "tool status contradicting the observation",
			edit: func(d map[string]any) {
				d["toolVersions"] = map[string]any{"kubectl": map[string]any{"status": "observed", "observed": false}}
			},
			inconsistent: true,
		},
		{
			name: "unknown tool status",
			edit: func(d map[string]any) {
				d["toolVersions"] = map[string]any{"kubectl": map[string]any{"status": "guessed"}}
			},
			inconsistent: true,
		},
		{
			name:         "tool without any observation status",
			edit:         func(d map[string]any) { d["toolVersions"] = map[string]any{"kubectl": map[string]any{"version": "v1"}} },
			inconsistent: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			owner := newProject(t)
			document, path := savedDocument(t, owner)
			tt.edit(document)
			rewrite(t, path, document)

			_, err := storage.Executions{}.Load(owner, runtimeID(t, "demo"))
			if err == nil || tt.inconsistent != errors.Is(err, microruntime.ErrInconsistentExecution) {
				t.Fatalf("Load() error = %v, want refusal (inconsistent %t)", err, tt.inconsistent)
			}
		})
	}
}

func TestExecutions_List(t *testing.T) {
	t.Parallel()

	owner := newProject(t)
	executions := storage.Executions{}

	if listed, err := executions.List(owner); err != nil || listed != nil {
		t.Fatalf("List() without runtimes = %v, %v; want none", listed, err)
	}

	for _, runtime := range []string{"demo", "other"} {
		if err := executions.Save(owner, runningExecution(t, owner, runtime)); err != nil {
			t.Fatal(err)
		}
	}

	listed, err := executions.List(owner)
	if err != nil || len(listed) != 2 || listed[0].Runtime.String() != "demo" || listed[1].Runtime.String() != "other" {
		t.Fatalf("List() = %v, %v; want both runtimes", listed, err)
	}
}

func TestExecutions_ListRefusesUnexpectedState(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		arrange func(t *testing.T, runtimes string)
	}{
		{name: "file among runtime directories", arrange: func(t *testing.T, runtimes string) { writeFile(t, filepath.Join(runtimes, "notes"), "") }},
		{name: "invalid runtime directory name", arrange: func(t *testing.T, runtimes string) {
			if err := os.MkdirAll(filepath.Join(runtimes, "Invalid_Name"), 0o700); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "symlinked runtime directory", arrange: func(t *testing.T, runtimes string) {
			if err := os.Symlink(tempDir(t), filepath.Join(runtimes, "linked")); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "runtime directory without a ledger", arrange: func(t *testing.T, runtimes string) {
			if err := os.MkdirAll(filepath.Join(runtimes, "partial"), 0o700); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "ledger in another runtime's directory", arrange: func(t *testing.T, runtimes string) {
			data, err := os.ReadFile(filepath.Join(runtimes, "demo", "state.json"))
			if err != nil {
				t.Fatal(err)
			}

			writeFile(t, filepath.Join(runtimes, "copied", "state.json"), string(data))
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			owner := newProject(t)
			if err := (storage.Executions{}).Save(owner, runningExecution(t, owner, "demo")); err != nil {
				t.Fatal(err)
			}

			tt.arrange(t, filepath.Dir(storage.LayoutOf(owner).RuntimeState(runtimeID(t, "demo"))))

			if _, err := (storage.Executions{}).List(owner); err == nil {
				t.Fatal("List() accepted unexpected runtime state")
			}
		})
	}
}
