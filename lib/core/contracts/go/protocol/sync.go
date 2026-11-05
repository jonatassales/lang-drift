package protocol

// Sync is what the agent reports on every poll.
type Sync struct {
	EnvironmentID string       `json:"environmentId"`
	Events        []Event      `json:"events"`
	State         AgentState   `json:"state"`
	ExecutionID   string       `json:"executionId,omitempty"`
	LastRevision  int          `json:"lastRevision"`
	PurgeEpoch    int          `json:"purgeEpoch"`
	Observation   *Observation `json:"observation,omitempty"`
}
