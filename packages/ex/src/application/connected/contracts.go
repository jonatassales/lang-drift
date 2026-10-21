// Package connected owns the long-running project agent and its bounded workers.
// Every operation is driven by authenticated console intent, never repository prose.
package connected

import (
	"context"
	"errors"
	"time"

	"github.com/jonatassales/tao/packages/ex/src/application"
	"github.com/jonatassales/tao/packages/ex/src/domain/connection"
	"github.com/jonatassales/tao/packages/ex/src/domain/project"
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

type Repository struct {
	Kind     string `json:"kind,omitempty"`
	Path     string `json:"path,omitempty"`
	FullName string `json:"fullName"`
	URL      string `json:"url"`
	Ref      string `json:"ref"`
}
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

type Binding struct {
	EnvironmentID string     `json:"environmentId"`
	Credential    string     `json:"credential"`
	RuntimeUID    string     `json:"runtimeUid"`
	RuntimeID     string     `json:"runtimeId"`
	ProjectID     string     `json:"projectId"`
	Repository    Repository `json:"repository"`
	Intent
}
type Claim struct {
	ExecutionID   string `json:"executionId"`
	IntentVersion int    `json:"intentVersion"`
	Revision      int    `json:"revision"`
}
type Sync struct {
	EnvironmentID string       `json:"environmentId"`
	Events        []Event      `json:"events"`
	State         string       `json:"state"`
	ExecutionID   string       `json:"executionId,omitempty"`
	LastRevision  int          `json:"lastRevision"`
	PurgeEpoch    int          `json:"purgeEpoch"`
	Observation   *Observation `json:"observation,omitempty"`
}

// Connection retains operational ownership separately from history. Credentials
// are loaded from the ignored runtime .env and excluded from JSON state.
type Connection struct {
	URL             string            `json:"url"`
	EnvironmentID   string            `json:"environmentId"`
	RuntimeUID      string            `json:"runtimeUid"`
	RuntimeID       string            `json:"runtimeId"`
	ProjectID       string            `json:"projectId"`
	Repository      Repository        `json:"repository"`
	Credential      string            `json:"-"`
	Secret          string            `json:"-"`
	State           string            `json:"state"`
	ExecutionID     string            `json:"executionId,omitempty"`
	Revision        int               `json:"revision"`
	IntentVersion   int               `json:"intentVersion"`
	LastAttempt     int               `json:"lastAttempt"`
	FailureRevision int               `json:"failureRevision"`
	ResumeEligible  bool              `json:"resumeEligible"`
	OperationID     string            `json:"operationId"`
	Sequence        uint64            `json:"sequence"`
	PurgeEpoch      int               `json:"purgeEpoch"`
	Markdown        string            `json:"markdown"`
	ConcreteHash    string            `json:"concreteHash"`
	Managed         map[string]string `json:"managed"`
	Outbox          []Event           `json:"outbox"`
	Observation     *Observation      `json:"observation,omitempty"`
}

// Console scopes every call to one immutable endpoint and credential.
type Console interface {
	Bind(context.Context, string, string, string) (Binding, error)
	Sync(context.Context, Connection, Sync) (Intent, error)
	Claim(context.Context, Connection, Intent) (Claim, error)
}

var ErrAccessRevoked = errors.New("console access is revoked")

// ProjectConsole discovers runtime membership without granting start permission.
type ProjectConsole interface {
	ProjectSync(context.Context, string, string) (ProjectMembership, error)
	ProjectBind(context.Context, string, string, string, string) (Binding, error)
}
type ProjectMembership struct {
	Repository Repository `json:"repository"`
	ProjectID  string     `json:"projectId"`
	RuntimeIDs []string   `json:"runtimeIds"`
}

// Store owns private operational state and the per-environment process lock.
type Store interface {
	List(project.Project) ([]Connection, error)
	Save(project.Project, Connection) error
	Bind(context.Context, project.Project, Connection) error
	Lock(project.Project, string) (func(), error)
}

// Native retains the existing lifecycle ownership rules. Materialize never
// changes live configuration. Start and Stop return only after actual outcomes.
type Native interface {
	VerifyRepository(context.Context, project.Project, Repository) error
	Materialize(context.Context, project.Project, Connection, connection.Configuration, string) (map[string]string, string, error)
	Start(context.Context, project.Project, Connection, connection.Configuration, func(string, string) error) error
	Stop(context.Context, project.Project, Connection) error
	Observe(context.Context, project.Project, Connection) (*Observation, error)
	Migrate(context.Context, project.Project, string) (string, error)
}

// Service is explicitly composed; no use case creates hidden global clients.
type Service struct {
	Projects         application.ProjectOpener
	Store            Store
	Console          Console
	ProjectConsole   ProjectConsole
	Native           Native
	PollInterval     time.Duration
	OperationTimeout time.Duration
	Now              func() time.Time
}
type ConnectRequest struct {
	Project application.ProjectSelection
	URL     string
	Secret  string
}
type RunRequest struct{ Project application.ProjectSelection }
type ProjectRunRequest struct {
	Project application.ProjectSelection
	URL     string
	Secret  string
}
type MigrateRequest struct {
	Project application.ProjectSelection
	Runtime string
	Recover bool
}
type Result struct {
	Runtime     string `json:"runtime,omitempty"`
	Environment string `json:"environment,omitempty"`
	Path        string `json:"path,omitempty"`
	Status      string `json:"status"`
}
