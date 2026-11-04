package connection

import (
	"errors"
	"fmt"
)

// ErrIllegalTransition marks a supervision step the state machine refuses.
// Hitting it means the supervisor's bookkeeping drifted, never the runtime.
var ErrIllegalTransition = errors.New("illegal supervision transition")

// AgentState is where one supervised runtime stands. The string values are
// persisted in connection state and synchronized with the console, so they
// never change.
type AgentState string

const (
	StateStopped  AgentState = "stopped"
	StateStarting AgentState = "starting"
	StateRunning  AgentState = "running"
	StateStopping AgentState = "stopping"
	// StateUnknown means the last operation did not settle; only a confirmed
	// stop leaves it.
	StateUnknown AgentState = "unknown"
)

// AgentAction is one bounded native operation the supervisor can own at a time.
type AgentAction string

const (
	ActionStart   AgentAction = "start"
	ActionStop    AgentAction = "stop"
	ActionObserve AgentAction = "observe"
)

// Cause names why an operation ran, exactly as events synchronize it.
type Cause string

const (
	CauseUserStart      Cause = "user_start"
	CauseUserStop       Cause = "user_stop"
	CauseDeleted        Cause = "deleted"
	CauseRevoked        Cause = "revoked"
	CauseReconnect      Cause = "reconnect"
	CauseConnectionLost Cause = "connection_lost"
	CauseAgentRestart   Cause = "agent_restart"
	CauseAgentShutdown  Cause = "agent_shutdown"
	CauseStartFailed    Cause = "start_failed"
	CauseRuntimeFailure Cause = "runtime_failure"
	CauseObservation    Cause = "observation"
	CauseConfiguration  Cause = "configuration"
)

// Launch is the state an action moves into when it begins, or a refusal:
// a start leaves stopped only, an observation needs a running runtime, and a
// stop settles a running or unsettled runtime.
func (state AgentState) Launch(action AgentAction) (AgentState, error) {
	switch {
	case action == ActionStart && state == StateStopped:
		return StateStarting, nil

	case action == ActionStop && (state == StateRunning || state == StateUnknown):
		return StateStopping, nil

	case action == ActionObserve && state == StateRunning:
		return StateRunning, nil
	}

	return "", fmt.Errorf("%w: cannot launch %s from %s", ErrIllegalTransition, action, state)
}

// Settle is the state a completed action lands in. A failed start or stop is
// unsettled, never silently stopped; a failed observation marks the runtime
// unsettled so a stop reconciles it.
func (state AgentState) Settle(action AgentAction, succeeded bool) (AgentState, error) {
	legal := (action == ActionStart && state == StateStarting) ||
		(action == ActionStop && state == StateStopping) ||
		(action == ActionObserve && state == StateRunning)
	if !legal {
		return "", fmt.Errorf("%w: cannot settle %s in %s", ErrIllegalTransition, action, state)
	}

	if !succeeded {
		return StateUnknown, nil
	}

	switch action {
	case ActionStart:
		return StateRunning, nil

	case ActionStop:
		return StateStopped, nil

	default:
		return StateRunning, nil
	}
}
