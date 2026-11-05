package execution

// Phase is what a run does next. The run loop asks NextPhase after each
// phase and executes the one it names through the existing ports.
type Phase string

// The phases of a run, in the order a run that meets no failure and no
// person goes through them.
const (
	// PhaseReserve opens the gate and reserves the lap's money. Start does
	// it before the loop; a loop that is asked for it holds no reservation.
	PhaseReserve Phase = "reserve"
	// PhaseBringUp materializes, builds and starts the runtime, then waits
	// for it to be ready.
	PhaseBringUp Phase = "bringUp"
	// PhasePlay plays the next stage through the driver.
	PhasePlay Phase = "play"
	// PhaseWait asks a person to confirm the next stage; without a yes the
	// run stays waiting with the runtime up.
	PhaseWait Phase = "wait"
	// PhaseFinish observes delivery, takes the runtime down, releases the
	// unused reservation and remembers.
	PhaseFinish Phase = "finish"
	// PhaseFail takes the runtime down, ends the lap as failed and releases
	// what the reservation did not spend.
	PhaseFail Phase = "fail"
	// PhaseDone ends the loop: the run waits for a person.
	PhaseDone Phase = "done"
)
