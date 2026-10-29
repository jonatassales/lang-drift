package collection_test

import (
	"os"
	"path/filepath"
	"testing"

	domain "github.com/jonatassales/tao/packages/ex/src/domain/collection"
	"github.com/jonatassales/tao/packages/ex/src/domain/project"
	"github.com/jonatassales/tao/packages/ex/src/infrastructure/collection"
)

func tempProject(t *testing.T) project.Project {
	t.Helper()

	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	return project.Project{Root: root}
}

func write(t *testing.T, root, relative string, content []byte) {
	t.Helper()

	path := filepath.Join(root, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestCollectorFollowsLinksAndRecordsWhatItRefuses(t *testing.T) {
	t.Parallel()

	owner := tempProject(t)
	write(t, owner.Root, "tickets/0001-demo.ticket.md", []byte("# 0001: Demo\n\n**Spec:** specs/demo.spec.md\n\n**Status:** ready-for-agent\n\n- [the design](/docs/design.md)\n- [the env](/.env)\n- [a link](/linked.md)\n- [outside](../outside.md)\n- [missing](/docs/missing.md)\n- [an image](/docs/logo.png)\n- [the site](https://example.com/)\n"))
	write(t, owner.Root, "specs/demo.spec.md", []byte("# Demo spec\n\nBack to [the ticket](/tickets/0001-demo.ticket.md) and on to [deep](/docs/deep.md).\n"))
	write(t, owner.Root, "docs/design.md", []byte("# Design\n"))
	write(t, owner.Root, "docs/deep.md", []byte("# Deep\n\n[deeper](/docs/deeper.md)\n"))
	write(t, owner.Root, "docs/deeper.md", []byte("# Deeper\n"))
	write(t, owner.Root, "docs/logo.png", []byte("PNG\x00\x01\x02"))
	write(t, owner.Root, ".env", []byte("SECRET=x\n"))
	if err := os.Symlink(filepath.Join(owner.Root, "docs", "design.md"), filepath.Join(owner.Root, "linked.md")); err != nil {
		t.Fatal(err)
	}

	result, err := (collection.Collector{}).Collect(owner, "tickets/0001-demo.ticket.md", domain.Limits{MaxDepth: 2, MaxFiles: 64, MaxBytes: 1 << 20})
	if err != nil {
		t.Fatal(err)
	}

	states := map[string]string{}
	for _, source := range result.Sources {
		states[source.Path] = source.State
	}

	want := map[string]string{
		"tickets/0001-demo.ticket.md": "read",
		"specs/demo.spec.md":          "read",
		"docs/design.md":              "read",
		".env":                        "refused",
		"linked.md":                   "refused",
		"../outside.md":               "refused",
		"docs/missing.md":             "unresolved",
		"docs/logo.png":               "binary",
		"docs/deep.md":                "read",
	}
	for path, state := range want {
		if states[path] != state {
			t.Errorf("%s is %q, want %q", path, states[path], state)
		}
	}

	if _, reached := states["docs/deeper.md"]; reached {
		t.Error("a file beyond the depth limit was visited")
	}

	if result.Sources[0].Path != "tickets/0001-demo.ticket.md" || result.Sources[0].LinkSyntax != domain.NamedInCLI || result.Sources[1].Path != "specs/demo.spec.md" || result.Sources[1].LinkSyntax != domain.FieldLine {
		t.Errorf("sources are not in document order: %v, %v", result.Sources[0], result.Sources[1])
	}

	if result.Files != 5 || result.Depth != 2 {
		t.Errorf("reached %d files to depth %d, want 5 files to depth 2", result.Files, result.Depth)
	}

	for _, source := range result.Sources {
		if source.State == "read" && len(source.Content) == 0 {
			t.Errorf("%s was read without content", source.Path)
		}
		if source.State != "read" && source.State != "binary" && (source.SHA256 != "" || source.Bytes != 0) {
			t.Errorf("%s records a hash although it was not read", source.Path)
		}
	}

	ticket := result.Sources[0]
	if ticket.Bytes == 0 || ticket.SHA256 == "" || ticket.MediaType != "text/markdown" {
		t.Errorf("the ticket is described as %+v", ticket)
	}
}

func TestCollectorStopsAtTheFileLimit(t *testing.T) {
	t.Parallel()

	owner := tempProject(t)
	write(t, owner.Root, "ticket.md", []byte("[a](/a.md) [b](/b.md) [c](/c.md)\n"))
	for _, name := range []string{"a.md", "b.md", "c.md"} {
		write(t, owner.Root, name, []byte("# "+name+"\n"))
	}

	result, err := (collection.Collector{}).Collect(owner, "ticket.md", domain.Limits{MaxDepth: 3, MaxFiles: 2, MaxBytes: 1 << 20})
	if err != nil {
		t.Fatal(err)
	}

	if result.Files != 2 || len(result.Sources) != 4 || result.Sources[3].State != "unresolved" {
		t.Fatalf("Collect() read %d files of %d sources; the third link should be known but unread: %+v", result.Files, len(result.Sources), result.Sources)
	}
}
