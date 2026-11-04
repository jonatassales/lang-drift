package records

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"strings"

	"github.com/jonatassales/tao/lib/cli/src/application"
	"github.com/jonatassales/tao/lib/cli/src/application/plan"
	"github.com/jonatassales/tao/lib/cli/src/domain/connection"
	"github.com/jonatassales/tao/lib/cli/src/domain/finops"
	"github.com/jonatassales/tao/lib/cli/src/domain/microruntime"
	"github.com/jonatassales/tao/lib/cli/src/domain/project"
	"github.com/jonatassales/tao/lib/cli/src/domain/run"
	"github.com/jonatassales/tao/lib/cli/src/infrastructure/files"
	"github.com/jonatassales/tao/lib/core/contracts/go/configuration"
)

// ReadPlan reads the runtime file: the configuration fence the engine
// executes and the plan fence.
func (runtimes Runtimes) ReadPlan(owner project.Project, runtime microruntime.Identifier) (configuration.Configuration, run.Plan, error) {
	markdown, err := runtimes.ReadRuntime(owner, runtime)
	if err != nil {
		return configuration.Configuration{}, run.Plan{}, err
	}

	config, err := connection.Parse(string(markdown))
	if err != nil {
		return configuration.Configuration{}, run.Plan{}, fmt.Errorf("the runtime.md of %s does not parse: %w", runtime, err)
	}

	var planned run.Plan
	if err := DecodeFence(string(markdown), "langdrift-plan", &planned); err != nil {
		return configuration.Configuration{}, run.Plan{}, fmt.Errorf("the runtime.md of %s has no readable plan: %w", runtime, err)
	}

	return config, planned, nil
}

// WritePlan renders the runtime file again with a changed plan; the prose
// is rebuilt from the plan's own records.
func (runtimes Runtimes) WritePlan(owner project.Project, runtime microruntime.Identifier, config configuration.Configuration, planned run.Plan) error {
	document := plan.Document{
		Name:          config.Name,
		TicketName:    path.Base(planned.Ticket.Path),
		Run:           planned.Run,
		ContextSHA256: planned.Collection.Context.SHA256,
		Record:        planned.Inputs.CPUMillis.Record,
		By:            planned.Inputs.CPUMillis.SetBy,
		Date:          planned.Inputs.CPUMillis.SetAt.Format("2006-01-02"),
		Config:        config,
		Plan:          planned,
	}

	markdown, err := RenderRuntime(document)
	if err != nil {
		return err
	}

	return runtimes.WriteRuntime(owner, runtime, markdown)
}

// ReadRun reads the execution record, or reports application.ErrNotFound.
func (Runtimes) ReadRun(owner project.Project, runtime microruntime.Identifier) (run.Run, error) {
	path, err := runtimeFile(owner, runtime, RunFile)
	if err != nil {
		return run.Run{}, err
	}

	var record run.Run

	err = files.ReadJSON(path, &record)
	if errors.Is(err, fs.ErrNotExist) {
		return run.Run{}, application.NotFound("runtime %s has no run record; plan it first", runtime)
	}

	return record, err
}

// ReadLedger reads the runtime's ledger, or reports application.ErrNotFound.
func (Runtimes) ReadLedger(owner project.Project, runtimeID string) (finops.Ledger, error) {
	path, err := files.Within(owner.Root, filepath.Join(".langdrift", "finops", "ledger", runtimeID+".ledger.json"))
	if err != nil {
		return finops.Ledger{}, err
	}

	var ledger finops.Ledger

	err = files.ReadJSON(path, &ledger)
	if errors.Is(err, fs.ErrNotExist) {
		return finops.Ledger{}, application.NotFound("runtime %s has no ledger; plan it first", runtimeID)
	}

	return ledger, err
}

// DecodeFence reads one JSON fence of a Markdown document, strictly.
func DecodeFence(markdown, language string, value any) error {
	_, after, found := strings.Cut(markdown, "```"+language+"\n")
	if !found {
		return errors.New("no " + language + " fence")
	}

	body, _, closed := strings.Cut(after, "\n```")
	if !closed {
		return errors.New("the " + language + " fence is not closed")
	}

	decoder := json.NewDecoder(strings.NewReader(body))
	decoder.DisallowUnknownFields()

	return decoder.Decode(value)
}
