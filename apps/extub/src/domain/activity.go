package domain

import (
	"fmt"
	"strconv"
	"time"
)

type Event struct {
	ID            string `json:"id"`
	RuntimeID     string `json:"runtimeId"`
	EnvironmentID string `json:"environmentId,omitempty"`
	ExecutionID   string `json:"executionId,omitempty"`
	Kind          string `json:"kind"`
	Cause         string `json:"cause"`
	Actor         string `json:"actor"`
	OccurredAt    string `json:"occurredAt"`
	ReceivedAt    string `json:"receivedAt"`
}
type ExecutionRecord struct {
	Position        int64     `json:"position"`
	RuntimeName     string    `json:"runtimeName"`
	EnvironmentName string    `json:"environmentName"`
	Execution       Execution `json:"execution"`
}
type EventRecord struct {
	Position        int64  `json:"position"`
	RuntimeName     string `json:"runtimeName"`
	EnvironmentName string `json:"environmentName"`
	Event           Event  `json:"event"`
}

// RecordActivity stores only synthetic lifecycle evidence in bounded fixture memory.
func (catalog *Catalog) RecordActivity(id, action string, now time.Time) {
	index := catalog.RuntimeIndex(id)
	if index < 0 {
		return
	}
	runtime := catalog.Runtimes[index]
	timestamp := now.UTC().Format(time.RFC3339Nano)
	for i := range catalog.ExecutionRecords {
		record := &catalog.ExecutionRecords[i]
		if record.Execution.RuntimeID == id && record.Execution.State == "running" {
			record.Execution.State = "stopped"
			record.Execution.UpdatedAt = timestamp
		}
	}
	if action == "start" {
		for _, widget := range catalog.Widgets(now, 0) {
			if widget.Runtime.ID != id || len(widget.Executions) == 0 {
				continue
			}
			catalog.ActivitySequence++
			catalog.ExecutionRecords = append(catalog.ExecutionRecords, ExecutionRecord{catalog.ActivitySequence, runtime.Name, "Seed machine 1", widget.Executions[0]})
		}
	}
	kinds := map[string][2]string{"start": {"start_requested", "started"}, "stop": {"stop_requested", "stopped"}, "revoke": {"credential_revoked", "stopped"}, "delete": {"delete_requested", "deleted"}}
	cause := map[string]string{"start": "user_start", "stop": "user_stop", "revoke": "credential_revocation", "delete": "user_delete"}[action]
	for _, kind := range kinds[action] {
		catalog.ActivitySequence++
		event := Event{ID: fmt.Sprintf("seed-event-%d", catalog.ActivitySequence), RuntimeID: id, Kind: kind, Cause: cause, Actor: "seed-user", OccurredAt: timestamp, ReceivedAt: timestamp}
		catalog.EventRecords = append(catalog.EventRecords, EventRecord{catalog.ActivitySequence, runtime.Name, "Seed machine 1", event})
	}
	if len(catalog.ExecutionRecords) > 2048 {
		catalog.ExecutionRecords = catalog.ExecutionRecords[len(catalog.ExecutionRecords)-2048:]
	}
	if len(catalog.EventRecords) > 2048 {
		catalog.EventRecords = catalog.EventRecords[len(catalog.EventRecords)-2048:]
	}
}

// Activity uses monotonically increasing positions, so inserts do not shift cursors.
func (catalog Catalog) Activity(product, kind, cursor string, now time.Time) (map[string]any, error) {
	if _, err := catalog.Product(product); err != nil {
		return nil, err
	}
	before := int64(1<<53 - 1)
	if cursor != "" {
		value, err := strconv.ParseInt(cursor, 10, 64)
		if err != nil || value < 1 || value > before || strconv.FormatInt(value, 10) != cursor {
			return nil, ErrInvalid
		}
		before = value
	}
	items := []any{}
	var next *string
	positions := []int64{}
	belongs := func(runtime string) bool {
		index := catalog.RuntimeIndex(runtime)
		return index >= 0 && catalog.Runtimes[index].ProductID == product
	}
	if kind == "history" {
		for i := len(catalog.EventRecords) - 1; i >= 0; i-- {
			record := catalog.EventRecords[i]
			if record.Position >= before || !belongs(record.Event.RuntimeID) {
				continue
			}
			items = append(items, record)
			positions = append(positions, record.Position)
			if len(items) == 51 {
				break
			}
		}
	} else {
		live := map[string]Execution{}
		for _, widget := range catalog.Widgets(now, 0) {
			for _, execution := range widget.Executions {
				live[execution.ID] = execution
			}
		}
		for i := len(catalog.ExecutionRecords) - 1; i >= 0; i-- {
			record := catalog.ExecutionRecords[i]
			if record.Position >= before || !belongs(record.Execution.RuntimeID) {
				continue
			}
			if current, ok := live[record.Execution.ID]; ok {
				record.Execution = current
			}
			items = append(items, record)
			positions = append(positions, record.Position)
			if len(items) == 51 {
				break
			}
		}
	}
	if len(items) > 50 {
		value := strconv.FormatInt(positions[49], 10)
		next = &value
		items = items[:50]
	}
	return map[string]any{"source": "stub", "items": items, "nextCursor": next}, nil
}
