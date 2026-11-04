package inputs_test

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jonatassales/tao/lib/cli/src/application"
	app "github.com/jonatassales/tao/lib/cli/src/application/inputs"
	"github.com/jonatassales/tao/lib/cli/src/domain/harness"
	"github.com/jonatassales/tao/lib/cli/src/domain/project"
	"github.com/jonatassales/tao/lib/cli/src/infrastructure/inputs"
)

func tempProject(t *testing.T) project.Project {
	t.Helper()

	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	return project.Project{Root: root}
}

func TestStore_KeepsOneRecordPerKindPrivately(t *testing.T) {
	t.Parallel()

	owner := tempProject(t)
	store := inputs.Store{}

	if _, err := store.Load(owner, "envelope"); !errors.Is(err, application.ErrNotFound) {
		t.Fatalf("Load() of an unset kind = %v, want ErrNotFound", err)
	}

	first := app.Record{Kind: "envelope", ID: "envelope-7c2a", SetBy: "owner", Via: "cli", SetAt: harness.At(time.Date(2026, 10, 27, 18, 4, 0, 0, time.UTC)), Value: json.RawMessage(`{"cpuMillis":4000}`)}
	if err := store.Save(owner, first); err != nil {
		t.Fatal(err)
	}

	second := first
	second.ID, second.Previous = "envelope-9d4e", &first
	if err := store.Save(owner, second); err != nil {
		t.Fatal(err)
	}

	loaded, err := store.Load(owner, "envelope")
	if err != nil {
		t.Fatal(err)
	}

	var value struct {
		CPUMillis int `json:"cpuMillis"`
	}
	if err := json.Unmarshal(loaded.Value, &value); err != nil {
		t.Fatal(err)
	}

	if loaded.ID != "envelope-9d4e" || loaded.Previous == nil || loaded.Previous.ID != "envelope-7c2a" || value.CPUMillis != 4000 {
		t.Fatalf("Load() = %+v, want the latest record with its previous one", loaded)
	}

	path := filepath.Join(owner.Root, ".langdrift", "agent", "state", "inputs", "envelope.json")
	info, err := os.Stat(path)
	if err != nil || info.Mode().Perm() != 0o600 {
		t.Fatalf("record mode = %v, %v; want 0600", info.Mode().Perm(), err)
	}
}

func TestStore_RefusesAKindThatIsNotAWord(t *testing.T) {
	t.Parallel()

	if _, err := (inputs.Store{}).Load(tempProject(t), "../escape"); err == nil {
		t.Fatal("Load() accepted a path as a kind")
	}
}
