package connected

import (
	"context"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jonatassales/tao/lib/cli/src/domain/connection"
)

// scriptedConsole walks one sync response per call and records what it saw.
type scriptedConsole struct {
	responses []Intent
	errors    []error
	calls     atomic.Int64
	claim     Claim
	claimErr  error
	seen      []Event
	cancel    context.CancelFunc
}

func (c *scriptedConsole) Sync(_ context.Context, _ Connection, sync Sync) (Intent, error) {
	index := int(c.calls.Add(1)) - 1
	c.seen = append(c.seen, sync.Events...)
	if index >= len(c.responses) {
		c.cancel()
		index = len(c.responses) - 1
	}
	var err error
	if index < len(c.errors) {
		err = c.errors[index]
	}
	response := c.responses[index]
	for _, event := range sync.Events {
		response.AcceptedEventIDs = append(response.AcceptedEventIDs, event.ID)
	}
	return response, err
}

func (c *scriptedConsole) Claim(context.Context, Connection, Intent) (Claim, error) {
	return c.claim, c.claimErr
}

func runScenario(t *testing.T, store *storeFake, console *scriptedConsole, native *nativeFake) error {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	console.cancel = cancel
	service := Service{Projects: projectFake{}, Store: store, Console: console, Native: native, PollInterval: time.Millisecond, OperationTimeout: time.Second}
	return service.Run(ctx, RunRequest{})
}

func TestService_InconsistentStartClaimFailsLoudly(t *testing.T) {
	markdown := configForTest().Markdown()
	store := &storeFake{record: Connection{RuntimeUID: "runtime-uid", RuntimeID: "runtime", ProjectID: "console-project", State: connection.StateStopped}}
	console := &scriptedConsole{
		responses: []Intent{
			{DesiredState: "running", Revision: 1, IntentVersion: 1, Markdown: markdown, Operation: Operation{ID: "op-1", Cause: "user_start"}},
		},
		claim: Claim{ExecutionID: "", Revision: 1, IntentVersion: 1},
	}

	err := runScenario(t, store, console, &nativeFake{})
	if err == nil || !strings.Contains(err.Error(), "inconsistent start claim") {
		t.Fatalf("err = %v", err)
	}
	var kinds []string
	for _, event := range console.seen {
		kinds = append(kinds, event.Kind)
	}
	if len(kinds) == 0 || kinds[0] != "configured" {
		t.Fatalf("events = %v, want the configuration applied before the claim", kinds)
	}
}

func TestService_AdoptsAnUnresolvedExecutionAndReconciles(t *testing.T) {
	store := &storeFake{record: Connection{RuntimeUID: "runtime-uid", RuntimeID: "runtime", ProjectID: "console-project", State: connection.StateStopped}}
	console := &scriptedConsole{
		responses: []Intent{
			{DesiredState: "stopped", Revision: 1, IntentVersion: 1, UnresolvedExecution: &UnresolvedExecution{ExecutionID: "execution-9", Revision: 1, IntentVersion: 1}, Operation: Operation{ID: "op-1"}},
			{DesiredState: "stopped", Revision: 1, IntentVersion: 1},
		},
	}
	native := &nativeFake{}

	if err := runScenario(t, store, console, native); err != nil {
		t.Fatal(err)
	}
	if native.stops == 0 {
		t.Fatal("the adopted execution was never reconciled by a stop")
	}
	if store.record.ExecutionID != "execution-9" {
		t.Fatalf("record execution = %q, want the adopted identity", store.record.ExecutionID)
	}
}

func TestService_RefusesAnInvalidReconciliationIdentity(t *testing.T) {
	store := &storeFake{record: Connection{RuntimeUID: "runtime-uid", RuntimeID: "runtime", State: connection.StateStopped}}
	console := &scriptedConsole{
		responses: []Intent{
			{DesiredState: "stopped", Revision: 1, IntentVersion: 1, UnresolvedExecution: &UnresolvedExecution{ExecutionID: "execution-9", Revision: 0, IntentVersion: 1}},
		},
	}

	err := runScenario(t, store, console, &nativeFake{})
	if err == nil || !strings.Contains(err.Error(), "invalid reconciliation identity") {
		t.Fatalf("err = %v", err)
	}
}

func TestService_RunRequiresAConnectedRuntime(t *testing.T) {
	service := Service{Projects: projectFake{}, Store: &emptyStore{}}
	err := service.Run(context.Background(), RunRequest{})
	if err == nil || !strings.Contains(err.Error(), "no connected runtimes") {
		t.Fatalf("err = %v", err)
	}
}

func TestService_ShutdownStopsARunningRuntime(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	store := &storeFake{record: Connection{RuntimeUID: "runtime-uid", RuntimeID: "runtime", ProjectID: "console-project", State: connection.StateStopped, Revision: 1, Markdown: configForTest().Markdown()}}
	console := &scriptedConsole{
		responses: []Intent{
			{DesiredState: "running", Revision: 1, IntentVersion: 1, Operation: Operation{ID: "op-1", Cause: "user_start"}},
		},
		claim: Claim{ExecutionID: "execution-1", Revision: 1, IntentVersion: 1},
	}
	console.cancel = func() {}
	native := &nativeFake{}

	service := Service{Projects: projectFake{}, Store: store, Console: console, Native: native, PollInterval: time.Millisecond, OperationTimeout: time.Second}
	done := make(chan error, 1)
	go func() { done <- service.Run(ctx, RunRequest{}) }()

	for deadline := time.Now().Add(time.Second); time.Now().Before(deadline); time.Sleep(5 * time.Millisecond) {
		store.mu.Lock()
		state := store.record.State
		store.mu.Unlock()
		if state == connection.StateRunning {
			break
		}
	}
	cancel()
	if err := <-done; err != nil {
		t.Fatal(err)
	}

	if native.stops == 0 {
		t.Fatal("shutdown left the runtime running")
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	if store.record.State != connection.StateStopped {
		t.Fatalf("state = %s, want stopped after shutdown", store.record.State)
	}
}
