package microruntime_test

import (
	"slices"
	"testing"

	"github.com/exphere/ex/src/domain/microruntime"
)

func TestContainerObservation_ImageStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		observation microruntime.ContainerObservation
		want        string
	}{
		{name: "reported image ID", observation: microruntime.ContainerObservation{PodUID: "uid-pod", ImageID: "sha256:abc"}, want: microruntime.ImageObserved},
		{name: "no pod yet", observation: microruntime.ContainerObservation{Image: "web:dev"}, want: microruntime.ImageNotObserved},
		{name: "pod without an image ID", observation: microruntime.ContainerObservation{PodUID: "uid-pod"}, want: microruntime.ImageUnavailable},
	}

	for _, tt := range tests {
		if got := tt.observation.ImageStatus(); got != tt.want {
			t.Errorf("%s: ImageStatus() = %q, want %q", tt.name, got, tt.want)
		}
	}
}

func TestConfiguredContainers(t *testing.T) {
	t.Parallel()

	manifests := validManifests(t)
	job := manifests[2].Object
	job.PodSpec()["initContainers"] = []any{map[string]any{"name": "migrate", "image": "migrate:dev"}}

	var got []string
	for _, container := range microruntime.ConfiguredContainers(manifests) {
		got = append(got, container.Unit+" "+container.Workload+" "+container.Container+" "+container.Image)
	}

	want := []string{
		"web Deployment/web web web:dev",
		"batch Job/batch batch batch:dev",
		"batch Job/batch migrate migrate:dev",
	}
	if !slices.Equal(got, want) {
		t.Fatalf("ConfiguredContainers() = %q, want %q", got, want)
	}

	if empty := microruntime.ConfiguredContainers(nil); empty == nil || len(empty) != 0 {
		t.Fatalf("ConfiguredContainers(nil) = %#v, want an empty, non-nil list", empty)
	}
}

func TestMergeContainers(t *testing.T) {
	t.Parallel()

	configured := []microruntime.ContainerObservation{
		{Unit: "web", Workload: "Deployment/web", Container: "web", Image: "web:dev"},
		{Unit: "batch", Workload: "Job/batch", Container: "batch", Image: "batch:dev"},
	}
	observed := []microruntime.ContainerObservation{
		{Unit: "web", Workload: "Deployment/web", Container: "web", Pod: "web-1", PodUID: "uid-1", ImageID: "sha256:1"},
		{Unit: "web", Workload: "Deployment/web", Container: "web", Pod: "web-2", PodUID: "uid-2", ImageID: "sha256:1"},
	}

	merged := microruntime.MergeContainers(configured, observed)

	var got []string
	for _, observation := range merged {
		got = append(got, observation.Container+"@"+observation.Pod)
	}

	if !slices.Equal(got, []string{"web@web-1", "web@web-2", "batch@"}) {
		t.Fatalf("MergeContainers() = %v, want live observations first and unreported configured containers kept", got)
	}

	observed[0].Pod = "edited"
	if merged[0].Pod != "web-1" {
		t.Fatal("MergeContainers() shares its backing array with the observed input")
	}
}

func TestExecution_ObserveContainers(t *testing.T) {
	t.Parallel()

	execution := runningExecution(t)

	claimed := func(unit string) map[string]any {
		return map[string]any{microruntime.LabelOwnershipToken: token, microruntime.LabelUnit: unit}
	}
	controlledBy := func(object microruntime.Object, uid string) microruntime.Object {
		object["metadata"].(map[string]any)["ownerReferences"] = []any{map[string]any{"uid": uid, "controller": true}}

		return object
	}
	pod := func(name, uid, unit, controller string, containers []any, status map[string]any) microruntime.Object {
		var labels map[string]any
		if unit != "" {
			labels = claimed(unit)
		}

		object := controlledBy(liveObject("v1", "Pod", name, uid, labels), controller)
		object["spec"], object["status"] = map[string]any{"containers": containers}, status

		return object
	}
	container := func(name, image string) map[string]any { return map[string]any{"name": name, "image": image} }

	replicaSets := []microruntime.Object{
		controlledBy(liveObject("apps/v1", "ReplicaSet", "web-6f6", "uid-rs", claimed("web")), "uid-Deployment-web"),
		controlledBy(liveObject("apps/v1", "ReplicaSet", "web-foreign", "uid-foreign-rs", nil), "uid-Deployment-web"),
	}
	pods := []microruntime.Object{
		pod("web-6f6-a", "uid-pod-web", "web", "uid-rs", []any{container("web", "web:dev")}, map[string]any{
			"phase": "Running",
			"containerStatuses": []any{map[string]any{
				"name": "web", "imageID": "docker.io/library/web@sha256:abc", "ready": true, "restartCount": 2.0,
			}},
		}),
		pod("batch-x", "uid-pod-batch", "batch", "uid-Job-batch", []any{container("batch", "batch:dev")}, map[string]any{
			"phase":             "Succeeded",
			"containerStatuses": []any{map[string]any{"name": "batch", "ready": false}},
		}),
		pod("stranger", "uid-stranger", "", "uid-rs", []any{container("web", "web:dev")}, nil),
		pod("orphan", "uid-orphan", "web", "uid-other", []any{container("web", "web:dev")}, nil),
		pod("web-foreign-a", "uid-pod-foreign", "web", "uid-foreign-rs", []any{container("web", "web:dev")}, nil),
	}

	observations, owned := execution.ObserveContainers(pods, replicaSets)

	if len(observations) != 2 {
		t.Fatalf("ObserveContainers() reported %d containers, want the two owned pods: %+v", len(observations), observations)
	}

	web, batch := observations[0], observations[1]
	if web.Unit != "web" || web.Workload != "Deployment/web" || web.Pod != "web-6f6-a" || web.PodUID != "uid-pod-web" ||
		web.ImageID != "docker.io/library/web@sha256:abc" || !web.Ready || web.Restarts != 2 || web.Phase != "Running" {
		t.Errorf("web observation = %+v", web)
	}

	if batch.Workload != "Job/batch" || batch.ImageStatus() != microruntime.ImageUnavailable || batch.Phase != "Succeeded" {
		t.Errorf("batch observation = %+v, want the Job's pod with an unavailable image ID", batch)
	}

	for _, uid := range []string{"uid-Deployment-web", "uid-rs", "uid-pod-web", "uid-pod-batch"} {
		if !owned[uid] {
			t.Errorf("owned UIDs lack %s: %v", uid, owned)
		}
	}

	for _, uid := range []string{"uid-foreign-rs", "uid-stranger", "uid-orphan", "uid-pod-foreign"} {
		if owned[uid] {
			t.Errorf("owned UIDs include foreign %s", uid)
		}
	}
}
