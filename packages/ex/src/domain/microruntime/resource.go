package microruntime

// Resource identifies one object an execution created. The UID guards against
// acting on a different object that later reuses the name.
type Resource struct {
	APIVersion string
	Kind       string
	Name       string
	Namespace  string
	UID        string
	Unit       Identifier
}

// ResourceOf records an object's identity as the API server reports it.
func ResourceOf(object Object, unit Identifier) Resource {
	return Resource{
		APIVersion: object.APIVersion(),
		Kind:       object.Kind(),
		Name:       object.Name(),
		Namespace:  object.Namespace(),
		UID:        object.UID(),
		Unit:       unit,
	}
}

// IsWorkload reports whether the resource runs pods.
func (resource Resource) IsWorkload() bool {
	return resource.Kind == "Deployment" || resource.Kind == "Job"
}

// isSupported reports whether the resource is one EX creates, in namespace,
// and may therefore delete.
func (resource Resource) isSupported(namespace Identifier) bool {
	expected, supported := apiVersions[resource.Kind]
	_, nameErr := ParseIdentifier(resource.Name)

	return supported &&
		resource.APIVersion == expected &&
		nameErr == nil &&
		resource.Namespace == namespace.String() &&
		!resource.Unit.IsZero()
}

// Readiness reports whether the live object for this resource is ready: a
// Service once it exists, a Job once it completes, and a Deployment once every
// desired replica is updated and available.
func (resource Resource) Readiness(live Object) (bool, error) {
	switch {
	case live == nil:
		return false, reject(ErrWorkloadFailed, "%s %s disappeared during readiness", resource.Kind, resource.Name)

	case live.UID() != resource.UID:
		return false, reject(ErrOwnershipChanged, "%s %s was replaced during readiness", resource.Kind, resource.Name)

	case resource.Kind == "Job":
		return jobReadiness(resource, live)

	case resource.Kind == "Deployment":
		return deploymentReady(live), nil
	}

	return true, nil
}

func jobReadiness(resource Resource, live Object) (bool, error) {
	for _, condition := range live.Status().Children("conditions") {
		if condition.Text("status") != "True" {
			continue
		}

		switch condition.Text("type") {
		case "Complete":
			return true, nil
		case "Failed":
			return false, reject(ErrWorkloadFailed, "Job %s failed: %s", resource.Name, condition.Text("reason"))
		}
	}

	return false, nil
}

func deploymentReady(live Object) bool {
	spec, status := live.Spec(), live.Status()

	desired := 1.0
	if spec["replicas"] != nil {
		desired = spec.Number("replicas")
	}

	return status.Number("observedGeneration") >= live.Metadata().Number("generation") &&
		status.Number("updatedReplicas") == desired &&
		status.Number("availableReplicas") == desired &&
		status.Number("replicas") == desired
}

// Removed reports whether the object recorded under uid is gone. Finding a
// different object under the same name means it was replaced meanwhile.
func Removed(live Object, uid string) (bool, error) {
	switch {
	case live == nil:
		return true, nil

	case live.UID() != uid:
		return false, reject(ErrOwnershipChanged, "%s %s was replaced during cleanup", live.Kind(), live.Name())
	}

	return false, nil
}
