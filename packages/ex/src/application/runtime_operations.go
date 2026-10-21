package application

import (
	"context"
	"fmt"
	"io"

	"github.com/exphere/ex/src/domain/microruntime"
	"github.com/exphere/ex/src/domain/project"
)

// RuntimeReferences reads the descriptor and checks its local references.
// Diagnostic operations do not require reparsing native configuration that
// the consumer may currently be repairing.
type RuntimeReferences interface {
	Load(owner project.Project, runtime microruntime.Identifier) (microruntime.Descriptor, error)
	CheckReferences(owner project.Project, descriptor microruntime.Descriptor) error
}

// LogStreamer streams a verified workload's logs until completion or
// cancellation. Non-following reads must bound the native log output.
type LogStreamer interface {
	Logs(ctx context.Context, target Target, resource microruntime.Resource, follow bool, output, diagnostics io.Writer) error
}

// PortForwarder owns a foreground localhost-only tunnel. Returning means the
// native process and its children have stopped.
type PortForwarder interface {
	Forward(ctx context.Context, target Target, resource microruntime.Resource, localPort, remotePort int, diagnostics io.Writer) error
}

// RuntimeUnitRequest selects one unit in an existing runtime.
type RuntimeUnitRequest struct {
	RuntimeRequest
	Unit string
}

// recordedOperation resolves a single target for diagnostics. The caller
// holds the project lock for selection and the entire foreground operation.
type recordedOperation struct {
	configuration RuntimeReferences
	clusters      ClusterReader
	observer      ClusterObserver
	executions    ExecutionReader
	objects       ObjectReader
}

func loadRuntimeReferences(configuration RuntimeReferences, owner project.Project, runtime microruntime.Identifier) (microruntime.Descriptor, error) {
	descriptor, err := configuration.Load(owner, runtime)
	if err != nil {
		return microruntime.Descriptor{}, err
	}

	if descriptor.ID != runtime {
		return microruntime.Descriptor{}, fmt.Errorf("%w: descriptor ID must match the runtime directory %s", microruntime.ErrInvalidDescriptor, runtime)
	}

	if err := descriptor.Validate(); err != nil {
		return microruntime.Descriptor{}, err
	}

	if err := configuration.CheckReferences(owner, descriptor); err != nil {
		return microruntime.Descriptor{}, err
	}

	return descriptor, nil
}

func (operation recordedOperation) selectTarget(ctx context.Context, owner project.Project, runtime, unit microruntime.Identifier, workload bool) (Target, microruntime.Resource, microruntime.Object, error) {
	execution, err := loadExecution(operation.executions, owner, runtime)
	if err != nil {
		return Target{}, microruntime.Resource{}, nil, err
	}

	descriptor, err := loadRuntimeReferences(operation.configuration, owner, runtime)
	if err != nil {
		return Target{}, microruntime.Resource{}, nil, err
	}

	record, err := (clusterGuard{records: operation.clusters, observer: operation.observer}).ready(ctx, owner)
	if err != nil {
		return Target{}, microruntime.Resource{}, nil, err
	}

	if record.Context != descriptor.Target.Context ||
		record.Context != execution.Context ||
		record.UID != execution.ClusterUID ||
		descriptor.Target.Namespace != execution.Namespace {
		return Target{}, microruntime.Resource{}, nil, fmt.Errorf("%w: target differs from recorded runtime", ErrTargetMismatch)
	}

	resource, err := execution.SelectResource(unit, workload)
	if err != nil {
		return Target{}, microruntime.Resource{}, nil, err
	}

	target := Target{ClusterAccess: accessOf(record), Namespace: execution.Namespace.String()}
	live, err := operation.objects.Get(ctx, target, resource.Kind, resource.Name)
	if err != nil {
		return Target{}, microruntime.Resource{}, nil, err
	}

	if live == nil || !live.IsOwnedBy(resource.UID, execution.Token) {
		return Target{}, microruntime.Resource{}, nil, fmt.Errorf("%w: recorded resource is absent or has changed ownership", microruntime.ErrOwnershipChanged)
	}

	return target, resource, live, nil
}

func parseUnit(text string) (microruntime.Identifier, error) {
	unit, err := microruntime.ParseIdentifier(text)
	if err != nil {
		return microruntime.Identifier{}, fmt.Errorf("unit: %w", err)
	}

	return unit, nil
}
