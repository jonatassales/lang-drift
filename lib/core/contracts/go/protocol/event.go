package protocol

import "time"

// Event is an immutable, durable observation. Its ID and sequence survive retries.
type Event struct {
	ID             string    `json:"id"`
	Sequence       uint64    `json:"sequence"`
	PurgeEpoch     int       `json:"purgeEpoch"`
	ExecutionID    string    `json:"executionId,omitempty"`
	Kind           string    `json:"kind"`
	Cause          string    `json:"cause"`
	OccurredAt     time.Time `json:"occurredAt"`
	Revision       int       `json:"revision"`
	IntentVersion  int       `json:"intentVersion,omitempty"`
	Detail         string    `json:"detail,omitempty"`
	OperationID    string    `json:"operationId,omitempty"`
	SourceRevision string    `json:"sourceRevision,omitempty"`
	ConcreteHash   string    `json:"concreteHash,omitempty"`
	MarkdownHash   string    `json:"markdownHash,omitempty"`
}
