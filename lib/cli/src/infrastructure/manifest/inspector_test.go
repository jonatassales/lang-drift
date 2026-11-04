package manifest_test

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	domain "github.com/jonatassales/tao/lib/cli/src/domain/manifest"
	"github.com/jonatassales/tao/lib/cli/src/domain/project"
	"github.com/jonatassales/tao/lib/cli/src/infrastructure/manifest"
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

func TestInspectorReadsAFileTheComponentNames(t *testing.T) {
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
	if err := os.WriteFile(filepath.Join(component, "Dockerfile"), []byte("FROM node\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	content, err := (manifest.Inspector{}).File(owner, "apps/signup/Dockerfile")
	if err != nil || string(content) != "FROM node\n" {
		t.Fatalf("File() = %q, %v", content, err)
	}

	if _, err := (manifest.Inspector{}).File(owner, "apps/signup/Dockerfile.dev"); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("a missing file = %v, want fs.ErrNotExist", err)
	}

	if _, err := (manifest.Inspector{}).File(owner, "../outside.env"); err == nil {
		t.Fatal("a file outside the project was read")
	}
}
