package application

import (
	"context"
	"fmt"

	"github.com/exphere/ex/src/domain/microruntime"
)

// LoadClusterImage copies an already-built local image into the verified
// project cluster. It never builds an image or pulls from a registry.
type LoadClusterImage struct {
	Projects ProjectOpener
	Locks    ProjectLocker
	Clusters ClusterReader
	Observer ClusterObserver
	Images   ImageLoader
}

// LoadClusterImageRequest names the image to load.
type LoadClusterImageRequest struct {
	Project ProjectSelection
	Image   string
}

// ImageResult describes a loaded image.
type ImageResult struct {
	Image  string `json:"image"`
	Status string `json:"status"`
}

// Execute loads the image into the ready cluster.
func (useCase LoadClusterImage) Execute(ctx context.Context, request LoadClusterImageRequest) (ImageResult, error) {
	if !microruntime.IsImageReference(request.Image) {
		return ImageResult{}, fmt.Errorf("%q is not one local image reference", request.Image)
	}

	owner, err := useCase.Projects.Open(request.Project)
	if err != nil {
		return ImageResult{}, err
	}

	unlock, err := useCase.Locks.Lock(owner)
	if err != nil {
		return ImageResult{}, err
	}
	defer unlock()

	record, err := clusterGuard{records: useCase.Clusters, observer: useCase.Observer}.ready(ctx, owner)
	if err != nil {
		return ImageResult{}, err
	}

	operationContext, cancel := context.WithTimeout(ctx, clusterOperationTimeout)
	defer cancel()

	if err := useCase.Images.LoadImage(operationContext, record.Name, request.Image); err != nil {
		return ImageResult{}, err
	}

	return ImageResult{Image: request.Image, Status: "loaded"}, nil
}
