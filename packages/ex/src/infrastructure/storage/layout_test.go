package storage_test

import (
	"path/filepath"
	"testing"

	"github.com/jonatassales/tao/packages/ex/src/domain/project"
	"github.com/jonatassales/tao/packages/ex/src/infrastructure/storage"
)

func TestLayout_PrivateStateAndRuntimeFolder(t *testing.T) {
	t.Parallel()

	owner := project.Project{Root: "/workspace"}
	runtime := runtimeID(t, "demo")
	layout := storage.LayoutOf(owner)

	if got := layout.RuntimeState(runtime); got != filepath.FromSlash("/workspace/.langdrift/agent/state/runtimes/demo") {
		t.Fatalf("RuntimeState() = %q", got)
	}

	if got := layout.Runtime(runtime); got != filepath.FromSlash("/workspace/.langdrift/microruntimes/runtimes/demo") {
		t.Fatalf("Runtime() = %q", got)
	}
}
