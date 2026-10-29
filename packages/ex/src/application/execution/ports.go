// Package execution holds the use cases that run a planned runtime: the
// money gate before a lap, the stages a lap plays through the owner's coding
// agent, the confirmations a person gives, and what the run records. The
// runtime is the execution: one runtime folder, one run, its laps.
package execution

import (
	"context"
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
	WriteStageEvidence(owner project.Project, runtime microruntime.Identifier, evidence StageEvidence) error
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
	Number          string                 `json:"number"`
	Run             string                 `json:"run"`
	Choice          run.Choice             `json:"choice"`
	ModelName       string                 `json:"modelName"`
	Estimate        finops.TokensThousands `json:"estimate"`
	EstimateCents   int                    `json:"estimateCents"`
	Requested       int                    `json:"requestedCents"`
	Reserved        int                    `json:"reservedCents"`
	Gates           []finops.GateResult    `json:"gates"`
	Outcome         string                 `json:"outcome"`
	Gate            string                 `json:"gate,omitempty"`
	Execution       string                 `json:"execution,omitempty"`
	Steps           []Step                 `json:"steps"`
	InputTokens     int64                  `json:"inputTokens"`
	OutputTokens    int64                  `json:"outputTokens"`
	ObservedCents   int                    `json:"observedCents"`
	Unmeasured      bool                   `json:"unmeasured"`
	StagesCompleted *harness.Timestamp     `json:"stagesCompleted"`
	Runtime         *run.RuntimeFence      `json:"runtime,omitempty"`
	Observed        *run.Observed          `json:"observed,omitempty"`
	Started         harness.Timestamp      `json:"started"`
	Ended           *harness.Timestamp     `json:"ended"`
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

// NothingPaid reports whether the lap paid nothing: a gate stopped it before
// its first paid call. FG-04 runs only after a call, so a lap it paused paid.
func (lap Lap) NothingPaid() bool {
	if lap.Outcome != PausedAtGate && lap.Outcome != BlockedAt {
		return false
	}

	return lap.Gate != duringRunGate
}

// Over reports whether the lap is over.
func (lap Lap) Over() bool { return lap.Ended != nil }

// duringRunGate is the gate that pauses the paid calls of an open lap.
const duringRunGate = "FG-04"

// ErrGate says a money gate did not pass: the lap is recorded as it stands.
var ErrGate = errors.New("a money gate did not pass")

// ErrRefused says the run cannot take the request as it stands.
var ErrRefused = errors.New("the run refuses")

// Confirmer answers whether the owner confirms a human stage now, as an
// interactive run asks; without one the run waits for ex agent confirm.
type Confirmer func(ctx context.Context, stage string) (bool, error)

// StageCall is one skill of one stage, run through the owner's coding agent.
type StageCall struct {
	ProjectRoot string
	Runtime     string
	Lap         string
	Stage       string
	Skill       string
	Ticket      string
	Execution   string
}

// StageResult is what the driver reported: the exit code, the usage when
// the result carried it, a one-line summary and where the log went.
type StageResult struct {
	ExitCode     int
	InputTokens  *int64
	OutputTokens *int64
	Summary      string
	LogPath      string
	ElapsedMs    int64
}

// HookCall is one committed script a stage runs first.
type HookCall struct {
	ProjectRoot string
	Runtime     string
	Lap         string
	Stage       string
	Script      string
}

// StageDriver runs a stage's skills through the owner's coding agent, and a
// stage's hooks as committed scripts, as trusted processes. Run reports a
// process that ran through its exit code; its error is a process that did
// not run or did not finish.
type StageDriver interface {
	Run(ctx context.Context, driver harness.Driver, call StageCall) (StageResult, error)
	Hook(ctx context.Context, call HookCall) (int, error)
}

// StageEvidence is the token-free record of one stage call: what ran, how
// it ended, what it used and where its log is; never what it said.
type StageEvidence struct {
	Lap          string            `json:"lap"`
	Stage        string            `json:"stage"`
	Skill        string            `json:"skill"`
	At           harness.Timestamp `json:"at"`
	ExitCode     int               `json:"exitCode"`
	InputTokens  *int64            `json:"inputTokens"`
	OutputTokens *int64            `json:"outputTokens"`
	Summary      string            `json:"summary"`
	Log          string            `json:"log"`
	ElapsedMs    int64             `json:"elapsedMs"`
}
