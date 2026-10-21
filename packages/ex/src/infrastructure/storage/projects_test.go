package storage_test

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/domain/project"
	"github.com/exphere/ex/src/infrastructure/storage"
)

const projectIDText = "0123456789abcdef0123456789abcdef"

// tempDir returns a fresh directory with symlinks resolved, because the
// storage adapters refuse symlinked ancestors such as macOS's /var link.
func tempDir(t *testing.T) string {
	t.Helper()

	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	return dir
}

func projectID(t *testing.T) project.ID {
	t.Helper()

	id, err := project.ParseID(projectIDText)
	if err != nil {
		t.Fatal(err)
	}

	return id
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

// newProject initializes a project in a fresh directory.
func newProject(t *testing.T) project.Project {
	t.Helper()

	created, err := storage.Projects{}.Create(tempDir(t), projectID(t))
	if err != nil {
		t.Fatal(err)
	}

	return created
}

func TestProjects_Create(t *testing.T) {
	t.Parallel()

	root := tempDir(t)

	created, err := storage.Projects{}.Create(root, projectID(t))
	if err != nil {
		t.Fatal(err)
	}

	if created.Root != root || created.ID != projectID(t) {
		t.Errorf("Create() = %+v, want the canonical root and the given ID", created)
	}

	metadata, err := os.ReadFile(filepath.Join(root, ".exphere", "config.json"))
	if err != nil || !strings.Contains(string(metadata), `"version": 1`) || !strings.Contains(string(metadata), projectIDText) {
		t.Errorf("metadata = %s, %v; want version 1 and the ID", metadata, err)
	}

	ignore, err := os.ReadFile(filepath.Join(root, ".exphere", "state", ".gitignore"))
	if err != nil || string(ignore) != "*\n" {
		t.Errorf("state ignore = %q, %v; want every private file ignored", ignore, err)
	}

	info, err := os.Stat(filepath.Join(root, ".exphere", "state"))
	if err != nil || info.Mode().Perm() != 0o700 {
		t.Errorf("state directory mode = %v, %v; want 0700", info.Mode().Perm(), err)
	}

	if _, err := (storage.Projects{}).Create(root, projectID(t)); !errors.Is(err, application.ErrAlreadyExists) {
		t.Errorf("second Create() error = %v, want ErrAlreadyExists", err)
	}
}

func TestProjects_CreateKeepsAnExistingStateIgnore(t *testing.T) {
	t.Parallel()

	root := tempDir(t)
	writeFile(t, filepath.Join(root, ".exphere", "state", ".gitignore"), "*\n")

	if _, err := (storage.Projects{}).Create(root, projectID(t)); err != nil {
		t.Fatalf("Create() with a matching state ignore = %v", err)
	}

	other := tempDir(t)
	writeFile(t, filepath.Join(other, ".exphere", "state", ".gitignore"), "*.log\n")

	if _, err := (storage.Projects{}).Create(other, projectID(t)); err == nil {
		t.Fatal("Create() accepted a state ignore file that does not ignore every private file")
	}
}

func TestProjects_CreateRefusesSymlinkedMetadata(t *testing.T) {
	t.Parallel()

	root := tempDir(t)
	if err := os.Symlink(tempDir(t), filepath.Join(root, ".exphere")); err != nil {
		t.Fatal(err)
	}

	if _, err := (storage.Projects{}).Create(root, projectID(t)); err == nil {
		t.Fatal("Create() wrote through a symlinked .exphere directory")
	}
}

func TestProjects_Open(t *testing.T) {
	t.Parallel()

	created := newProject(t)
	nested := filepath.Join(created.Root, "services", "web")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}

	t.Run("discovers the nearest project", func(t *testing.T) {
		t.Parallel()

		opened, err := storage.Projects{}.Open(application.ProjectSelection{Directory: nested})
		if err != nil || opened != created {
			t.Fatalf("Open(nested) = %+v, %v; want %+v", opened, err, created)
		}
	})

	t.Run("explicit selection requires the root itself", func(t *testing.T) {
		t.Parallel()

		if opened, err := (storage.Projects{}).Open(application.ProjectSelection{Directory: created.Root, Explicit: true}); err != nil || opened != created {
			t.Fatalf("Open(root) = %+v, %v", opened, err)
		}

		if _, err := (storage.Projects{}).Open(application.ProjectSelection{Directory: nested, Explicit: true}); !errors.Is(err, application.ErrNotFound) {
			t.Fatalf("Open(nested, explicit) error = %v, want ErrNotFound", err)
		}
	})

	t.Run("no project anywhere above", func(t *testing.T) {
		t.Parallel()

		if _, err := (storage.Projects{}).Open(application.ProjectSelection{Directory: tempDir(t)}); !errors.Is(err, application.ErrNotFound) {
			t.Fatalf("Open() error = %v, want ErrNotFound", err)
		}
	})
}

func TestProjects_OpenRefusesUntrustedMetadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		metadata string
	}{
		{name: "unsupported version", metadata: `{"version": 2, "id": "` + projectIDText + `"}`},
		{name: "invalid identity", metadata: `{"version": 1, "id": "not-an-id"}`},
		{name: "unreadable document", metadata: `{"version": 1,`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			parent := tempDir(t)
			writeFile(t, filepath.Join(parent, ".exphere", "config.json"), `{"version": 1, "id": "`+projectIDText+`"}`)
			root := filepath.Join(parent, "child")
			writeFile(t, filepath.Join(root, ".exphere", "config.json"), tt.metadata)

			_, err := storage.Projects{}.Open(application.ProjectSelection{Directory: root})
			if err == nil || errors.Is(err, application.ErrNotFound) {
				t.Fatalf("Open() error = %v, want invalid metadata reported instead of selecting the parent", err)
			}
		})
	}
}

func TestProjects_OpenRefusesSymlinkedMetadataFile(t *testing.T) {
	t.Parallel()

	root := tempDir(t)
	elsewhere := filepath.Join(tempDir(t), "config.json")
	writeFile(t, elsewhere, `{"version": 1, "id": "`+projectIDText+`"}`)

	if err := os.MkdirAll(filepath.Join(root, ".exphere"), 0o755); err != nil {
		t.Fatal(err)
	}

	if err := os.Symlink(elsewhere, filepath.Join(root, ".exphere", "config.json")); err != nil {
		t.Fatal(err)
	}

	if _, err := (storage.Projects{}).Open(application.ProjectSelection{Directory: root, Explicit: true}); err == nil {
		t.Fatal("Open() followed a symlinked metadata file")
	}
}
