package microruntime

import (
	"time"

	"github.com/exphere/ex/src/domain/project"
)

// ExecutionStatus describes where the last operation on a runtime stands.
type ExecutionStatus string

const (
	StatusPreparing   ExecutionStatus = "preparing"
	StatusApplying    ExecutionStatus = "applying"
	StatusRunning     ExecutionStatus = "running"
	StatusCleaning    ExecutionStatus = "cleaning"
	StatusDown        ExecutionStatus = "down"
	StatusFailed      ExecutionStatus = "failed"
	StatusInterrupted ExecutionStatus = "interrupted"
)

// UnitStatus describes how far an operation got with one unit.
type UnitStatus string

const (
	UnitNotAttempted UnitStatus = "not-attempted"
	UnitRunning      UnitStatus = "running"
	UnitCompleted    UnitStatus = "completed"
	UnitFailed       UnitStatus = "failed"
)

// Operation names the lifecycle operation an execution record last ran.
type Operation string

const (
	OperationUp   Operation = "up"
	OperationDown Operation = "down"
)

// Execution is a runtime's ownership ledger: what an up created and what a
// later down must remove. It is recorded apart from the editable configuration.
// Settled means the operation's subprocesses have returned, not that its
// resources are healthy.
type Execution struct {
	ProjectID    project.ID
	Runtime      Identifier
	Status       ExecutionStatus
	Settled      bool
	Token        string
	ClusterUID   string
	Context      string
	Namespace    Identifier
	NamespaceUID string
	ConfigHash   string
	RunID        string
	Operation    Operation
	UpdatedAt    time.Time
	ToolVersions map[string]ToolVersion
	Images       []ContainerObservation
	Resources    []Resource
	Units        map[Identifier]UnitStatus
}

// Start describes a new up before anything exists in the cluster.
type Start struct {
	ProjectID  project.ID
	Runtime    Identifier
	Namespace  Identifier
	ClusterUID string
	Context    string
	ConfigHash string
	Token      string
	RunID      string
	Tools      map[string]ToolVersion
	Manifests  []Manifest
}

// StartExecution claims every manifest for a new up and records the
// obligations before any object is created.
func StartExecution(start Start) *Execution {
	execution := &Execution{
		ProjectID:    start.ProjectID,
		Runtime:      start.Runtime,
		Status:       StatusPreparing,
		Token:        start.Token,
		ClusterUID:   start.ClusterUID,
		Context:      start.Context,
		Namespace:    start.Namespace,
		ConfigHash:   start.ConfigHash,
		RunID:        start.RunID,
		Operation:    OperationUp,
		ToolVersions: start.Tools,
		Images:       ConfiguredContainers(start.Manifests),
		Units:        map[Identifier]UnitStatus{},
	}

	claim := Claim{
		Runtime:   start.Runtime,
		Namespace: start.Namespace,
		Token:     start.Token,
		RunID:     start.RunID,
	}

	for _, manifest := range start.Manifests {
		claim.Apply(manifest)
		execution.Resources = append(execution.Resources, ResourceOf(manifest.Object, manifest.Unit))
		execution.Units[manifest.Unit] = UnitNotAttempted
	}

	return execution
}

// Validate rejects a record whose identity or targets are inconsistent,
// because cleanup acts on exactly what is recorded.
func (execution *Execution) Validate(owner project.ID, runtime Identifier) error {
	if execution.ProjectID != owner ||
		execution.Runtime != runtime ||
		execution.Token == "" ||
		execution.Namespace.IsZero() ||
		execution.ClusterUID == "" {
		return reject(ErrInconsistentExecution, "runtime %s record identity is inconsistent", runtime)
	}

	for _, resource := range execution.Resources {
		if !resource.isSupported(execution.Namespace) {
			return reject(
				ErrInconsistentExecution,
				"runtime %s records an unsupported target %s/%s",
				runtime, resource.Kind, resource.Name,
			)
		}
	}

	return nil
}

// BeginApplying advances an active startup only after preparation recorded
// the namespace it owns. The caller must have validated every native object
// with a dry run.
func (execution *Execution) BeginApplying() error {
	if execution.Operation != OperationUp || execution.Status != StatusPreparing || execution.Settled {
		return reject(ErrInconsistentExecution, "runtime %s is not preparing an active startup", execution.Runtime)
	}

	if execution.NamespaceUID == "" {
		return reject(ErrInconsistentExecution, "runtime %s has no recorded namespace ownership", execution.Runtime)
	}

	execution.Status = StatusApplying

	return nil
}

// CompleteUp records successful startup only when every declared resource
// has a recorded UID and every unit completed. Readiness is observed by the
// caller; these checks prevent incomplete bookkeeping from claiming success.
func (execution *Execution) CompleteUp() error {
	if execution.Operation != OperationUp || execution.Status != StatusApplying || execution.Settled {
		return reject(ErrInconsistentExecution, "runtime %s is not applying an active startup", execution.Runtime)
	}

	if execution.NamespaceUID == "" || len(execution.Resources) == 0 || len(execution.Units) == 0 {
		return reject(ErrInconsistentExecution, "runtime %s startup has incomplete ownership records", execution.Runtime)
	}

	for _, resource := range execution.Resources {
		if resource.UID == "" || execution.Units[resource.Unit] != UnitCompleted {
			return reject(ErrInconsistentExecution, "runtime %s resource %s/%s has not completed startup", execution.Runtime, resource.Kind, resource.Name)
		}
	}

	for unit, status := range execution.Units {
		if status != UnitCompleted {
			return reject(ErrInconsistentExecution, "runtime %s unit %s has not completed startup", execution.Runtime, unit)
		}
	}

	execution.Status = StatusRunning

	return nil
}

