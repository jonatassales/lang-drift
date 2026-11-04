package application

import (
	"fmt"

	"github.com/jonatassales/tao/lib/cli/src/domain/microruntime"
)

// addEnvironments puts each unit's Secret first among the manifests, so it
// is claimed, created and recorded before the workload that reads it, and
// cleanup removes it with the rest of what the execution owns.
func (run *lifecycle) addEnvironments(environments []microruntime.Environment) error {
	secrets := make([]microruntime.Manifest, 0, len(environments))

	for _, environment := range environments {
		if _, declared := run.descriptor.Units[environment.Unit]; !declared {
			return fmt.Errorf("%w: the environment names unit %s, which the runtime does not declare", ErrConfigurationChanged, environment.Unit)
		}

		secrets = append(secrets, environment.Manifest())
	}

	run.snapshot.Manifests = append(secrets, run.snapshot.Manifests...)

	return nil
}
