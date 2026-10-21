package files_test

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/exphere/ex/src/infrastructure/files"
)

// tempDir returns a fresh directory with symlinks resolved, because these
// adapters refuse symlinked ancestors such as macOS's /var link.
func tempDir(t testing.TB) string {
	t.Helper()

	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	return dir
}

func mkdir(t *testing.T, path string) string {
	t.Helper()

	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatal(err)
	}

	return path
}

func symlink(t *testing.T, target, link string) string {
	t.Helper()

	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}

	return link
}

func TestCanonicalDirectory(t *testing.T) {
	t.Parallel()

	root := tempDir(t)
	real := mkdir(t, filepath.Join(root, "real"))
	link := symlink(t, real, filepath.Join(root, "link"))

	if got, err := files.CanonicalDirectory(link); err != nil || got != real {
		t.Errorf("CanonicalDirectory(link) = %q, %v; want %q", got, err, real)
	}

	if got, err := files.CanonicalDirectory(real + "/."); err != nil || got != real {
		t.Errorf("CanonicalDirectory(unclean) = %q, %v; want %q", got, err, real)
	}

	file := filepath.Join(root, "file")
	if err := os.WriteFile(file, nil, 0o600); err != nil {
		t.Fatal(err)
	}

	if _, err := files.CanonicalDirectory(file); err == nil {
		t.Error("CanonicalDirectory accepted a regular file")
	}

	if _, err := files.CanonicalDirectory(filepath.Join(root, "missing")); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("CanonicalDirectory(missing) error = %v, want fs.ErrNotExist", err)
	}
}

func TestWithinAndExisting(t *testing.T) {
	t.Parallel()

	root := tempDir(t)
	mkdir(t, filepath.Join(root, "runtimes", "demo"))
	outside := mkdir(t, filepath.Join(root, "outside"))
	symlink(t, outside, filepath.Join(root, "runtimes", "linked"))

	tests := []struct {
		name         string
		relative     string
		wantWithin   bool
		wantExisting bool
	}{
		{name: "existing path", relative: "runtimes/demo", wantWithin: true, wantExisting: true},
		{name: "path not created yet", relative: "runtimes/demo/state", wantWithin: true},
		{name: "symlinked component", relative: "runtimes/linked/file"},
		{name: "escaping path", relative: "../escape"},
		{name: "absolute path", relative: "/etc"},
		{name: "empty path", relative: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			want := filepath.Join(root, filepath.FromSlash(tt.relative))

			got, err := files.Within(root, tt.relative)
			if tt.wantWithin != (err == nil) || (err == nil && got != want) {
				t.Errorf("Within(%q) = %q, %v; want success %t", tt.relative, got, err, tt.wantWithin)
			}

			got, err = files.Existing(root, tt.relative)
			if tt.wantExisting != (err == nil) || (err == nil && got != want) {
				t.Errorf("Existing(%q) = %q, %v; want success %t", tt.relative, got, err, tt.wantExisting)
			}
		})
	}

	if _, err := files.Existing(root, "runtimes/missing"); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("Existing(missing) error = %v, want fs.ErrNotExist", err)
	}

	link := symlink(t, root, filepath.Join(tempDir(t), "project"))
	if _, err := files.Within(link, "runtimes"); err == nil {
		t.Error("Within accepted a symlinked base")
	}
}

func TestRequireNoSymlinks(t *testing.T) {
	t.Parallel()

	root := tempDir(t)
	real := mkdir(t, filepath.Join(root, "real"))
	link := symlink(t, real, filepath.Join(root, "link"))

	if err := files.RequireNoSymlinks(filepath.Join(real, "not", "created", "yet")); err != nil {
		t.Errorf("RequireNoSymlinks(missing descendant) = %v", err)
	}

	for _, path := range []string{link, filepath.Join(link, "child")} {
		if err := files.RequireNoSymlinks(path); err == nil {
			t.Errorf("RequireNoSymlinks(%q) accepted a symlink", path)
		}
	}
}
