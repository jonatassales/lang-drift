package execution

import "testing"

// runStateFields is how many booleans RunState holds; every combination of
// them is one bit pattern below.
const runStateFields = 7

// TestNextPhase covers every transition of the run's table, then every
// combination of outcomes, so that no sequence plays a stage without a
// reservation and a runtime that came up.
func TestNextPhase(t *testing.T) {
	t.Parallel()

	transitions := []struct {
		name  string
		state RunState
		want  Phase
	}{
		{"a lap with no reservation reserves first", RunState{StagesLeft: true}, PhaseReserve},
		{"a reserved lap brings the runtime up before any stage", RunState{Reserved: true, StagesLeft: true, NextStageWaits: true}, PhaseBringUp},
		{"a bring-up that failed fails the lap", RunState{Reserved: true, BringUpFailed: true, StagesLeft: true}, PhaseFail},
		{"a human stage waits with the runtime up", RunState{Reserved: true, RuntimeUp: true, StagesLeft: true, NextStageWaits: true}, PhaseWait},
		{"a person who did not answer leaves the run waiting", RunState{Reserved: true, RuntimeUp: true, StagesLeft: true, NextStageWaits: true, Waiting: true}, PhaseDone},
		{"a confirmed stage plays", RunState{Reserved: true, RuntimeUp: true, StagesLeft: true}, PhasePlay},
		{"a stage that failed fails the lap", RunState{Reserved: true, RuntimeUp: true, StageFailed: true}, PhaseFail},
		{"no stage left finishes", RunState{Reserved: true, RuntimeUp: true}, PhaseFinish},
		{"a lap with no stage at all brings the runtime up before finishing", RunState{Reserved: true}, PhaseBringUp},
	}

	for _, transition := range transitions {
		t.Run(transition.name, func(t *testing.T) {
			t.Parallel()

			if got := NextPhase(transition.state); got != transition.want {
				t.Errorf("NextPhase(%+v) = %s, want %s", transition.state, got, transition.want)
			}
		})
	}

	t.Run("no combination plays or finishes without a reservation and a runtime", func(t *testing.T) {
		t.Parallel()

		for bits := range 1 << runStateFields {
			state := RunState{
				Reserved:       bits&1 != 0,
				RuntimeUp:      bits&2 != 0,
				BringUpFailed:  bits&4 != 0,
				StagesLeft:     bits&8 != 0,
				NextStageWaits: bits&16 != 0,
				Waiting:        bits&32 != 0,
				StageFailed:    bits&64 != 0,
			}

			phase := NextPhase(state)
			acting := phase == PhasePlay || phase == PhaseWait || phase == PhaseFinish
			if acting && (!state.Reserved || !state.RuntimeUp) {
				t.Errorf("NextPhase(%+v) = %s without a reservation and a runtime", state, phase)
			}
			if phase == PhaseBringUp && !state.Reserved {
				t.Errorf("NextPhase(%+v) brings the runtime up before reserving", state)
			}
			if (state.BringUpFailed || state.StageFailed) && !state.Waiting && phase != PhaseFail {
				t.Errorf("NextPhase(%+v) = %s hides a failure", state, phase)
			}
		}
	})
}
