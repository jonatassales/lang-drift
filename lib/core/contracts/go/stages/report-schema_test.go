package stages_test

import (
	"encoding/json"
	"slices"
	"testing"
)

// reportSchema is the part of a published report schema the Go check
// mirrors.
type reportSchema struct {
	ID                   string   `json:"$id"`
	AdditionalProperties bool     `json:"additionalProperties"`
	Required             []string `json:"required"`
}

// Each stage's schema names itself, defines no field beyond the report's,
// and requires exactly its sample's fields.
func TestReportSchemasNameThemselvesAndRequireTheSamplesFields(t *testing.T) {
	t.Parallel()

	for _, stage := range []string{"planning", "spec", "tickets", "implement", "review"} {
		var schema reportSchema
		if err := json.Unmarshal([]byte(readSample(t, stage+".report.schema.json")), &schema); err != nil {
			t.Fatal(err)
		}

		var sample map[string]json.RawMessage
		if err := json.Unmarshal([]byte(readSample(t, stage+".report.sample.json")), &sample); err != nil {
			t.Fatal(err)
		}

		fields := make([]string, 0, len(sample))
		for field := range sample {
			fields = append(fields, field)
		}

		slices.Sort(fields)
		required := slices.Sorted(slices.Values(schema.Required))
		named := schema.ID == "https://langdrift.md/schemas/stages/"+stage+".report.schema.json"

		if !named || schema.AdditionalProperties || !slices.Equal(fields, required) {
			t.Errorf("%s: $id %q, required %v, sample %v", stage, schema.ID, schema.Required, fields)
		}
	}
}
