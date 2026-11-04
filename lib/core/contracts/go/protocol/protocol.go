// Package protocol is the wire contract between the Microruntimes connected
// agent and the server it supervises runtimes for: an MCP endpoint at
// <console-url>/api/mcp, called with JSON-RPC 2.0 tools/call requests.
//
// Every operation is driven by authenticated console intent, never by
// repository prose. A field that was not measured is null, never zero.
package protocol

import "time"

// Version is the MCP protocol version the agent and the server agree on.
const Version = "2025-11-25"

// The tools the server exposes to the agent.
const (
	// ToolProjectSync lists a project's runtimes without granting a start.
	ToolProjectSync = "langdrift.project.sync"

	// ToolProjectBind issues one environment credential for one runtime.
	ToolProjectBind = "langdrift.project.bind"

	// ToolBind binds one runtime with a runtime secret.
	ToolBind = "langdrift.bind"

	// ToolSync exchanges events and observations for the current intent.
	ToolSync = "langdrift.sync"

	// ToolClaim claims a start for one intent version and revision.
	ToolClaim = "langdrift.claim"
)

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

// Observation contains measured usage only; nil values mean unavailable.
type Observation struct {
	// Nil until a consumer integration supplies measured cumulative execution cost.
	LLMCostUSD         *float64            `json:"llmCostUsd,omitempty"`
	State              string              `json:"-"`
	CPUUsed            *float64            `json:"cpuUsed"`
	MemoryUsedMiB      *float64            `json:"memoryUsedMiB"`
	ObservedAt         time.Time           `json:"observedAt"`
	MeasuredFrom       *time.Time          `json:"measuredFrom,omitempty"`
	MeasurementWindows []MeasurementWindow `json:"measurementWindows,omitempty"`
}

// MeasurementWindow preserves one native pod sample interval. Aggregates use
// the oldest sample timestamp, so a stale pod cannot make evidence look fresh.
type MeasurementWindow struct {
	Pod        string    `json:"pod"`
	ObservedAt time.Time `json:"observedAt"`
	Window     string    `json:"window"`
}

// Repository names the source a runtime builds from: a local folder or a
// GitHub repository at a revision.
type Repository struct {
	Kind     string `json:"kind,omitempty"`
	Path     string `json:"path,omitempty"`
	FullName string `json:"fullName"`
	URL      string `json:"url"`
	Ref      string `json:"ref"`
}

// Operation identifies who requested the current intent, and why.
type Operation struct {
	ID          string `json:"id"`
	Actor       string `json:"actor"`
	Cause       string `json:"cause"`
	RequestedAt string `json:"requestedAt"`
}

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

// UnresolvedExecution identifies work requiring reconciliation. It never
// authorizes a new start or replaces the console claim check.
type UnresolvedExecution struct {
	ExecutionID   string `json:"executionId"`
	IntentVersion int    `json:"intentVersion"`
	Revision      int    `json:"revision"`
	State         string `json:"state"`
}

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

// Claim is the server's grant of one start for one intent version and revision.
type Claim struct {
	ExecutionID   string `json:"executionId"`
	IntentVersion int    `json:"intentVersion"`
	Revision      int    `json:"revision"`
}

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

// ProjectMembership lists a project's runtimes. Discovery never grants a start.
type ProjectMembership struct {
	Repository Repository `json:"repository"`
	ProjectID  string     `json:"projectId"`
	RuntimeIDs []string   `json:"runtimeIds"`
}
