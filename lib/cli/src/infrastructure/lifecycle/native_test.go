package lifecycle_test

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"github.com/jonatassales/tao/lib/cli/src/application"
	connectedapp "github.com/jonatassales/tao/lib/cli/src/application/connected"
	"github.com/jonatassales/tao/lib/cli/src/domain/connection"
	"github.com/jonatassales/tao/lib/cli/src/domain/microruntime"
	"github.com/jonatassales/tao/lib/cli/src/domain/project"
	"github.com/jonatassales/tao/lib/cli/src/infrastructure/files"
	"github.com/jonatassales/tao/lib/cli/src/infrastructure/lifecycle"
	"github.com/jonatassales/tao/lib/core/contracts/go/configuration"
)

// engine is a native engine that keeps the records it was handed.
type engine struct {
	materialized connectedapp.Connection
	started      connectedapp.Connection
	starts       int
	usage        *connectedapp.Observation
}

func (engine *engine) Materialize(_ context.Context, _ project.Project, record connectedapp.Connection, _ connection.Configuration, _ string) (map[string]string, string, error) {
	engine.materialized = record

	return map[string]string{"runtime.md": "a1", "schema.json": "b2"}, "hash-1", nil
}

func (engine *engine) Start(_ context.Context, _ project.Project, record connectedapp.Connection, _ connection.Configuration, report func(string, string) error) error {
	engine.started = record
	engine.starts++
	_ = report("build_started", "rev")
	_ = report("build_succeeded", "build_completed")

	return nil
}

func (engine *engine) Observe(context.Context, project.Project, connectedapp.Connection) (*connectedapp.Observation, error) {
	return engine.usage, nil
}

type stopper struct {
	outcome application.ExecutionOutcome
	stops   int
}

func (stopper *stopper) Execute(context.Context, application.StopRuntimeRequest) (application.ExecutionOutcome, error) {
	stopper.stops++

	return stopper.outcome, nil
}

type inspector struct{ inspection application.Inspection }

func (inspector inspector) Execute(context.Context, application.RuntimeRequest) (application.Inspection, error) {
	return inspector.inspection, nil
}

type reader struct{ execution *microruntime.Execution }

func (reader reader) Load(project.Project, microruntime.Identifier) (*microruntime.Execution, error) {
	return reader.execution, nil
}

func owner(t *testing.T) (project.Project, microruntime.Identifier) {
	t.Helper()

	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	id, err := project.ParseID("a759267ef752f98d6c25c4114ad65d0c")
	if err != nil {
		t.Fatal(err)
	}
	runtime, err := microruntime.ParseIdentifier("console-microruntime")
	if err != nil {
		t.Fatal(err)
	}

	return project.Project{ID: id, Root: root}, runtime
}

func closed() *microruntime.Execution {
	return &microruntime.Execution{RunID: "exec-5c1e9a", Settled: true, Status: microruntime.StatusDown}
}

func TestUpMaterializesAndStartsAsTheLocalOperator(t *testing.T) {
	t.Parallel()

	project, runtime := owner(t)
	engine := &engine{}
	native := lifecycle.Native{Engine: engine, Stop: &stopper{}, Inspect: inspector{}, Executions: reader{execution: closed()}}
	config := configuration.Configuration{RuntimeID: "console-microruntime", RuntimeUID: "701706be7e70d681d608826d3ad548b9", Revision: 1}

	var steps []string
	report := func(step, note string) { steps = append(steps, step+": "+note) }

	execution, err := native.Up(context.Background(), project, runtime, config, []byte("# runtime\n"), report)
	if err != nil {
		t.Fatal(err)
	}

	if execution.ID != "exec-5c1e9a" {
		t.Errorf("execution = %+v", execution)
	}
	record := engine.materialized
	if record.Repository.Kind != "local" || record.Repository.Path != project.Root || record.RuntimeUID != config.RuntimeUID || record.RuntimeID != "console-microruntime" || record.ProjectID != project.ID.String() || len(record.Managed) != 0 {
		t.Errorf("the first materialization got %+v, want a fresh local operator record", record)
	}
	if engine.started.ConcreteHash != "hash-1" || engine.started.Managed["schema.json"] != "b2" {
		t.Errorf("the start got %+v, want what the materialization fixed", engine.started)
	}
	if !slices.Equal(steps, []string{"Configured: `runtime.md` revision 1 materialized", "Built: images built from the components' Dockerfiles"}) {
		t.Errorf("steps = %v", steps)
	}

	var kept struct {
		Managed      map[string]string `json:"managed"`
		ConcreteHash string            `json:"concreteHash"`
	}
	if err := files.ReadJSON(filepath.Join(project.Root, ".langdrift", "agent", "state", "runs", "console-microruntime", "operator.json"), &kept); err != nil || kept.ConcreteHash != "hash-1" || kept.Managed["runtime.md"] != "a1" {
		t.Errorf("operator record = %+v, %v", kept, err)
	}

	if _, err := native.Up(context.Background(), project, runtime, config, []byte("# runtime\n"), report); err != nil {
		t.Fatal(err)
	}
	if engine.materialized.Managed["runtime.md"] != "a1" || engine.materialized.ConcreteHash != "hash-1" {
		t.Errorf("the second materialization got %+v, want the kept operator record", engine.materialized)
	}
}

