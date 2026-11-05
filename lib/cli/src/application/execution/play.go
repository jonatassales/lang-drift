package execution

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/jonatassales/tao/lib/cli/src/application/plan"
	"github.com/jonatassales/tao/lib/cli/src/application/session"
	"github.com/jonatassales/tao/lib/cli/src/domain/finops"
	"github.com/jonatassales/tao/lib/cli/src/domain/harness"
	"github.com/jonatassales/tao/lib/cli/src/domain/microruntime"
	"github.com/jonatassales/tao/lib/cli/src/domain/project"
	"github.com/jonatassales/tao/lib/cli/src/domain/run"
	"github.com/jonatassales/tao/lib/core/contracts/go/configuration"
)

// allowUnmeasured is the driver setting that lets a call without usage pass FG-04.
const allowUnmeasured = "allow"

// playing is one play of the stages: the run's records and the lap as they
// stand, written at every change.
type playing struct {
	service    Service
	owner      project.Project
	runtime    microruntime.Identifier
	guards     plan.Guardrails
	driver     *harness.Driver
	config     configuration.Configuration
	planned    run.Plan
	record     run.Run
	ledger     finops.Ledger
	lap        Lap
	sink       session.Sink
	by         string
	confirm    Confirmer
	bounds     finops.Bounds
	checksWait time.Duration
	pushed     bool
	now        func() harness.Timestamp
}

// play drives the lap through its phases in the order NextPhase decides:
// the runtime comes up before the first stage, a stage that waits for a
// person pauses with the runtime running, the stages play in order from the
// first not done, and the lap finishes, or fails, from there. The owner's
// stage, and everything after it, is never played.
func (playing *playing) play(ctx context.Context) (StartResult, error) {
	playing.seedStages()

	var outcome RunState
	var failure error
	for {
		state := playing.state(outcome)
		stage, _ := playing.nextStage()

		switch NextPhase(state) {
		case PhaseReserve:
			return playing.result(), fmt.Errorf("%w: lap %s of %s holds no reservation", ErrRefused, playing.lap.Number, playing.runtime)
		case PhaseBringUp:
			if err := playing.bringUp(ctx); err != nil {
				outcome.BringUpFailed, failure = true, err
			}
		case PhaseWait:
			confirmed, err := playing.wait(ctx, stage)
			if err != nil {
				return playing.result(), err
			}
			outcome.Waiting = !confirmed
		case PhasePlay:
			err := playing.perform(ctx, stage)
			if errors.Is(err, ErrGate) {
				return playing.pausedAtGate(ctx, err)
			}
			if err != nil {
				outcome.StageFailed, failure = true, playing.markFailed(stage, err)
			}
		case PhaseFinish:
			return playing.finish(ctx)
		case PhaseFail:
			return playing.failLap(ctx, failure)
		case PhaseDone:
			return playing.result(), nil
		}
	}
}

// state is where the lap stands for NextPhase: what its records say, and
// what the phase the loop just executed returned.
func (playing *playing) state(outcome RunState) RunState {
	stage, left := playing.nextStage()
	waits := left && plan.IsHuman(stage, playing.guards.Workflow.Defaults) && playing.stageState(stage.ID).ConfirmedBy == nil

	reserved := playing.lap.Outcome == Open && playing.lap.Reserved > 0

	return RunState{
		Reserved:       reserved,
		RuntimeUp:      playing.lap.Ready != nil,
		BringUpFailed:  outcome.BringUpFailed,
		StagesLeft:     left,
		NextStageWaits: waits,
		Waiting:        outcome.Waiting,
		StageFailed:    outcome.StageFailed,
	}
}

// nextStage is the first stage the run plays that is not done, before the
// owner's stage.
func (playing *playing) nextStage() (harness.Stage, bool) {
	for _, stage := range playing.guards.Workflow.Stages {
		if stage.Human == harness.OwnerOnly {
			break
		}
		if !plays(stage) || playing.stageState(stage.ID).State == run.StageDone {
			continue
		}

		return stage, true
	}

	return harness.Stage{}, false
}

// plays reports whether the run has anything to do in a stage.
func plays(stage harness.Stage) bool {
	return len(stage.Skills) > 0 || len(stage.Hooks) > 0
}

