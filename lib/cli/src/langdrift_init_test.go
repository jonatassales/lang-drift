package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"
)

// ex init writes the project identity and the Microruntimes execution rules,
// and keeps its private state apart. Nothing LangDrift does not prescribe
// lands in the project: no template, skill, policy text or workflow.
func TestInitWritesTheIdentityAndTheExecutionRulesOnly(t *testing.T) {
	t.Parallel()

	p := newProject(t)
	if status, _, stderr := p.ex("", "init"); status != 0 {
		t.Fatal(stderr)
	}

	folder := filepath.Join(p.root, ".langdrift")

	if got := visibleFiles(t, folder); !slices.Equal(got, []string{"config.json", "microruntimes/AGENTS.md"}) {
		t.Fatalf("ex init wrote %v outside the private state, want config.json and microruntimes/AGENTS.md", got)
	}

	rules, err := os.ReadFile(filepath.Join(folder, "microruntimes", "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	reference, err := os.ReadFile(filepath.Join("testdata", "langdrift", "sample", "microruntimes", "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(rules) != string(reference) {
		t.Error("the execution rules differ from the sample's AGENTS.md")
	}

	identity, err := os.ReadFile(filepath.Join(folder, "config.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !regexp.MustCompile(`^\{\n\t"version": 1,\n\t"id": "[0-9a-f]{32}"\n\}\n$`).Match(identity) {
		t.Errorf("config.json = %q, want the sample's tab layout with a fresh id", identity)
	}

	ignore, err := os.ReadFile(filepath.Join(folder, "agent", "state", ".gitignore"))
	if err != nil || string(ignore) != "*\n" {
		t.Errorf("private state ignore = %q, %v", ignore, err)
	}

	if status, stdout, stderr := p.ex("", "init"); status != 0 || !strings.Contains(stdout, "already initialized") {
		t.Fatalf("second ex init = %d, %s%s", status, stdout, stderr)
	}

	if again, _ := os.ReadFile(filepath.Join(folder, "config.json")); string(again) != string(identity) {
		t.Error("a second ex init changed the identity")
	}
}

// visibleFiles lists the folder's files outside its private state, as
// slash-separated paths in sorted order.
func visibleFiles(t *testing.T, folder string) []string {
	t.Helper()

	var found []string

	err := filepath.WalkDir(folder, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, err := filepath.Rel(folder, path)
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if filepath.ToSlash(relative) == "agent/state" {
				return filepath.SkipDir
			}

			return nil
		}
		found = append(found, filepath.ToSlash(relative))

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	slices.Sort(found)

	return found
}
