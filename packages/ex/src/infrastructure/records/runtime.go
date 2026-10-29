package records

import (
	"encoding/json"
	"errors"
	"io/fs"
	"path/filepath"

	"github.com/jonatassales/tao/packages/contracts/configuration"
	"github.com/jonatassales/tao/packages/ex/src/application"
	"github.com/jonatassales/tao/packages/ex/src/application/plan"
	"github.com/jonatassales/tao/packages/ex/src/domain/collection"
	"github.com/jonatassales/tao/packages/ex/src/domain/finops"
	"github.com/jonatassales/tao/packages/ex/src/domain/microruntime"
	"github.com/jonatassales/tao/packages/ex/src/domain/project"
	"github.com/jonatassales/tao/packages/ex/src/domain/run"
	"github.com/jonatassales/tao/packages/ex/src/infrastructure/files"
)

// The names of a runtime's records inside its folder.
const (
	RuntimeFile = "runtime.md"
	RunFile     = "run.json"
	MemoryFile  = "run.md"
)

// The fixed prose of a runtime file.
const (
	runtimeSentence  = "Configuration is interpreted by the connected project agent. Prose does not execute commands."
	runtimeParagraph = "The planner read the ticket and the files it links, and the owner confirmed the components through the CLI. Nothing collected set a limit: the envelope and the budget are CLI input. The model and effort configure the suggestion."
)

// ConfigurationFence renders the langdrift fence's body as the contracts do.
func ConfigurationFence(config configuration.Configuration) ([]byte, error) {
	return json.MarshalIndent(config, "", "  ")
}

// RenderRuntime renders runtime.md: the prose, the langdrift fence the
// engine executes and the langdrift-plan fence that holds the analysis and
// the plan.
func RenderRuntime(document plan.Document) ([]byte, error) {
	config, err := ConfigurationFence(document.Config)
	if err != nil {
		return nil, err
	}

	plan, err := Format(document.Plan, FenceProfile)
	if err != nil {
		return nil, err
	}

	markdown := "# " + document.Name + "\n\n" +
		runtimeSentence + "\n\n" +
		"**Ticket:** " + document.TicketName + " · run `" + document.Run + "`\n" +
		"**Context:** sha256 " + document.ContextSHA256 + "\n" +
		"**Envelope and budget:** CLI input, record `" + document.Record + "`, set by " + document.By + " on " + document.Date + "\n\n" +
		runtimeParagraph + "\n\n" +
		"```langdrift\n" + string(config) + "\n```\n\n" +
		"```langdrift-plan\n" + string(plan) + "```\n"

	return []byte(markdown), nil
}

// Runtimes reads and writes a runtime's records in its folder.
type Runtimes struct{}

func runtimeFile(owner project.Project, runtime microruntime.Identifier, name string) (string, error) {
	return files.Within(owner.Root, filepath.Join(filepath.FromSlash(microruntime.ConfigurationDirectory(runtime).String()), name))
}

// ReadRuntime reads the runtime file, or reports application.ErrNotFound.
func (Runtimes) ReadRuntime(owner project.Project, runtime microruntime.Identifier) ([]byte, error) {
	path, err := runtimeFile(owner, runtime, RuntimeFile)
	if err != nil {
		return nil, err
	}

	content, err := files.ReadBounded(path, 1<<20)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, application.NotFound("runtime %s has no runtime.md", runtime)
	}

	return content, err
}

// WriteRuntime writes the runtime file.
func (Runtimes) WriteRuntime(owner project.Project, runtime microruntime.Identifier, markdown []byte) error {
	path, err := runtimeFile(owner, runtime, RuntimeFile)
	if err != nil {
		return err
	}

	return files.WriteDocument(path, markdown, documentMode)
}

// WriteRun writes the execution record.
func (Runtimes) WriteRun(owner project.Project, runtime microruntime.Identifier, record run.Run) error {
	path, err := runtimeFile(owner, runtime, RunFile)
	if err != nil {
		return err
	}

	data, err := Format(record, FileProfile)
	if err != nil {
		return err
	}

	return files.WriteDocument(path, data, documentMode)
}

// WriteLedger writes the runtime's ledger under finops/ledger.
func (Runtimes) WriteLedger(owner project.Project, runtimeID string, ledger finops.Ledger) error {
	path, err := files.Within(owner.Root, filepath.Join(".langdrift", "finops", "ledger", runtimeID+".ledger.json"))
	if err != nil {
		return err
	}

	data, err := Format(ledger, FileProfile)
	if err != nil {
		return err
	}

	return files.WriteDocument(path, data, documentMode)
}

// Render renders the runtime file for the plan use case.
func (Runtimes) Render(document plan.Document) ([]byte, error) { return RenderRuntime(document) }

// Fence renders the configuration fence and hashes it.
func (Runtimes) Fence(config configuration.Configuration) ([]byte, string, error) {
	body, err := ConfigurationFence(config)
	if err != nil {
		return nil, "", err
	}

	return body, collection.SHA256(body), nil
}
