package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jonatassales/tao/packages/ex/src/domain/project"
)

func TestConnected_PreparePreservesConsumerFilesAndAttribution(t *testing.T) {
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	owner := project.Project{Root: root}
	bootstrap := Connected{}
	if err = bootstrap.Prepare(owner); err != nil {
		t.Fatal(err)
	}
	skills, err := filepath.Glob(filepath.Join(root, ".langdrift/workflow/skills/*/SKILL.md"))
	if err != nil || len(skills) != 26 {
		t.Fatalf("skill count=%d err=%v", len(skills), err)
	}
	license, err := os.ReadFile(filepath.Join(root, ".langdrift/workflow/skills/UPSTREAM-LICENSE"))
	if err != nil || !strings.Contains(string(license), "MIT") {
		t.Fatal("upstream license missing")
	}
	target := filepath.Join(root, ".langdrift/workflow/README.md")
	if err = os.WriteFile(target, []byte("consumer customization\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err = bootstrap.Prepare(owner); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(target)
	if string(data) != "consumer customization\n" {
		t.Fatal("customized workflow replaced")
	}
	if _, err = os.Stat(filepath.Join(root, "AGENTS.md")); !os.IsNotExist(err) {
		t.Fatal("consumer root guidance was replaced")
	}
}
