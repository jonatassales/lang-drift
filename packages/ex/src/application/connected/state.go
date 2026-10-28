package connected

import (
	"encoding/json"
	"errors"
)

// MaxConnectionStateBytes is the shared serialized write/read bound. Reserved
// space lets a saturated event outbox still record bounded cleanup and shutdown.
const MaxConnectionStateBytes = 4 << 20

const (
	cleanupReserveBytes  = 64 << 10
	shutdownReserveBytes = 32 << 10
	maxEventBytes        = 8 << 10
	maxBufferedEvents    = 10000
)

// MarshalConnection produces the exact private representation accepted by Store.
// Credentials remain excluded by Connection's JSON contract.
func MarshalConnection(record Connection) ([]byte, error) {
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return nil, err
	}
	if len(data)+1 > MaxConnectionStateBytes {
		return nil, errors.New("connection state exceeds its persistent byte limit")
	}
	return append(data, '\n'), nil
}

// addEvent validates before mutating the ledger. Stop outcomes have reserved
// capacity; final agent shutdown has a second reserve if repeated cleanup fails.
func addEvent(record *Connection, event Event) error {
	encoded, err := json.Marshal(event)
	if err != nil {
		return err
	}
	if len(encoded) > maxEventBytes {
		return errors.New("event exceeds its persistent byte limit")
	}
	limit, count := MaxConnectionStateBytes-cleanupReserveBytes, maxBufferedEvents-4
	if event.Kind == "stop_started" || event.Kind == "stopped" {
		limit, count = MaxConnectionStateBytes-shutdownReserveBytes, maxBufferedEvents-2
	}
	if event.Cause == "agent_shutdown" {
		limit, count = MaxConnectionStateBytes, maxBufferedEvents
	}
	if len(record.Outbox) >= count {
		return errors.New("event outbox is full; synchronize before more operations")
	}
	candidate := *record
	candidate.Sequence = event.Sequence
	candidate.Outbox = append(append([]Event(nil), record.Outbox...), event)
	data, err := MarshalConnection(candidate)
	if err != nil {
		return err
	}
	if len(data) > limit {
		return errors.New("event outbox byte limit reached; reserved space is available only for cleanup")
	}
	*record = candidate
	return nil
}
