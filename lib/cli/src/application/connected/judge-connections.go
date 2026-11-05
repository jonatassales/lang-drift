package connected

import (
	"context"

	"github.com/jonatassales/tao/lib/cli/src/domain/project"
)

// JudgeConnections keeps the console a project's planner asks, with its
// deployment secret, in the project's private state.
type JudgeConnections interface {
	Record(ctx context.Context, owner project.Project, consoleURL string, secret string) error
}
