package manifest_test

import (
	"os"
	"path/filepath"
	"testing"

	domain "github.com/jonatassales/tao/packages/ex/src/domain/manifest"
	"github.com/jonatassales/tao/packages/ex/src/domain/project"
	"github.com/jonatassales/tao/packages/ex/src/infrastructure/manifest"
)

func TestInspectorReadsTheManifestsUnderAComponent(t *testing.T) {
	t.Parallel()

	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	owner := project.Project{Root: root}

	component := filepath.Join(root, "apps", "signup")
	if err := os.MkdirAll(component, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(component, "package.json"), []byte(`{"scripts":{"dev":"next dev -p 3000"},"dependencies":{"next":"16"}}`), 0o644); err != nil {
		t.Fatal(err)
	}

	inspection, err := (manifest.Inspector{}).Inspect(owner, "apps/signup")
	if err != nil {
		t.Fatal(err)
	}

	if inspection.Manifest != "package.json" || inspection.Kind != domain.Web || inspection.Port != 3000 {
		t.Fatalf("Inspect() = %+v", inspection)
	}

	empty, err := (manifest.Inspector{}).Inspect(owner, "apps")
	if err != nil || empty.Manifest != "" {
		t.Fatalf("a path without manifests = %+v, %v", empty, err)
	}

	if _, err := (manifest.Inspector{}).Inspect(owner, "apps/missing"); err == nil {
		t.Fatal("a missing path was inspected")
	}
}
