package storage_test

import (
	"path/filepath"
	"testing"

	"github.com/jonatassales/tao/packages/ex/src/domain/project"
	"github.com/jonatassales/tao/packages/ex/src/infrastructure/storage"
)

func TestLayout_RuntimeState(t *testing.T) {
	t.Parallel()

	owner := project.Project{Root: "/workspace"}
	runtime := runtimeID(t, "demo")

	if got := storage.LayoutOf(owner).RuntimeState(runtime); got != filepath.FromSlash("/workspace/.langdrift/state/runtimes/demo") {
		t.Fatalf("RuntimeState() = %q", got)
	}
}
