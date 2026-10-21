package application

import (
	"context"

	"github.com/exphere/ex/src/domain/cluster"
	"github.com/exphere/ex/src/domain/microruntime"
	"github.com/exphere/ex/src/domain/project"
)

// The ports below are everything the application needs from the outside
// world. Each is as small as its consumers allow. Ports backed by local files
// take no context; ports that run native tools or reach a cluster do.

// --- Projects ---------------------------------------------------------------

// ProjectSelection says where to find the consumer project.
type ProjectSelection struct {
	// Directory is where the search starts.
	Directory string

	// Explicit requires Directory to be the project root itself. Otherwise the
	// nearest project at or above Directory is used.
	Explicit bool
}

// ProjectOpener opens an existing project, or returns ErrNotFound.
type ProjectOpener interface {
	Open(selection ProjectSelection) (project.Project, error)
}

// ProjectCreator records a new project's identity in directory. It returns
// ErrAlreadyExists rather than replacing existing metadata.
type ProjectCreator interface {
	Create(directory string, id project.ID) (project.Project, error)
}

// ProjectLocker serializes every change to a project's cluster and runtimes,
// across processes.
type ProjectLocker interface {
	Lock(owner project.Project) (unlock func(), err error)
}

// --- Consumer-authored configuration ------------------------------------------

// RuntimeConfiguration reads the configuration a consumer authored for a runtime.
type RuntimeConfiguration interface {
	// Load reads the runtime's descriptor exactly as authored, or returns ErrNotFound.
	Load(owner project.Project, runtime microruntime.Identifier) (microruntime.Descriptor, error)

	// Snapshot reads the native files a valid descriptor references and
	// fingerprints them together with the descriptor.
	Snapshot(owner project.Project, descriptor microruntime.Descriptor) (Snapshot, error)
}

// Snapshot is the native configuration one execution runs against. Hash ties
// cleanup to the exact configuration that created the resources.
type Snapshot struct {
	Manifests            []microruntime.Manifest
	PreparationDirectory string
	Hash                 string
	// CompatibleHashes contains fully recomputed fingerprints from earlier
	// encodings of the same configuration. It never contains recorded values.
	CompatibleHashes []string
}

// RuntimeScaffolder publishes a new runtime configuration directory. It never
// writes into an existing one.
type RuntimeScaffolder interface {
	// Scaffold writes a descriptor for the consumer to complete by hand.
	Scaffold(owner project.Project, descriptor microruntime.Descriptor) (location string, err error)

	// Generate writes complete native files from a validated blueprint.
	Generate(owner project.Project, runtime microruntime.Identifier, blueprint microruntime.Blueprint) (location string, err error)
}

// AgentScaffolder publishes an editable agent project. It never builds,
// installs, or runs anything.
type AgentScaffolder interface {
	Scaffold(owner project.Project, name microruntime.Identifier) (location string, err error)
}

// --- Cluster ----------------------------------------------------------------

// ClusterAccess is how to reach a cluster: explicit credentials and context,
// never the user's ambient configuration.
type ClusterAccess struct {
	Kubeconfig string
	Context    string
}

// ClusterReader reads the project's cluster record, or returns ErrNotFound.
type ClusterReader interface {
	Load(owner project.Project) (cluster.Record, error)
}

// ClusterWriter also persists the record.
type ClusterWriter interface {
	ClusterReader
	Save(owner project.Project, record cluster.Record) error
}

// ClusterCredentials manages the private kubeconfig reserved for the project's cluster.
type ClusterCredentials interface {
	Credentials(owner project.Project) (location string, exists bool, err error)
	DiscardCredentials(owner project.Project) error
}

// ClusterProvisioner creates and deletes local clusters by name.
type ClusterProvisioner interface {
	Clusters(ctx context.Context) ([]string, error)
	Create(ctx context.Context, name, kubeconfig string) (contextName string, err error)
	Delete(ctx context.Context, name string) error
}

// ImageLoader copies an already-built local image into a cluster.
type ImageLoader interface {
	LoadImage(ctx context.Context, clusterName, image string) error
}

// ImageChecker fails unless every image is present on every node of a cluster.
type ImageChecker interface {
	RequireImages(ctx context.Context, clusterName string, images []string) error
}

// ClusterObserver observes a cluster's identity and contents without changing it.
type ClusterObserver interface {
	Namespaces(ctx context.Context, access ClusterAccess) (map[string]string, error)
	ClusterInventory(ctx context.Context, access ClusterAccess) (cluster.Inventory, error)
}

