package application

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/exphere/ex/src/domain/cluster"
	"github.com/exphere/ex/src/domain/project"
)

// CreateCluster creates the project's local cluster and records its identity
// and baseline. A ready cluster is reused after verification. Intent is
// recorded before creation, so an interrupted creation is never forgotten.
type CreateCluster struct {
	Projects    ProjectOpener
	Locks       ProjectLocker
	Clusters    ClusterWriter
	Credentials ClusterCredentials
	Provisioner ClusterProvisioner
	Observer    ClusterObserver
}

// Execute creates or reuses the cluster.
func (useCase CreateCluster) Execute(ctx context.Context, request ClusterRequest) (ClusterResult, error) {
	owner, err := useCase.Projects.Open(request.Project)
	if err != nil {
		return ClusterResult{}, err
	}

	unlock, err := useCase.Locks.Lock(owner)
	if err != nil {
		return ClusterResult{}, err
	}
	defer unlock()

	prior, err := useCase.prior(owner)
	if err != nil {
		return ClusterResult{}, err
	}

	reuse, err := cluster.PlanCreation(prior)
	if err != nil {
		return ClusterResult{}, err
	}

	if reuse {
		if err := useCase.guard().verify(ctx, *prior); err != nil {
			return ClusterResult{}, err
		}

		return clusterResultOf(*prior), nil
	}

	if err := useCase.requireUnusedName(ctx, owner.ClusterName()); err != nil {
		return ClusterResult{}, err
	}

	kubeconfig, err := useCase.reserveCredentials(owner, prior)
	if err != nil {
		return ClusterResult{}, err
	}

	record := cluster.Begin(owner, kubeconfig)
	if err := saveCluster(useCase.Clusters, owner, &record); err != nil {
		return ClusterResult{}, err
	}

	operationContext, cancel := context.WithTimeout(ctx, clusterOperationTimeout)
	defer cancel()

	if err := useCase.provision(operationContext, owner, &record); err != nil {
		return ClusterResult{}, recordFailure(useCase.Clusters, owner, &record, err)
	}

	return clusterResultOf(record), nil
}

func (useCase CreateCluster) guard() clusterGuard {
	return clusterGuard{records: useCase.Clusters, observer: useCase.Observer}
}

// prior returns the existing record, or nil when the project never had one.
func (useCase CreateCluster) prior(owner project.Project) (*cluster.Record, error) {
	record, err := useCase.guard().load(owner)

	switch {
	case errors.Is(err, ErrNotFound):
		return nil, nil
	case err != nil:
		return nil, err
	}

	return &record, nil
}

func (useCase CreateCluster) requireUnusedName(ctx context.Context, name string) error {
	existing, err := useCase.Provisioner.Clusters(ctx)
	if err != nil {
		return err
	}

	if slices.Contains(existing, name) {
		return fmt.Errorf("%w: a cluster named %s already exists", ErrAdoptionRefused, name)
	}

	return nil
}

// reserveCredentials returns where the new cluster's kubeconfig goes. Leftover
// credentials are discarded only when the record says their cluster was deleted.
func (useCase CreateCluster) reserveCredentials(owner project.Project, prior *cluster.Record) (string, error) {
	location, exists, err := useCase.Credentials.Credentials(owner)
	if err != nil {
		return "", err
	}

	if !exists {
		return location, nil
	}

	if prior == nil || prior.Status != cluster.StatusDeleted {
		return "", fmt.Errorf("%w: an unrecorded kubeconfig exists at %s", ErrAdoptionRefused, location)
	}

	if err := useCase.Credentials.DiscardCredentials(owner); err != nil {
		return "", err
	}

	return location, nil
}

// provision creates the cluster, then records its identity before the
// optional inventory can fail, then its baseline.
func (useCase CreateCluster) provision(ctx context.Context, owner project.Project, record *cluster.Record) error {
	contextName, err := useCase.Provisioner.Create(ctx, record.Name, record.Kubeconfig)
	if err != nil {
		return err
	}

	record.Context = contextName

	namespaces, err := useCase.Observer.Namespaces(ctx, accessOf(*record))
	if err != nil {
		return err
	}

	if err := record.Identify(namespaces); err != nil {
		return err
	}

	if err := saveCluster(useCase.Clusters, owner, record); err != nil {
		return err
	}

	bootstrap, err := useCase.Observer.ClusterInventory(ctx, accessOf(*record))
	if err != nil {
		return err
	}

	if err := record.Complete(bootstrap); err != nil {
		return err
	}

	return saveCluster(useCase.Clusters, owner, record)
}