// RequireClosed refuses to proceed while the runtime may still own anything
// or an operation on it never settled.
func (execution *Execution) RequireClosed() error {
	switch {
	case !execution.Settled:
		return reject(ErrExistingObligations, "runtime %s has an unresolved operation", execution.Runtime)

	case execution.Status != StatusDown:
		return reject(ErrExistingObligations, "runtime %s is %s", execution.Runtime, execution.Status)

	case execution.NamespaceUID != "" || len(execution.Resources) != 0:
		return reject(ErrExistingObligations, "runtime %s still records owned resources", execution.Runtime)
	}

	return nil
}

// AuthorizeDown decides whether a down may start from this record. A partial
// operation needs an explicit recovery after inspection.
func (execution *Execution) AuthorizeDown(recoverPartial bool) error {
	switch {
	case !execution.Settled:
		return reject(
			ErrCleanupBlocked,
			"previous operation did not settle; reconcile any surviving processes before cleanup",
		)

	case execution.Status != StatusRunning && !recoverPartial:
		return reject(
			ErrCleanupBlocked,
			"runtime %s is %s; a partial operation requires explicit recovery after inspection",
			execution.Runtime, execution.Status,
		)
	}

	return nil
}

// BeginDown starts a cleanup run: every unit is pending again until proven otherwise.
func (execution *Execution) BeginDown(runID string, tools map[string]ToolVersion) {
	execution.RunID = runID
	execution.Operation = OperationDown
	execution.ToolVersions = tools
	execution.Status = StatusCleaning
	execution.Settled = false
	execution.markAllUnits(UnitNotAttempted)
}

// CompleteDown records that the runtime no longer owns anything.
func (execution *Execution) CompleteDown() {
	execution.Resources = nil
	execution.NamespaceUID = ""
	execution.Status = StatusDown
	execution.markAllUnits(UnitCompleted)
}

// MarkUnit records how far the current operation got with one unit.
func (execution *Execution) MarkUnit(unit Identifier, status UnitStatus) {
	execution.Units[unit] = status
}

func (execution *Execution) markAllUnits(status UnitStatus) {
	for unit := range execution.Units {
		execution.Units[unit] = status
	}
}

// Settle records how an operation ended. It runs on every exit path, so the
// record always says whether the operation's subprocesses have returned.
func (execution *Execution) Settle(failure error, interrupted bool) {
	if failure != nil {
		execution.Status = StatusFailed
		if interrupted {
			execution.Status = StatusInterrupted
		}

		for unit, status := range execution.Units {
			if status == UnitRunning {
				execution.Units[unit] = UnitFailed
			}
		}
	}

	execution.Settled = true
}

// RecordNamespace records the namespace that preparation created; later
// cleanup is conditioned on its UID.
func (execution *Execution) RecordNamespace(live Object) error {
	if live == nil {
		return reject(ErrOwnershipChanged, "preparation did not create namespace %s", execution.Namespace)
	}

	execution.NamespaceUID = live.UID()

	if !live.IsOwnedBy(execution.NamespaceUID, execution.Token) {
		return reject(ErrOwnershipChanged, "namespace %s does not carry the runtime's ownership", execution.Namespace)
	}

	return nil
}

// ConfirmNamespace requires the live namespace to belong to this runtime. A
// partial up may never have recorded its UID; it is learned when the token
// matches.
func (execution *Execution) ConfirmNamespace(live Object) error {
	if execution.NamespaceUID == "" && live.IsOwnedBy(live.UID(), execution.Token) {
		execution.NamespaceUID = live.UID()
	}

	if !live.IsOwnedBy(execution.NamespaceUID, execution.Token) {
		return reject(ErrOwnershipChanged, "namespace %s ownership changed", execution.Namespace)
	}

	return nil
}

// RecordCreated records the identity the API server gave a created object.
func (execution *Execution) RecordCreated(index int, created Object) error {
	resource := ResourceOf(created, execution.Resources[index].Unit)
	execution.Resources[index] = resource

	if !created.IsOwnedBy(resource.UID, execution.Token) {
		return reject(ErrOwnershipChanged, "%s %s ownership was not observed", resource.Kind, resource.Name)
	}

	return nil
}

// ConfirmResource requires a recorded object that still exists to carry the
// runtime's ownership, learning a UID that a partial up never recorded.
func (execution *Execution) ConfirmResource(index int, live Object) error {
	if live == nil {
		return nil
	}

	resource := &execution.Resources[index]

	if resource.UID == "" {
		if !live.IsOwnedBy(live.UID(), execution.Token) {
			return reject(ErrOwnershipChanged, "unrecorded %s %s does not carry runtime ownership", resource.Kind, resource.Name)
		}

		resource.UID = live.UID()
	}

	if !live.IsOwnedBy(resource.UID, execution.Token) {
		return reject(ErrOwnershipChanged, "ownership changed for %s %s", resource.Kind, resource.Name)
	}

	return nil
}

// UnitProgress counts the resources each unit still has to process, so a unit
// completes when its last resource does.
type UnitProgress map[Identifier]int

// Progress starts counting over the recorded resources.
func (execution *Execution) Progress() UnitProgress {
	progress := UnitProgress{}
	for _, resource := range execution.Resources {
		progress[resource.Unit]++
	}

	return progress
}

// Finish records one processed resource and reports whether its unit is done.
func (progress UnitProgress) Finish(unit Identifier) bool {
	progress[unit]--

	return progress[unit] == 0
}
