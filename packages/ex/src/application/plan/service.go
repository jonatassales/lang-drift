package plan

import (
	"context"
	"fmt"
	"strings"

	"github.com/jonatassales/tao/packages/contracts/configuration"
	"github.com/jonatassales/tao/packages/ex/src/application"
	"github.com/jonatassales/tao/packages/ex/src/application/inputs"
	"github.com/jonatassales/tao/packages/ex/src/application/session"
	"github.com/jonatassales/tao/packages/ex/src/domain/collection"
	"github.com/jonatassales/tao/packages/ex/src/domain/harness"
	"github.com/jonatassales/tao/packages/ex/src/domain/rag"
	"github.com/jonatassales/tao/packages/ex/src/domain/run"
)

// Service plans tickets.
type Service struct {
	Projects  application.ProjectOpener
	Collector Collector
	Inspector Inspector
	Inputs    inputs.Store
	Records   Records
	Index     Index
	Embedder  rag.Embedder
	// Judge answers the planner's questions; nil plans without one.
	Judge Judge

	// Injectable for tests; nil uses the real sources.
	IDs   application.IDSource
	Clock application.Clock
}

// ComponentInput is one component the owner names on the command line.
type ComponentInput struct {
	Path          string
	Name          string
	Kind          string
	Image         string
	Port          int
	ReadinessPath string
	Build         configuration.Build
	Technologies  []run.Technology
}

// Request plans one ticket.
type Request struct {
	Project    application.ProjectSelection
	Ticket     string
	Runtime    string
	Name       string
	Components []ComponentInput
	CPUMillis  int
	MemoryMiB  int
	Limits     collection.Limits
	NoJudge    bool
	By         string
	Session    session.Sink
}

// Result says what was written.
type Result struct {
	Runtime  string `json:"runtime"`
	Revision int    `json:"revision"`
	Path     string `json:"path"`
	State    string `json:"state"`
}

// The states a plan may be in.
const (
	Proposed = "proposed"
	Blocked  = "blocked"
)

// Plan reads the ticket, follows its links, judges what they say, retrieves
// earlier runs, resolves and sizes the components, checks the draft and
// writes the runtime file, the run record and the ledger's planning line,
// printing the session as it goes.
func (service Service) Plan(ctx context.Context, request Request) (Result, error) {
	if err := validate(request); err != nil {
		return Result{}, err
	}

	owner, err := service.Projects.Open(request.Project)
	if err != nil {
		return Result{}, err
	}

	guards, err := LoadGuardrails(service.Inputs, owner)
	if err != nil {
		return Result{}, err
	}

	ticketPath, err := ticketPathOf(owner, request.Ticket)
	if err != nil {
		return Result{}, err
	}

	planning := &planning{
		service: service,
		request: request,
		owner:   owner,
		sink:    sinkOf(request.Session),
		guards:  guards,
		by:      orDefault(request.By, inputs.DefaultBy),
		now:     harness.At(service.Clock.Now()),
	}

	if err := planning.collect(ticketPath); err != nil {
		return Result{}, err
	}

	if err := planning.openJudge(); err != nil {
		return Result{}, err
	}

	if err := planning.analyze(ctx); err != nil {
		return Result{}, err
	}

	if err := planning.retrieve(ctx); err != nil {
		return Result{}, err
	}

	planning.suggest(ctx)

	if err := planning.resolve(); err != nil {
		return Result{}, err
	}

	if err := planning.check(); err != nil {
		return Result{}, err
	}

	if err := planning.write(); err != nil {
		return Result{}, err
	}

	if planning.state == Blocked {
		return Result{}, fmt.Errorf("%w: %s failed; the plan is written as blocked", ErrBlocked, strings.Join(failingChecks(planning.checks), ", "))
	}

	return Result{Runtime: planning.runtime.String(), Revision: planning.config.Revision, Path: runtimePath(planning.runtime), State: planning.state}, nil
}

func validate(request Request) error {
	if strings.TrimSpace(request.Ticket) == "" {
		return invalid("name the ticket to plan")
	}

	if len(request.Components) == 0 {
		return invalid("name at least one component with --component")
	}

	for _, component := range request.Components {
		if component.Path == "" {
			return invalid("every --component needs a path")
		}
	}

	if request.Limits.MaxDepth < 1 || request.Limits.MaxFiles < 1 || request.Limits.MaxBytes < 1 {
		return invalid("the collection limits must be positive")
	}

	if overridesEnvelopeHalfway(request) {
		return invalid("override the envelope with both --cpu-millis and --memory-mib, or neither")
	}

	return nil
}

func invalid(format string, arguments ...any) error {
	return fmt.Errorf("%w: %s", inputs.ErrInvalid, fmt.Sprintf(format, arguments...))
}

func sinkOf(sink session.Sink) session.Sink {
	if sink == nil {
		return session.Discard{}
	}

	return sink
}

// overridesEnvelopeHalfway reports an envelope override that names only one
// of its two figures, or a negative one.
func overridesEnvelopeHalfway(request Request) bool {
	if request.CPUMillis < 0 || request.MemoryMiB < 0 {
		return true
	}

	return (request.CPUMillis > 0) != (request.MemoryMiB > 0)
}
