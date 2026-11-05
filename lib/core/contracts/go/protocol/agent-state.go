package protocol

// AgentState is where one supervised runtime stands, as the agent reports it.
type AgentState string

// The agent states. Their values are persisted and synchronized, so they
// never change.
const (
	StateStopped  AgentState = "stopped"
	StateStarting AgentState = "starting"
	StateRunning  AgentState = "running"
	StateStopping AgentState = "stopping"
	StateUnknown  AgentState = "unknown"
)
