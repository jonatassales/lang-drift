package microruntime_test

import (
	"errors"
	"testing"

	"github.com/exphere/ex/src/domain/microruntime"
)

// quantity parses text that the test knows to be valid.
func quantity(t testing.TB, text string) microruntime.Quantity {
	t.Helper()

	parsed, err := microruntime.ParseQuantity(text)
	if err != nil {
		t.Fatal(err)
	}

	return parsed
}

func TestParseQuantity(t *testing.T) {
	t.Parallel()

	valid := []string{"1", "0", "100m", "0.5", ".5", "128Mi", "1Gi", "2k", "1M", "1e3", "1E-2", "4Pi"}
	for _, text := range valid {
		got, err := microruntime.ParseQuantity(text)
		if err != nil || got.String() != text || got.IsZero() {
			t.Errorf("ParseQuantity(%q) = %q, %v", text, got, err)
		}
	}

	invalid := []string{"", "-1", "+1", "1K", "1Mb", "100 m", " 100m", "1.5.0", "m", "Mi", "1e", "one"}
	for _, text := range invalid {
		got, err := microruntime.ParseQuantity(text)
		if !errors.Is(err, microruntime.ErrInvalidQuantity) || !got.IsZero() {
			t.Errorf("ParseQuantity(%q) = %q, %v; want the zero value and ErrInvalidQuantity", text, got, err)
		}
	}
}

func TestQuotaResource_Supported(t *testing.T) {
	t.Parallel()

	supported := []microruntime.QuotaResource{
		microruntime.QuotaRequestsCPU, microruntime.QuotaLimitsCPU,
		microruntime.QuotaRequestsMemory, microruntime.QuotaLimitsMemory,
		microruntime.QuotaRequestsEphemeralStorage, microruntime.QuotaLimitsEphemeralStorage,
		microruntime.QuotaPods, microruntime.QuotaJobs, microruntime.QuotaDeployments, microruntime.QuotaServices,
	}
	for _, resource := range supported {
		if !resource.Supported() {
			t.Errorf("%s.Supported() = false", resource)
		}
	}

	for _, resource := range []microruntime.QuotaResource{"", "count/secrets", "requests.nvidia.com/gpu", "services.loadbalancers"} {
		if resource.Supported() {
			t.Errorf("%q.Supported() = true", resource)
		}
	}
}
