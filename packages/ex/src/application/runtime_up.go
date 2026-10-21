package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/exphere/ex/src/domain/microruntime"
	"github.com/exphere/ex/src/domain/project"
)

// StartRuntime validates a runtime, prepares its namespace, creates its
// workloads, and waits for each to become ready. Partial effects stay
// recorded for inspection and recovery.
type StartRuntime struct {
	Projects      ProjectOpener
	Locks         ProjectLocker
	Configuration RuntimeConfiguration
	Clusters      ClusterReader
	Observer      ClusterObserver
	Images        ImageChecker
	Executions    ExecutionLedger
	Objects       WorkloadCreator
	Preparation   Preparer
	Evidence      ExecutionEvidence
	Tools         ToolVersions
}

// Execute checks everything that can be checked without side effects, then
// records the obligations and runs the startup.
func (useCase StartRuntime) Execute(ctx context.Context, request RuntimeRequest) (ExecutionOutcome, error) {
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

	run, err := useCase.sources().open(ctx, owner, runtime)
	if err != nil {
		return ExecutionOutcome{}, err
	}

	if err := useCase.requireImages(ctx, run); err != nil {
		return ExecutionOutcome{}, err
	}

	if err := useCase.requireNoObligations(owner, runtime); err != nil {
		return ExecutionOutcome{}, err
	}

	if err := useCase.requireFreeNamespace(ctx, run); err != nil {
		return ExecutionOutcome{}, err
	}

	startup := &startup{lifecycle: run, objects: useCase.Objects, preparer: useCase.Preparation}

	return startup.run(ctx, useCase.Tools.Versions(ctx))
}

func (useCase StartRuntime) sources() lifecycleSources {
	return lifecycleSources{
		configuration: useCase.Configuration,
		guard:         clusterGuard{records: useCase.Clusters, observer: useCase.Observer},
		ledger:        useCase.Executions,
		evidence:      useCase.Evidence,
	}
}

func (useCase StartRuntime) requireImages(ctx context.Context, run lifecycle) error {
	images := microruntime.Images(run.snapshot.Manifests)

	if err := useCase.Images.RequireImages(ctx, run.cluster.Name, images); err != nil {
		return fmt.Errorf("load prebuilt images into the project cluster before up: %w", err)
	}

	return nil
}

// requireNoObligations refuses to start while an earlier up still owns anything.
func (useCase StartRuntime) requireNoObligations(owner project.Project, runtime microruntime.Identifier) error {
	previous, err := loadExecution(useCase.Executions, owner, runtime)

	switch {
	case errors.Is(err, ErrNotFound):
		return nil
	case err != nil:
		return err
	}

	if err := previous.RequireClosed(); err != nil {
		return fmt.Errorf("%w; inspect it and run down before a new up", err)
	}

	return nil
}

func (useCase StartRuntime) requireFreeNamespace(ctx context.Context, run lifecycle) error {
	namespace, err := run.namespace(ctx, useCase.Objects)
	if err != nil {
		return err
	}

	if namespace != nil {
		return fmt.Errorf("%w: namespace %s already exists", ErrAdoptionRefused, run.target.Namespace)
	}

	return nil
}

// startup is one locked up.
type startup struct {
	lifecycle
	objects  WorkloadCreator
	preparer Preparer
}

// run records the obligations, then performs the phases. Once the record is
// on disk, the run always settles it and writes the evidence, whatever the
// phases returned.
func (run *startup) run(ctx context.Context, tools map[string]microruntime.ToolVersion) (ExecutionOutcome, error) {
	if err := run.begin(tools); err != nil {
		return ExecutionOutcome{}, err
	}

	failure := run.execute(ctx)
	if failure != nil {
		run.keepObservedImages(ctx)
	}

	if err := run.conclude(ctx, failure); err != nil {
		return ExecutionOutcome{}, err
	}

	return outcomeOf(run.execution), nil
}

// begin claims every manifest and writes the obligations before anything is created.
func (run *startup) begin(tools map[string]microruntime.ToolVersion) error {
	token, err := newID()
	if err != nil {
		return err
	}

	runID, err := newID()
	if err != nil {
		return err
	}

	run.execution = microruntime.StartExecution(microruntime.Start{
		ProjectID:  run.owner.ID,
		Runtime:    run.descriptor.ID,
		Namespace:  run.descriptor.Target.Namespace,
		ClusterUID: run.cluster.UID,
		Context:    run.cluster.Context,
		ConfigHash: run.snapshot.Hash,
		Token:      token,
		RunID:      runID,
		Tools:      tools,
		Manifests:  run.snapshot.Manifests,
	})

	return run.save()
}

