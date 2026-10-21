package application_test

import (
	"context"
	"errors"
	"slices"
	"testing"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/domain/cluster"
)

func (w *world) loadImage() application.LoadClusterImage {
	return application.LoadClusterImage{Projects: w.projects, Locks: w.locks, Clusters: w.clusters, Observer: w.observer, Images: w.images}
}

func (w *world) imageRequest(image string) application.LoadClusterImageRequest {
	return application.LoadClusterImageRequest{Project: w.clusterRequest().Project, Image: image}
}

func TestLoadClusterImage_LoadsIntoTheVerifiedCluster(t *testing.T) {
	t.Parallel()

	w := newWorld(t)

	result, err := w.loadImage().Execute(context.Background(), w.imageRequest("web:dev"))
	if err != nil || result != (application.ImageResult{Image: "web:dev", Status: "loaded"}) {
		t.Fatalf("Execute() = %+v, %v", result, err)
	}

	if !slices.Equal(w.images.loaded, []string{w.owner.ClusterName() + " web:dev"}) {
		t.Fatalf("loaded %v, want the image in the project's cluster", w.images.loaded)
	}

	w.requireLockBalanced(t)
}

func TestLoadClusterImage_Refuses(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		image   string
		arrange func(*world)
		wantErr error
	}{
		{name: "not one image reference", image: "web:dev --all", arrange: func(*world) {}},
		{name: "empty image reference", image: "", arrange: func(*world) {}},
		{name: "cluster not ready", image: "web:dev", arrange: func(w *world) { w.clusters.record.Status = cluster.StatusCreating }, wantErr: cluster.ErrNotReady},
		{name: "replaced cluster", image: "web:dev", arrange: func(w *world) { w.observer.namespaces["kube-system"] = "uid-new" }, wantErr: cluster.ErrUntrustedIdentity},
		{name: "lock held elsewhere", image: "web:dev", arrange: func(w *world) { w.locks.err = errInjected }, wantErr: errInjected},
		{name: "unknown project", image: "web:dev", arrange: func(w *world) { w.projects.err = errInjected }, wantErr: errInjected},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := newWorld(t)
			tt.arrange(w)

			_, err := w.loadImage().Execute(context.Background(), w.imageRequest(tt.image))
			if err == nil || (tt.wantErr != nil && !errors.Is(err, tt.wantErr)) {
				t.Fatalf("Execute() error = %v, want %v", err, tt.wantErr)
			}

			if len(w.images.loaded) != 0 {
				t.Error("a refused request loaded an image")
			}

			w.requireLockBalanced(t)
		})
	}
}

func TestLoadClusterImage_ReportsTheLoaderFailure(t *testing.T) {
	t.Parallel()

	w := newWorld(t)
	w.images.loadErr = errInjected

	if _, err := w.loadImage().Execute(context.Background(), w.imageRequest("web:dev")); !errors.Is(err, errInjected) {
		t.Fatalf("Execute() error = %v, want the loader failure", err)
	}
}
