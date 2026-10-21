package microruntime

// Blueprint holds the explicit consumer choices from which a runtime's native
// files are generated. It is the domain's view of guided setup: how the choices
// are collected (prompts, flags, a file) is not its concern.
type Blueprint struct {
	Owner     string
	Target    KubernetesTarget
	Workloads []Workload
	Quotas    ResourceQuotas
}

// Validate enforces the rules every blueprint must satisfy before generation.
// Whether each workload's component directory exists is checked elsewhere.
func (blueprint Blueprint) Validate() error {
	if !isPlainText(blueprint.Owner) ||
		!isPlainText(blueprint.Target.Context) ||
		blueprint.Target.Namespace.IsZero() {
		return reject(
			ErrInvalidBlueprint,
			"owner, context, and a valid namespace are required",
		)
	}

	if len(blueprint.Workloads) == 0 || len(blueprint.Quotas) == 0 {
		return reject(
			ErrInvalidBlueprint,
			"at least one workload and explicit resource quotas are required",
		)
	}

	for resource, quantity := range blueprint.Quotas {
		if !resource.Supported() || quantity.IsZero() {
			return reject(
				ErrInvalidBlueprint,
				"unsupported quota or quantity %s=%s",
				resource, quantity,
			)
		}
	}

	names := make(map[Identifier]bool, len(blueprint.Workloads))

	for _, workload := range blueprint.Workloads {
		if err := workload.Validate(); err != nil {
			return err
		}

		if names[workload.Name] {
			return reject(
				ErrInvalidBlueprint,
				"duplicate workload %q",
				workload.Name,
			)
		}

		names[workload.Name] = true
	}

	return nil
}

// Descriptor derives the descriptor that references the native files generated
// for this blueprint. Each unit's manifests live under units/<name>/kubernetes
// inside the runtime's configuration directory. Call Validate first.
func (blueprint Blueprint) Descriptor(id Identifier) Descriptor {
	units := make(map[Identifier]Unit, len(blueprint.Workloads))

	for _, workload := range blueprint.Workloads {
		units[workload.Name] = Unit{
			Component: workload.Component,
			Manifests: mustRelativePath(
				"units/" + workload.Name.String() + "/kubernetes",
			),
		}
	}

	return Descriptor{
		Version:            SupportedVersion,
		ID:                 id,
		Environment:        LocalEnvironment,
		Owner:              blueprint.Owner,
		Target:             blueprint.Target,
		TerraformDirectory: defaultTerraformDirectory,
		Units:              units,
	}
}
