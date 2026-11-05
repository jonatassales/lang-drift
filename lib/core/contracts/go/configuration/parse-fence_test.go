package configuration_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/jonatassales/tao/lib/core/contracts/go/configuration"
)

// The simulation's fence reads as version 3, with the five components of the
// Plan card's COMPONENTS rows in their order.
func TestParseFenceReadsTheSimulationsFence(t *testing.T) {
	t.Parallel()

	fence, err := configuration.ParseFence(fenceMarkdown(t, "runtime.sample.json"))
	if err != nil {
		t.Fatal(err)
	}

	bound := fence.Version == configuration.FenceVersion && fence.CPUMillis == 2000 && fence.MemoryMiB == 3072
	if !bound || fence.RuntimeID != "book-a-demo-runtime" || fence.Revision != 1 || len(fence.Components) != 5 {
		t.Fatalf("fence = %+v", fence)
	}

	want := configuration.FenceComponent{
		Path:      "services/signup-api",
		Stack:     "Go",
		Role:      "api",
		Port:      8080,
		Image:     "book-a-demo-runtime/signup-api:r1",
		CPUMillis: 250,
		MemoryMiB: 256,
		Replicas:  1,
	}
	if fence.Components[0] != want {
		t.Errorf("first component = %+v, want %+v", fence.Components[0], want)
	}
}

// A field the version does not define, and a fence of another version, are
// refused rather than misread.
func TestParseFenceRefusesAnotherShape(t *testing.T) {
	t.Parallel()

	sample := fenceMarkdown(t, "runtime.sample.json")

	for name, markdown := range map[string]string{
		"an unknown field": strings.Replace(sample, `"version": 3,`, `"version": 3, "runtimeUid": "uid",`, 1),
		"version 2":        strings.Replace(sample, `"version": 3,`, `"version": 2,`, 1),
		"two fences":       sample + sample,
	} {
		_, err := configuration.ParseFence(markdown)
		if !errors.Is(err, configuration.ErrInvalidConfiguration) {
			t.Errorf("%s: err = %v", name, err)
		}
	}

	_, err := configuration.ParseFence(strings.Replace(sample, `"version": 3,`, `"version": 2,`, 1))
	if err == nil || !strings.Contains(err.Error(), "re-plan the runtime") {
		t.Errorf("a version-2 fence gave %v", err)
	}
}
