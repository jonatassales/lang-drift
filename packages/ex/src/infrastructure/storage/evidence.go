package storage

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/jonatassales/tao/packages/ex/src/application"
	"github.com/jonatassales/tao/packages/ex/src/domain/microruntime"
	"github.com/jonatassales/tao/packages/ex/src/domain/project"
	"github.com/jonatassales/tao/packages/ex/src/infrastructure/files"
)

// executionEvidence is the retained record of one lifecycle operation.
type executionEvidence struct {
	Operation    application.ExecutionOutcome `json:"operation"`
	Billing      string                       `json:"billing"`
	Verification string                       `json:"verification"`
}

// runRecordVersion is the run.json format this version writes.
const runRecordVersion = 1

// executionRun is the run record of an execution nobody planned: the
// connected agent started it from console intent, so the runtime gets a
// record of kind execution and its evidence a home.
type executionRun struct {
	Version int    `json:"version"`
	ID      string `json:"id"`
	Kind    string `json:"kind"`
	Runtime string `json:"runtime"`
}

// Evidence keeps token-free records of every operation and inspection inside
// the runtime's folder, under evidence/<execution>.
type Evidence struct{}

// RecordExecution writes the evidence of one up or down.
func (Evidence) RecordExecution(owner project.Project, outcome application.ExecutionOutcome) error {
	record := executionEvidence{
		Operation:    outcome,
		Billing:      "unavailable",
		Verification: "observations are local; no production or savings claim",
	}

	return writeEvidence(owner, outcome.Runtime, outcome.RunID, "runtime.json", record)
}

// RecordInspection writes one inspection report beside the others of its execution.
func (Evidence) RecordInspection(owner project.Project, report application.Inspection) error {
	name := filepath.Join("inspections", report.ObservationID+".json")

	return writeEvidence(owner, report.Runtime, report.ExecutionRunID, name, report)
}

// writeEvidence resolves the execution's evidence directory inside its
// runtime, makes sure the runtime has a run record, and writes one document.
func writeEvidence(owner project.Project, runtimeText, executionRunID, name string, document any) error {
	runtime, err := microruntime.ParseIdentifier(runtimeText)
	if err != nil {
		return fmt.Errorf("evidence for runtime %q: %w", runtimeText, err)
	}

	if executionRunID == "" {
		return errors.New("evidence needs the execution it belongs to")
	}

	layout := LayoutOf(owner)

	if err := ensureRunRecord(layout, runtime); err != nil {
		return err
	}

	return files.WriteJSON(filepath.Join(layout.evidence(runtime, executionRunID), name), document)
}

// ensureRunRecord writes an execution-kind run.json when the runtime has none,
// so an execution started without a plan still has a run of its own.
func ensureRunRecord(layout Layout, runtime microruntime.Identifier) error {
	path := layout.runRecord(runtime)

	_, err := os.Lstat(path)
	if err == nil {
		return nil
	}

	if !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	record := executionRun{Version: runRecordVersion, ID: runtime.String(), Kind: "execution", Runtime: runtime.String()}

	return files.WriteJSON(path, record)
}
