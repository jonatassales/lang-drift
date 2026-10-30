package main

import (
	"time"

	"github.com/jonatassales/tao/packages/ex/src/application"
	connectedapp "github.com/jonatassales/tao/packages/ex/src/application/connected"
	"github.com/jonatassales/tao/packages/ex/src/application/execution"
	gatesapp "github.com/jonatassales/tao/packages/ex/src/application/gates"
	"github.com/jonatassales/tao/packages/ex/src/application/inputs"
	"github.com/jonatassales/tao/packages/ex/src/application/plan"
	"github.com/jonatassales/tao/packages/ex/src/application/setup"
	"github.com/jonatassales/tao/packages/ex/src/domain/rag"
	collectioninfra "github.com/jonatassales/tao/packages/ex/src/infrastructure/collection"
	"github.com/jonatassales/tao/packages/ex/src/infrastructure/configuration"
	connectedinfra "github.com/jonatassales/tao/packages/ex/src/infrastructure/connected"
	deliveryinfra "github.com/jonatassales/tao/packages/ex/src/infrastructure/delivery"
	"github.com/jonatassales/tao/packages/ex/src/infrastructure/history"
	"github.com/jonatassales/tao/packages/ex/src/infrastructure/hooks"
	inputsinfra "github.com/jonatassales/tao/packages/ex/src/infrastructure/inputs"
	judgeinfra "github.com/jonatassales/tao/packages/ex/src/infrastructure/judge"
	"github.com/jonatassales/tao/packages/ex/src/infrastructure/kind"
	"github.com/jonatassales/tao/packages/ex/src/infrastructure/kubectl"
	lifecycleinfra "github.com/jonatassales/tao/packages/ex/src/infrastructure/lifecycle"
	manifestinfra "github.com/jonatassales/tao/packages/ex/src/infrastructure/manifest"
	"github.com/jonatassales/tao/packages/ex/src/infrastructure/process"
	raginfra "github.com/jonatassales/tao/packages/ex/src/infrastructure/rag"
	"github.com/jonatassales/tao/packages/ex/src/infrastructure/records"
	"github.com/jonatassales/tao/packages/ex/src/infrastructure/scaffold"
	stagesinfra "github.com/jonatassales/tao/packages/ex/src/infrastructure/stages"
	"github.com/jonatassales/tao/packages/ex/src/infrastructure/storage"
	"github.com/jonatassales/tao/packages/ex/src/infrastructure/terraform"
	"github.com/jonatassales/tao/packages/ex/src/infrastructure/toolchain"
	"github.com/jonatassales/tao/packages/ex/src/presentation/cli"
)

// dependencies are the sources a test scripts; nil means the real ones.
type dependencies struct {
	IDs       application.IDSource
	Clock     application.Clock
	Embedder  rag.Embedder
	Judge     plan.Judge
	Driver    execution.StageDriver
	Lifecycle execution.Lifecycle
	Delivery  execution.DeliveryObserver
	// Poll is how often the checks are looked at; zero takes the default.
	Poll time.Duration
}

// useCases is the production composition.
func useCases() (cli.UseCases, lifecycles) { return useCasesWith(dependencies{}) }

