package microruntime

// ToolVersion records whether a native tool's version could be observed.
type ToolVersion struct {
	Version  string
	Observed bool
}

// Image observation states. A configured container that no live pod has
// reported yet is "not-observed"; a live pod without an image ID is
// "unavailable". Missing data is never interpreted as a value.
const (
	ImageObserved    = "observed"
	ImageUnavailable = "unavailable"
	ImageNotObserved = "not-observed"
)

// ContainerObservation reports a configured image separately from the
// immutable image identity a live pod reports. It deliberately excludes
// arguments, environment, and logs.
type ContainerObservation struct {
	Unit      string
	Workload  string
	Pod       string
	PodUID    string
	Container string
	Image     string
	ImageID   string
	Ready     bool
	Restarts  int
	Phase     string
}

// ImageStatus classifies how much is known about the container's image.
func (observation ContainerObservation) ImageStatus() string {
	switch {
	case observation.ImageID != "":
		return ImageObserved
	case observation.PodUID == "":
		return ImageNotObserved
	}

	return ImageUnavailable
}

// ConfiguredContainers lists the containers the manifests declare, before any
// pod exists.
func ConfiguredContainers(manifests []Manifest) []ContainerObservation {
	observations := []ContainerObservation{}

	for _, manifest := range manifests {
		object := manifest.Object
		if object.Kind() == "Service" {
			continue
		}

		for _, container := range object.PodSpec().Containers() {
			observations = append(observations, ContainerObservation{
				Unit:      manifest.Unit.String(),
				Workload:  object.Kind() + "/" + object.Name(),
				Container: container.Text("name"),
				Image:     container.Text("image"),
			})
		}
	}

	return observations
}

// MergeContainers prefers live observations and keeps configured containers
// that no live pod reported.
func MergeContainers(configured, observed []ContainerObservation) []ContainerObservation {
	key := func(observation ContainerObservation) string {
		return observation.Unit + "/" + observation.Workload + "/" + observation.Container
	}

	merged := append([]ContainerObservation(nil), observed...)

	reported := make(map[string]bool, len(observed))
	for _, observation := range observed {
		reported[key(observation)] = true
	}

	for _, observation := range configured {
		if !reported[key(observation)] {
			merged = append(merged, observation)
		}
	}

	return merged
}

// ObserveContainers reports the containers of the pods that the execution's
// workloads control, directly or through a ReplicaSet. It also returns the UIDs
// of every owned object it traversed.
func (execution *Execution) ObserveContainers(pods, replicaSets []Object) ([]ContainerObservation, map[string]bool) {
	owned := map[string]bool{}
	workloads := map[string]string{}

	for _, resource := range execution.Resources {
		if resource.UID != "" {
			owned[resource.UID] = true
			workloads[resource.UID] = resource.Kind + "/" + resource.Name
		}
	}

	for _, replicaSet := range replicaSets {
		uid := replicaSet.UID()
		if !replicaSet.IsOwnedBy(uid, execution.Token) {
			continue
		}

		for _, controller := range replicaSet.ControllerUIDs() {
			if owned[controller] {
				owned[uid] = true
				workloads[uid] = workloads[controller]
			}
		}
	}

	observations := []ContainerObservation{}

	for _, pod := range pods {
		uid := pod.UID()
		if !pod.IsOwnedBy(uid, execution.Token) {
			continue
		}

		workload, controlled := "", false

		for _, controller := range pod.ControllerUIDs() {
			if owned[controller] {
				workload, controlled = workloads[controller], true
			}
		}

		if !controlled {
			continue
		}

		owned[uid] = true
		observations = append(observations, podContainers(pod, workload)...)
	}

	return observations, owned
}

func podContainers(pod Object, workload string) []ContainerObservation {
	status := pod.Status()

	statuses := map[string]Object{}
	for _, containerStatus := range append(status.Children("containerStatuses"), status.Children("initContainerStatuses")...) {
		statuses[containerStatus.Text("name")] = containerStatus
	}

	observations := []ContainerObservation{}

	for _, container := range pod.Spec().Containers() {
		name := container.Text("name")
		observed := statuses[name]

		observations = append(observations, ContainerObservation{
			Unit:      pod.Labels().Text(LabelUnit),
			Workload:  workload,
			Pod:       pod.Name(),
			PodUID:    pod.UID(),
			Container: name,
			Image:     container.Text("image"),
			ImageID:   observed.Text("imageID"),
			Ready:     observed.Flag("ready"),
			Restarts:  int(observed.Number("restartCount")),
			Phase:     status.Text("phase"),
		})
	}

	return observations
}
