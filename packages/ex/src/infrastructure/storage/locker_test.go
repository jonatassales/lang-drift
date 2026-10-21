package storage_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/exphere/ex/src/infrastructure/storage"
)

func TestLocker_Lock(t *testing.T) {
	t.Parallel()

	owner := newProject(t)

	unlock, err := storage.Locker{}.Lock(owner)
	if err != nil {
		t.Fatal(err)
	}
	defer unlock()

	if _, err := (storage.Locker{}).Lock(owner); err == nil {
		t.Fatal("a second operation acquired the project lock")
	}

	if _, err := os.Stat(filepath.Join(owner.Root, ".exphere", "state", "cluster.lock")); err != nil {
		t.Fatalf("project lock file: %v", err)
	}
}
