package microruntime_test

import (
	"encoding/json"
	"testing"

	"github.com/jonatassales/tao/lib/cli/src/domain/microruntime"
)

func TestEnvironment_Manifest(t *testing.T) {
	t.Parallel()

	environment := microruntime.Environment{
		Unit:   identifier(t, "demo-requests"),
		Secret: identifier(t, "demo-requests-env"),
		Variables: map[string]string{
			"LOG_URL":   "http://log:3500",
			"API_TOKEN": "from-the-agent",
		},
	}

	manifest := environment.Manifest()

	rendered, err := json.Marshal(manifest.Object)
	if err != nil {
		t.Fatal(err)
	}

	want := `{"apiVersion":"v1","kind":"Secret","metadata":{"name":"demo-requests-env"},"stringData":{"API_TOKEN":"from-the-agent","LOG_URL":"http://log:3500"},"type":"Opaque"}`
	if manifest.Unit.String() != "demo-requests" || string(rendered) != want {
		t.Fatalf("Manifest() = %s for %s, want %s", rendered, manifest.Unit, want)
	}

	claim := microruntime.Claim{
		Runtime:   identifier(t, "demo"),
		Namespace: identifier(t, "demo"),
		Token:     "token",
		RunID:     "run",
	}
	claim.Apply(manifest)

	if manifest.Object.Labels().Text(microruntime.LabelUnit) != "demo-requests" || manifest.Object.Namespace() != "demo" {
		t.Fatalf("claimed Secret metadata = %v", manifest.Object.Metadata())
	}

	if _, found := manifest.Object["spec"]; found {
		t.Fatalf("claiming a Secret gave it a pod template: %v", manifest.Object)
	}
}
