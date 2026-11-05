package protocol

// Binding is the server's answer to a bind: the environment credential and
// the runtime's current intent.
type Binding struct {
	EnvironmentID string     `json:"environmentId"`
	Credential    string     `json:"credential"`
	RuntimeUID    string     `json:"runtimeUid"`
	RuntimeID     string     `json:"runtimeId"`
	ProjectID     string     `json:"projectId"`
	Repository    Repository `json:"repository"`
	Intent
}
