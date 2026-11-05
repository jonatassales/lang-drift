package handoff

// RunOutcome is where the run stands (spec 0065).
type RunOutcome string

// A run is running, paused at the gate, stopped, or completed once the ready
// line prints.
const (
	RunRunning   RunOutcome = "running"
	RunPaused    RunOutcome = "paused"
	RunStopped   RunOutcome = "stopped"
	RunCompleted RunOutcome = "completed"
)