// execute runs the phases in order and stops at the first failure.
func (run *startup) execute(ctx context.Context) error {
	phases := []func(context.Context) error{
		run.prepareNamespace,
		run.validateObjects,
		run.createObjects,
		run.awaitWorkloads,
		run.observeImages,
	}

	for _, phase := range phases {
		if err := phase(ctx); err != nil {
			return err
		}
	}

	return run.execution.CompleteUp()
}

// prepareNamespace creates the namespace and quota, then records the
// namespace UID that later cleanup is conditioned on.
func (run *startup) prepareNamespace(ctx context.Context) error {
	if err := run.preparer.Apply(ctx, run.preparation()); err != nil {
		return err
	}

	namespace, err := run.namespace(ctx, run.objects)
	if err != nil {
		return err
	}

	if err := run.execution.RecordNamespace(namespace); err != nil {
		return err
	}

	return run.save()
}

// validateObjects dry-runs every object against the live namespace, so a bad
// manifest fails before any workload exists.
func (run *startup) validateObjects(ctx context.Context) error {
	for _, manifest := range run.snapshot.Manifests {
		if err := run.objects.DryRun(ctx, run.target, manifest.Object); err != nil {
			run.execution.MarkUnit(manifest.Unit, microruntime.UnitFailed)

			return err
		}
	}

	if err := run.execution.BeginApplying(); err != nil {
		return err
	}

	return run.save()
}

// createObjects creates the objects one at a time, recording each UID as it
// appears so an interrupted run can still be cleaned up precisely.
func (run *startup) createObjects(ctx context.Context) error {
	for index, manifest := range run.snapshot.Manifests {
		if err := run.createObject(ctx, index, manifest); err != nil {
			return err
		}
	}

	return nil
}

func (run *startup) createObject(ctx context.Context, index int, manifest microruntime.Manifest) error {
	existing, err := run.objects.Get(ctx, run.target, manifest.Object.Kind(), manifest.Object.Name())
	if err != nil {
		return err
	}

	if existing != nil {
		return fmt.Errorf("%w: %s already exists", ErrAdoptionRefused, manifest.Key())
	}

	run.execution.MarkUnit(manifest.Unit, microruntime.UnitRunning)
	if err := run.save(); err != nil {
		return err
	}

	created, err := run.objects.Create(ctx, run.target, manifest.Object)
	if err != nil {
		run.execution.MarkUnit(manifest.Unit, microruntime.UnitFailed)

		return err
	}

	if err := run.execution.RecordCreated(index, created); err != nil {
		return err
	}

	return run.save()
}

// awaitWorkloads waits for every recorded workload and marks a unit completed
// once its last resource is ready.
func (run *startup) awaitWorkloads(ctx context.Context) error {
	progress := run.execution.Progress()

	for _, resource := range run.execution.Resources {
		if resource.IsWorkload() {
			if err := run.awaitReady(ctx, resource); err != nil {
				run.execution.MarkUnit(resource.Unit, microruntime.UnitFailed)

				return err
			}
		}

		if progress.Finish(resource.Unit) {
			run.execution.MarkUnit(resource.Unit, microruntime.UnitCompleted)

			if err := run.save(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (run *startup) awaitReady(ctx context.Context, resource microruntime.Resource) error {
	return poll(ctx, func() (bool, error) {
		live, err := run.objects.Get(ctx, run.target, resource.Kind, resource.Name)
		if err != nil {
			return false, err
		}

		return resource.Readiness(live)
	})
}

func (run *startup) observeImages(ctx context.Context) error {
	observed, _, err := observeContainers(ctx, run.objects, run.target, run.execution)
	if err != nil {
		return fmt.Errorf("workloads became ready but image observation failed: %w", err)
	}

	run.execution.Images = microruntime.MergeContainers(run.execution.Images, observed)

	return nil
}

// keepObservedImages enriches a failed run's record on a best-effort basis.
// It never hides the original outcome or fabricates image IDs.
func (run *startup) keepObservedImages(ctx context.Context) {
	if ctx.Err() != nil || run.execution.NamespaceUID == "" {
		return
	}

	if observed, _, err := observeContainers(ctx, run.objects, run.target, run.execution); err == nil {
		run.execution.Images = microruntime.MergeContainers(run.execution.Images, observed)
	}
}