// wait prints that the stage waits for a person, and asks the confirmer when
// the run has one; without a yes the run stays waiting.
func (playing *playing) wait(ctx context.Context, stage harness.Stage) (bool, error) {
	emit(playing.sink, session.Waiting, session.Data{"stage": stage.ID})
	playing.setStage(stage.ID, func(state *run.StageState) { state.State = run.StageWaiting })
	playing.record.State, playing.record.Outcome = run.StateWaiting, text(run.StateRunning)

	if playing.confirm == nil {
		return false, playing.persist()
	}

	confirmed, err := playing.confirm(ctx, stage.ID)
	if err != nil {
		return false, err
	}
	if !confirmed {
		return false, playing.persist()
	}

	confirmStage(&playing.record, &playing.lap, stage.ID, playing.by, playing.now())
	emit(playing.sink, session.Confirmed, session.Data{"stage": stage.ID, "by": playing.by})

	return true, nil
}

// confirmStage records the owner's confirmation on the run and the lap.
func confirmStage(record *run.Run, lap *Lap, stageID, by string, at harness.Timestamp) {
	for index := range record.Stages {
		if record.Stages[index].ID != stageID {
			continue
		}
		confirmer := by
		record.Stages[index].ConfirmedBy = &confirmer
		record.Stages[index].State = run.StageActive
	}

	lap.Steps = append(lap.Steps, Step{Name: "Confirmed", At: at, Note: stageID + " · by " + by})
}

// perform runs a stage: its hooks first, then its skills, each paid call
// under FG-04.
func (playing *playing) perform(ctx context.Context, stage harness.Stage) error {
	playing.setStage(stage.ID, func(state *run.StageState) { state.State = run.StageActive })
	playing.record.State, playing.record.Outcome = run.StateRunning, text(run.StateRunning)

	for _, script := range stage.Hooks {
		if err := playing.hook(ctx, stage, script); err != nil {
			return err
		}
	}

	summary := ""
	if state := playing.stageState(stage.ID); state.Summary != nil {
		summary = *state.Summary
	}
	for _, skill := range stage.Skills {
		if slices.Contains(playing.stageState(stage.ID).SkillsDone, skill) {
			continue
		}

		result, err := playing.skill(ctx, stage, skill)
		if err != nil {
			return err
		}
		summary = result.Summary
		playing.setStage(stage.ID, func(state *run.StageState) {
			state.SkillsDone = append(state.SkillsDone, skill)
			state.Summary = text(summary)
		})
	}

	at := playing.now()
	playing.setStage(stage.ID, func(state *run.StageState) {
		state.State = run.StageDone
		state.At = &at
		if summary != "" {
			state.Summary = text(summary)
		}
	})

	return playing.persist()
}

// hook runs one committed script of a stage; a script that fails stops the
// stage before any paid call.
func (playing *playing) hook(ctx context.Context, stage harness.Stage, script string) error {
	call := HookCall{ProjectRoot: playing.owner.Root, Runtime: playing.runtime.String(), Lap: playing.lap.Number, Stage: stage.ID, Script: script}

	exitCode, err := playing.service.Driver.Hook(ctx, call)
	if err != nil {
		return playing.fail(stage, fmt.Errorf("hook %s of stage %s: %w", script, stage.ID, err))
	}

	summary := script + " · exit " + fmt.Sprint(exitCode)
	emit(playing.sink, session.Hook, session.Data{"summary": summary})
	playing.lap.Steps = append(playing.lap.Steps, Step{Name: "Hook", At: playing.now(), Note: stage.ID + " · " + summary})

	if exitCode != 0 {
		return playing.fail(stage, fmt.Errorf("hook %s of stage %s exited %d", script, stage.ID, exitCode))
	}

	return nil
}

// skill runs one skill through the driver, after FG-04, and records what it
// used and what it said it did.
func (playing *playing) skill(ctx context.Context, stage harness.Stage, skill string) (StageResult, error) {
	if err := playing.gate(stage); err != nil {
		return StageResult{}, err
	}

	call := StageCall{ProjectRoot: playing.owner.Root, Runtime: playing.runtime.String(), Lap: playing.lap.Number, Stage: stage.ID, Skill: skill, Ticket: playing.planned.Ticket.Path, Execution: playing.lap.Execution}

	result, err := playing.service.Driver.Run(ctx, *playing.driver, call)
	if err != nil && !errors.Is(err, ErrNotStarted) {
		// The agent ran and was cut off: what it used is unknown, never zero.
		playing.lap.Unmeasured = true
	}
	if err != nil {
		return StageResult{}, playing.fail(stage, fmt.Errorf("skill %s of stage %s: %w", skill, stage.ID, err))
	}

	at := playing.now()
	playing.account(stage.ID, result)

	summary := orDefault(result.Summary, "exit "+fmt.Sprint(result.ExitCode)+" · no summary")
	if result.InputTokens == nil || result.OutputTokens == nil {
		summary += " · no usage reported"
	}
	emit(playing.sink, session.Stage, session.Data{"skill": skill, "summary": summary})
	playing.lap.Steps = append(playing.lap.Steps, Step{Name: "Stage", At: at, Note: stage.ID + " · " + skill + " · " + summary})

	evidence := StageEvidence{Lap: playing.lap.Number, Stage: stage.ID, Skill: skill, At: at, ExitCode: result.ExitCode, InputTokens: result.InputTokens, OutputTokens: result.OutputTokens, Summary: result.Summary, Log: result.LogPath, ElapsedMs: result.ElapsedMs}
	if err := playing.service.Records.WriteStageEvidence(playing.owner, playing.runtime, evidence); err != nil {
		return StageResult{}, err
	}

	if result.ExitCode != 0 {
		return StageResult{}, playing.fail(stage, fmt.Errorf("skill %s of stage %s exited %d; see %s", skill, stage.ID, result.ExitCode, result.LogPath))
	}

	return result, nil
}

