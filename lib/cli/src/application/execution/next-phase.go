package execution

// NextPhase decides what the run does next from what has happened in the
// lap. A run that asked a person and got no answer is done for now; a failed
// bring-up or stage fails the lap; before anything else the gate reserves;
// the runtime comes up before any stage; a stage that waits for a person
// waits with the runtime up; a stage that does not plays; with none left,
// the run finishes.
func NextPhase(state RunState) Phase {
	switch {
	case state.Waiting:
		return PhaseDone
	case state.BringUpFailed || state.StageFailed:
		return PhaseFail
	case !state.Reserved:
		return PhaseReserve
	case !state.RuntimeUp:
		return PhaseBringUp
	case state.StagesLeft && state.NextStageWaits:
		return PhaseWait
	case state.StagesLeft:
		return PhasePlay
	}

	return PhaseFinish
}
