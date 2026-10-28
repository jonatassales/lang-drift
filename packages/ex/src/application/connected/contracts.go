// Package connected owns the long-running project agent and its bounded workers.
// Every operation is driven by authenticated console intent, never repository prose.
package connected

import (
	"context"
	"errors"
	"time"

	"github.com/jonatassales/tao/exphere/cli/src/application"
	"github.com/jonatassales/tao/exphere/cli/src/domain/connection"
	"github.com/jonatassales/tao/exphere/cli/src/domain/project"
	"github.com/jonatassales/tao/exphere/contracts/protocol"
)

// The wire types the agent exchanges with the server belong to the Exphere
// contracts; these names keep the agent's code reading as before.
type (
	Event               = protocol.Event
	Observation         = protocol.Observation
	MeasurementWindow   = protocol.MeasurementWindow
	Repository          = protocol.Repository
	Operation           = protocol.Operation
	Intent              = protocol.Intent
	UnresolvedExecution = protocol.UnresolvedExecution
	Binding             = protocol.Binding
	Claim               = protocol.Claim
	Sync                = protocol.Sync
	ProjectMembership   = protocol.ProjectMembership
)

// Connection retains operational ownership separately from history. Credentials
// are loaded from the ignored runtime .env and excluded from JSON state.
type Connection struct {
	URL             string                `json:"url"`
	EnvironmentID   string                `json:"environmentId"`
	RuntimeUID      string                `json:"runtimeUid"`
	RuntimeID       string                `json:"runtimeId"`
	ProjectID       string                `json:"projectId"`
	Repository      Repository            `json:"repository"`
	Credential      string                `json:"-"`
	Secret          string                `json:"-"`
	State           connection.AgentState `json:"state"`
	ExecutionID     string                `json:"executionId,omitempty"`
	Revision        int                   `json:"revision"`
	IntentVersion   int                   `json:"intentVersion"`
	LastAttempt     int                   `json:"lastAttempt"`
	FailureRevision int                   `json:"failureRevision"`
	ResumeEligible  bool                  `json:"resumeEligible"`
	OperationID     string                `json:"operationId"`
	Sequence        uint64                `json:"sequence"`
	PurgeEpoch      int                   `json:"purgeEpoch"`
	Markdown        string                `json:"markdown"`
	ConcreteHash    string                `json:"concreteHash"`
	Managed         map[string]string     `json:"managed"`
	Outbox          []Event               `json:"outbox"`
	Observation     *Observation          `json:"observation,omitempty"`
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
	IDs              func() (string, error)
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
