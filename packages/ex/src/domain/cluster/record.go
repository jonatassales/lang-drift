// Package cluster is the domain model of the one local cluster a project owns.
//
// A cluster Record is an ownership ledger. It says which cluster EX created
// for the project and what existed in it right after creation, so EX can later
// prove that deleting the cluster only removes what EX itself created.
package cluster

import (
	"errors"
	"fmt"
	"time"

	"github.com/exphere/ex/src/domain/project"
)

// SystemNamespace is the namespace whose UID identifies a cluster.
const SystemNamespace = "kube-system"

// Status describes the last recorded operation, not live readiness.
type Status string

const (
	StatusCreating    Status = "creating"
	StatusReady       Status = "ready"
	StatusFailed      Status = "failed"
	StatusInterrupted Status = "interrupted"
	StatusDeleting    Status = "deleting"
	StatusDeleted     Status = "deleted"
)

// Domain errors classify why EX refused to act on a cluster.
var (
	ErrNotOwned          = errors.New("cluster record does not belong to the project")
	ErrUnresolved        = errors.New("cluster has an unresolved operation")
	ErrNotReady          = errors.New("cluster is not ready")
	ErrUntrustedIdentity = errors.New("cluster identity cannot be trusted")
	ErrUnsafeDeletion    = errors.New("cluster deletion is unsafe")
)

// ObjectIdentity is a persisted object's identity, without payload or secrets.
// Controllers holds the UIDs of owner references marked as controller.
type ObjectIdentity struct {
	UID         string
	Controllers []string
}

// Inventory is a complete observation of every persisted object in a cluster,
// keyed by resource type, namespace, and name.
type Inventory struct {
	Objects       map[string]ObjectIdentity
	ResourceTypes []string
}

// Baseline is what existed in the cluster right after EX created it.
type Baseline struct {
	Namespaces map[string]string
	Inventory  Inventory
	// InventoryVersion identifies the ownership evidence contract, not the
	// record format. Zero means no complete versioned proof was recorded.
	InventoryVersion int
}

// Record binds a cluster's observed identity to the project that owns it.
// Kubeconfig locates the private credentials EX uses to reach the cluster;
// the domain treats it as opaque.
type Record struct {
	ProjectID  project.ID
	Name       string
	Context    string
	Kubeconfig string
	UID        string
	Status     Status
	Baseline   Baseline
	UpdatedAt  time.Time
	Error      string
}

// Begin records the intent to create the project's cluster, before anything
// exists, so an interrupted creation is never forgotten.
func Begin(owner project.Project, kubeconfig string) Record {
	return Record{
		ProjectID:  owner.ID,
		Name:       owner.ClusterName(),
		Context:    "kind-" + owner.ClusterName(),
		Kubeconfig: kubeconfig,
		Status:     StatusCreating,
	}
}

// PlanCreation decides what creating the cluster means given the prior record,
// which is nil when none exists. It reports whether the existing ready cluster
// should be reused instead of creating a new one.
func PlanCreation(prior *Record) (reuse bool, err error) {
	switch {
	case prior == nil || prior.Status == StatusDeleted:
		return false, nil

	case prior.Status == StatusReady:
		return true, nil

	default:
		return false, fmt.Errorf(
			"%w: status %q; inspect it before recovery",
			ErrUnresolved, prior.Status,
		)
	}
}

// BelongsTo verifies that a persisted record was written for owner.
func (record *Record) BelongsTo(owner project.Project) error {
	if record.ProjectID != owner.ID || record.Name != owner.ClusterName() {
		return ErrNotOwned
	}
	if record.Context != "kind-"+record.Name && !(record.Context == "" && record.UID == "") {
		return ErrNotOwned
	}

	return nil
}

// RequireReady refuses to host runtimes on a cluster that is not ready.
func (record *Record) RequireReady() error {
	if record.Status != StatusReady {
		return fmt.Errorf("%w: %s", ErrNotReady, record.Status)
	}

	return nil
}

// Identify records the cluster's identity from its observed namespaces, which
// also become the deletion baseline.
func (record *Record) Identify(namespaces map[string]string) error {
	uid := namespaces[SystemNamespace]
	if uid == "" {
		return fmt.Errorf("%w: the cluster reported no %s UID", ErrUntrustedIdentity, SystemNamespace)
	}

	record.UID = uid
	record.Baseline.Namespaces = namespaces

	return nil
}

// Complete records the bootstrap inventory and marks the cluster ready. An
// empty inventory is refused, because a later deletion could not be proven safe.
func (record *Record) Complete(bootstrap Inventory) error {
	if len(bootstrap.Objects) == 0 || len(bootstrap.ResourceTypes) == 0 {
		return fmt.Errorf("%w: an empty bootstrap inventory cannot establish ownership", ErrUnsafeDeletion)
	}

	record.Baseline.Inventory = bootstrap
	record.Baseline.InventoryVersion = 1
	record.Status, record.Error = StatusReady, ""

	return nil
}

// Fail records why the last operation stopped.
func (record *Record) Fail(cause error, interrupted bool) {
	record.Status = StatusFailed
	if interrupted {
		record.Status = StatusInterrupted
	}

	record.Error = cause.Error()
}

// MarkDeleting records the intent to delete before the cluster is touched.
func (record *Record) MarkDeleting() {
	record.Status, record.Error = StatusDeleting, ""
}

// MarkDeleted records that the cluster no longer exists.
func (record *Record) MarkDeleted() {
	record.Status, record.Error = StatusDeleted, ""
}

// VerifyIdentity refuses to operate on a cluster that replaced the recorded one.
func (record *Record) VerifyIdentity(namespaces map[string]string) error {
	switch {
	case record.UID == "":
		return fmt.Errorf(
			"%w: its UID was never recorded, so adopting or deleting it is unsafe",
			ErrUntrustedIdentity,
		)

	case namespaces[SystemNamespace] != record.UID:
		return fmt.Errorf(
			"%w: the cluster was replaced; refusing to operate on it",
			ErrUntrustedIdentity,
		)
	}

	return nil
}
