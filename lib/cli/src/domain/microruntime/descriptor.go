// Package microruntime is the domain model of a microruntime: how it is
// defined, which Kubernetes objects it may own, and the ledger that records
// what one execution created so cleanup removes exactly that.
//
// Like the rest of the domain, it depends only on the standard library, the
// Microruntimes contracts, whose lexical rules it shares, and other domain
// packages. It knows nothing about Cobra, configuration files,
// the filesystem, Terraform, Kubernetes clients, or process execution; the
// outer layers depend on it, never the other way around.
package microruntime

import "sort"

const (
	// SupportedVersion is the only descriptor version EX accepts.
	SupportedVersion = 2

	// LocalEnvironment is the only environment EX can run in.
	LocalEnvironment = "local"
)

var (
	configurationRoot         = mustRelativePath(".langdrift/microruntimes/runtimes")
	defaultTerraformDirectory = mustRelativePath("terraform")
)

// ConfigurationDirectory returns where a runtime's configuration lives,
// relative to the project root. The id must come from ParseIdentifier.
func ConfigurationDirectory(id Identifier) RelativePath {
	return mustRelativePath(configurationRoot.String() + "/" + id.String())
}

// KubernetesTarget selects the exact context and namespace; no ambient target
// is ever used.
type KubernetesTarget struct {
	Context   string
	Namespace Identifier
}

// Unit references existing component source and its native manifests.
// Component is relative to the project root; Manifests is relative to the
// runtime's configuration directory.
type Unit struct {
	Component RelativePath
	Manifests RelativePath
}

// Descriptor identifies a runtime and references consumer-owned native
// configuration. Resource values belong to the native Kubernetes files, not to
// the descriptor. TerraformDirectory and Hooks are relative to the runtime's
// configuration directory.
type Descriptor struct {
	Version            int
	ID                 Identifier
	Environment        string
	Owner              string
	Target             KubernetesTarget
	TerraformDirectory RelativePath
	Units              map[Identifier]Unit
	Hooks              map[Identifier]RelativePath
}

// Scaffold returns the intentionally incomplete descriptor for manual setup.
// Owner, target, and units are left to the consumer, so Validate rejects it
// until they are filled in.
func Scaffold(id Identifier) Descriptor {
	return Descriptor{
		Version:            SupportedVersion,
		ID:                 id,
		Environment:        LocalEnvironment,
		TerraformDirectory: defaultTerraformDirectory,
		Units:              map[Identifier]Unit{},
	}
}

// UnitNames returns the unit names in a stable order.
func (descriptor Descriptor) UnitNames() []Identifier {
	names := make([]Identifier, 0, len(descriptor.Units))
	for name := range descriptor.Units {
		names = append(names, name)
	}

	sort.Slice(names, func(first, second int) bool {
		return names[first].String() < names[second].String()
	})

	return names
}

// Validate enforces the invariants that need no filesystem or cluster access.
// Whether referenced files exist is verified by the layers that can read them.
func (descriptor Descriptor) Validate() error {
	switch {
	case descriptor.Version != SupportedVersion:
		return reject(
			ErrInvalidDescriptor,
			"version %d is unsupported; expected version %d with explicit migration",
			descriptor.Version, SupportedVersion,
		)

	case descriptor.Environment != LocalEnvironment:
		return reject(
			ErrInvalidDescriptor,
			"only the %s environment is supported",
			LocalEnvironment,
		)

	case descriptor.ID.IsZero() ||
		descriptor.Target.Namespace.IsZero() ||
		!isPlainText(descriptor.Owner) ||
		!isPlainText(descriptor.Target.Context):
		return reject(
			ErrInvalidDescriptor,
			"runtime ID, owner, context, and namespace are required",
		)

	case descriptor.TerraformDirectory.IsZero():
		return reject(ErrInvalidDescriptor, "a Terraform directory is required")

	case len(descriptor.Units) == 0:
		return reject(ErrInvalidDescriptor, "at least one unit is required")
	}

	if err := descriptor.validateUnits(); err != nil {
		return err
	}

	return descriptor.validateHooks()
}

func (descriptor Descriptor) validateUnits() error {
	unitByManifests := make(map[RelativePath]Identifier, len(descriptor.Units))

	for _, name := range descriptor.UnitNames() {
		unit := descriptor.Units[name]

		switch {
		case name.IsZero() || unit.Component.IsZero() || unit.Manifests.IsZero():
			return reject(
				ErrInvalidDescriptor,
				"unit %q requires a component and manifests",
				name,
			)

		case unit.Component.IsWithin(configurationRoot):
			return reject(
				ErrInvalidDescriptor,
				"unit %s source must remain outside runtime configuration",
				name,
			)
		}

		if previous, shared := unitByManifests[unit.Manifests]; shared {
			return reject(
				ErrInvalidDescriptor,
				"units %s and %s share the same manifest path",
				previous, name,
			)
		}

		unitByManifests[unit.Manifests] = name
	}

	return nil
}

func (descriptor Descriptor) validateHooks() error {
	for name, script := range descriptor.Hooks {
		if name.IsZero() || !script.HasSuffix(".sh") {
			return reject(
				ErrInvalidDescriptor,
				"hook %q must reference a .sh file",
				name,
			)
		}
	}

	return nil
}
