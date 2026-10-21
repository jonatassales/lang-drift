package files_test

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/exphere/ex/src/infrastructure/files"
)

func writeFile(t *testing.T, path, content string) string {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	return path
}

func TestReadBounded(t *testing.T) {
	t.Parallel()

	root := tempDir(t)
	exact := writeFile(t, filepath.Join(root, "exact"), "12345")
	link := symlink(t, exact, filepath.Join(root, "link"))
	linkedDirectory := symlink(t, root, filepath.Join(tempDir(t), "linked"))

	if data, err := files.ReadBounded(exact, 5); err != nil || string(data) != "12345" {
		t.Errorf("ReadBounded(at limit) = %q, %v", data, err)
	}

	if _, err := files.ReadBounded(exact, 4); err == nil || !strings.Contains(err.Error(), "exceeds 4 bytes") {
		t.Errorf("ReadBounded(over limit) error = %v, want a size refusal rather than truncation", err)
	}

	for name, path := range map[string]string{
		"symlinked file":      link,
		"directory":           root,
		"symlinked directory": filepath.Join(linkedDirectory, "exact"),
	} {
		if _, err := files.ReadBounded(path, 1024); err == nil {
			t.Errorf("ReadBounded(%s) succeeded", name)
		}
	}

	if _, err := files.ReadBounded(filepath.Join(root, "missing"), 1024); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("ReadBounded(missing) error = %v, want fs.ErrNotExist", err)
	}
}

func TestWriteJSONAndReadJSON(t *testing.T) {
	t.Parallel()

	type record struct {
		Version int    `json:"version"`
		ID      string `json:"id"`
	}

	root := tempDir(t)
	path := filepath.Join(root, "state", "nested", "record.json")

	if err := files.WriteJSON(path, record{Version: 1, ID: "first"}); err != nil {
		t.Fatal(err)
	}

	if err := files.WriteJSON(path, record{Version: 2, ID: "second"}); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "{\n  \"version\": 2,\n  \"id\": \"second\"\n}\n" {
		t.Errorf("written document = %q, want indented JSON ending in a newline", data)
	}

	var read struct {
		ID string `json:"id"`
	}
	if err := files.ReadJSON(path, &read); err != nil || read.ID != "second" {
		t.Errorf("ReadJSON() = %+v, %v; want the replaced document with unknown fields allowed", read, err)
	}

	assertMode(t, path, 0o600)
	assertMode(t, filepath.Dir(path), 0o700)

	entries, err := os.ReadDir(filepath.Dir(path))
	if err != nil {
		t.Fatal(err)
	}

	if len(entries) != 1 {
		t.Errorf("directory holds %d entries, want no leftover temporary files", len(entries))
	}
}

func assertMode(t *testing.T, path string, want fs.FileMode) {
	t.Helper()

	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}

	if got := info.Mode().Perm(); got != want {
		t.Errorf("%s mode = %o, want %o", path, got, want)
	}
}

func TestWriteJSON_RefusesLinks(t *testing.T) {
	t.Parallel()

	root := tempDir(t)
	target := writeFile(t, filepath.Join(root, "target.json"), "{}")
	link := symlink(t, target, filepath.Join(root, "record.json"))
	linkedDirectory := symlink(t, mkdir(t, filepath.Join(root, "real")), filepath.Join(root, "linked"))

	if err := files.WriteJSON(link, map[string]int{"version": 1}); err == nil {
		t.Error("WriteJSON replaced a symlink")
	}

	if err := files.WriteJSON(filepath.Join(linkedDirectory, "record.json"), map[string]int{"version": 1}); err == nil {
		t.Error("WriteJSON wrote through a symlinked directory")
	}

	if data, _ := os.ReadFile(target); string(data) != "{}" {
		t.Errorf("the symlink target changed to %q", data)
	}
}

func TestReadJSON_ReportsInvalidDocuments(t *testing.T) {
	t.Parallel()

	path := writeFile(t, filepath.Join(tempDir(t), "record.json"), "{not json")

	var value map[string]any
	if err := files.ReadJSON(path, &value); err == nil || !strings.Contains(err.Error(), path) {
		t.Fatalf("ReadJSON() error = %v, want a decoding error naming the file", err)
	}
}

// BenchmarkWriteJSON measures the durable replacement every recorded ledger
// step pays: write, file sync, rename, and directory sync.
func BenchmarkWriteJSON(b *testing.B) {
	path := filepath.Join(tempDir(b), "state.json")

	resources := make([]map[string]string, 30)
	for index := range resources {
		resources[index] = map[string]string{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"name":       "unit",
			"namespace":  "runtime",
			"uid":        "0b6f6c1e-7a52-4c39-9a8e-1f2d3c4b5a69",
		}
	}

	document := map[string]any{"status": "running", "settled": true, "resources": resources}

	b.ReportAllocs()

	for range b.N {
		if err := files.WriteJSON(path, document); err != nil {
			b.Fatal(err)
		}
	}
}