// gate is FG-04 before one paid call: the paid calls pause when what the lap
// observed reached its reservation, or when a call reported no usage and the
// driver does not allow that. A pause ends the lap; the next one plays on.
func (playing *playing) gate(stage harness.Stage) error {
	known := !playing.lap.Unmeasured || playing.driver.Unmeasured == allowUnmeasured
	result := finops.DuringRun(playing.lap.ObservedCents, known, playing.lap.Reserved)
	playing.recordGate(result)
	if result.Result == finops.NotReached {
		return nil
	}

	detail := result.ID + " " + result.Result
	if result.Detail != nil {
		detail = *result.Detail
	}
	emit(playing.sink, session.GateExhausted, session.Data{"detail": detail})

	now := playing.now()
	playing.lap.Steps = append(playing.lap.Steps, Step{Name: "Paused", At: now, Note: "FG-04: " + detail})
	playing.lap.Gate, playing.lap.Outcome, playing.lap.Ended = duringRunGate, PausedAtGate, &now
	playing.setStage(stage.ID, func(state *run.StageState) { state.State = run.StagePending })
	playing.record.State, playing.record.Outcome = run.StatePaused, text(run.StatePaused)

	if err := playing.persist(); err != nil {
		return err
	}

	return fmt.Errorf("%w: FG-04 paused the paid calls · %s · lap %s ended with %s observed", ErrGate, detail, playing.lap.Number, finops.Dollars(playing.lap.ObservedCents))
}

// recordGate keeps one result per gate on the lap: the latest.
func (playing *playing) recordGate(result finops.GateResult) {
	for index := range playing.lap.Gates {
		if playing.lap.Gates[index].ID == result.ID {
			playing.lap.Gates[index] = result

			return
		}
	}

	playing.lap.Gates = append(playing.lap.Gates, result)
}

// account adds a call's usage to the lap and to its stage, at the lap's
// list price; a call without usage marks the lap unmeasured.
func (playing *playing) account(stageID string, result StageResult) {
	if result.InputTokens == nil || result.OutputTokens == nil {
		playing.lap.Unmeasured = true

		return
	}

	price, _ := playing.guards.Price(playing.lap.Choice.Model)
	playing.lap.InputTokens += *result.InputTokens
	playing.lap.OutputTokens += *result.OutputTokens
	playing.lap.ObservedCents = finops.ObservedCents(playing.lap.InputTokens, playing.lap.OutputTokens, price)

	playing.setStage(stageID, func(state *run.StageState) {
		usage := run.StageUsage{}
		if state.Usage != nil {
			usage = *state.Usage
		}
		usage.InputTokens += *result.InputTokens
		usage.OutputTokens += *result.OutputTokens
		usage.Cents = finops.ObservedCents(usage.InputTokens, usage.OutputTokens, price)
		state.Usage = &usage
	})
}

// markFailed records the stage a failure stopped, once: the loop then fails
// the lap.
func (playing *playing) markFailed(stage harness.Stage, cause error) error {
	if playing.stageState(stage.ID).State == run.StageFailed {
		return cause
	}

	return playing.fail(stage, cause)
}

// pausedAtGate takes the runtime down after FG-04 ended the lap: the lap is
// paused, not failed, and the next one brings the runtime up again.
func (playing *playing) pausedAtGate(ctx context.Context, cause error) (StartResult, error) {
	if playing.lap.Execution == "" {
		return playing.result(), cause
	}

	if err := playing.takeDown(ctx); err != nil {
		return playing.result(), errors.Join(cause, err)
	}

	return playing.result(), cause
}

// fail records a stage the run could not finish; the next run plays it again.
func (playing *playing) fail(stage harness.Stage, cause error) error {
	playing.setStage(stage.ID, func(state *run.StageState) { state.State = run.StageFailed })
	playing.record.State, playing.record.Outcome = run.StateFailed, text(run.StateFailed)

	if err := playing.persist(); err != nil {
		return errors.Join(cause, err)
	}

	return cause
}

