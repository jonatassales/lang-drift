package application

import (
	"context"
	"fmt"
	"slices"

	"github.com/exphere/ex/src/domain/cluster"
	"github.com/exphere/ex/src/domain/project"
)

// DeleteCluster removes the project's cluster once every runtime is cleaned
// up and the cluster holds nothing EX did not create. Configuration and
// evidence stay.
type DeleteCluster struct {
	Projects    ProjectOpener
	Locks       ProjectLocker
	Clusters    ClusterWriter
	Credentials ClusterCredentials
	Executions  ExecutionLister
	Provisioner ClusterProvisioner
	Observer    ClusterObserver
}

// Execute deletes the cluster, or reconciles a record whose cluster is already gone.
func (useCase DeleteCluster) Execute(ctx context.Context, request ClusterRequest) (ClusterResult, error) {
	owner, err := useCase.Projects.Open(request.Project)
	if err != nil {
		return ClusterResult{}, err
	}

	unlock, err := useCase.Locks.Lock(owner)
	if err != nil {
		return ClusterResult{}, err
	}
	defer unlock()

	guard := clusterGuard{records: useCase.Clusters, observer: useCase.Observer}

	record, err := guard.load(owner)
	if err != nil {
		return ClusterResult{}, err
	}

	if record.Status == cluster.StatusDeleted {
		return clusterResultOf(record), nil
	}

	if err := useCase.requireClosedRuntimes(owner); err != nil {
		return ClusterResult{}, err
	}

	present, err := useCase.exists(ctx, record.Name)
	if err != nil {
		return ClusterResult{}, err
	}

	if !present {
		// A failed creation the provisioner cleaned up, or an interrupted
		// deletion. Nothing live is touched.
		return useCase.forget(owner, &record)
	}

	if err := guard.verify(ctx, record); err != nil {
		return ClusterResult{}, err
	}

	if err := useCase.authorize(ctx, &record); err != nil {
		return ClusterResult{}, err
	}

	record.MarkDeleting()
	if err := saveCluster(useCase.Clusters, owner, &record); err != nil {
		return ClusterResult{}, err
	}

	operationContext, cancel := context.WithTimeout(ctx, clusterOperationTimeout)
	defer cancel()

	if err := useCase.destroy(operationContext, record.Name); err != nil {
		return ClusterResult{}, recordFailure(useCase.Clusters, owner, &record, err)
	}

	return useCase.forget(owner, &record)
}

// requireClosedRuntimes refuses deletion while any runtime may still own resources.
func (useCase DeleteCluster) requireClosedRuntimes(owner project.Project) error {
	executions, err := useCase.Executions.List(owner)
	if err != nil {
		return err
	}

	for _, execution := range executions {
		if err := execution.Validate(owner.ID, execution.Runtime); err != nil {
			return err
		}

		if err := execution.RequireClosed(); err != nil {
			return fmt.Errorf("%w; clean it up before deleting the cluster", err)
		}
	}

	return nil
}

func (useCase DeleteCluster) authorize(ctx context.Context, record *cluster.Record) error {
	access := accessOf(*record)

	namespaces, err := useCase.Observer.Namespaces(ctx, access)
	if err != nil {
		return err
	}

	inventory, err := useCase.Observer.ClusterInventory(ctx, access)
	if err != nil {
		return err
	}

	return record.AuthorizeDeletion(namespaces, inventory)
}

func (useCase DeleteCluster) destroy(ctx context.Context, name string) error {
	if err := useCase.Provisioner.Delete(ctx, name); err != nil {
		return err
	}

	stillPresent, err := useCase.exists(ctx, name)
	if err != nil {
		return err
	}

	if stillPresent {
		return fmt.Errorf("cluster %s still exists after deletion", name)
	}

	return nil
}

func (useCase DeleteCluster) exists(ctx context.Context, name string) (bool, error) {
	names, err := useCase.Provisioner.Clusters(ctx)
	if err != nil {
		return false, err
	}

	return slices.Contains(names, name), nil
}

// forget records the deletion, then discards the credentials it made useless.
func (useCase DeleteCluster) forget(owner project.Project, record *cluster.Record) (ClusterResult, error) {
	record.MarkDeleted()

	if err := saveCluster(useCase.Clusters, owner, record); err != nil {
		return ClusterResult{}, err
	}

	if err := useCase.Credentials.DiscardCredentials(owner); err != nil {
		return ClusterResult{}, err
	}

	return clusterResultOf(*record), nil
}
