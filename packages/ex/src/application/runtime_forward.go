package application

import (
	"context"
	"fmt"
	"io"
)

// ForwardRuntime owns a foreground tunnel to one recorded Service. The
// project lock and caller cancellation cover its complete lifetime.
type ForwardRuntime struct {
	Projects      ProjectOpener
	Locks         ProjectLocker
	Configuration RuntimeReferences
	Clusters      ClusterReader
	Observer      ClusterObserver
	Executions    ExecutionReader
	Objects       ObjectReader
	Forwarder     PortForwarder
}

// ForwardRuntimeRequest selects local and remote ports. A zero remote port
// selects the Service's sole port; ambiguous Services require an explicit one.
type ForwardRuntimeRequest struct {
	RuntimeUnitRequest
	LocalPort   int
	RemotePort  int
	Diagnostics io.Writer
}

// Execute verifies ownership and waits until the foreground tunnel ends.
func (useCase ForwardRuntime) Execute(ctx context.Context, request ForwardRuntimeRequest) error {
	if request.LocalPort < 1 || request.LocalPort > 65535 || request.RemotePort < 0 || request.RemotePort > 65535 {
		return fmt.Errorf("ports must be between 1 and 65535")
	}
	runtime, err := parseRuntime(request.Runtime)
	if err != nil {
		return err
	}
	unit, err := parseUnit(request.Unit)
	if err != nil {
		return err
	}
	owner, err := useCase.Projects.Open(request.Project)
	if err != nil {
		return err
	}
	unlock, err := useCase.Locks.Lock(owner)
	if err != nil {
		return err
	}
	defer unlock()

	operation := recordedOperation{
		configuration: useCase.Configuration,
		clusters:      useCase.Clusters, observer: useCase.Observer,
		executions: useCase.Executions, objects: useCase.Objects,
	}
	target, resource, live, err := operation.selectTarget(ctx, owner, runtime, unit, false)
	if err != nil {
		return err
	}
	remotePort, err := live.ForwardPort(request.RemotePort)
	if err != nil {
		return err
	}

	return useCase.Forwarder.Forward(ctx, target, resource, request.LocalPort, remotePort, request.Diagnostics)
}
