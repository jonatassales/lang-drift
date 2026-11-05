package execution

// RunState is what has happened in the lap so far, as NextPhase reads it.
// The loop builds it from the lap, the run record and what the phase it
// just executed returned.
type RunState struct {
	// Reserved says the gate opened and the lap holds its reservation.
	Reserved bool
	// RuntimeUp says the runtime came up and was ready in this lap.
	RuntimeUp bool
	// BringUpFailed says the runtime did not come up or was not ready.
	BringUpFailed bool
	// StagesLeft says a stage the run plays, before the owner's stage, is
	// not done.
	StagesLeft bool
	// NextStageWaits says that stage waits for a person and nobody has
	// confirmed it.
	NextStageWaits bool
	// Waiting says the run asked a person and got no confirmation.
	Waiting bool
	// StageFailed says the stage the run played failed.
	StageFailed bool
}
