// Package plan holds the use case that plans a ticket: it reads the ticket
// and the files its links reach, takes the guardrail inputs, and writes the
// runtime file, the run record and the ledger's planning line, printing the
// session as it goes. The plan file is the runtime file.
package plan

import (
	"errors"

	"github.com/jonatassales/tao/packages/contracts/configuration"
	"github.com/jonatassales/tao/packages/ex/src/domain/collection"
	"github.com/jonatassales/tao/packages/ex/src/domain/finops"
	"github.com/jonatassales/tao/packages/ex/src/domain/harness"
	"github.com/jonatassales/tao/packages/ex/src/domain/manifest"
	"github.com/jonatassales/tao/packages/ex/src/domain/microruntime"
	"github.com/jonatassales/tao/packages/ex/src/domain/project"
	"github.com/jonatassales/tao/packages/ex/src/domain/rag"
	"github.com/jonatassales/tao/packages/ex/src/domain/run"
)

// Collector follows a ticket's links under the limits.
type Collector interface {
	Collect(owner project.Project, ticketPath string, limits collection.Limits) (collection.Result, error)
}

// Document is what the runtime file holds beyond the fences: its prose, and
// the two fences' contents.
type Document struct {
	Name          string
	TicketName    string
	Run           string
	ContextSHA256 string
	Record        string
	By            string
	Date          string
	Config        configuration.Configuration
	Plan          run.Plan
}

// Records reads and writes a runtime's records.
type Records interface {
	ReadRuntime(owner project.Project, runtime microruntime.Identifier) ([]byte, error)
	WriteRuntime(owner project.Project, runtime microruntime.Identifier, markdown []byte) error
	WriteRun(owner project.Project, runtime microruntime.Identifier, record run.Run) error
	WriteLedger(owner project.Project, runtimeID string, ledger finops.Ledger) error
	ReadRun(owner project.Project, runtime microruntime.Identifier) (run.Run, error)
	ReadLedger(owner project.Project, runtimeID string) (finops.Ledger, error)
	Render(document Document) ([]byte, error)
	Fence(config configuration.Configuration) (body []byte, sha256 string, err error)
}

// ErrBlocked says a check failed and the plan was not written.
var ErrBlocked = errors.New("the plan is blocked")

// ErrRunning says the runtime has a lap open, so it cannot be planned again
// until that lap ended.
var ErrRunning = errors.New("the runtime is running")

// Spend keeps the planner's spend privately beside the laps, and sums what a
// period spent from the private state alone: no gate reads the ledger mirror.
type Spend interface {
	SavePlanning(owner project.Project, runtime microruntime.Identifier, at harness.Timestamp, line finops.PlanningLine) error
	Period(owner project.Project, period string) (int, error)
	// Lock takes the runtime's run lock, which a run and a confirmation take too.
	Lock(owner project.Project, runtime microruntime.Identifier) (func(), error)
}

// Inspector reads a component's manifests and the files it names: its
// Dockerfile and its env file. An absent file is fs.ErrNotExist.
type Inspector interface {
	Inspect(owner project.Project, componentPath string) (manifest.Inspection, error)
	File(owner project.Project, relativePath string) ([]byte, error)
}

// Index holds the project's episodes, rebuilt from the runtimes' memories.
type Index interface {
	Episodes(owner project.Project, embedder rag.Embedder) ([]rag.Episode, error)
	Rebuild(owner project.Project, embedder rag.Embedder) ([]rag.Episode, error)
}
