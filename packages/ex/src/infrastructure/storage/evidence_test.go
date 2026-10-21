package storage_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/infrastructure/storage"
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

func TestEvidence_RecordsExecutionsAndInspectionsUnderTheirRuns(t *testing.T) {
	t.Parallel()

	owner := newProject(t)
	evidence := storage.Evidence{}
	runs := filepath.Join(owner.Root, ".audits", "evidence", "runs")

	if err := evidence.RecordExecution(owner, application.ExecutionOutcome{Runtime: "demo", RunID: "run-1", Status: "running"}); err != nil {
		t.Fatal(err)
	}

	execution := readEvidence(t, filepath.Join(runs, "run-1", "runtime.json"))
	operation, _ := execution["operation"].(map[string]any)
	if operation["status"] != "running" || execution["billing"] != "unavailable" || execution["verification"] == "" {
		t.Errorf("execution evidence = %v, want the outcome with explicit billing and verification limits", execution)
	}

	if err := evidence.RecordInspection(owner, application.Inspection{Runtime: "demo", ObservationID: "observation-1"}); err != nil {
		t.Fatal(err)
	}

	if inspection := readEvidence(t, filepath.Join(runs, "observation-1", "inspection.json")); inspection["runtimeId"] != "demo" {
		t.Errorf("inspection evidence = %v", inspection)
	}
}
