package microruntime

// SelectResource requires exactly one recorded resource for a unit. Workload
// selects a Deployment or Job; otherwise it selects a Service. A missing UID
// still requires a later ownership check and is never permission to adopt.
func (execution *Execution) SelectResource(unit Identifier, workload bool) (Resource, error) {
	kind := "Service"
	if workload {
		kind = "workload"
	}

	var selected Resource
	matches := 0
	for _, resource := range execution.Resources {
		if resource.Unit == unit && ((workload && resource.IsWorkload()) || (!workload && resource.Kind == "Service")) {
			selected = resource
			matches++
		}
	}

	if matches != 1 {
		return Resource{}, reject(ErrInconsistentExecution, "unit %q must have exactly one %s for this operation", unit, kind)
	}

	return selected, nil
}

// ForwardPort selects an explicit remote port or the Service's sole port.
// It does not infer a target for an ambiguous or malformed Service.
func (object Object) ForwardPort(explicit int) (int, error) {
	if explicit != 0 {
		if explicit < 1 || explicit > 65535 {
			return 0, reject(ErrUnsupportedManifest, "ports must be between 1 and 65535")
		}
		return explicit, nil
	}

	ports := object.Spec().Children("ports")
	if len(ports) != 1 {
		return 0, reject(ErrUnsupportedManifest, "select --remote-port for a service with multiple ports")
	}

	port := ports[0].Number("port")
	if port < 1 || port > 65535 || port != float64(int(port)) {
		return 0, reject(ErrUnsupportedManifest, "invalid service port")
	}

	return int(port), nil
}
