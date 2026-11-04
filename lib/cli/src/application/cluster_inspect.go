package application

import "context"

// InspectCluster reports the recorded cluster after verifying the live
// identity. It never updates the record or adopts anything.
type InspectCluster struct {
	Projects ProjectOpener
	Clusters ClusterReader
	Observer ClusterObserver
}

// Execute verifies and reports the cluster.
func (useCase InspectCluster) Execute(ctx context.Context, request ClusterRequest) (ClusterResult, error) {
	owner, err := useCase.Projects.Open(request.Project)
	if err != nil {
		return ClusterResult{}, err
	}

	record, err := clusterGuard{records: useCase.Clusters, observer: useCase.Observer}.verified(ctx, owner)
	if err != nil {
		return ClusterResult{}, err
	}

	return clusterResultOf(record), nil
}
