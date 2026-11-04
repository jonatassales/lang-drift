// Package lifecycle brings a planned runtime up inside a run and takes it
// down, through the native lifecycle the connected agent uses, with a local
// operator record instead of a console's: the run is the operator.
package lifecycle

import (
	"context"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"

	"github.com/jonatassales/tao/lib/cli/src/application"
	connectedapp "github.com/jonatassales/tao/lib/cli/src/application/connected"
	"github.com/jonatassales/tao/lib/cli/src/application/execution"
	"github.com/jonatassales/tao/lib/cli/src/domain/connection"
	"github.com/jonatassales/tao/lib/cli/src/domain/microruntime"
	"github.com/jonatassales/tao/lib/cli/src/domain/project"
	"github.com/jonatassales/tao/lib/cli/src/infrastructure/files"
	"github.com/jonatassales/tao/lib/core/contracts/go/configuration"
	"github.com/jonatassales/tao/lib/core/contracts/go/protocol"
)

// Engine is the native lifecycle: it materializes a configuration, starts
// it and observes what it uses.
type Engine interface {
	Materialize(ctx context.Context, owner project.Project, record connectedapp.Connection, config connection.Configuration, markdown string) (map[string]string, string, error)
	Start(ctx context.Context, owner project.Project, record connectedapp.Connection, config connection.Configuration, report func(string, string) error) error
	Observe(ctx context.Context, owner project.Project, record connectedapp.Connection) (*connectedapp.Observation, error)
}

// Stopper takes a runtime down and cleans up what it owned.
type Stopper interface {
	Execute(context.Context, application.StopRuntimeRequest) (application.ExecutionOutcome, error)
}

// Inspector reports the runtime's objects and containers.
type Inspector interface {
	Execute(context.Context, application.RuntimeRequest) (application.Inspection, error)
}

// Native is the run's lifecycle over the native engine.
type Native struct {
	Engine     Engine
	Stop       Stopper
	Inspect    Inspector
	Executions application.ExecutionReader
}

// operator is the run's own operator record, kept privately between laps:
// what the last materialization wrote and the configuration it fixed.
type operator struct {
	RuntimeUID   string            `json:"runtimeUid"`
	Managed      map[string]string `json:"managed"`
	ConcreteHash string            `json:"concreteHash"`
}

// Up materializes the runtime, builds and starts it, and reports the
// execution it claimed. An execution an earlier lap left open is cleaned
// up first.
func (native Native) Up(ctx context.Context, owner project.Project, runtime microruntime.Identifier, config configuration.Configuration, markdown []byte, report execution.Reporter) (execution.Execution, error) {
	record, err := native.record(owner, runtime, config)
	if err != nil {
		return execution.Execution{}, err
	}

	if err := native.ensureStopped(ctx, owner, runtime); err != nil {
		return execution.Execution{}, err
	}

	managed, hash, err := native.Engine.Materialize(ctx, owner, record, config, string(markdown))
	if err != nil {
		return execution.Execution{}, fmt.Errorf("materialize: %w", err)
	}

	record.Managed, record.ConcreteHash = managed, hash
	if err := native.save(owner, runtime, record); err != nil {
		return execution.Execution{}, err
	}
	report("Configured", fmt.Sprintf("`runtime.md` revision %d materialized", config.Revision))

	err = native.Engine.Start(ctx, owner, record, config, func(step, _ string) error {
		if step == "build_succeeded" {
			report("Built", "images built from the components' Dockerfiles")
		}

		return nil
	})
	if err != nil {
		return execution.Execution{}, fmt.Errorf("start: %w", err)
	}

	ledger, err := native.Executions.Load(owner, runtime)
	if err != nil {
		return execution.Execution{}, fmt.Errorf("the execution ledger after the start: %w", err)
	}

	return execution.Execution{ID: ledger.RunID}, nil
}

