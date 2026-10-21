package application

import (
	"context"
	"fmt"
	"slices"

	"github.com/exphere/ex/src/domain/microruntime"
)

// StopRuntime removes a runtime's recorded workloads, then its namespace and
// preparation, re-verifying ownership before each destructive step. It never
// removes the project cluster, authored files, or evidence.
type StopRuntime struct {
	Projects      ProjectOpener
	Locks         ProjectLocker
	Configuration RuntimeConfiguration
	Clusters      ClusterReader
	Observer      ClusterObserver
	Executions    ExecutionLedger
	Objects       WorkloadRemover
	Preparation   PreparationCleaner
	Evidence      ExecutionEvidence
	Tools         ToolVersions
}

// StopRuntimeRequest names the runtime to stop.
type StopRuntimeRequest struct {
	RuntimeRequest

	// Recover permits cleaning up after a failed or interrupted operation,
	// once it has been inspected.
	Recover bool
}

// Execute verifies that the configuration and cluster still match what the
// recorded execution used, then cleans up.
func (useCase StopRuntime) Execute(ctx context.Context, request StopRuntimeRequest) (ExecutionOutcome, error) {
	runtime, err := parseRuntime(request.Runtime)
	if err != nil {
		return ExecutionOutcome{}, err
	}

	owner, err := useCase.Projects.Open(request.Project)
	if err != nil {
		return ExecutionOutcome{}, err
	}

	unlock, err := useCase.Locks.Lock(owner)
	if err != nil {
		return ExecutionOutcome{}, err
	}
	defer unlock()

	execution, err := loadExecution(useCase.Executions, owner, runtime)
	if err != nil {
		return ExecutionOutcome{}, err
	}

	if execution.Status == microruntime.StatusDown {
		return outcomeOf(execution), nil
	}

	if err := execution.AuthorizeDown(request.Recover); err != nil {
		return ExecutionOutcome{}, err
	}

	run, err := useCase.sources().open(ctx, owner, runtime)
	if err != nil {
		return ExecutionOutcome{}, err
	}

	if err := requireRecordedTarget(run, execution); err != nil {
		return ExecutionOutcome{}, err
	}

	run.execution = execution
	shutdown := &shutdown{lifecycle: run, objects: useCase.Objects, cleaner: useCase.Preparation}

	return shutdown.run(ctx, useCase.Tools.Versions(ctx))
}

func (useCase StopRuntime) sources() lifecycleSources {
	return lifecycleSources{
		configuration: useCase.Configuration,
		guard:         clusterGuard{records: useCase.Clusters, observer: useCase.Observer},
		ledger:        useCase.Executions,
		evidence:      useCase.Evidence,
	}
}

// requireRecordedTarget ties cleanup to the configuration and cluster that
// created the resources.
func requireRecordedTarget(run lifecycle, execution *microruntime.Execution) error {
	if run.snapshot.Hash != execution.ConfigHash && !slices.Contains(run.snapshot.CompatibleHashes, execution.ConfigHash) {
		return fmt.Errorf("%w; restore the recorded configuration before cleanup", ErrConfigurationChanged)
	}

	if run.cluster.UID != execution.ClusterUID ||
		run.target.Context != execution.Context ||
		run.target.Namespace != execution.Namespace.String() {
		return fmt.Errorf("%w: cleanup must target the cluster and namespace that execution used", ErrTargetMismatch)
	}

	return nil
}

// shutdown is one locked down.
type shutdown struct {
	lifecycle
	objects WorkloadRemover
	cleaner PreparationCleaner
}

// run records the start of the cleanup, then performs it. Once that is on
// disk, the run always settles the record and writes the evidence, whatever
// the cleanup returned.
func (run *shutdown) run(ctx context.Context, tools map[string]microruntime.ToolVersion) (ExecutionOutcome, error) {
	runID, err := newID()
	if err != nil {
		return ExecutionOutcome{}, err
	}

	run.execution.BeginDown(runID, tools)
	if err := run.save(); err != nil {
		return ExecutionOutcome{}, err
	}

	failure := run.execute(ctx)

	if err := run.conclude(ctx, failure); err != nil {
		return ExecutionOutcome{}, err
	}

	return outcomeOf(run.execution), nil
}

// execute removes the workloads, then the namespace, then the preparation.
func (run *shutdown) execute(ctx context.Context) error {
	namespace, err := run.namespace(ctx, run.objects)
	if err != nil {
		return err
	}

	if namespace != nil {
		if err := run.removeNamespaceContents(ctx, namespace); err != nil {
			return err
		}
	}

	if err := run.destroyPreparation(ctx); err != nil {
		return err
	}

	if err := run.requireNamespaceGone(ctx); err != nil {
		return err
	}

	run.execution.CompleteDown()

	return nil
}

