package microruntime

import (
	"fmt"
	"sort"
	"strings"
)

// Ownership labels mark every object an execution creates.
const (
	LabelRuntime        = "exphere.dev/runtime"
	LabelUnit           = "exphere.dev/unit"
	LabelOwnershipToken = "exphere.dev/ownership-token"
)

// Job containers receive these variables so their output can be correlated.
const (
	EnvironmentRuntimeID = "EXPHERE_RUNTIME_ID"
	EnvironmentRunID     = "EXPHERE_RUN_ID"
)

// apiVersions lists the kinds EX creates and the one API version accepted for
// each: stateless Deployments, ClusterIP Services, and finite Jobs.
var apiVersions = map[string]string{
	"Deployment": "apps/v1",
	"Service":    "v1",
	"Job":        "batch/v1",
}

// supportedVolumeSources excludes persistent volumes, host paths, generic
// ephemeral claims, and inline network storage.
var supportedVolumeSources = map[string]bool{
	"emptyDir":    true,
	"configMap":   true,
	"secret":      true,
	"projected":   true,
	"downwardAPI": true,
}

// Manifest binds one native Kubernetes object to the unit it belongs to.
type Manifest struct {
	Unit   Identifier
	Object Object
}

// Key identifies the manifest's object within one runtime.
func (manifest Manifest) Key() string {
	return manifest.Object.Kind() + "/" + manifest.Object.Name()
}

// ValidateManifests checks the native objects read for this descriptor: each
// must satisfy the manifest policy, no object may be declared twice, and every
// unit must contribute at least one object.
func (descriptor Descriptor) ValidateManifests(manifests []Manifest) error {
	declared := make(map[string]bool, len(manifests))
	perUnit := make(map[Identifier]int, len(descriptor.Units))

	for _, manifest := range manifests {
		if _, known := descriptor.Units[manifest.Unit]; !known {
			return reject(ErrUnsupportedManifest, "%s belongs to undeclared unit %s", manifest.Key(), manifest.Unit)
		}

		if err := manifest.validate(descriptor.Target.Namespace); err != nil {
			return fmt.Errorf("unit %s: %w", manifest.Unit, err)
		}

		if declared[manifest.Key()] {
			return reject(ErrUnsupportedManifest, "duplicate object %s", manifest.Key())
		}

		declared[manifest.Key()] = true
		perUnit[manifest.Unit]++
	}

	for _, name := range descriptor.UnitNames() {
		if perUnit[name] == 0 {
			return reject(ErrUnsupportedManifest, "unit %s has no supported manifests", name)
		}
	}

	return nil
}

// validate enforces the stateless pilot policy on one object.
func (manifest Manifest) validate(namespace Identifier) error {
	object := manifest.Object
	metadata := object.Metadata()

	if _, err := ParseIdentifier(metadata.Text("name")); err != nil {
		return reject(ErrUnsupportedManifest, "invalid explicit resource name %q", metadata.Text("name"))
	}

	if declared := metadata.Text("namespace"); declared != "" && declared != namespace.String() {
		return reject(ErrUnsupportedManifest, "manifest namespace %q differs from the descriptor", declared)
	}

	if metadata.Text("generateName") != "" || metadata.Text("uid") != "" {
		return reject(ErrUnsupportedManifest, "generated names and authored UIDs are unsupported")
	}

	expected, supported := apiVersions[object.Kind()]
	if !supported {
		return reject(ErrUnsupportedManifest, "unsupported resource kind %q", object.Kind())
	}

	if object.APIVersion() != expected {
		return reject(ErrUnsupportedManifest, "%s apiVersion must be %s", object.Kind(), expected)
	}

	switch object.Kind() {
	case "Service":
		return validateService(object.Spec())

	case "Job":
		if err := validateJob(object.Spec()); err != nil {
			return err
		}
	}

	return validatePod(object.PodSpec())
}

func validateService(spec Object) error {
	if serviceType := spec.Text("type"); serviceType != "" && serviceType != "ClusterIP" {
		return reject(ErrUnsupportedManifest, "only ClusterIP Services are supported")
	}

	if len(asSlice(spec["externalIPs"])) > 0 || spec.Text("externalName") != "" {
		return reject(ErrUnsupportedManifest, "external service targets are unsupported")
	}

	return nil
}

func validateJob(spec Object) error {
	if !isPositiveInteger(spec["activeDeadlineSeconds"]) {
		return reject(ErrUnsupportedManifest, "finite Jobs require a positive activeDeadlineSeconds")
	}

	if spec.Flag("suspend") {
		return reject(ErrUnsupportedManifest, "suspended Jobs cannot complete")
	}

	return nil
}