func TestUpStopsAnExecutionAnEarlierLapLeftOpen(t *testing.T) {
	t.Parallel()

	project, runtime := owner(t)
	stop := &stopper{outcome: application.ExecutionOutcome{Settled: true, Status: "down"}}
	open := &microruntime.Execution{RunID: "exec-old", Settled: false, Status: microruntime.StatusRunning}
	native := lifecycle.Native{Engine: &engine{}, Stop: stop, Inspect: inspector{}, Executions: reader{execution: open}}

	if _, err := native.Up(context.Background(), project, runtime, configuration.Configuration{RuntimeID: "console-microruntime"}, nil, func(string, string) {}); err != nil {
		t.Fatal(err)
	}
	if stop.stops != 1 {
		t.Errorf("the open execution was stopped %d times, want once before materializing", stop.stops)
	}
}

func TestObserveReadsReadinessAndUsage(t *testing.T) {
	t.Parallel()

	project, runtime := owner(t)
	cpu, memory := 0.41, 690.4
	at := time.Date(2026, 10, 28, 14, 31, 2, 0, time.UTC)
	engine := &engine{usage: &connectedapp.Observation{CPUUsed: &cpu, MemoryUsedMiB: &memory, ObservedAt: at, MeasurementWindows: []connectedapp.MeasurementWindow{{Pod: "web", ObservedAt: at, Window: "30s"}}}}
	inspection := application.Inspection{
		Execution:  application.ExecutionOutcome{Status: "running"},
		Resources:  []application.ResourceObservation{{Resource: application.ResourceResult{Kind: "Deployment"}, Observed: "present"}, {Resource: application.ResourceResult{Kind: "Service"}, Observed: "present"}},
		Containers: []application.ContainerResult{{Ready: true}},
	}
	native := lifecycle.Native{Engine: engine, Stop: &stopper{}, Inspect: inspector{inspection: inspection}, Executions: reader{execution: closed()}}

	observation, err := native.Observe(context.Background(), project, runtime)
	if err != nil {
		t.Fatal(err)
	}
	if !observation.Ready || observation.Deployments != 1 || observation.Containers != 1 || observation.ReadyContainers != 1 || observation.CPUMillis == nil || *observation.CPUMillis != 410 || observation.MemoryMiB == nil || *observation.MemoryMiB != 690 || observation.Window != "30s" || !observation.ObservedAt.Equal(at) {
		t.Errorf("observation = %+v", observation)
	}

	inspection.Containers[0].Ready = false
	native.Inspect = inspector{inspection: inspection}
	if observation, _ := native.Observe(context.Background(), project, runtime); observation.Ready {
		t.Error("a container that is not ready read as ready")
	}
}

func TestDownCountsWhatTheCleanupRemoved(t *testing.T) {
	t.Parallel()

	project, runtime := owner(t)
	before := &microruntime.Execution{Resources: []microruntime.Resource{{Kind: "Namespace"}, {Kind: "Deployment"}, {Kind: "Service"}, {Kind: "ResourceQuota"}}}
	stop := &stopper{outcome: application.ExecutionOutcome{Settled: true, Status: "down"}}
	native := lifecycle.Native{Engine: &engine{}, Stop: stop, Inspect: inspector{}, Executions: reader{execution: before}}

	cleanup, err := native.Down(context.Background(), project, runtime)
	if err != nil || !cleanup.Confirmed || cleanup.Removed != 4 || cleanup.Total != 4 || !cleanup.ClusterKept {
		t.Errorf("Down() = %+v, %v", cleanup, err)
	}

	stop.outcome = application.ExecutionOutcome{Settled: false, Status: "stopping", Resources: []application.ResourceResult{{Kind: "Namespace"}, {Kind: "Deployment"}}}
	cleanup, err = native.Down(context.Background(), project, runtime)
	if err != nil || cleanup.Confirmed || cleanup.Removed != 2 || cleanup.Total != 4 {
		t.Errorf("an unsettled Down() = %+v, %v", cleanup, err)
	}

	if _, err := os.Stat(filepath.Join(project.Root, ".langdrift")); !os.IsNotExist(err) {
		t.Error("Down wrote under the project")
	}
}
