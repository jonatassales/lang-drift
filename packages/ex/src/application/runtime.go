package application

import (
	"fmt"

	"github.com/exphere/ex/src/domain/microruntime"
	"github.com/exphere/ex/src/domain/project"
)

// RuntimeRequest names one runtime of a project.
type RuntimeRequest struct {
	Project ProjectSelection
	Runtime string
}

func parseRuntime(text string) (microruntime.Identifier, error) {
	runtime, err := microruntime.ParseIdentifier(text)
	if err != nil {
		return microruntime.Identifier{}, fmt.Errorf("runtime ID: %w", err)
	}

	return runtime, nil
}

// loadConfiguration is the single definition of a valid runtime, shared by
// validate, up, and down: an authored descriptor that satisfies the domain
// rules, whose native files exist and satisfy the manifest policy.
func loadConfiguration(
	configuration RuntimeConfiguration,
	owner project.Project,
	runtime microruntime.Identifier,
) (microruntime.Descriptor, Snapshot, error) {
	descriptor, err := configuration.Load(owner, runtime)
	if err != nil {
		return microruntime.Descriptor{}, Snapshot{}, err
	}

	if descriptor.ID != runtime {
		return microruntime.Descriptor{}, Snapshot{}, fmt.Errorf(
			"%w: descriptor ID must match the runtime directory %s",
			microruntime.ErrInvalidDescriptor, runtime,
		)
	}

	if err := descriptor.Validate(); err != nil {
		return microruntime.Descriptor{}, Snapshot{}, err
	}

	snapshot, err := configuration.Snapshot(owner, descriptor)
	if err != nil {
		return microruntime.Descriptor{}, Snapshot{}, err
	}

	if err := descriptor.ValidateManifests(snapshot.Manifests); err != nil {
		return microruntime.Descriptor{}, Snapshot{}, err
	}

	return descriptor, snapshot, nil
}

// loadExecution reads a runtime's ledger and refuses an inconsistent one,
// because cleanup acts on exactly what it records.
func loadExecution(
	reader ExecutionReader,
	owner project.Project,
	runtime microruntime.Identifier,
) (*microruntime.Execution, error) {
	execution, err := reader.Load(owner, runtime)
	if err != nil {
		return nil, err
	}

	if err := execution.Validate(owner.ID, runtime); err != nil {
		return nil, err
	}

	return execution, nil
}