func validatePod(pod Object) error {
	for _, key := range []string{"hostNetwork", "hostPID", "hostIPC"} {
		if pod.Flag(key) {
			return reject(ErrUnsupportedManifest, "host sharing is unsupported")
		}
	}

	if err := validateVolumes(pod); err != nil {
		return err
	}

	return validateContainers(pod)
}

func validateVolumes(pod Object) error {
	volumes, ok := strictSlice(pod["volumes"])
	if !ok {
		return reject(ErrUnsupportedManifest, "pod volumes must be an array")
	}

	for _, value := range volumes {
		volume := asObject(value)
		if _, err := ParseIdentifier(volume.Text("name")); err != nil {
			return reject(ErrUnsupportedManifest, "pod volumes require an explicit valid name")
		}

		sources := 0

		for key, source := range volume {
			if key == "name" {
				continue
			}

			if !supportedVolumeSources[key] {
				return reject(ErrUnsupportedManifest, "volume source %s is outside the stateless pilot", key)
			}

			if asObject(source) == nil {
				return reject(ErrUnsupportedManifest, "volume source %s must be an object", key)
			}

			sources++
		}

		if sources != 1 {
			return reject(ErrUnsupportedManifest, "pod volumes require exactly one supported source")
		}
	}

	return nil
}

func validateContainers(pod Object) error {
	if len(pod.Children("containers")) == 0 {
		return reject(ErrUnsupportedManifest, "workload requires containers")
	}

	for _, container := range pod.Containers() {
		switch {
		case strings.TrimSpace(container.Text("image")) == "":
			return reject(ErrUnsupportedManifest, "container image is required")

		case container.Child("securityContext").Flag("privileged"):
			return reject(ErrUnsupportedManifest, "privileged containers are unsupported")

		case container.Text("imagePullPolicy") != "Never":
			return reject(ErrUnsupportedManifest, "prebuilt local images require imagePullPolicy: Never")
		}
	}

	return nil
}

// isPositiveInteger accepts the numeric types JSON and YAML decode to.
func isPositiveInteger(value any) bool {
	switch number := value.(type) {
	case int:
		return number > 0
	case int64:
		return number > 0
	case uint64:
		return number > 0
	case float64:
		return number > 0 && number == float64(int64(number))
	}

	return false
}

// Claim marks the objects of one execution. Every claimed object, and the pod
// template of every workload, carries the runtime, unit, and ownership-token
// labels; Job containers also receive the runtime and run IDs.
type Claim struct {
	Runtime   Identifier
	Namespace Identifier
	Token     string
	RunID     string
}

// Apply rewrites the manifest's object in place so it belongs to the claim.
func (claim Claim) Apply(manifest Manifest) {
	object := manifest.Object

	metadata := object.ensureChild("metadata")
	metadata["namespace"] = claim.Namespace.String()
	claim.label(metadata, manifest.Unit)

	if object.Kind() == "Service" {
		return
	}

	template := object.ensureChild("spec").ensureChild("template")
	claim.label(template.ensureChild("metadata"), manifest.Unit)

	if object.Kind() == "Job" {
		claim.correlate(template.Child("spec"))
	}
}

func (claim Claim) label(metadata Object, unit Identifier) {
	labels := metadata.ensureChild("labels")
	labels[LabelRuntime] = claim.Runtime.String()
	labels[LabelOwnershipToken] = claim.Token
	labels[LabelUnit] = unit.String()
}

// correlate replaces any authored correlation variables in each container.
func (claim Claim) correlate(pod Object) {
	for _, container := range pod.Children("containers") {
		if container == nil {
			continue
		}

		environment := []any{}

		for _, variable := range asSlice(container["env"]) {
			name := asObject(variable).Text("name")
			if name != EnvironmentRuntimeID && name != EnvironmentRunID {
				environment = append(environment, variable)
			}
		}

		container["env"] = append(
			environment,
			map[string]any{"name": EnvironmentRuntimeID, "value": claim.Runtime.String()},
			map[string]any{"name": EnvironmentRunID, "value": claim.RunID},
		)
	}
}

// Images returns the distinct images the workloads need, sorted so a failing
// preflight names the same image on every run.
func Images(manifests []Manifest) []string {
	unique := map[string]bool{}
	for _, container := range ConfiguredContainers(manifests) {
		unique[container.Image] = true
	}

	images := make([]string, 0, len(unique))
	for image := range unique {
		images = append(images, image)
	}

	sort.Strings(images)

	return images
}
