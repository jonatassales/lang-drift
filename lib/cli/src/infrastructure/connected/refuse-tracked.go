package connected

import (
	"context"
	"errors"
	"fmt"

	"github.com/jonatassales/tao/lib/cli/src/infrastructure/process"
)

// refuseTracked refuses to write a private file Git tracks, which would put
// the secret it holds in the repository's history. A checkout outside Git
// has nothing to refuse.
func refuseTracked(ctx context.Context, git runner, root string, privatePath string) error {
	inGit, err := hasGitRoot(root)
	if err != nil {
		return err
	}

	if !inGit {
		return nil
	}

	if git == nil {
		return errors.New("private storage requires a Git index reader")
	}

	command := process.Command{
		Name: "git",
		Args: []string{"ls-files", "-z", "--", privatePath},
		Dir:  root,
	}

	tracked, err := git.Run(ctx, command)
	if err != nil {
		return fmt.Errorf("cannot verify the private path against the Git index: %w", err)
	}

	if tracked.Stdout != "" {
		return errors.New(privatePath + " is tracked by Git; remove it from the index before connecting")
	}

	return nil
}
