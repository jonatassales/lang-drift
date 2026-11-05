package judge_test

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// requestSchemaViolations checks a request body against the contract's own
// request.schema.json, the file the console's tests read too, found through
// the module that ships it.
func requestSchemaViolations(t *testing.T, body []byte) []string {
	t.Helper()

	located, err := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", "github.com/jonatassales/tao/lib/core/contracts/go").Output()
	if err != nil {
		t.Fatalf("locating the contracts module: %v", err)
	}

	raw, err := os.ReadFile(filepath.Join(strings.TrimSpace(string(located)), "judge", "schema", "request.schema.json"))
	if err != nil {
		t.Fatal(err)
	}

	var schema requestSchema
	if err := json.Unmarshal(raw, &schema); err != nil {
		t.Fatal(err)
	}

	var request map[string]any
	if err := json.Unmarshal(body, &request); err != nil {
		return []string{"the body is not a JSON object: " + err.Error()}
	}

	var violations []string

	for _, field := range schema.Required {
		if _, present := request[field]; !present {
			violations = append(violations, field+" is missing")
		}
	}

	for field, value := range request {
		property, known := schema.Properties[field]
		if !known {
			violations = append(violations, field+" is not a field of the schema")

			continue
		}

		violations = append(violations, propertyViolations(field, property, value)...)
	}

	return violations
}
