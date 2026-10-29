// Package execution holds the use cases that run a planned runtime: the
// money gate before a lap, the lap's attempts, and what the run records. The
// runtime is the execution: one runtime folder, one run, its laps.
package execution

import (
	"errors"

	"github.com/jonatassales/tao/packages/contracts/configuration"
	"github.com/jonatassales/tao/packages/ex/src/domain/finops"
	"github.com/jonatassales/tao/packages/ex/src/domain/harness"
	"github.com/jonatassales/tao/packages/ex/src/domain/microruntime"
	"github.com/jonatassales/tao/packages/ex/src/domain/project"
	"github.com/jonatassales/tao/packages/ex/src/domain/run"
)

// Records reads and writes the runtime's records.
type Records interface {
	ReadPlan(owner project.Project, runtime microruntime.Identifier) (configuration.Configuration, run.Plan, error)
	WritePlan(owner project.Project, runtime microruntime.Identifier, config configuration.Configuration, plan run.Plan) error
	ReadRun(owner project.Project, runtime microruntime.Identifier) (run.Run, error)
	WriteRun(owner project.Project, runtime microruntime.Identifier, record run.Run) error
	ReadLedger(owner project.Project, runtimeID string) (finops.Ledger, error)
	WriteLedger(owner project.Project, runtimeID string, ledger finops.Ledger) error
	WriteLap(owner project.Project, runtime microruntime.Identifier, lap Lap) error
}

// Laps keeps the private state of each lap, from which the lap file is
// rendered at every change.
type Laps interface {
	Load(owner project.Project, runtime microruntime.Identifier, number string) (Lap, error)
	Save(owner project.Project, runtime microruntime.Identifier, lap Lap) error
	Open(owner project.Project, runtime microruntime.Identifier) (Lap, bool, error)
}

// Lap is one attempt of the run, as the private state keeps it.
type Lap struct {
	Number        string                 `json:"number"`
	Run           string                 `json:"run"`
	Choice        run.Choice             `json:"choice"`
	ModelName     string                 `json:"modelName"`
	Estimate      finops.TokensThousands `json:"estimate"`
	EstimateCents int                    `json:"estimateCents"`
	Requested     int                    `json:"requestedCents"`
	Reserved      int                    `json:"reservedCents"`
	Gates         []finops.GateResult    `json:"gates"`
	Outcome       string                 `json:"outcome"`
	Gate          string                 `json:"gate,omitempty"`
	Execution     string                 `json:"execution,omitempty"`
	Steps         []Step                 `json:"steps"`
	InputTokens   int64                  `json:"inputTokens"`
	OutputTokens  int64                  `json:"outputTokens"`
	Unmeasured    bool                   `json:"unmeasured"`
	Runtime       *run.RuntimeFence      `json:"runtime,omitempty"`
	Observed      *run.Observed          `json:"observed,omitempty"`
	Started       harness.Timestamp      `json:"started"`
	Ended         *harness.Timestamp     `json:"ended"`
}

// Step is one row of the lap's steps table.
type Step struct {
	Name string            `json:"name"`
	At   harness.Timestamp `json:"at"`
	Note string            `json:"note"`
}

// The outcomes a lap may have.
const (
	PausedAtGate = "paused-at-gate"
	BlockedAt    = "blocked"
	Open         = "open"
	Completed    = "completed"
	Failed       = "failed"
	Unknown      = "unknown"
)

// NothingPaid reports whether the lap paid nothing: it never got past a gate.
func (lap Lap) NothingPaid() bool {
	return lap.Outcome == PausedAtGate || lap.Outcome == BlockedAt
}

// Ended reports whether the lap is over.
func (lap Lap) Over() bool { return lap.Ended != nil }

// ErrGate says a money gate did not pass: the lap is recorded, nothing paid.
var ErrGate = errors.New("a money gate did not pass")

// ErrRefused says the run cannot take the request as it stands.
var ErrRefused = errors.New("the run refuses")
