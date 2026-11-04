package storage_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/jonatassales/tao/lib/cli/src/application"
	"github.com/jonatassales/tao/lib/cli/src/infrastructure/storage"
)

func readEvidence(t *testing.T, path string) map[string]any {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	var document map[string]any
	if err := json.Unmarshal(data, &document); err != nil {
		t.Fatal(err)
	}

	return document
}

// runtimeFolder is the runtime's own folder under the project's microruntimes.
func runtimeFolder(root, runtime string) string {
	return filepath.Join(root, ".langdrift", "microruntimes", "runtimes", runtime)
}

func TestEvidence_RecordsExecutionsAndInspectionsInsideTheRuntime(t *testing.T) {
	t.Parallel()

	owner := newProject(t)
	evidence := storage.Evidence{}
	executions := filepath.Join(runtimeFolder(owner.Root, "demo"), "evidence")

	if err := evidence.RecordExecution(owner, application.ExecutionOutcome{Runtime: "demo", RunID: "run-1", Status: "running"}); err != nil {
		t.Fatal(err)
	}

	execution := readEvidence(t, filepath.Join(executions, "run-1", "runtime.json"))
	operation, _ := execution["operation"].(map[string]any)
	if operation["status"] != "running" || execution["billing"] != "unavailable" || execution["verification"] == "" {
		t.Errorf("execution evidence = %v, want the outcome with explicit billing and verification limits", execution)
	}

	report := application.Inspection{Runtime: "demo", ExecutionRunID: "run-1", ObservationID: "observation-1"}
	if err := evidence.RecordInspection(owner, report); err != nil {
		t.Fatal(err)
	}

	if inspection := readEvidence(t, filepath.Join(executions, "run-1", "inspections", "observation-1.json")); inspection["runtimeId"] != "demo" {
		t.Errorf("inspection evidence = %v", inspection)
	}

	if _, err := os.Stat(filepath.Join(owner.Root, ".audits")); !os.IsNotExist(err) {
		t.Error("evidence was written outside the runtime")
	}
}

func TestEvidence_GivesAnUnplannedExecutionARunOfItsOwn(t *testing.T) {
	t.Parallel()

	owner := newProject(t)

	if err := (storage.Evidence{}).RecordExecution(owner, application.ExecutionOutcome{Runtime: "demo", RunID: "run-1", Status: "running"}); err != nil {
		t.Fatal(err)
	}

	run := readEvidence(t, filepath.Join(runtimeFolder(owner.Root, "demo"), "run.json"))
	if run["kind"] != "execution" || run["id"] != "demo" || run["version"] != 1.0 {
		t.Errorf("run record = %v, want an execution-kind run for the runtime", run)
	}
}

func TestEvidence_RefusesARecordWithoutItsExecution(t *testing.T) {
	t.Parallel()

	owner := newProject(t)

	err := (storage.Evidence{}).RecordInspection(owner, application.Inspection{Runtime: "demo", ObservationID: "observation-1"})
	if err == nil {
		t.Fatal("an inspection without an execution was recorded")
	}
}
