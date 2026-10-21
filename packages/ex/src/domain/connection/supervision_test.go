package connection

import (
	"errors"
	"testing"
)

func TestAgentState_Launch(t *testing.T) {
	tests := []struct {
		state  AgentState
		action AgentAction
		next   AgentState
		fail   bool
	}{
		{StateStopped, ActionStart, StateStarting, false},
		{StateRunning, ActionStop, StateStopping, false},
		{StateUnknown, ActionStop, StateStopping, false},
		{StateRunning, ActionObserve, StateRunning, false},

		{StateRunning, ActionStart, "", true},
		{StateStarting, ActionStart, "", true},
		{StateUnknown, ActionStart, "", true},
		{StateStopped, ActionStop, "", true},
		{StateStarting, ActionStop, "", true},
		{StateStopping, ActionStop, "", true},
		{StateStopped, ActionObserve, "", true},
		{StateUnknown, ActionObserve, "", true},
		{AgentState("bogus"), ActionStart, "", true},
	}
	for _, test := range tests {
		next, err := test.state.Launch(test.action)
		if (err != nil) != test.fail || (!test.fail && next != test.next) {
			t.Errorf("%s.Launch(%s) = %s, %v; want %s fail=%t", test.state, test.action, next, err, test.next, test.fail)
		}
		if test.fail && !errors.Is(err, ErrIllegalTransition) {
			t.Errorf("%s.Launch(%s) error %v does not match ErrIllegalTransition", test.state, test.action, err)
		}
	}
}

func TestAgentState_Settle(t *testing.T) {
	tests := []struct {
		state     AgentState
		action    AgentAction
		succeeded bool
		next      AgentState
		fail      bool
	}{
		{StateStarting, ActionStart, true, StateRunning, false},
		{StateStarting, ActionStart, false, StateUnknown, false},
		{StateStopping, ActionStop, true, StateStopped, false},
		{StateStopping, ActionStop, false, StateUnknown, false},
		{StateRunning, ActionObserve, true, StateRunning, false},
		{StateRunning, ActionObserve, false, StateUnknown, false},

		{StateStopped, ActionStart, true, "", true},
		{StateRunning, ActionStart, true, "", true},
		{StateRunning, ActionStop, true, "", true},
		{StateStopping, ActionObserve, true, "", true},
	}
	for _, test := range tests {
		next, err := test.state.Settle(test.action, test.succeeded)
		if (err != nil) != test.fail || (!test.fail && next != test.next) {
			t.Errorf("%s.Settle(%s, %t) = %s, %v; want %s fail=%t", test.state, test.action, test.succeeded, next, err, test.next, test.fail)
		}
	}
}

func TestAgentState_Vocabulary(t *testing.T) {
	// The persisted and synchronized values are a contract with the console
	// and with existing state files; the types must not change them.
	states := map[AgentState]string{
		StateStopped:  "stopped",
		StateStarting: "starting",
		StateRunning:  "running",
		StateStopping: "stopping",
		StateUnknown:  "unknown",
	}
	for state, text := range states {
		if string(state) != text {
			t.Errorf("state %q changed its persisted value", text)
		}
	}
	causes := map[Cause]string{
		CauseUserStart:      "user_start",
		CauseUserStop:       "user_stop",
		CauseDeleted:        "deleted",
		CauseRevoked:        "revoked",
		CauseReconnect:      "reconnect",
		CauseConnectionLost: "connection_lost",
		CauseAgentRestart:   "agent_restart",
		CauseAgentShutdown:  "agent_shutdown",
		CauseStartFailed:    "start_failed",
		CauseRuntimeFailure: "runtime_failure",
		CauseObservation:    "observation",
		CauseConfiguration:  "configuration",
	}
	for cause, text := range causes {
		if string(cause) != text {
			t.Errorf("cause %q changed its synchronized value", text)
		}
	}
}
