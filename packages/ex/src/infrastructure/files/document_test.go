package files_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jonatassales/tao/packages/ex/src/infrastructure/files"
)

func TestWriteDocument_WritesAReadableFileInAReadableDirectory(t *testing.T) {
	t.Parallel()

	root := tempDir(t)
	path := filepath.Join(root, "microruntimes", "runtimes", "demo", "run.json")

	if err := files.WriteDocument(path, []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(path)
	if err != nil || info.Mode().Perm() != 0o644 {
		t.Fatalf("document mode = %v, %v; want 0644", info.Mode().Perm(), err)
	}

	directory, err := os.Stat(filepath.Dir(path))
	if err != nil || directory.Mode().Perm() != 0o755 {
		t.Fatalf("directory mode = %v, %v; want 0755", directory.Mode().Perm(), err)
	}

	if err := files.WriteDocument(path, []byte("{\"version\": 1}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if content, _ := os.ReadFile(path); string(content) != "{\"version\": 1}\n" {
		t.Fatalf("document = %q after the second write", content)
	}
}

func TestWriteDocument_RefusesALinkedTarget(t *testing.T) {
	t.Parallel()

	root := tempDir(t)
	elsewhere := filepath.Join(tempDir(t), "elsewhere.json")
	if err := os.WriteFile(elsewhere, []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(root, "run.json")
	if err := os.Symlink(elsewhere, path); err != nil {
		t.Fatal(err)
	}

	if err := files.WriteDocument(path, []byte("{}\n"), 0o644); err == nil {
		t.Fatal("WriteDocument wrote through a symlink")
	}
}
