package connected

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jonatassales/tao/exphere/cli/src/application"
	app "github.com/jonatassales/tao/exphere/cli/src/application/connected"
	"github.com/jonatassales/tao/exphere/cli/src/domain/connection"
	"github.com/jonatassales/tao/exphere/cli/src/infrastructure/storage"
)

// A bind that died between writing .env and saving the record must not block
// every later attempt: the orphaned .env is the interrupted bind's own file.
func TestStore_BindRepairsAnOrphanedEnv(t *testing.T) {
	owner := ownerForTest(t)
	store := Store{Runner: untrackedRunner{}}
	record := app.Connection{RuntimeUID: "runtime-uid", RuntimeID: "runtime", State: connection.StateStopped, Secret: "second-secret", Credential: "second-credential"}

	orphan := filepath.Join(owner.Root, ".langdrift/runtimes/runtime/.env")
	if err := os.MkdirAll(filepath.Dir(orphan), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(orphan, []byte("RUNTIME_SECRET=first-secret\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := store.Bind(context.Background(), owner, record); err != nil {
		t.Fatalf("orphan repair refused: %v", err)
	}
	env, err := os.ReadFile(orphan)
	if err != nil || !strings.Contains(string(env), "second-secret") {
		t.Fatalf("env = %q, %v; want the repaired credentials", env, err)
	}

	// With the record saved, the runtime is connected and a replacement is refused.
	if err := store.Bind(context.Background(), owner, record); err == nil || !strings.Contains(err.Error(), "refusing to replace") {
		t.Fatalf("a connected runtime's .env was replaced: %v", err)
	}
}

// A materialization that died between file writes must be adopted by the next
// attempt, not reported as an unmanaged conflict forever.
func TestNative_MaterializeAdoptsItsInterruptedWrites(t *testing.T) {
	owner := ownerForTest(t)
	native := Native{Runner: &fakeRunner{}, Executions: storage.Executions{}}
	record := app.Connection{RuntimeID: "runtime", Repository: app.Repository{FullName: "owner/repository"}}

	first := nativeConfig()
	if _, _, err := native.Materialize(context.Background(), owner, record, first, first.Markdown()); err != nil {
		t.Fatal(err)
	}
	// The crash: files and intent exist, but the record never learned Managed.

	second := nativeConfig()
	second.Revision = 2
	second.Components[0].Port = 9090
	managed, _, err := native.Materialize(context.Background(), owner, record, second, second.Markdown())
	if err != nil {
		t.Fatalf("the retry refused its own interrupted writes: %v", err)
	}
	if len(managed) == 0 {
		t.Fatal("the retry reported nothing managed")
	}
}

// A genuinely foreign file keeps failing loudly even with an intent recorded.
func TestNative_MaterializeStillRefusesForeignFiles(t *testing.T) {
	owner := ownerForTest(t)
	native := Native{Runner: &fakeRunner{}, Executions: storage.Executions{}}
	record := app.Connection{RuntimeID: "runtime", Repository: app.Repository{FullName: "owner/repository"}}

	config := nativeConfig()
	if _, _, err := native.Materialize(context.Background(), owner, record, config, config.Markdown()); err != nil {
		t.Fatal(err)
	}

	foreign := filepath.Join(owner.Root, ".langdrift/runtimes/runtime/runtime.md")
	if err := os.WriteFile(foreign, []byte("hand edited\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	changed := nativeConfig()
	changed.Revision = 2
	if _, _, err := native.Materialize(context.Background(), owner, record, changed, changed.Markdown()); err == nil || !strings.Contains(err.Error(), "conflict") {
		t.Fatalf("a hand-edited file was overwritten: %v", err)
	}
}

func TestStore_LockSerializesAndRejectsTraversal(t *testing.T) {
	owner := ownerForTest(t)
	store := Store{}

	if _, err := store.Lock(owner, "../escape"); err == nil {
		t.Fatal("a traversal lock id was accepted")
	}

	unlock, err := store.Lock(owner, "runtime-uid")
	if err != nil {
		t.Fatal(err)
	}
	defer unlock()
	if _, err := store.Lock(owner, "runtime-uid"); err == nil {
		t.Fatal("the same lock was taken twice")
	}
}

func TestNative_StopWithoutARecordIsANoOp(t *testing.T) {
	owner := ownerForTest(t)
	native := Native{Runner: &fakeRunner{}, Executions: storage.Executions{}, StopRuntime: application.StopRuntime{
		Projects:   storage.Projects{},
		Locks:      storage.Locker{},
		Clusters:   storage.Clusters{},
		Executions: storage.Executions{},
	}}

	if err := native.Stop(context.Background(), owner, app.Connection{RuntimeID: "runtime"}); err != nil {
		t.Fatalf("stop without an execution record must reconcile to nothing: %v", err)
	}
}

func TestNative_MigrateWritesOneDraftAndPreservesIt(t *testing.T) {
	owner := ownerForTest(t)
	native := Native{Runner: &fakeRunner{}, Executions: storage.Executions{}}
	record := app.Connection{RuntimeID: "runtime", Repository: app.Repository{FullName: "owner/repository"}}
	config := nativeConfig()
	if _, _, err := native.Materialize(context.Background(), owner, record, config, config.Markdown()); err != nil {
		t.Fatal(err)
	}

	path, err := native.Migrate(context.Background(), owner, "runtime")
	if err != nil {
		t.Fatal(err)
	}
	draft, err := os.ReadFile(path)
	if err != nil || !strings.Contains(string(draft), "Import this existing JSON v2 runtime ID") {
		t.Fatalf("draft = %q, %v", draft, err)
	}

	if _, err := native.Migrate(context.Background(), owner, "runtime"); err == nil || !strings.Contains(err.Error(), "preserving your draft") {
		t.Fatalf("an existing draft was replaced: %v", err)
	}
}
