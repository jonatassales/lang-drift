package handoff_test

import (
	"encoding/json"
	"testing"

	"github.com/jonatassales/tao/lib/core/contracts/go/handoff"
)

// Before the run has an id and an outcome, its record says so with null;
// sections the run has not reached are left out, and every section of spec
// 0065's run.json round trips under its key, the dev ticket among the
// stages.
func TestRunRecordKeepsSpec0065Keys(t *testing.T) {
	t.Parallel()

	encoded, err := json.Marshal(handoff.RunRecord{})
	if err != nil || string(encoded) != `{"run":null,"outcome":null}` {
		t.Errorf("encoded = %s, %v", encoded, err)
	}

	sections := `{"run":"book-a-demo-runtime-0001","intent":{},"stages":[{"report":{}},{}],"effort":{},"safeSeat":{},"runtime":{},"delivery":{"pullRequest":{}},"outcome":"completed"}`

	record, err := handoff.Decode[handoff.RunRecord](sections)
	if err != nil {
		t.Fatal(err)
	}

	completed := record.Outcome != nil && *record.Outcome == handoff.RunCompleted
	if !completed {
		t.Errorf("outcome = %v", record.Outcome)
	}

	encoded, err = json.Marshal(record)
	if err != nil || string(encoded) != sections {
		t.Errorf("round trip = %s, %v; want %s", encoded, err, sections)
	}
}