// Observe looks at the runtime: its Deployments and containers, and what it
// uses when the metrics are there.
func (native Native) Observe(ctx context.Context, owner project.Project, runtime microruntime.Identifier) (execution.Observation, error) {
	inspection, err := native.Inspect.Execute(ctx, request(owner, runtime))
	if err != nil {
		return execution.Observation{}, err
	}

	observation := execution.Observation{ObservedAt: inspection.ObservedAt}
	for _, resource := range inspection.Resources {
		if resource.Resource.Kind == "Deployment" && resource.Observed != "absent" {
			observation.Deployments++
		}
	}
	for _, container := range inspection.Containers {
		observation.Containers++
		if container.Ready {
			observation.ReadyContainers++
		}
	}
	observation.Ready = inspection.Execution.Status == string(microruntime.StatusRunning) && observation.Containers > 0 && observation.ReadyContainers == observation.Containers

	usage, err := native.Engine.Observe(ctx, owner, connectedapp.Connection{RuntimeID: runtime.String()})
	if err != nil || usage == nil {
		return observation, nil
	}

	if usage.CPUUsed != nil {
		millis := int(math.Round(*usage.CPUUsed * 1000))
		observation.CPUMillis = &millis
	}
	if usage.MemoryUsedMiB != nil {
		mib := int(math.Round(*usage.MemoryUsedMiB))
		observation.MemoryMiB = &mib
	}
	if len(usage.MeasurementWindows) > 0 {
		observation.Window = usage.MeasurementWindows[0].Window
		observation.ObservedAt = usage.ObservedAt
	}

	return observation, nil
}

// Down stops the runtime and says what the cleanup removed, when the
// execution ledger confirms it.
func (native Native) Down(ctx context.Context, owner project.Project, runtime microruntime.Identifier) (execution.Cleanup, error) {
	total := 0
	if before, err := native.Executions.Load(owner, runtime); err == nil {
		total = len(before.Resources)
	}

	outcome, err := native.Stop.Execute(ctx, application.StopRuntimeRequest{RuntimeRequest: request(owner, runtime), Recover: true})
	if err != nil {
		return execution.Cleanup{Total: total}, err
	}

	return execution.Cleanup{
		Confirmed:   outcome.Settled && outcome.Status == string(microruntime.StatusDown),
		Removed:     total - len(outcome.Resources),
		Total:       total,
		ClusterKept: true,
	}, nil
}

// ensureStopped cleans up an execution an earlier lap left open, so the
// materialization finds the runtime stopped.
func (native Native) ensureStopped(ctx context.Context, owner project.Project, runtime microruntime.Identifier) error {
	ledger, err := native.Executions.Load(owner, runtime)
	if errors.Is(err, application.ErrNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	if ledger.RequireClosed() == nil {
		return nil
	}

	if _, err := native.Stop.Execute(ctx, application.StopRuntimeRequest{RuntimeRequest: request(owner, runtime), Recover: true}); err != nil {
		return fmt.Errorf("an earlier execution of %s is open and could not be stopped: %w", runtime, err)
	}

	return nil
}

// record is the operator record the engine reads: the run itself, on the
// project's own folder, with what the last lap materialized.
func (native Native) record(owner project.Project, runtime microruntime.Identifier, config configuration.Configuration) (connectedapp.Connection, error) {
	record := connectedapp.Connection{
		RuntimeUID: config.RuntimeUID,
		RuntimeID:  runtime.String(),
		ProjectID:  owner.ID.String(),
		Repository: protocol.Repository{Kind: "local", Path: owner.Root},
		State:      connection.StateStopped,
		Managed:    map[string]string{},
		Outbox:     []connectedapp.Event{},
	}

	path, err := operatorPath(owner, runtime)
	if err != nil {
		return record, err
	}

	var kept operator
	err = files.ReadJSON(path, &kept)
	if errors.Is(err, os.ErrNotExist) {
		return record, nil
	}
	if err != nil {
		return record, fmt.Errorf("the operator record of %s is unreadable: %w", runtime, err)
	}

	if kept.Managed != nil {
		record.Managed = kept.Managed
	}
	record.ConcreteHash = kept.ConcreteHash

	return record, nil
}

func (native Native) save(owner project.Project, runtime microruntime.Identifier, record connectedapp.Connection) error {
	path, err := operatorPath(owner, runtime)
	if err != nil {
		return err
	}

	return files.WriteJSON(path, operator{RuntimeUID: record.RuntimeUID, Managed: record.Managed, ConcreteHash: record.ConcreteHash})
}

func operatorPath(owner project.Project, runtime microruntime.Identifier) (string, error) {
	return files.Within(owner.Root, filepath.Join(".langdrift", "agent", "state", "runs", runtime.String(), "operator.json"))
}

func request(owner project.Project, runtime microruntime.Identifier) application.RuntimeRequest {
	return application.RuntimeRequest{Project: application.ProjectSelection{Directory: owner.Root, Explicit: true}, Runtime: runtime.String()}
}
