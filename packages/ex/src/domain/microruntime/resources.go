package microruntime

import "regexp"

var quantityPattern = regexp.MustCompile(
	`^(?:[0-9]+(?:\.[0-9]+)?|\.[0-9]+)(?:[eE][+-]?[0-9]+|[numkMGTPE]|[KMGTPE]i)?$`,
)

// Quantity is a Kubernetes resource quantity such as "100m" or "128Mi".
// The zero Quantity means "not set"; it is not the quantity zero.
type Quantity struct{ value string }

// ParseQuantity validates text as a Quantity.
func ParseQuantity(text string) (Quantity, error) {
	if !quantityPattern.MatchString(text) {
		return Quantity{}, reject(
			ErrInvalidQuantity,
			"%q is not a quantity such as 100m or 128Mi",
			text,
		)
	}

	return Quantity{value: text}, nil
}

func (quantity Quantity) String() string { return quantity.value }

// IsZero reports whether the quantity was never set.
func (quantity Quantity) IsZero() bool { return quantity.value == "" }

// Resources are the explicit CPU and memory bounds of one workload.
type Resources struct {
	CPURequest    Quantity
	CPULimit      Quantity
	MemoryRequest Quantity
	MemoryLimit   Quantity
}

func (resources Resources) isComplete() bool {
	return !resources.CPURequest.IsZero() &&
		!resources.CPULimit.IsZero() &&
		!resources.MemoryRequest.IsZero() &&
		!resources.MemoryLimit.IsZero()
}

// QuotaResource names a namespace quota that EX supports.
type QuotaResource string

const (
	QuotaRequestsCPU              QuotaResource = "requests.cpu"
	QuotaLimitsCPU                QuotaResource = "limits.cpu"
	QuotaRequestsMemory           QuotaResource = "requests.memory"
	QuotaLimitsMemory             QuotaResource = "limits.memory"
	QuotaRequestsEphemeralStorage QuotaResource = "requests.ephemeral-storage"
	QuotaLimitsEphemeralStorage   QuotaResource = "limits.ephemeral-storage"
	QuotaPods                     QuotaResource = "pods"
	QuotaJobs                     QuotaResource = "count/jobs.batch"
	QuotaDeployments              QuotaResource = "count/deployments.apps"
	QuotaServices                 QuotaResource = "services"
)

// Supported reports whether EX can enforce this quota.
func (resource QuotaResource) Supported() bool {
	switch resource {
	case QuotaRequestsCPU, QuotaLimitsCPU,
		QuotaRequestsMemory, QuotaLimitsMemory,
		QuotaRequestsEphemeralStorage, QuotaLimitsEphemeralStorage,
		QuotaPods, QuotaJobs, QuotaDeployments, QuotaServices:
		return true
	}

	return false
}

// ResourceQuotas bound what a runtime namespace may consume.
// They limit resources, not monetary spend.
type ResourceQuotas map[QuotaResource]Quantity
