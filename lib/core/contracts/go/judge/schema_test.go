package judge_test

import (
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// TestSchemasAndSamplesAgree checks the files the console reads: every
// schema parses and names itself, and has a sample that carries exactly the
// top-level fields the schema requires. Whether a sample's values
// are an answer is ParseAnswer's to decide.
func TestSchemasAndSamplesAgree(t *testing.T) {
	t.Parallel()

	schemas, err := filepath.Glob(filepath.Join("schema", "*.schema.json"))
	if err != nil {
		t.Fatal(err)
	}

	if len(schemas) != 6 {
		t.Fatalf("schemas = %v, want the request, the response and four answers", schemas)
	}

	for _, path := range schemas {
		schema := readObject(t, path)

		identifier, isText := schema["$id"].(string)
		namesItself := isText && strings.HasSuffix(identifier, "/"+filepath.Base(path))

		if !namesItself {
			t.Errorf("%s names itself %q", path, identifier)
		}

		sample := strings.TrimSuffix(path, ".schema.json") + ".sample.json"
		required := requiredFields(t, schema, path)
		fields := keys(readObject(t, sample))
		if !slices.Equal(fields, required) {
			t.Errorf("%s has %v, want %v", sample, fields, required)
		}
	}
}
