package connected

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/jonatassales/tao/lib/cli/src/application/plan"
	"github.com/jonatassales/tao/lib/cli/src/domain/project"
	"github.com/jonatassales/tao/lib/cli/src/infrastructure/files"
)

// consoleConnectionPath is the private file that holds the console a project
// is connected to: its URL and the deployment secret, never a JSON record.
const consoleConnectionPath = credentialsDirectory + "/console.env"

// consoleConnectionLimit bounds what is read of that file: a URL and a secret.
const consoleConnectionLimit = 8 << 10

// JudgeConnections keeps the console a project's planner asks in the
// project's private state, in a file the state directory's ignore file keeps
// out of Git and that Git must not track.
type JudgeConnections struct{ Runner runner }

// Record keeps the console the project is connected to, replacing an earlier
// one, so connecting again rotates the secret.
func (connections JudgeConnections) Record(ctx context.Context, owner project.Project, consoleURL string, secret string) error {
	if strings.ContainsAny(consoleURL+secret, "\n\r\x00") {
		return errors.New("the console URL or RUNTIME_SECRET holds a control character")
	}

	if err := refuseTracked(ctx, connections.Runner, owner.Root, consoleConnectionPath); err != nil {
		return err
	}

	path, err := files.Within(owner.Root, consoleConnectionPath)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}

	content := "CONSOLE_URL=" + consoleURL + "\nRUNTIME_SECRET=" + secret + "\n"

	return atomicWrite(path, []byte(content), 0o600)
}

// Read reads the console the project is connected to. A project that never
// connected reads as no connection.
func (JudgeConnections) Read(owner project.Project) (plan.Connection, error) {
	path, err := files.Within(owner.Root, consoleConnectionPath)
	if err != nil {
		return plan.Connection{}, err
	}

	content, err := files.ReadBounded(path, consoleConnectionLimit)
	if errors.Is(err, fs.ErrNotExist) {
		return plan.Connection{}, nil
	}

	if err != nil {
		return plan.Connection{}, err
	}

	return consoleConnectionOf(string(content))
}
