package judgment_test

import (
	"encoding/json"
	"slices"
	"testing"

	"github.com/jonatassales/tao/lib/core/contracts/go/configuration"
	"github.com/jonatassales/tao/lib/core/contracts/go/judgment"
)

// answerSchema is the part of a published answer schema the Go contract
// mirrors.
type answerSchema struct {
	ID         string                     `json:"$id"`
	Required   []string                   `json:"required"`
	Properties map[string]json.RawMessage `json:"properties"`
}

// The three schemas name themselves, require exactly their sample's fields,
// define no price, and close the effort, the verdicts and the reasons to
// the Go sets.
func TestAnswerSchemasAgreeWithTheGoContract(t *testing.T) {
	t.Parallel()

	for _, kind := range []string{"criteria", "episodes", "suggestion"} {
		var schema answerSchema
		if err := json.Unmarshal([]byte(readSample(t, kind+".answer.schema.json")), &schema); err != nil {
			t.Fatal(err)
		}

		var sample map[string]json.RawMessage
		if err := json.Unmarshal([]byte(readSample(t, kind+".answer.sample.json")), &sample); err != nil {
			t.Fatal(err)
		}

		names := make([]string, 0, len(sample))
		for name := range sample {
			names = append(names, name)
		}

		slices.Sort(names)
		required := slices.Sorted(slices.Values(schema.Required))

		if schema.ID != "https://langdrift.md/schemas/judgment/"+kind+".answer.schema.json" || !slices.Equal(names, required) {
			t.Errorf("%s: $id %q, required %v, sample %v", kind, schema.ID, schema.Required, names)
		}

		for _, priced := range []string{"cost", "usage", "provenance", "judgeVersion"} {
			if _, found := schema.Properties[priced]; found {
				t.Errorf("%s defines %s; the plan pays for nothing", kind, priced)
			}
		}
	}

	var suggestion struct {
		Properties struct {
			Effort struct {
				Enum []string `json:"enum"`
			} `json:"effort"`
		} `json:"properties"`
	}
	if err := json.Unmarshal([]byte(readSample(t, "suggestion.answer.schema.json")), &suggestion); err != nil {
		t.Fatal(err)
	}

	if !slices.Equal(suggestion.Properties.Effort.Enum, configuration.FenceEfforts) {
		t.Errorf("efforts = %v", suggestion.Properties.Effort.Enum)
	}

	var episodes struct {
		Properties struct {
			Episodes struct {
				Items struct {
					Properties struct {
						Verdict struct {
							Enum []string `json:"enum"`
						} `json:"verdict"`
						Reason struct {
							Enum []string `json:"enum"`
						} `json:"reason"`
					} `json:"properties"`
				} `json:"items"`
			} `json:"episodes"`
		} `json:"properties"`
	}
	if err := json.Unmarshal([]byte(readSample(t, "episodes.answer.schema.json")), &episodes); err != nil {
		t.Fatal(err)
	}

	item := episodes.Properties.Episodes.Items.Properties
	reasons := append([]string{judgment.KeptReason}, judgment.DroppedReasons...)
	if !slices.Equal(item.Verdict.Enum, []string{judgment.Kept, judgment.Dropped}) || !slices.Equal(item.Reason.Enum, reasons) {
		t.Errorf("verdicts %v, reasons %v", item.Verdict.Enum, item.Reason.Enum)
	}
}
