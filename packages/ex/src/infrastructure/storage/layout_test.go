package storage_test

import (
	"path/filepath"
	"testing"

	"github.com/exphere/ex/src/domain/project"
	"github.com/exphere/ex/src/infrastructure/storage"
)

func TestLayout_RuntimeState(t *testing.T) {
	t.Parallel()

	owner := project.Project{Root: "/workspace"}
	runtime := runtimeID(t, "demo")

	if got := storage.LayoutOf(owner).RuntimeState(runtime); got != filepath.FromSlash("/workspace/.exphere/state/runtimes/demo") {
		t.Fatalf("RuntimeState() = %q", got)
	}
}
