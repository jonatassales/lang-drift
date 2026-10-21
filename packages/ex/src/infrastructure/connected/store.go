// Package connected implements the project's private persistence, MCP transport,
// and native lifecycle adapters. It never changes the shared cluster ownership.
package connected

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	app "github.com/exphere/ex/src/application/connected"
	"github.com/exphere/ex/src/domain/microruntime"
	"github.com/exphere/ex/src/domain/project"
	"github.com/exphere/ex/src/infrastructure/files"
	"github.com/exphere/ex/src/infrastructure/process"
)

// Store persists operational state separately from the runtime's ignored .env.
type Store struct{ Runner runner }

func statePath(owner project.Project, id string) (string, error) {
	if _, err := microruntime.ParseIdentifier(id); err != nil {
		return "", err
	}
	return files.Within(owner.Root, ".exphere/state/connections/"+id+".json")
}
func (Store) List(owner project.Project) ([]app.Connection, error) {
	directory, err := files.Within(owner.Root, ".exphere/state/connections")
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(directory)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var records []app.Connection
	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		var record app.Connection
		data, err := files.ReadBounded(filepath.Join(directory, entry.Name()), app.MaxConnectionStateBytes)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(data, &record); err != nil {
			return nil, err
		}
		if entry.Name() != record.RuntimeID+".json" {
			return nil, errors.New("connection identity differs from filename")
		}
		envPath, err := files.Existing(owner.Root, ".exphere/runtimes/"+record.RuntimeID+"/.env")
		if err != nil {
			return nil, err
		}
		data, err = files.ReadBounded(envPath, 16384)
		if err != nil {
			return nil, err
		}
		for _, line := range strings.Split(string(data), "\n") {
			key, value, ok := strings.Cut(line, "=")
			if !ok {
				continue
			}
			switch key {
			case "RUNTIME_SECRET":
				record.Secret = value
			case "AGENT_CREDENTIAL":
				record.Credential = value
			}
		}
		if record.Credential == "" {
			return nil, fmt.Errorf("runtime %s has no agent credential", record.RuntimeID)
		}
		records = append(records, record)
	}
	return records, nil
}
func (Store) Save(owner project.Project, record app.Connection) error {
	path, err := statePath(owner, record.RuntimeID)
	if err != nil {
		return err
	}
	data, err := app.MarshalConnection(record)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	return atomicWrite(path, data, 0o600)
}
func (store Store) Bind(ctx context.Context, owner project.Project, record app.Connection) error {
	if strings.ContainsAny(record.Credential+record.Secret, "\n\r\x00") {
		return errors.New("credential contains unsupported control characters")
	}
	privatePath := ".exphere/runtimes/" + record.RuntimeID + "/.env"
	if _, err := files.Within(owner.Root, privatePath); err != nil {
		return err
	}
	hasGit, err := hasGitRoot(owner.Root)
	if err != nil {
		return err
	}
	if hasGit || record.Repository.Kind != "local" {
		if store.Runner == nil {
			return errors.New("credential storage requires a Git index reader")
		}
		tracked, err := store.Runner.Run(ctx, process.Command{Name: "git", Args: []string{"ls-files", "-z", "--", privatePath}, Dir: owner.Root, UnsetEnv: []string{"GIT_DIR", "GIT_WORK_TREE", "GIT_INDEX_FILE", "GIT_COMMON_DIR", "GIT_CONFIG_COUNT", "GIT_NAMESPACE"}})
		if err != nil {
			return errors.New("cannot verify the private credential path against the Git index")
		}
		if tracked.Stdout != "" {
			return errors.New("runtime .env is tracked by Git; remove it from the index before connecting (existing index and files were preserved)")
		}
	}
	directory, err := files.Within(owner.Root, ".exphere/runtimes/"+record.RuntimeID)
	if err != nil {
		return err
	}
	if err = os.MkdirAll(directory, 0o755); err != nil {
		return err
	}
	env := filepath.Join(directory, ".env")
	if _, err = os.Lstat(env); err == nil {
		return errors.New("runtime .env already exists; refusing to replace private settings")
	} else if !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	ignore := filepath.Join(directory, ".gitignore")
	existing, err := files.ReadBounded(ignore, 1<<20)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	if !strings.Contains("\n"+string(existing)+"\n", "\n.env\n") {
		existing = append(existing, []byte("\n.env\n")...)
		if err = atomicWrite(ignore, existing, 0o644); err != nil {
			return err
		}
	}
	if err = atomicWrite(env, []byte("RUNTIME_SECRET="+record.Secret+"\nAGENT_CREDENTIAL="+record.Credential+"\n"), 0o600); err != nil {
		return err
	}
	return store.Save(owner, record)
}
func (Store) Lock(owner project.Project, id string) (func(), error) {
	if strings.ContainsAny(id, "/\\.") {
		return nil, errors.New("invalid runtime UID")
	}
	path, err := files.Within(owner.Root, ".exphere/state/connections/"+id+".lock")
	if err != nil {
		return nil, err
	}
	if err = os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil, err
	}
	return files.Lock(path)
}

func atomicWrite(path string, data []byte, mode fs.FileMode) error {
	if err := files.RequireNoSymlinks(path); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".connected-*")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if err = tmp.Chmod(mode); err == nil {
		_, err = tmp.Write(data)
	}
	if err == nil {
		err = tmp.Sync()
	}
	err = errors.Join(err, tmp.Close())
	if err != nil {
		return err
	}
	if err := os.Rename(tmp.Name(), path); err != nil {
		return err
	}
	directory, err := os.Open(filepath.Dir(path))
	if err != nil {
		return err
	}
	defer directory.Close()
	return directory.Sync()
}
