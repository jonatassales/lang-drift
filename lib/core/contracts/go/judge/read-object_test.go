package judge_test

import (
	"encoding/json"
	"os"
	"testing"
)

// readObject reads a JSON file into a map, failing the test when it cannot.
func readObject(t *testing.T, path string) map[string]any {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	var object map[string]any
	if err := json.Unmarshal(data, &object); err != nil {
		t.Fatalf("%s: %v", path, err)
	}

	return object
}
