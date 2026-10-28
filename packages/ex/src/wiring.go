package main

import (
	"github.com/jonatassales/tao/packages/ex/src/application"
	connectedapp "github.com/jonatassales/tao/packages/ex/src/application/connected"
	"github.com/jonatassales/tao/packages/ex/src/application/inputs"
	"github.com/jonatassales/tao/packages/ex/src/application/plan"
	"github.com/jonatassales/tao/packages/ex/src/domain/rag"
	"github.com/jonatassales/tao/packages/ex/src/infrastructure/configuration"
	collectioninfra "github.com/jonatassales/tao/packages/ex/src/infrastructure/collection"
	connectedinfra "github.com/jonatassales/tao/packages/ex/src/infrastructure/connected"
	"github.com/jonatassales/tao/packages/ex/src/infrastructure/hooks"
	inputsinfra "github.com/jonatassales/tao/packages/ex/src/infrastructure/inputs"
	"github.com/jonatassales/tao/packages/ex/src/infrastructure/kind"
	"github.com/jonatassales/tao/packages/ex/src/infrastructure/kubectl"
	manifestinfra "github.com/jonatassales/tao/packages/ex/src/infrastructure/manifest"
	"github.com/jonatassales/tao/packages/ex/src/infrastructure/process"
	raginfra "github.com/jonatassales/tao/packages/ex/src/infrastructure/rag"
	"github.com/jonatassales/tao/packages/ex/src/infrastructure/records"
	"github.com/jonatassales/tao/packages/ex/src/infrastructure/scaffold"
	"github.com/jonatassales/tao/packages/ex/src/infrastructure/storage"
	"github.com/jonatassales/tao/packages/ex/src/infrastructure/terraform"
	"github.com/jonatassales/tao/packages/ex/src/infrastructure/toolchain"
	"github.com/jonatassales/tao/packages/ex/src/presentation/cli"
)

// dependencies are the sources a test scripts; nil means the real ones.
type dependencies struct {
	IDs      application.IDSource
	Clock    application.Clock
	Embedder rag.Embedder
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

	cases := cli.UseCases{
		ProjectInit: application.InitializeProject{
			Projects:  projects,
			Creator:   projects,
			Bootstrap: scaffold.LangDrift{},
			IDs:       deps.IDs,
		},
		Inputs: inputs.Service{
			Projects: projects,
			Store:    inputsinfra.Store{},
			Mirrors:  records.Mirrors{},
			IDs:      deps.IDs,
			Clock:    deps.Clock,
		},
		Plan: plan.Service{
			Projects:  projects,
			Collector: collectioninfra.Collector{},
			Inspector: manifestinfra.Inspector{},
			Inputs:    inputsinfra.Store{},
			Records:   records.Runtimes{},
			Index:     raginfra.Index{},
			Embedder:  deps.Embedder,
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
	}
	projectConsole := &connectedinfra.MCP{}
	cases.ProjectAgent = connectedapp.Service{Projects: projects, Store: connectedinfra.Store{Runner: runner}, Console: projectConsole, ProjectConsole: projectConsole, Native: connectedinfra.Native{
		Create: cases.ClusterCreate, Load: cases.ClusterLoadImage, StartRuntime: native.Up, StopRuntime: native.Down, Inspect: native.Inspect, Executions: executions, Runner: runner,
	}}

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
