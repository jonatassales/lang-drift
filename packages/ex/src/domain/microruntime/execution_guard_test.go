package microruntime

import (
	"errors"
	"testing"
)

func TestExecution_RefusesTransitionsItNeverClaimed(t *testing.T) {
	execution := StartExecution(Start{Runtime: mustIdentifier(t, "demo"), Token: "token", RunID: "run"})

	if err := execution.MarkUnit(mustIdentifier(t, "ghost"), UnitFailed); !errors.Is(err, ErrInconsistentExecution) {
		t.Fatalf("marking an unclaimed unit: %v", err)
	}
	if err := execution.RecordCreated(0, nil); !errors.Is(err, ErrInconsistentExecution) {
		t.Fatalf("recording at an unclaimed position: %v", err)
	}
	if err := execution.BeginDown("", nil); !errors.Is(err, ErrInconsistentExecution) {
		t.Fatalf("cleanup without a run identity: %v", err)
	}
}

func mustIdentifier(t *testing.T, name string) Identifier {
	t.Helper()
	id, err := ParseIdentifier(name)
	if err != nil {
		t.Fatal(err)
	}
	return id
}
