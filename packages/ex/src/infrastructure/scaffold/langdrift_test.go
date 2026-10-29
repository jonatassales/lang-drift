package scaffold_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jonatassales/tao/packages/ex/src/domain/project"
	"github.com/jonatassales/tao/packages/ex/src/infrastructure/scaffold"
)

func TestLangDrift_WritesTheExecutionRulesOnce(t *testing.T) {
	t.Parallel()

	root := tempDir(t)
	owner := project.Project{Root: root}
	rules := filepath.Join(root, ".langdrift", "microruntimes", "AGENTS.md")

	if err := (scaffold.LangDrift{}).Prepare(owner); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(rules)
	if err != nil {
		t.Fatal(err)
	}
	if len(content) == 0 || string(content[:len("# Microruntimes execution rules")]) != "# Microruntimes execution rules" {
		t.Fatalf("execution rules start with %q", content[:min(len(content), 40)])
	}

	entries, err := os.ReadDir(filepath.Join(root, ".langdrift"))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Name() != "microruntimes" {
		t.Fatalf(".langdrift holds %d entries after Prepare, want only microruntimes", len(entries))
	}

	if err := os.WriteFile(rules, []byte("# The team's own rules\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := (scaffold.LangDrift{}).Prepare(owner); err != nil {
		t.Fatal(err)
	}

	if kept, _ := os.ReadFile(rules); string(kept) != "# The team's own rules\n" {
		t.Fatalf("a second Prepare replaced the team's file: %q", kept)
	}
}

func TestLangDrift_RefusesASymlinkedFolder(t *testing.T) {
	t.Parallel()

	root := tempDir(t)
	if err := os.Symlink(tempDir(t), filepath.Join(root, ".langdrift")); err != nil {
		t.Fatal(err)
	}

	if err := (scaffold.LangDrift{}).Prepare(project.Project{Root: root}); err == nil {
		t.Fatal("Prepare wrote through a symlinked .langdrift")
	}
}
