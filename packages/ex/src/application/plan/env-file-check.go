package plan

import (
	"strings"

	"github.com/jonatassales/tao/packages/ex/src/domain/run"
)

// envFileCheck names the env files the components read, and fails naming
// one that is absent, malformed or overlaps the pass-through names.
func envFileCheck(components []resolved) run.Check {
	var named, refused []string
	for _, component := range components {
		if component.input.EnvFile == "" {
			continue
		}

		if component.envErr != nil {
			refused = append(refused, component.input.EnvFile+": "+component.envErr.Error())

			continue
		}

		named = append(named, component.input.EnvFile)
	}

	if len(refused) > 0 {
		return run.Check{
			ID:     "env-file-exists",
			Result: failed,
			Detail: text(strings.Join(refused, " · ")),
		}
	}

	if len(named) == 0 {
		return run.Check{
			ID:     "env-file-exists",
			Result: passed,
			Detail: text("no env file named"),
		}
	}

	return run.Check{
		ID:     "env-file-exists",
		Result: passed,
		Detail: text(strings.Join(named, " · ")),
	}
}
