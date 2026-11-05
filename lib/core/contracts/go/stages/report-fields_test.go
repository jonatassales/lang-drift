package stages

import (
	"encoding/json"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

// fieldsOfSchema is the part of a published report schema whose fields the
// Go check must hold.
type fieldsOfSchema struct {
	Required   []string                   `json:"required"`
	Properties map[string]json.RawMessage `json:"properties"`
}

// The fields the Go check requires of each stage are the fields its
// published schema requires and defines, so the two refuse the same missing
// and unknown fields.
func TestReportFieldsAreTheSchemasFields(t *testing.T) {
	t.Parallel()

	for stage, fields := range reportFields {
		text, err := os.ReadFile(filepath.Join("schema", string(stage)+".report.schema.json"))
		if err != nil {
			t.Fatal(err)
		}

		var schema fieldsOfSchema
		if err := json.Unmarshal(text, &schema); err != nil {
			t.Fatal(err)
		}

		want := slices.Sorted(slices.Values(fields))
		required := slices.Sorted(slices.Values(schema.Required))
		defined := slices.Sorted(maps.Keys(schema.Properties))

		agree := slices.Equal(want, required) && slices.Equal(want, defined)
		if !agree {
			t.Errorf("%s: Go %v, schema requires %v and defines %v", stage, want, required, defined)
		}
	}
}
