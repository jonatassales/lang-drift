package stages_test

import (
	"os"
	"path/filepath"
	"testing"
)

// readSample reads a published schema or sample.
func readSample(t *testing.T, name string) string {
	t.Helper()

	raw, err := os.ReadFile(filepath.Join("schema", name))
	if err != nil {
		t.Fatal(err)
	}

	return string(raw)
}
