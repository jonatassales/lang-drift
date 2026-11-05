package handoff_test

import (
	"encoding/json"
	"testing"

	"github.com/jonatassales/tao/lib/core/contracts/go/handoff"
)

// A value not known yet decodes as null and encodes as null, never as zero.
func TestEvidenceKeepsAnUnknownValueNull(t *testing.T) {
	t.Parallel()

	text := `{"file":"runtime.md","state":"draft","revision":null,"proposedTo":"console"}`

	evidence, err := handoff.Decode[handoff.Evidence](text)
	if err != nil || evidence.Revision != nil {
		t.Fatalf("evidence = %+v, %v; want an unknown revision", evidence, err)
	}

	encoded, err := json.Marshal(evidence)
	if err != nil || string(encoded) != text {
		t.Errorf("encoded = %s, %v; want %s", encoded, err, text)
	}
}
