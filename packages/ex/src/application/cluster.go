package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/exphere/ex/src/domain/cluster"
	"github.com/exphere/ex/src/domain/project"
)

// clusterOperationTimeout bounds creating, loading into, or deleting a cluster,
// independently of the caller's own deadline.
const clusterOperationTimeout = 5 * time.Minute

// ClusterRequest selects the project whose cluster an operation acts on.
type ClusterRequest struct {
	Project ProjectSelection
}

// ClusterResult is the externally visible part of a cluster record.
type ClusterResult struct {
	ProjectID        string                           `json:"projectId"`
	Name             string                           `json:"name"`
	Context          string                           `json:"context"`
	Kubeconfig       string                           `json:"kubeconfig"`
	UID              string                           `json:"uid"`
	Status           string                           `json:"status"`
	NamespaceUIDs    map[string]string                `json:"namespaceUIDs,omitempty"`
	InventoryVersion int                              `json:"inventoryVersion"`
	Bootstrap        map[string]ClusterObjectIdentity `json:"bootstrap,omitempty"`
	ResourceTypes    []string                         `json:"resourceTypes,omitempty"`
	UpdatedAt        time.Time                        `json:"updatedAt"`
	Error            string                           `json:"error,omitempty"`
}

// ClusterObjectIdentity exposes only an object's UID and controller ancestry.
type ClusterObjectIdentity struct {
	UID         string   `json:"uid"`
	Controllers []string `json:"controllers,omitempty"`
}

func clusterResultOf(record cluster.Record) ClusterResult {
	bootstrap := make(map[string]ClusterObjectIdentity, len(record.Baseline.Inventory.Objects))
	for key, identity := range record.Baseline.Inventory.Objects {
		bootstrap[key] = ClusterObjectIdentity(identity)
	}
	return ClusterResult{
		ProjectID:        record.ProjectID.String(),
		Name:             record.Name,
		Context:          record.Context,
		Kubeconfig:       record.Kubeconfig,
		UID:              record.UID,
		Status:           string(record.Status),
		NamespaceUIDs:    record.Baseline.Namespaces,
		InventoryVersion: record.Baseline.InventoryVersion,
		Bootstrap:        bootstrap,
		ResourceTypes:    record.Baseline.Inventory.ResourceTypes,
		UpdatedAt:        record.UpdatedAt,
		Error:            record.Error,
	}
}

func accessOf(record cluster.Record) ClusterAccess {
	return ClusterAccess{Kubeconfig: record.Kubeconfig, Context: record.Context}
}

// clusterGuard loads the project's cluster record and checks that the live
// cluster is still the one it records. Every use case that touches the
// cluster goes through it, so none can act on a replacement target.
type clusterGuard struct {
	records  ClusterReader
	observer ClusterObserver
}

// load reads the record and checks that it belongs to owner.
func (guard clusterGuard) load(owner project.Project) (cluster.Record, error) {
	record, err := guard.records.Load(owner)
	if err != nil {
		return cluster.Record{}, err
	}

	if err := record.BelongsTo(owner); err != nil {
		return cluster.Record{}, err
	}

	return record, nil
}

// verified loads the record and verifies the live identity. A deleted cluster
// has nothing live to verify.
func (guard clusterGuard) verified(ctx context.Context, owner project.Project) (cluster.Record, error) {
	record, err := guard.load(owner)
	if err != nil {
		return cluster.Record{}, err
	}

	if record.Status == cluster.StatusDeleted {
		return record, nil
	}

	if err := guard.verify(ctx, record); err != nil {
		return cluster.Record{}, err
	}

	return record, nil
}

// ready returns a verified record only when the cluster can host runtimes.
func (guard clusterGuard) ready(ctx context.Context, owner project.Project) (cluster.Record, error) {
	record, err := guard.load(owner)
	if err != nil {
		return cluster.Record{}, err
	}

	if err := record.RequireReady(); err != nil {
		return cluster.Record{}, err
	}

	if err := guard.verify(ctx, record); err != nil {
		return cluster.Record{}, err
	}

	return record, nil
}

func (guard clusterGuard) verify(ctx context.Context, record cluster.Record) error {
	if record.UID == "" {
		return record.VerifyIdentity(nil)
	}
	namespaces, err := guard.observer.Namespaces(ctx, accessOf(record))
	if err != nil {
		return err
	}

	return record.VerifyIdentity(namespaces)
}

// recordFailure persists why a cluster operation stopped without hiding the cause.
func recordFailure(store ClusterWriter, owner project.Project, record *cluster.Record, cause error) error {
	record.Fail(cause, isInterruption(cause))

	if err := saveCluster(store, owner, record); err != nil {
		return errors.Join(cause, fmt.Errorf("record cluster failure: %w", err))
	}

	return cause
}

func saveCluster(store ClusterWriter, owner project.Project, record *cluster.Record) error {
	record.UpdatedAt = time.Now().UTC()

	return store.Save(owner, *record)
}
