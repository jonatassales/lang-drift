package storage_test

import (
	"path/filepath"
	"testing"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/infrastructure/storage"
)

func TestEvidence_RecordHook(t *testing.T) {
	t.Parallel()

	owner := newProject(t)
	exitCode := 3

	running := application.HookOutcome{Runtime: "demo", RunID: "run-2", Hook: "seed", Status: "running"}
	if err := (storage.Evidence{}).RecordHook(owner, running); err != nil {
		t.Fatal(err)
	}

	finished := running
	finished.Status, finished.ExitCode = "failed", &exitCode
	if err := (storage.Evidence{}).RecordHook(owner, finished); err != nil {
		t.Fatal(err)
	}

	hook := readEvidence(t, filepath.Join(owner.Root, ".audits", "evidence", "runs", "run-2", "hook.json"))
	if hook["status"] != "failed" || hook["exitCode"] != 3.0 || hook["hook"] != "seed" {
		t.Errorf("hook evidence = %v, want the outcome replacing the recorded intent", hook)
	}

	if _, present := hook["stdout"]; present {
		t.Error("hook evidence retains script output")
	}
}
