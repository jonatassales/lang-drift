package storage

import (
	"path/filepath"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/domain/project"
	"github.com/exphere/ex/src/infrastructure/files"
)

// executionEvidence is the retained record of one lifecycle operation.
type executionEvidence struct {
	Operation    application.ExecutionOutcome `json:"operation"`
	Billing      string                       `json:"billing"`
	Verification string                       `json:"verification"`
}

// Evidence keeps token-free records of every operation and inspection under
// .audits/evidence, one directory per run.
type Evidence struct{}

// RecordExecution writes the evidence of one up or down.
func (Evidence) RecordExecution(owner project.Project, outcome application.ExecutionOutcome) error {
	record := executionEvidence{
		Operation:    outcome,
		Billing:      "unavailable",
		Verification: "observations are local; no production or savings claim",
	}

	return files.WriteJSON(filepath.Join(LayoutOf(owner).evidence(outcome.RunID), "runtime.json"), record)
}

// RecordInspection writes one inspection report.
func (Evidence) RecordInspection(owner project.Project, report application.Inspection) error {
	return files.WriteJSON(filepath.Join(LayoutOf(owner).evidence(report.ObservationID), "inspection.json"), report)
}
