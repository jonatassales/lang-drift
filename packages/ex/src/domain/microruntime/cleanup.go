package microruntime

// endpointControllerLabel marks Endpoints the legacy endpoint controller derives
// from a Service. They carry no owner reference to that Service.
const endpointControllerLabel = "endpoints.kubernetes.io/managed-by"

// AuthorizeCleanup checks a complete inventory of the runtime's namespace.
//
// Before deletion, every object must be recorded, descend from a recorded
// object through owner references, or be the Endpoints of a recorded Service.
// After deletion, only the namespace's defaults and the runtime's own quota
// may remain.
func (execution *Execution) AuthorizeCleanup(inventory []Object, deleted bool) error {
	byUID := make(map[string]Object, len(inventory))
	for _, object := range inventory {
		if uid := object.UID(); uid != "" {
			byUID[uid] = object
		}
	}

	owned := execution.descendants(inventory)

	for _, object := range inventory {
		switch {
		case isNamespaceDefault(object):
			continue

		case object.Kind() == "ResourceQuota" &&
			object.Labels().Text(LabelOwnershipToken) == execution.Token:
			continue

		case !deleted && owned[object.UID()]:
			continue

		case !deleted && execution.isServiceEndpoints(object, byUID, owned):
			continue
		}

		return reject(ErrCleanupBlocked, "remaining or foreign %s/%s", object.Kind(), object.Name())
	}

	return nil
}

// isNamespaceDefault recognizes objects Kubernetes creates in every namespace.
func isNamespaceDefault(object Object) bool {
	return (object.Kind() == "ServiceAccount" && object.Name() == "default") ||
		(object.Kind() == "ConfigMap" && object.Name() == "kube-root-ca.crt")
}

// descendants returns the UIDs of recorded objects and of every object that
// reaches one through owner references.
func (execution *Execution) descendants(inventory []Object) map[string]bool {
	owned := map[string]bool{}
	for _, resource := range execution.Resources {
		if resource.UID != "" {
			owned[resource.UID] = true
		}
	}

	for changed := true; changed; {
		changed = false

		for _, object := range inventory {
			uid := object.UID()
			if uid == "" || owned[uid] {
				continue
			}

			for _, owner := range object.OwnerUIDs() {
				if owned[owner] {
					owned[uid] = true
					changed = true

					break
				}
			}
		}
	}

	return owned
}

// isServiceEndpoints recognizes the endpoint controller's object for a
// recorded Service. A name or label alone is not enough: the live Service and
// every Pod the Endpoints advertise must be correlated with this execution.
// Empty Endpoints are normal before Pods exist or after a failed readiness.
func (execution *Execution) isServiceEndpoints(endpoints Object, byUID map[string]Object, owned map[string]bool) bool {
	metadata := endpoints.Metadata()

	if endpoints.APIVersion() != "v1" ||
		endpoints.Kind() != "Endpoints" ||
		metadata.Text("namespace") != execution.Namespace.String() ||
		endpoints.UID() == "" ||
		len(asSlice(metadata["ownerReferences"])) != 0 ||
		endpoints.Labels().Text(endpointControllerLabel) != "endpoint-controller" {
		return false
	}

	service, unit, found := execution.recordedService(endpoints.Name(), byUID)
	if !found || !execution.carriesClaim(endpoints.Labels(), unit) || !execution.carriesClaim(service.Labels(), unit) {
		return false
	}

	selector := service.Spec().Child("selector")
	if len(selector) == 0 {
		return false
	}

	for _, value := range selector {
		if _, isText := value.(string); !isText {
			return false
		}
	}

	return execution.advertisesOwnedPods(endpoints, selector, byUID, owned)
}

// recordedService finds the live Service this execution recorded under name.
func (execution *Execution) recordedService(name string, byUID map[string]Object) (Object, Identifier, bool) {
	namespace := execution.Namespace.String()

	for _, resource := range execution.Resources {
		if resource.Kind != "Service" || resource.APIVersion != "v1" ||
			resource.Namespace != namespace || resource.Name != name || resource.UID == "" {
			continue
		}

		service := byUID[resource.UID]
		if !service.IsOwnedBy(resource.UID, execution.Token) ||
			service.APIVersion() != "v1" || service.Kind() != "Service" ||
			service.Namespace() != namespace || service.Name() != name {
			return nil, Identifier{}, false
		}

		return service, resource.Unit, true
	}

	return nil, Identifier{}, false
}

// carriesClaim reports whether labels carry this execution's claim for unit.
func (execution *Execution) carriesClaim(labels Object, unit Identifier) bool {
	return labels.Text(LabelRuntime) == execution.Runtime.String() &&
		labels.Text(LabelUnit) == unit.String() &&
		labels.Text(LabelOwnershipToken) == execution.Token
}

func (execution *Execution) advertisesOwnedPods(endpoints, selector Object, byUID map[string]Object, owned map[string]bool) bool {
	subsets, ok := strictSlice(endpoints["subsets"])
	if !ok {
		return false
	}

	for _, value := range subsets {
		subset := asObject(value)
		if subset == nil {
			return false
		}

		for _, field := range []string{"addresses", "notReadyAddresses"} {
			addresses, ok := strictSlice(subset[field])
			if !ok {
				return false
			}

			for _, address := range addresses {
				reference := asObject(address).Child("targetRef")
				if !execution.isOwnedSelectedPod(reference, selector, byUID, owned) {
					return false
				}
			}
		}
	}

	return true
}

func (execution *Execution) isOwnedSelectedPod(reference, selector Object, byUID map[string]Object, owned map[string]bool) bool {
	namespace := execution.Namespace.String()
	uid := reference.Text("uid")
	pod := byUID[uid]
	labels := pod.Labels()

	unit, err := ParseIdentifier(labels.Text(LabelUnit))
	_, recordedUnit := execution.Units[unit]

	if err != nil || !recordedUnit || !owned[uid] ||
		reference.Text("kind") != "Pod" ||
		(reference.Text("apiVersion") != "" && reference.Text("apiVersion") != "v1") ||
		pod.APIVersion() != "v1" || pod.Kind() != "Pod" ||
		reference.Text("namespace") != namespace || pod.Namespace() != namespace ||
		reference.Text("name") != pod.Name() ||
		labels.Text(LabelRuntime) != execution.Runtime.String() ||
		!pod.IsOwnedBy(uid, execution.Token) {
		return false
	}

	for label, expected := range selector {
		actual, isText := labels[label].(string)
		if !isText || actual != expected {
			return false
		}
	}

	return true
}
