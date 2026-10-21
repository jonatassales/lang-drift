package scaffold_test

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/exphere/ex/src/infrastructure/scaffold"
)

// tempDir returns a fresh directory with symlinks resolved, because the
// adapter refuses symlinked ancestors such as macOS's /var link.
func tempDir(t *testing.T) string {
	t.Helper()

	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	return dir
}

func entries(t *testing.T, root string) []string {
	t.Helper()

	var found []string

	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil || path == root {
			return err
		}

		relative, err := filepath.Rel(root, path)
		if entry.IsDir() {
			relative += "/"
		}

		found = append(found, filepath.ToSlash(relative))

		return err
	})
	if err != nil {
		t.Fatal(err)
	}

	slices.Sort(found)

	return found
}

func buildTwoFiles(stage string) error {
	if err := scaffold.WriteFile(stage, "README.md", []byte("# demo\n")); err != nil {
		return err
	}

	return scaffold.WriteJSONFile(stage, "units/web/workload.json", map[string]string{"kind": "Deployment"})
}

func TestPublish(t *testing.T) {
	t.Parallel()

	root := tempDir(t)

	location, err := scaffold.Publish(root, ".exphere/runtimes/demo", buildTwoFiles)
	if err != nil {
		t.Fatal(err)
	}

	if location != filepath.Join(root, ".exphere", "runtimes", "demo") {
		t.Fatalf("Publish() location = %q", location)
	}

	if got := entries(t, location); !slices.Equal(got, []string{"README.md", "units/", "units/web/", "units/web/workload.json"}) {
		t.Fatalf("published entries = %v", got)
	}

	workload, err := os.ReadFile(filepath.Join(location, "units", "web", "workload.json"))
	if err != nil || string(workload) != "{\n  \"kind\": \"Deployment\"\n}\n" {
		t.Fatalf("workload.json = %q, %v; want indented JSON with a trailing newline", workload, err)
	}

	if staged := entries(t, filepath.Dir(location)); !slices.Equal(staged, []string{"demo/", "demo/README.md", "demo/units/", "demo/units/web/", "demo/units/web/workload.json"}) {
		t.Fatalf("parent entries = %v, want no staging directory left behind", staged)
	}
}

func TestPublish_NeverReplacesADestination(t *testing.T) {
	t.Parallel()

	root := tempDir(t)
	existing := filepath.Join(root, ".exphere", "agents", "researcher")
	if err := os.MkdirAll(existing, 0o755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(existing, "keep.txt"), []byte("consumer edits"), 0o644); err != nil {
		t.Fatal(err)
	}

	built := false
	_, err := scaffold.Publish(root, ".exphere/agents/researcher", func(stage string) error {
		built = true

		return buildTwoFiles(stage)
	})
	if err == nil || !strings.Contains(err.Error(), "already exists") || built {
		t.Fatalf("Publish() = %v, built %t; want a refusal before building", err, built)
	}

	if got := entries(t, existing); !slices.Equal(got, []string{"keep.txt"}) {
		t.Fatalf("existing destination = %v, want it unchanged", got)
	}
}

func TestPublish_LeavesNothingWhenTheBuildFails(t *testing.T) {
	t.Parallel()

	root := tempDir(t)
	failure := errors.New("template failed")

	_, err := scaffold.Publish(root, ".exphere/runtimes/demo", func(stage string) error {
		if err := scaffold.WriteFile(stage, "README.md", []byte("partial")); err != nil {
			return err
		}

		return failure
	})
	if !errors.Is(err, failure) {
		t.Fatalf("Publish() error = %v, want the build failure", err)
	}

	if got := entries(t, filepath.Join(root, ".exphere", "runtimes")); len(got) != 0 {
		t.Fatalf("runtimes directory = %v, want no destination or staging left", got)
	}
}

func TestPublish_RefusesUnsafeDestinations(t *testing.T) {
	t.Parallel()

	for _, relative := range []string{"../escape", "/absolute", ".exphere//runtimes", ".exphere/./runtimes", `.exphere\runtimes`, ""} {
		if _, err := scaffold.Publish(tempDir(t), relative, buildTwoFiles); err == nil {
			t.Errorf("Publish(%q) succeeded", relative)
		}
	}

	root := tempDir(t)
	if err := os.Symlink(tempDir(t), filepath.Join(root, ".exphere")); err != nil {
		t.Fatal(err)
	}

	if _, err := scaffold.Publish(root, ".exphere/runtimes/demo", buildTwoFiles); err == nil {
		t.Error("Publish() wrote through a symlinked parent")
	}
}

func TestWriteFile(t *testing.T) {
	t.Parallel()

	stage := tempDir(t)

	if err := scaffold.WriteFile(stage, "nested/file.txt", []byte("content")); err != nil {
		t.Fatal(err)
	}

	if err := scaffold.WriteFile(stage, "nested/file.txt", []byte("again")); err == nil {
		t.Error("WriteFile() replaced an existing file")
	}

	for _, relative := range []string{"../escape", "/absolute", ".", "nested/../file", `nested\file`, ""} {
		if err := scaffold.WriteFile(stage, relative, nil); err == nil {
			t.Errorf("WriteFile(%q) succeeded", relative)
		}
	}

	if err := scaffold.WriteJSONFile(stage, "invalid.json", func() {}); err == nil {
		t.Error("WriteJSONFile() wrote a value JSON cannot encode")
	}
}

func TestPublish_RollsBackAPartialPublication(t *testing.T) {
	t.Parallel()

	root := tempDir(t)

	_, err := scaffold.Publish(root, ".exphere/runtimes/demo", func(stage string) error {
		if err := buildTwoFiles(stage); err != nil {
			return err
		}

		// A symlink cannot be published, so publication stops after the
		// regular files before it in walk order.
		return os.Symlink("README.md", filepath.Join(stage, "zz-link"))
	})
	if err == nil || !strings.Contains(err.Error(), "unsupported staged scaffold entry") {
		t.Fatalf("Publish() error = %v, want the unsupported entry reported", err)
	}

	if got := entries(t, filepath.Join(root, ".exphere", "runtimes")); len(got) != 0 {
		t.Fatalf("runtimes directory = %v, want the partial publication removed", got)
	}
}
