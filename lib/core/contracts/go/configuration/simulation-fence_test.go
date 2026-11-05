package configuration_test

import (
	"testing"

	"github.com/jonatassales/tao/lib/core/contracts/go/configuration"
)

// simulationFence is the fence the simulation's plan drafts, read from the
// published sample.
func simulationFence(t *testing.T) configuration.Fence {
	t.Helper()

	fence, err := configuration.ParseFence(fenceMarkdown(t, "runtime.sample.json"))
	if err != nil {
		t.Fatal(err)
	}

	return fence
}
