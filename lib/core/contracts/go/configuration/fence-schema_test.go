package configuration_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/jonatassales/tao/lib/core/contracts/go/configuration"
)

// fenceSchema is the part of the published schema the Go contract mirrors.
type fenceSchema struct {
	ID         string   `json:"$id"`
	Required   []string `json:"required"`
	Properties struct {
		Components struct {
			Items struct {
				Required   []string `json:"required"`
				Properties struct {
					Stack struct {
						Enum []string `json:"enum"`
					} `json:"stack"`
					Role struct {
						Enum []string `json:"enum"`
					} `json:"role"`
				} `json:"properties"`
			} `json:"items"`
		} `json:"components"`
	} `json:"properties"`
}

// The published schema and the Go contract agree on the same sample: the
// schema names itself, requires exactly the sample's fields at both levels,
// closes the stacks and roles to the Go sets, and the sample validates in Go
// with no problem.
func TestFenceSchemaAgreesWithTheGoContract(t *testing.T) {
	t.Parallel()

	raw, err := os.ReadFile(filepath.Join("schema", "runtime.schema.json"))
	if err != nil {
		t.Fatal(err)
	}

	var schema fenceSchema
	if err := json.Unmarshal(raw, &schema); err != nil {
		t.Fatal(err)
	}

	if schema.ID != "https://langdrift.md/schemas/configuration/runtime.schema.json" {
		t.Errorf("$id = %q", schema.ID)
	}

	sampleRaw, err := os.ReadFile(filepath.Join("schema", "runtime.sample.json"))
	if err != nil {
		t.Fatal(err)
	}

	var sample struct {
		Fields     map[string]json.RawMessage
		Components []map[string]json.RawMessage `json:"components"`
	}
	if err := json.Unmarshal(sampleRaw, &sample.Fields); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(sampleRaw, &sample); err != nil {
		t.Fatal(err)
	}

	if !sameKeys(schema.Required, sample.Fields) || !sameKeys(schema.Properties.Components.Items.Required, sample.Components[0]) {
		t.Errorf("the schema requires %v and %v; the sample has other fields", schema.Required, schema.Properties.Components.Items.Required)
	}

	stacks := schema.Properties.Components.Items.Properties.Stack.Enum
	roles := schema.Properties.Components.Items.Properties.Role.Enum
	if !slices.Equal(stacks, configuration.FenceStacks) || !slices.Equal(roles, configuration.FenceRoles) {
		t.Errorf("the schema's stacks %v and roles %v differ from the Go sets", stacks, roles)
	}

	fence, err := configuration.ParseFence(fenceMarkdown(t, "runtime.sample.json"))
	if err != nil {
		t.Fatal(err)
	}

	if problems := fence.Validate(); len(problems) != 0 {
		t.Errorf("the sample breaks %v", problems)
	}
}
