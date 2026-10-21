package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/exphere/ex/src/domain/cluster"
	"github.com/exphere/ex/src/domain/microruntime"
	"github.com/exphere/ex/src/domain/project"
)

// pollInterval paces every wait for the cluster to reach a state.
const pollInterval = time.Second

// lifecycleSources loads what both up and down act on.
type lifecycleSources struct {
	configuration RuntimeConfiguration
	guard         clusterGuard
	ledger        ExecutionLedger
	evidence      ExecutionEvidence
}

// lifecycle is the shared context of one locked up or down. Once its
// execution record is on disk, every exit path settles it and records evidence.
type lifecycle struct {
	owner      project.Project
	descriptor microruntime.Descriptor
	snapshot   Snapshot
	cluster    cluster.Record
	target     Target
	ledger     ExecutionLedger
	evidence   ExecutionEvidence
	execution  *microruntime.Execution
}

// open loads a valid configuration and the ready cluster it must target.
func (sources lifecycleSources) open(
	ctx context.Context,
	owner project.Project,
	runtime microruntime.Identifier,
) (lifecycle, error) {
	descriptor, snapshot, err := loadConfiguration(sources.configuration, owner, runtime)
	if err != nil {
		return lifecycle{}, err
	}

	record, err := sources.guard.ready(ctx, owner)
	if err != nil {
		return lifecycle{}, err
	}

	if descriptor.Target.Context != record.Context {
		return lifecycle{}, fmt.Errorf(
			"%w: runtime context must be the project cluster's context %q",
			ErrTargetMismatch, record.Context,
		)
	}

	return lifecycle{
		owner:      owner,
		descriptor: descriptor,
		snapshot:   snapshot,
		cluster:    record,
		target:     Target{ClusterAccess: accessOf(record), Namespace: descriptor.Target.Namespace.String()},
		ledger:     sources.ledger,
		evidence:   sources.evidence,
	}, nil
}

func (run *lifecycle) save() error {
	run.execution.UpdatedAt = time.Now().UTC()

	return run.ledger.Save(run.owner, run.execution)
}

// conclude settles the record and persists it with its evidence. Both writes
// are always attempted, and their errors are joined to the operation's own.
func (run *lifecycle) conclude(ctx context.Context, failure error) error {
	run.execution.Settle(failure, ctx.Err() != nil)

	return errors.Join(
		failure,
		run.save(),
		run.evidence.RecordExecution(run.owner, outcomeOf(run.execution)),
	)
}

func (run *lifecycle) namespace(ctx context.Context, reader ObjectReader) (microruntime.Object, error) {
	return reader.Get(ctx, run.target, "namespace", run.target.Namespace)
}

func (run *lifecycle) preparation() PreparationRequest {
	return PreparationRequest{
		Project:   run.owner,
		Runtime:   run.descriptor.ID,
		Directory: run.snapshot.PreparationDirectory,
		Access:    run.target.ClusterAccess,
		Namespace: run.target.Namespace,
		Owner:     run.descriptor.Owner,
		Token:     run.execution.Token,
	}
}

// observeContainers reports the containers of pods the execution's workloads
// control, and the UIDs of every owned object traversed to find them.
func observeContainers(
	ctx context.Context,
	reader ObjectReader,
	target Target,
	execution *microruntime.Execution,
) ([]microruntime.ContainerObservation, map[string]bool, error) {
	pods, err := reader.List(ctx, target, "pods")
	if err != nil {
		return nil, nil, err
	}

	replicaSets, err := reader.List(ctx, target, "replicasets")
	if err != nil {
		return nil, nil, err
	}

	observations, owned := execution.ObserveContainers(pods, replicaSets)

	return observations, owned, nil
}

// awaitRemoval waits until the object recorded under uid is gone.
func awaitRemoval(ctx context.Context, reader ObjectReader, target Target, kind, name, uid string) error {
	return poll(ctx, func() (bool, error) {
		live, err := reader.Get(ctx, target, kind, name)
		if err != nil {
			return false, err
		}

		return microruntime.Removed(live, uid)
	})
}

// poll calls check until it reports done, fails, or ctx ends.
func poll(ctx context.Context, check func() (done bool, err error)) error {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		done, err := check()
		if err != nil || done {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}
