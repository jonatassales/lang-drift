package connected

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	app "github.com/jonatassales/tao/lib/cli/src/application/connected"
)

func TestLocalSourceIncludesEditsAndExcludesPrivateIgnoredAndLinkedFiles(t *testing.T) {
	owner := ownerForTest(t)
	contents := map[string]string{"Dockerfile": "FROM scratch\n", "app.txt": "first version", ".gitignore": "ignored/\n*.log\n", "ignored/private.txt": "ignored", "debug.log": "ignored", ".env": "synthetic-secret", "node_modules/dependency/index.js": "dependency"}
	for name, content := range contents {
		path := filepath.Join(owner.Root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Symlink(filepath.Join(owner.Root, ".env"), filepath.Join(owner.Root, "linked")); err != nil {
		t.Fatal(err)
	}
	native := Native{}
	if err := native.VerifyRepository(context.Background(), owner, app.Repository{Kind: "local", Path: owner.Root}); err != nil {
		t.Fatal(err)
	}
	root, cleanup, first, err := native.localBuildTree(context.Background(), owner)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	for _, name := range []string{".env", "ignored", "debug.log", "node_modules", "linked", ".langdrift"} {
		if _, err := os.Stat(filepath.Join(root, name)); !os.IsNotExist(err) {
			t.Errorf("excluded source was copied: %s", name)
		}
	}
	if data, err := os.ReadFile(filepath.Join(root, "app.txt")); err != nil || string(data) != "first version" {
		t.Fatalf("source missing: %s %v", data, err)
	}
	if err := os.WriteFile(filepath.Join(owner.Root, "app.txt"), []byte("local edit"), 0o600); err != nil {
		t.Fatal(err)
	}
	_, cleanupNext, second, err := native.localBuildTree(context.Background(), owner)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanupNext()
	if first == second || len(first) != 64 {
		t.Fatal("snapshot fingerprint does not identify source changes")
	}
	if err := (Store{}).Bind(context.Background(), owner, app.Connection{RuntimeID: "local", RuntimeUID: "local-uid", Credential: "synthetic-recovery", Repository: app.Repository{Kind: "local"}}); err != nil {
		t.Fatalf("non-Git credential storage: %v", err)
	}
}

func TestLocalSourceRejectsUnsupportedIgnorePatternsBeforeSnapshot(t *testing.T) {
	for _, pattern := range []string{"secrets/[ab].json", `\#private`} {
		t.Run(pattern, func(t *testing.T) {
			owner := ownerForTest(t)
			if err := os.WriteFile(filepath.Join(owner.Root, ".gitignore"), []byte(pattern), 0o600); err != nil {
				t.Fatal(err)
			}
			root, _, _, err := (Native{}).localBuildTree(context.Background(), owner)
			if err == nil || root != "" {
				t.Fatal("unsupported ignore pattern did not fail closed")
			}
		})
	}
}
