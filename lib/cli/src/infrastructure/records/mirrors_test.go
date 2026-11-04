package records_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jonatassales/tao/lib/cli/src/domain/project"
	"github.com/jonatassales/tao/lib/cli/src/infrastructure/records"
)

func tempProject(t *testing.T) project.Project {
	t.Helper()

	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	return project.Project{Root: root}
}

func TestMirrors_WriteMirrorFormatsAReadableFile(t *testing.T) {
	t.Parallel()

	owner := tempProject(t)
	value := struct {
		Fields []string `json:"fields"`
	}{Fields: []string{"gpu", "disk"}}

	if err := (records.Mirrors{}).WriteMirror(owner, "harness/envelope.json", value); err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(owner.Root, ".langdrift", "harness", "envelope.json")
	content, err := os.ReadFile(path)
	if err != nil || string(content) != "{\n\t\"fields\": [\"gpu\", \"disk\"]\n}\n" {
		t.Fatalf("mirror = %q, %v", content, err)
	}

	if info, _ := os.Stat(path); info.Mode().Perm() != 0o644 {
		t.Fatalf("mirror mode = %v, want 0644", info.Mode().Perm())
	}
}

func TestMirrors_NextPriceNumberContinuesTheSequence(t *testing.T) {
	t.Parallel()

	owner := tempProject(t)

	number, err := (records.Mirrors{}).NextPriceNumber(owner)
	if err != nil || number != "0001" {
		t.Fatalf("first number = %q, %v", number, err)
	}

	prices := filepath.Join(owner.Root, ".langdrift", "finops", "prices")
	if err := os.MkdirAll(prices, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"0001-illustrative.prices.json", "0007-october.prices.json", "notes.md"} {
		if err := os.WriteFile(filepath.Join(prices, name), []byte("{}\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	number, err = (records.Mirrors{}).NextPriceNumber(owner)
	if err != nil || number != "0008" {
		t.Fatalf("next number = %q, %v; want 0008", number, err)
	}
}
