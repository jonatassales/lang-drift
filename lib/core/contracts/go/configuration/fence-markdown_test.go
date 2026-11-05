package configuration_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// fenceMarkdown wraps a fence's JSON as runtime.md carries it.
func fenceMarkdown(t *testing.T, sample string) string {
	t.Helper()

	body, err := os.ReadFile(filepath.Join("schema", sample))
	if err != nil {
		t.Fatal(err)
	}

	return "# runtime\n\n```langdrift\n" + strings.TrimSpace(string(body)) + "\n```\n"
}