// removeNamespaceContents deletes the recorded objects and then the
// namespace, refusing to act on anything the runtime cannot prove it owns.
func (run *shutdown) removeNamespaceContents(ctx context.Context, namespace microruntime.Object) error {
	if err := run.execution.ConfirmNamespace(namespace); err != nil {
		return err
	}

	if err := run.confirmResources(ctx); err != nil {
		return err
	}

	if err := run.save(); err != nil {
		return err
	}

	if err := run.authorizeCleanup(ctx, false); err != nil {
		return err
	}

	if err := run.deleteResources(ctx); err != nil {
		return err
	}

	if err := run.authorizeCleanup(ctx, true); err != nil {
		return err
	}

	return run.removeNamespace(ctx)
}

func (run *shutdown) confirmResources(ctx context.Context) error {
	for index, resource := range run.execution.Resources {
		live, err := run.objects.Get(ctx, run.target, resource.Kind, resource.Name)
		if err != nil {
			return err
		}

		if err := run.execution.ConfirmResource(index, live); err != nil {
			return err
		}
	}

	return nil
}

func (run *shutdown) authorizeCleanup(ctx context.Context, deleted bool) error {
	inventory, err := run.objects.NamespaceInventory(ctx, run.target)
	if err != nil {
		return err
	}

	return run.execution.AuthorizeCleanup(inventory, deleted)
}

// deleteResources removes the objects in recorded order and tracks progress
// per unit, so a partial cleanup can be inspected and resumed.
func (run *shutdown) deleteResources(ctx context.Context) error {
	progress := run.execution.Progress()

	for _, resource := range run.execution.Resources {
		run.execution.MarkUnit(resource.Unit, microruntime.UnitRunning)
		if err := run.save(); err != nil {
			return err
		}

		if err := run.deleteIfStillOwned(ctx, resource); err != nil {
			return err
		}

		if progress.Finish(resource.Unit) {
			run.execution.MarkUnit(resource.Unit, microruntime.UnitCompleted)
		}

		if err := run.save(); err != nil {
			return err
		}
	}

	return nil
}

// deleteIfStillOwned deletes the recorded object only if it is still the one
// recorded. An object that was never created has no UID and is skipped.
func (run *shutdown) deleteIfStillOwned(ctx context.Context, resource microruntime.Resource) error {
	if resource.UID == "" {
		return nil
	}

	live, err := run.objects.Get(ctx, run.target, resource.Kind, resource.Name)
	if err != nil || live == nil {
		return err
	}

	if !live.IsOwnedBy(resource.UID, run.execution.Token) {
		return fmt.Errorf("%w: %s %s was replaced before cleanup", microruntime.ErrOwnershipChanged, resource.Kind, resource.Name)
	}

	if err := run.objects.Delete(ctx, run.target, resource); err != nil {
		return err
	}

	return awaitRemoval(ctx, run.objects, run.target, resource.Kind, resource.Name, resource.UID)
}

// removeNamespace deletes the emptied namespace, provided it is still the one
// this runtime created.
func (run *shutdown) removeNamespace(ctx context.Context) error {
	live, err := run.namespace(ctx, run.objects)
	if err != nil || live == nil {
		return err
	}

	uid := run.execution.NamespaceUID
	if !live.IsOwnedBy(uid, run.execution.Token) {
		return fmt.Errorf("%w: namespace %s was replaced before cleanup", microruntime.ErrOwnershipChanged, run.target.Namespace)
	}

	if err := run.objects.DeleteNamespace(ctx, run.target, uid); err != nil {
		return err
	}

	return awaitRemoval(ctx, run.objects, run.target, "namespace", run.target.Namespace, uid)
}

// destroyPreparation lets the preparation tool remove what it created.
// Without preparation state there is nothing for it to destroy.
func (run *shutdown) destroyPreparation(ctx context.Context) error {
	prepared, err := run.cleaner.HasState(run.owner, run.descriptor.ID)
	if err != nil || !prepared {
		return err
	}

	return run.cleaner.Destroy(ctx, run.preparation())
}

func (run *shutdown) requireNamespaceGone(ctx context.Context) error {
	live, err := run.namespace(ctx, run.objects)
	if err != nil {
		return err
	}

	if live != nil {
		return fmt.Errorf("%w: namespace %s still exists after cleanup", microruntime.ErrCleanupBlocked, run.target.Namespace)
	}

	return nil
}
