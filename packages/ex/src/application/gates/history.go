// Package gates makes the CI Gates report of a pull request. It needs no
// project: it reads the repository's history and nothing under .langdrift/.
package gates

import (
	"context"

	"github.com/jonatassales/tao/packages/ex/src/domain/gates"
)

// History lists the commits a range holds, oldest first.
type History interface {
	Commits(ctx context.Context, directory, base, head string) ([]gates.Commit, error)
}
