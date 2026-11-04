package plan

import (
	"strings"

	"github.com/jonatassales/tao/packages/ex/src/domain/run"
)

// dockerfilesCheck names each component's Dockerfile, and fails naming the
// ones that are not there.
func dockerfilesCheck(components []resolved) run.Check {
	var found, missing []string
	for _, component := range components {
		if !component.dockerfileFound {
			missing = append(missing, component.dockerfile)

			continue
		}

		found = append(found, component.dockerfile)
	}

	if len(missing) > 0 {
		return run.Check{
			ID:     "dockerfiles-exist",
			Result: failed,
			Detail: text("missing " + strings.Join(missing, " · ")),
		}
	}

	return run.Check{
		ID:     "dockerfiles-exist",
		Result: passed,
		Detail: text(strings.Join(found, " · ")),
	}
}
