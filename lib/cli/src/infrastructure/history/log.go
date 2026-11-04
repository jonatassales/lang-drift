// Package history reads a repository's commits through git, as its own
// process, and changes nothing.
package history

import (
	"context"
	"strings"

	"github.com/jonatassales/tao/lib/cli/src/domain/gates"
	"github.com/jonatassales/tao/lib/cli/src/infrastructure/process"
)

// Runner is the direct native execution contract the log consumes.
type Runner interface {
	Run(context.Context, process.Command) (process.Result, error)
}

// Log lists commits with git log.
type Log struct {
	Runner Runner
}

// Commits lists the commits of base..head, oldest first, as their short SHA
// and subject.
func (log Log) Commits(ctx context.Context, directory, base, head string) ([]gates.Commit, error) {
	command := process.Command{
		Name: "git",
		Args: []string{"log", "--reverse", "--format=%h%x09%s", base + ".." + head, "--"},
		Dir:  directory,
	}

	result, err := log.Runner.Run(ctx, command)
	if err != nil {
		return nil, process.Explain("read the commit range "+base+".."+head, result, err)
	}

	var commits []gates.Commit

	for _, line := range strings.Split(strings.TrimSpace(result.Stdout), "\n") {
		sha, subject, found := strings.Cut(line, "\t")
		if !found {
			continue
		}

		commits = append(commits, gates.Commit{
			SHA:     sha,
			Subject: subject,
		})
	}

	return commits, nil
}
