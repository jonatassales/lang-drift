package execution

import (
	"context"

	"github.com/jonatassales/tao/lib/cli/src/application/session"
	"github.com/jonatassales/tao/lib/cli/src/domain/run"
)

// abandonLap lets the waiting lap go. The runtime stops and is cleaned up
// first; a cleanup nobody confirmed leaves the outcome unknown and the
// reservation held, as at the end of a run. Otherwise finops.Abandon settles
// the lap: it closes as abandoned with who and when, its ledger line releases
// what Abandon releases, and the run finishes abandoned and remembers.
func (playing *playing) abandonLap(ctx context.Context) (StopResult, error) {
	if err := playing.takeDown(ctx); err != nil {
		return StopResult{}, err
	}

	now := playing.now()
	lap := &playing.lap
	lap.Ended = &now
	lap.Observed = playing.observed()

	outcome := Abandoned
	released := 0

	if lap.Outcome == Unknown {
		outcome = Unknown
		lap.Steps = append(lap.Steps, Step{
			Name: "Held",
			At:   now,
			Note: "stopped by " + playing.by + "; the cleanup was not confirmed, so the reservation is held until the outcome is known",
		})
	} else {
		released = playing.settleAbandoned(now)
	}

	playing.record.State, playing.record.Outcome = run.StateFinished, text(outcome)
	playing.record.Window.Finished = &now
	playing.syncLapSummary()

	if err := playing.service.Records.WriteMemory(playing.owner, playing.runtime, playing.memory(outcome)); err != nil {
		return StopResult{}, err
	}

	if err := playing.persist(); err != nil {
		return StopResult{}, err
	}

	if _, err := playing.service.Index.Rebuild(playing.owner, playing.service.embedder()); err != nil {
		return StopResult{}, err
	}

	emit(playing.sink, session.LapAbandoned, session.Data{
		"lap":           lap.Number,
		"outcome":       outcome,
		"by":            playing.by,
		"releasedCents": released,
		"held":          lap.Unmeasured,
	})
	emit(playing.sink, session.NextPlan, session.Data{"runtime": playing.runtime.String()})

	result := StopResult{
		Runtime:       playing.runtime.String(),
		Lap:           lap.Number,
		By:            playing.by,
		ReleasedCents: released,
		Status:        outcome + "; plan it again with dft agent plan",
	}

	return result, playing.stopCause
}
