// Package inputs holds the use cases that record the guardrail inputs: the
// hardware envelope, the workflow stages, the model and effort, the driver,
// the budget, the gates, the prices and the forecasts. Every value that can
// block or allow a run is explicit input, held in a private store with who
// set it, when and the value before, and mirrored for reading. No gate reads
// a mirror.
package inputs

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jonatassales/tao/packages/ex/src/domain/harness"
	"github.com/jonatassales/tao/packages/ex/src/domain/project"
)

// The kinds the store holds, one record each.
const (
	KindEnvelope  = "envelope"
	KindWorkflow  = "workflow"
	KindChoice    = "choice"
	KindDriver    = "driver"
	KindBudget    = "budget"
	KindGates     = "gates"
	KindPrices    = "prices"
	KindForecasts = "forecasts"
)

// The mirrors, relative to the .langdrift folder.
const (
	EnvelopeMirror  = "harness/envelope.json"
	WorkflowMirror  = "harness/workflow.json"
	BudgetMirror    = "finops/budgets.json"
	GatesMirror     = "finops/gates.json"
	ForecastsMirror = "finops/forecasts.json"
	PricesDirectory = "finops/prices"
)

// Via is what every record names as the way it was set.
const Via = "cli"

// MirrorNote opens every mirror the CLI writes.
const MirrorNote = "Written by ex from its own record. No gate reads this file; an edit here is shown as an unauthorized change."

// Record is one input as the store keeps it: what it is, who set it, when,
// its value as the mirror shows it, and the record it replaced.
type Record struct {
	Kind     string            `json:"kind"`
	ID       string            `json:"id"`
	SetBy    string            `json:"setBy"`
	Via      string            `json:"via"`
	SetAt    harness.Timestamp `json:"setAt"`
	Value    json.RawMessage   `json:"value"`
	Previous *Record           `json:"previous"`
}

// Store keeps the records privately. Load reports application.ErrNotFound
// for a kind nobody set; Save keeps the record it replaces as the previous one.
type Store interface {
	Load(owner project.Project, kind string) (Record, error)
	Save(owner project.Project, record Record) error
}

// Mirrors writes the files the consumer and a console read, and numbers the
// price lists.
type Mirrors interface {
	WriteMirror(owner project.Project, relative string, value any) error
	NextPriceNumber(owner project.Project) (string, error)
}

// ErrInvalid classifies input the use cases refuse; the message says why.
var ErrInvalid = errors.New("invalid input")

func invalid(format string, arguments ...any) error {
	return fmt.Errorf("%w: %s", ErrInvalid, fmt.Sprintf(format, arguments...))
}

// Result says which record was set and where its mirror is.
type Result struct {
	Kind   string `json:"kind"`
	Record string `json:"record"`
	Mirror string `json:"mirror,omitempty"`
	Status string `json:"status"`
}