// useCasesWith is the composition root: the only place that knows every
// concrete adapter. It builds each adapter once and hands every use case
// exactly the ports it declares.
func useCasesWith(deps dependencies) (cli.UseCases, lifecycles) {
	runner := process.Local{}

	var (
		projects    = storage.Projects{}
		locks       = storage.Locker{}
		clusters    = storage.Clusters{}
		executions  = storage.Executions{}
		evidence    = storage.Evidence{}
		files       = configuration.Files{}
		scaffolder  = configuration.Scaffolder{}
		agents      = scaffold.Agents{}
		kubernetes  = kubectl.Client{Runner: runner}
		provisioner = kind.Provisioner{Runner: runner}
		preparation = terraform.Preparation{Runner: runner}
		tools       = toolchain.Toolchain{Runner: runner}
		extensions  = hooks.Executor{Runner: runner}
	)

	native := lifecycles{
		Init: application.InitializeRuntime{
			Projects:   projects,
			Scaffolder: scaffolder,
		},
		Validate: application.ValidateRuntime{
			Projects:      projects,
			Configuration: files,
		},
		Up: application.StartRuntime{
			Projects:      projects,
			Locks:         locks,
			Configuration: files,
			Clusters:      clusters,
			Observer:      kubernetes,
			Images:        provisioner,
			Executions:    executions,
			Objects:       kubernetes,
			Preparation:   preparation,
			Evidence:      evidence,
			Tools:         tools,
		},
		Down: application.StopRuntime{
			Projects:      projects,
			Locks:         locks,
			Configuration: files,
			Clusters:      clusters,
			Observer:      kubernetes,
			Executions:    executions,
			Objects:       kubernetes,
			Preparation:   preparation,
			Evidence:      evidence,
			Tools:         tools,
		},
		Inspect: application.InspectRuntime{
			Projects:   projects,
			Clusters:   clusters,
			Observer:   kubernetes,
			Executions: executions,
			Objects:    kubernetes,
			Evidence:   evidence,
			Tools:      tools,
		},
		Logs: application.ReadRuntimeLogs{
			Projects: projects, Locks: locks, Configuration: files,
			Clusters: clusters, Observer: kubernetes, Executions: executions,
			Objects: kubernetes, Logs: kubernetes,
		},
		Forward: application.ForwardRuntime{
			Projects: projects, Locks: locks, Configuration: files,
			Clusters: clusters, Observer: kubernetes, Executions: executions,
			Objects: kubernetes, Forwarder: kubernetes,
		},
		Hook: application.RunRuntimeHook{
			Projects: projects, Locks: locks, Configuration: files,
			Scripts: files, Executor: extensions, Evidence: evidence,
		},
	}

	inputsService := inputs.Service{
		Projects: projects,
		Store:    inputsinfra.Store{},
		Mirrors:  records.Mirrors{},
		IDs:      deps.IDs,
		Clock:    deps.Clock,
	}

	cases := cli.UseCases{
		ProjectInit: setup.SetUpProject{
			Initialize: application.InitializeProject{
				Projects:  projects,
				Creator:   projects,
				Bootstrap: scaffold.LangDrift{},
				IDs:       deps.IDs,
			},
			Projects:  projects,
			Inputs:    inputsService,
			Store:     inputsinfra.Store{},
			Artifacts: storage.Artifacts{},
			Clock:     deps.Clock,
		},
		Inputs: inputsService,
		Plan: plan.Service{
			Projects:  projects,
			Collector: collectioninfra.Collector{},
			Inspector: manifestinfra.Inspector{},
			Inputs:    inputsinfra.Store{},
			Records:   records.Runtimes{},
			Index:     raginfra.Index{},
			Spend:     records.LapStore{},
			Embedder:  deps.Embedder,
			Judge:     judgeOf(deps),
			IDs:       deps.IDs,
			Clock:     deps.Clock,
		},
		Doctor: application.Diagnose{
			Prerequisites: tools,
			Projects:      projects,
			Clusters:      clusters,
			Observer:      kubernetes,
		},
		ClusterCreate: application.CreateCluster{
			Projects:    projects,
			Locks:       locks,
			Clusters:    clusters,
			Credentials: clusters,
			Provisioner: provisioner,
			Observer:    kubernetes,
		},
		ClusterDelete: application.DeleteCluster{
			Projects:    projects,
			Locks:       locks,
			Clusters:    clusters,
			Credentials: clusters,
			Executions:  executions,
			Provisioner: provisioner,
			Observer:    kubernetes,
		},
		ClusterInspect: application.InspectCluster{
			Projects: projects,
			Clusters: clusters,
			Observer: kubernetes,
		},
		ClusterLoadImage: application.LoadClusterImage{
			Projects: projects,
			Locks:    locks,
			Clusters: clusters,
			Observer: kubernetes,
			Images:   provisioner,
		},
		AgentInit: application.InitializeAgent{
			Projects: projects,
			Agents:   agents,
		},
		Gates: gatesapp.Report{
			History: history.Log{
				Runner: runner,
			},
		},
	}
	engine := connectedinfra.Native{
		Create: cases.ClusterCreate, Load: cases.ClusterLoadImage, StartRuntime: native.Up, StopRuntime: native.Down, Inspect: native.Inspect, Executions: executions, Runner: runner,
	}
	projectConsole := &connectedinfra.MCP{}
	cases.ProjectAgent = connectedapp.Service{Projects: projects, Store: connectedinfra.Store{Runner: runner}, Console: projectConsole, ProjectConsole: projectConsole, Native: engine}
	cases.Execution = execution.Service{
		Projects:   projects,
		Inputs:     inputsinfra.Store{},
		Records:    records.Runtimes{},
		Laps:       records.LapStore{},
		Driver:     driverOf(deps, runner),
		Lifecycle:  lifecycleOf(deps, engine, native, executions),
		Delivery:   deliveryOf(deps, runner),
		Index:      raginfra.Index{},
		Embedder:   deps.Embedder,
		ChecksPoll: pollOf(deps),
		IDs:        deps.IDs,
		Clock:      deps.Clock,
	}

	return cases, native
}

// lifecycles are the native runtime use cases. The connected agent is their
// only production consumer; tests drive them through executeNativeCatalogs.
type lifecycles struct {
	Init     application.InitializeRuntime
	Validate application.ValidateRuntime
	Up       application.StartRuntime
	Down     application.StopRuntime
	Inspect  application.InspectRuntime
	Logs     application.ReadRuntimeLogs
	Forward  application.ForwardRuntime
	Hook     application.RunRuntimeHook
}

// judgeOf is the planner's judge: the scripted one, or Claude.
func judgeOf(deps dependencies) plan.Judge {
	if deps.Judge != nil {
		return deps.Judge
	}

	return judgeinfra.Anthropic{}
}

// driverOf is the stages' driver: the scripted one, or the owner's coding
// agent run as a process.
func driverOf(deps dependencies, runner process.Local) execution.StageDriver {
	if deps.Driver != nil {
		return deps.Driver
	}

	return stagesinfra.Driver{Runner: runner}
}

// lifecycleOf is the run's lifecycle: the scripted one, or the native engine
// with a local operator record.
func lifecycleOf(deps dependencies, engine connectedinfra.Native, native lifecycles, executions storage.Executions) execution.Lifecycle {
	if deps.Lifecycle != nil {
		return deps.Lifecycle
	}

	return lifecycleinfra.Native{Engine: engine, Stop: native.Down, Inspect: native.Inspect, Executions: executions}
}

// deliveryOf is the delivery observer: the scripted one, or git and gh.
func deliveryOf(deps dependencies, runner process.Local) execution.DeliveryObserver {
	if deps.Delivery != nil {
		return deps.Delivery
	}

	return deliveryinfra.Observer{Runner: runner}
}

// defaultChecksPoll is how often the pull request's checks are looked at.
const defaultChecksPoll = 15 * time.Second

func pollOf(deps dependencies) time.Duration {
	if deps.Poll > 0 {
		return deps.Poll
	}

	return defaultChecksPoll
}
