package connected

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestConnectionOutbox_ReservesCleanupCapacity(t *testing.T) {
	event := Event{ID: strings.Repeat("a", 32), ExecutionID: strings.Repeat("b", 36), Kind: "failed", Cause: "start_failed", OccurredAt: time.Date(2026, 10, 21, 12, 0, 0, 0, time.UTC), Revision: 10, IntentVersion: 100, Detail: strings.Repeat("operation_failed ", 24), OperationID: strings.Repeat("c", 36), ConcreteHash: strings.Repeat("d", 64), MarkdownHash: strings.Repeat("e", 64), SourceRevision: strings.Repeat("f", 40)}
	record := Connection{RuntimeID: "runtime", RuntimeUID: "runtime-uid", Markdown: strings.Repeat("# Authored configuration\n", 2048), State: "running"}
	// Find the largest realistic backlog admitted before the cleanup reserve.
	low, high := 0, maxBufferedEvents-4
	for low < high {
		count := (low + high + 1) / 2
		record.Outbox = make([]Event, count)
		for i := range record.Outbox {
			record.Outbox[i] = event
			record.Outbox[i].ID = fmt.Sprintf("%032d", i+1)
			record.Outbox[i].Sequence = uint64(i + 1)
		}
		data, err := MarshalConnection(record)
		if err == nil && len(data) <= MaxConnectionStateBytes-cleanupReserveBytes {
			low = count
		} else {
			high = count - 1
		}
	}
	record.Outbox = record.Outbox[:low]
	record.Sequence = uint64(low)
	event.Sequence = record.Sequence + 1
	if err := addEvent(&record, event); err == nil {
		t.Fatal("ordinary event consumed cleanup reserve")
	}
	if record.Sequence != uint64(low) || len(record.Outbox) != low {
		t.Fatal("rejected event changed the durable ledger")
	}
	for _, kind := range []string{"stop_started", "stopped"} {
		event.Kind, event.Cause, event.Sequence = kind, "connection_lost", record.Sequence+1
		if err := addEvent(&record, event); err != nil {
			t.Fatalf("reserved cleanup event failed: %v", err)
		}
	}
	data, err := MarshalConnection(record)
	if err != nil {
		t.Fatal(err)
	}
	var reopened Connection
	if err := json.Unmarshal(data, &reopened); err != nil || reopened.Sequence != uint64(low+2) || reopened.Outbox[len(reopened.Outbox)-1].Kind != "stopped" {
		t.Fatalf("cleanup ledger cannot be reopened: %v", err)
	}
}
