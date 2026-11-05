package protocol

// Intent is the console's current authorization, not evidence of local outcomes.
type Intent struct {
	AcceptedEventIDs    []string             `json:"acceptedEventIds"`
	DesiredState        string               `json:"desiredState"`
	IntentVersion       int                  `json:"intentVersion"`
	Revision            int                  `json:"revision"`
	Markdown            string               `json:"markdown"`
	Operation           Operation            `json:"operation"`
	LeaseSeconds        int                  `json:"leaseSeconds"`
	Revoked             bool                 `json:"revoked"`
	Deleted             bool                 `json:"deleted"`
	PurgeEpoch          int                  `json:"purgeEpoch"`
	UnresolvedExecution *UnresolvedExecution `json:"unresolvedExecution,omitempty"`
}