// --- Runtime execution ------------------------------------------------------

// ExecutionReader reads a runtime's ledger, or returns ErrNotFound.
type ExecutionReader interface {
	Load(owner project.Project, runtime microruntime.Identifier) (*microruntime.Execution, error)
}

// ExecutionLedger also persists the ledger after every recorded step.
type ExecutionLedger interface {
	ExecutionReader
	Save(owner project.Project, execution *microruntime.Execution) error
}

// ExecutionLister lists the ledgers of every runtime in a project.
type ExecutionLister interface {
	List(owner project.Project) ([]*microruntime.Execution, error)
}

// Target is one runtime namespace on the project cluster.
type Target struct {
	ClusterAccess
	Namespace string
}

// ObjectReader reads namespaced objects. Get returns nil when the object is absent.
type ObjectReader interface {
	Get(ctx context.Context, target Target, kind, name string) (microruntime.Object, error)
	List(ctx context.Context, target Target, resource string) ([]microruntime.Object, error)
}

// ObjectWriter validates objects against the live API and creates them.
// Create fails rather than adopting an object that already exists.
type ObjectWriter interface {
	DryRun(ctx context.Context, target Target, object microruntime.Object) error
	Create(ctx context.Context, target Target, object microruntime.Object) (microruntime.Object, error)
}

// ObjectDeleter requests deletion of exactly the recorded object, guarded by
// its UID. It returns once the API server accepts the request.
type ObjectDeleter interface {
	Delete(ctx context.Context, target Target, resource microruntime.Resource) error
	DeleteNamespace(ctx context.Context, target Target, uid string) error
}

// InventoryReader lists every listable object in the target namespace.
// A partial inventory is an error, never a shorter list.
type InventoryReader interface {
	NamespaceInventory(ctx context.Context, target Target) ([]microruntime.Object, error)
}

// UsageReader reads optional pod metrics.
type UsageReader interface {
	PodMetrics(ctx context.Context, target Target) ([]microruntime.Object, error)
}

// WorkloadCreator is what starting a runtime needs from the cluster API.
type WorkloadCreator interface {
	ObjectReader
	ObjectWriter
}

// WorkloadRemover is what stopping a runtime needs from the cluster API.
type WorkloadRemover interface {
	ObjectReader
	ObjectDeleter
	InventoryReader
}

// WorkloadObserver is what inspecting a runtime needs from the cluster API.
type WorkloadObserver interface {
	ObjectReader
	UsageReader
}

// PreparationRequest describes the namespace and quota preparation of one
// runtime execution.
type PreparationRequest struct {
	Project   project.Project
	Runtime   microruntime.Identifier
	Directory string
	Access    ClusterAccess
	Namespace string
	Owner     string
	Token     string
}

// Preparer creates a runtime's namespace and quota.
type Preparer interface {
	Apply(ctx context.Context, request PreparationRequest) error
}

// PreparationCleaner destroys what a Preparer created, when it created anything.
type PreparationCleaner interface {
	HasState(owner project.Project, runtime microruntime.Identifier) (bool, error)
	Destroy(ctx context.Context, request PreparationRequest) error
}

// ExecutionEvidence keeps a token-free record of every lifecycle operation.
type ExecutionEvidence interface {
	RecordExecution(owner project.Project, outcome ExecutionOutcome) error
}

// InspectionEvidence keeps every inspection report.
type InspectionEvidence interface {
	RecordInspection(owner project.Project, report Inspection) error
}

// --- Local tools --------------------------------------------------------------

// ToolVersions observes the versions of the native tools an execution uses.
type ToolVersions interface {
	Versions(ctx context.Context) map[string]microruntime.ToolVersion
}

// Prerequisites checks the native tools EX relies on. It never installs anything.
type Prerequisites interface {
	Check(ctx context.Context) PrerequisiteReport
}

// Tool states reported by Prerequisites.
const (
	ToolAvailable   = "available"
	ToolMissing     = "missing"
	ToolUnavailable = "unavailable"
)

// ToolStatus describes one native tool on the host.
type ToolStatus struct {
	Name    string `json:"name"`
	Path    string `json:"path,omitempty"`
	Version string `json:"version,omitempty"`
	Status  string `json:"status"`
	Install string `json:"install,omitempty"`
}

// PrerequisiteReport is the result of checking the host.
type PrerequisiteReport struct {
	Tools                    []ToolStatus
	ContainerEngineReachable bool
}
