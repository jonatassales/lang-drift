package microruntime

import (
	"net/url"
	"strings"

	"github.com/jonatassales/tao/lib/core/contracts/go/configuration"
)

// WorkloadKind selects how Kubernetes runs a workload.
type WorkloadKind string

const (
	DeploymentWorkload WorkloadKind = "Deployment"
	JobWorkload        WorkloadKind = "Job"
)

// IsSupported reports whether EX can run workloads of this kind.
func (kind WorkloadKind) IsSupported() bool {
	return kind == DeploymentWorkload || kind == JobWorkload
}

// ServiceSpec describes how a Deployment is exposed and kept ready.
// A Job has no service, so its ServiceSpec must stay empty.
type ServiceSpec struct {
	Port     int
	Replicas int
	Probe    Probe
}

// isZero reports whether the spec names nothing, as a Job's must.
func (service ServiceSpec) isZero() bool {
	return service.Port == 0 &&
		service.Replicas == 0 &&
		service.Probe.Kind == "" &&
		len(service.Probe.Command) == 0
}

// Workload is the initial definition of one unit; later edits belong to the
// native Kubernetes files. Component is relative to the project root.
// EnvSecret, when set, names the Secret the container reads its
// environment from.
type Workload struct {
	Name      Identifier
	Component RelativePath
	Image     string
	Kind      WorkloadKind
	Resources Resources
	Service   ServiceSpec
	EnvSecret Identifier
}

// Validate enforces the rules every workload must satisfy before generation.
func (workload Workload) Validate() error {
	switch {
	case workload.Name.IsZero() || workload.Component.IsZero():
		return reject(ErrInvalidBlueprint, "every workload requires a name and a component")

	case workload.Component.IsWithin(configurationRoot):
		return reject(
			ErrInvalidBlueprint,
			"workload %s source must remain outside runtime configuration",
			workload.Name,
		)

	case !IsImageReference(workload.Image):
		return reject(
			ErrInvalidBlueprint,
			"workload %s requires an existing image reference",
			workload.Name,
		)

	case !workload.Resources.isComplete():
		return reject(
			ErrInvalidBlueprint,
			"workload %s requires explicit CPU and memory quantities",
			workload.Name,
		)
	}

	switch workload.Kind {
	case DeploymentWorkload:
		return workload.validateDeployment()

	case JobWorkload:
		if !workload.Service.isZero() {
			return reject(
				ErrInvalidBlueprint,
				"job %s must not specify service port, replicas, or readiness",
				workload.Name,
			)
		}

		return nil

	default:
		return reject(
			ErrInvalidBlueprint,
			"unsupported workload kind %q; use %s or %s",
			workload.Kind, DeploymentWorkload, JobWorkload,
		)
	}
}

func (workload Workload) validateDeployment() error {
	service := workload.Service

	if !ValidPort(service.Port) || service.Replicas < 1 {
		return reject(
			ErrInvalidBlueprint,
			"workload %s requires a valid port and positive replicas",
			workload.Name,
		)
	}

	if err := service.Probe.Validate(); err != nil {
		return reject(ErrInvalidBlueprint, "workload %s: %v", workload.Name, err)
	}

	return nil
}

// IsImageReference reports whether text can name one already-built image.
func IsImageReference(text string) bool {
	return configuration.IsImageReference(text)
}

// isReadinessPath accepts an absolute HTTP path such as "/healthz".
func isReadinessPath(text string) bool {
	parsed, err := url.ParseRequestURI(text)

	return err == nil &&
		strings.HasPrefix(text, "/") &&
		!parsed.IsAbs() &&
		parsed.Host == ""
}
