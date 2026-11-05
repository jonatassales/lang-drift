package judge_test

import (
	"slices"
	"testing"
)

// requiredFields lists, sorted, the fields a schema requires.
func requiredFields(t *testing.T, schema map[string]any, path string) []string {
	t.Helper()

	listed, ok := schema["required"].([]any)
	if !ok {
		t.Fatalf("%s lists no required fields", path)
	}

	fields := make([]string, 0, len(listed))
	for _, field := range listed {
		name, isText := field.(string)
		if !isText {
			t.Fatalf("%s requires a field that is not named by a string", path)
		}

		fields = append(fields, name)
	}

	slices.Sort(fields)

	return fields
}
