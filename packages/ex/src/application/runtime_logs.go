package application

import (
	"context"
	"io"
)

// ReadRuntimeLogs reads one recorded workload while holding the project lock,
// so cleanup cannot replace its target during the foreground operation.
type ReadRuntimeLogs struct {
	Projects      ProjectOpener
	Locks         ProjectLocker
	Configuration RuntimeReferences
	Clusters      ClusterReader
	Observer      ClusterObserver
	Executions    ExecutionReader
	Objects       ObjectReader
	Logs          LogStreamer
}

// ReadRuntimeLogsRequest selects one workload and its output streams. Follow
// remains active until the caller cancels; finite reads use the caller deadline.
type ReadRuntimeLogsRequest struct {
	RuntimeUnitRequest
	Follow      bool
	Output      io.Writer
	Diagnostics io.Writer
}

// Execute verifies the recorded target and waits for the native log operation.
func (useCase ReadRuntimeLogs) Execute(ctx context.Context, request ReadRuntimeLogsRequest) error {
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
	target, resource, _, err := operation.selectTarget(ctx, owner, runtime, unit, true)
	if err != nil {
		return err
	}

	return useCase.Logs.Logs(ctx, target, resource, request.Follow, request.Output, request.Diagnostics)
}