// result says where the lap stands.
func (playing *playing) result() StartResult {
	return StartResult{Runtime: playing.runtime.String(), Lap: playing.lap.Number, Outcome: playing.lap.Outcome, State: playing.record.State, ReservedCents: playing.lap.Reserved, ObservedCents: playing.lap.ObservedCents}
}

// seedStages gives the run its stage states when the workflow was set after
// the plan was written.
func (playing *playing) seedStages() {
	if len(playing.record.Stages) > 0 {
		return
	}

	workflow := playing.guards.Workflow
	playing.record.HumanInTheLoop = run.HumanInTheLoop{Record: text(workflow.Mirror.Record), Mirror: text(".langdrift/" + playing.guards.WorkflowAt), Chosen: append([]string{}, workflow.Defaults.HumanAt...)}
	for _, stage := range workflow.Stages {
		playing.record.Stages = append(playing.record.Stages, run.StageState{ID: stage.ID, State: run.StagePending, Human: plan.IsHuman(stage, workflow.Defaults)})
	}
}

func (playing *playing) stageState(id string) run.StageState {
	state, _ := stageOf(playing.record, id)

	return state
}

func (playing *playing) setStage(id string, change func(*run.StageState)) {
	for index := range playing.record.Stages {
		if playing.record.Stages[index].ID == id {
			change(&playing.record.Stages[index])

			return
		}
	}
}

// stageOf finds a stage's state on the run record.
func stageOf(record run.Run, id string) (run.StageState, bool) {
	for _, state := range record.Stages {
		if state.ID == id {
			return state, true
		}
	}

	return run.StageState{ID: id, State: run.StagePending}, false
}

// persist keeps the lap's ledger line and the run's lap summary current,
// then writes the lap, the run record and the ledger.
func (playing *playing) persist() error {
	playing.settleLine()
	finops.SettleTotals(&playing.ledger, playing.guards.Budget.Bounded.PeriodCents)

	playing.syncLapSummary()

	return playing.service.persist(playing.owner, playing.runtime, playing.lap, playing.record, playing.ledger)
}

// settleLine writes on the lap's ledger line what its calls used so far and,
// once the lap ended, what the reservation did not spend. A lap in which a
// call reported no usage settles nothing: its tokens, its variance and what
// it recovered stay unknown. A lap whose outcome is unknown holds its
// reservation.
func (playing *playing) settleLine() {
	lap := playing.lap
	if !hasLine(playing.ledger, lap.Number) {
		// An earlier write of the lap did not reach the ledger.
		playing.ledger.Laps = append(playing.ledger.Laps, lineOf(lap))
	}
	for index := range playing.ledger.Laps {
		line := &playing.ledger.Laps[index]
		if line.Lap != lap.Number {
			continue
		}

		line.Cents.Observed = lap.ObservedCents
		line.NothingPaid = lap.NothingPaid()
		line.Gates = append([]finops.GateResult{}, lap.Gates...)
		line.TokensThousands.Observed = nil
		if !lap.Unmeasured && (lap.InputTokens > 0 || lap.OutputTokens > 0) {
			line.TokensThousands.Observed = &finops.TokensThousands{Input: thousands(lap.InputTokens), Output: thousands(lap.OutputTokens)}
		}

		if lap.Ended == nil || lap.Unmeasured {
			return
		}

		settlement := finops.Settle(lap.EstimateCents, lap.Reserved, lap.ObservedCents)
		line.Cents.Recovered, line.Cents.Overspend = settlement.Recovered, settlement.Overspend
		line.Cents.Variance = number(settlement.Variance)
		if lap.Outcome == Unknown {
			// The reservation is held until the outcome is known.
			line.Cents.Recovered = 0
		}

		return
	}
}

// thousands rounds raw tokens up to the thousands the ledger counts in.
func thousands(tokens int64) int {
	return int((tokens + 999) / 1000)
}

// syncLapSummary keeps the run record's entry for the lap current, adding
// it when an earlier write of the lap did not reach the record.
func (playing *playing) syncLapSummary() {
	for index := range playing.record.Laps {
		if playing.record.Laps[index].Number == playing.lap.Number {
			playing.record.Laps[index] = summaryOf(playing.lap)

			return
		}
	}

	playing.record.Laps = append(playing.record.Laps, summaryOf(playing.lap))
}

// hasLine reports whether the ledger has a line for the lap.
func hasLine(ledger finops.Ledger, number string) bool {
	for _, line := range ledger.Laps {
		if line.Lap == number {
			return true
		}
	}

	return false
}
