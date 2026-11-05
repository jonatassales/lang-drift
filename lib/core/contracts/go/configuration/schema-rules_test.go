package configuration_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/jonatassales/tao/lib/core/contracts/go/configuration"
)

// schemaRules is the part of the published schema that states the rules a
// JSON Schema can state.
type schemaRules struct {
	AdditionalProperties bool `json:"additionalProperties"`
	Properties           struct {
		CPUMillis struct {
			Maximum int `json:"maximum"`
		} `json:"cpuMillis"`
		MemoryMiB struct {
			Maximum int `json:"maximum"`
		} `json:"memoryMiB"`
		Effort struct {
			Enum []string `json:"enum"`
		} `json:"effort"`
		Components struct {
			MaxItems int `json:"maxItems"`
			Items    struct {
				AdditionalProperties bool `json:"additionalProperties"`
				Properties           struct {
					Replicas struct {
						Const int `json:"const"`
					} `json:"replicas"`
				} `json:"properties"`
			} `json:"items"`
		} `json:"components"`
	} `json:"properties"`
}

// The schema states the same rules Validate enforces wherever a JSON Schema
// can: no field beyond the version's, the efforts, the bound's limits, at
// most 16 components and one replica each. The rules a schema cannot state
// (unique paths, distinct ports and images, shares within the bound) are
// Validate's alone; the module stays on the standard library, so no schema
// validator runs the samples.
func TestFenceSchemaStatesTheRulesValidateEnforces(t *testing.T) {
	t.Parallel()

	raw, err := os.ReadFile(filepath.Join("schema", "runtime.schema.json"))
	if err != nil {
		t.Fatal(err)
	}

	var rules schemaRules
	if err := json.Unmarshal(raw, &rules); err != nil {
		t.Fatal(err)
	}

	closed := !rules.AdditionalProperties && !rules.Properties.Components.Items.AdditionalProperties
	bounds := rules.Properties.CPUMillis.Maximum == 64000 && rules.Properties.MemoryMiB.Maximum == 1048576
	components := rules.Properties.Components.MaxItems == 16 && rules.Properties.Components.Items.Properties.Replicas.Const == 1

	if !closed || !bounds || !components || !slices.Equal(rules.Properties.Effort.Enum, configuration.FenceEfforts) {
		t.Errorf("the schema's rules = %+v", rules)
	}

	fence := simulationFence(t)
	fence.CPUMillis = 64001
	fence.Effort = "max"
	fence.Components[0].Replicas = 2

	if problems := fence.Validate(); len(problems) != 3 {
		t.Errorf("Validate found %v where the schema refuses three fields", problems)
	}
}
