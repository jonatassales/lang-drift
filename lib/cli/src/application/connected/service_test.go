package connected

import (
	"context"
	"encoding/json"
	"errors"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/jonatassales/tao/lib/cli/src/application"
	"github.com/jonatassales/tao/lib/cli/src/domain/connection"
	"github.com/jonatassales/tao/lib/cli/src/domain/project"
)

type projectFake struct{ owner project.Project }

func (p projectFake) Open(application.ProjectSelection) (project.Project, error) { return p.owner, nil }

type storeFake struct {
	mu     sync.Mutex
	record Connection
	locks  []string
}

func (s *storeFake) List(project.Project) ([]Connection, error) { return []Connection{s.record}, nil }
func (s *storeFake) Save(_ project.Project, record Connection) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, _ := json.Marshal(record)
	_ = json.Unmarshal(data, &s.record)
	return nil
}
func (s *storeFake) Bind(_ context.Context, p project.Project, r Connection) error {
	return s.Save(p, r)
}
func (s *storeFake) Lock(_ project.Project, id string) (func(), error) {
	s.locks = append(s.locks, id)
	return func() {}, nil
}

type nativeFake struct {
	mu            sync.Mutex
	starts, stops int
	failStart     bool
	failObserve   bool
	slowStop      time.Duration
}

func (*nativeFake) VerifyRepository(context.Context, project.Project, Repository) error { return nil }
func (*nativeFake) Materialize(context.Context, project.Project, Connection, connection.Configuration, string) (map[string]string, string, error) {
	return map[string]string{}, "", nil
}
func (n *nativeFake) Start(ctx context.Context, _ project.Project, _ Connection, _ connection.Configuration, report func(string, string) error) error {
	n.mu.Lock()
	n.starts++
	n.mu.Unlock()
	if err := report("build_started", "build_started"); err != nil {
		return err
	}
	if n.failStart {
		return errors.New("build failed")
	}
	return nil
}
func (n *nativeFake) Stop(ctx context.Context, _ project.Project, _ Connection) error {
	n.mu.Lock()
	n.stops++
	delay := n.slowStop
	n.mu.Unlock()
	if delay > 0 {
		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}
func (n *nativeFake) Observe(context.Context, project.Project, Connection) (*Observation, error) {
	if n.failObserve {
		return nil, errors.New("cluster unavailable")
	}
	return &Observation{}, nil
}
func (*nativeFake) Migrate(context.Context, project.Project, string) (string, error) { return "", nil }

type consoleFake struct {
	intent         Intent
	now            time.Time
	claims         int
	seen           []Event
	disconnect     bool
	disconnected   bool
	denyResume     bool
	cancel         context.CancelFunc
	runningResults int
	lastSequence   uint64
}

func (*consoleFake) Bind(context.Context, string, string, string) (Binding, error) {
	return Binding{}, nil
}
func (c *consoleFake) Sync(_ context.Context, _ Connection, sync Sync) (Intent, error) {
	response := c.intent
	for _, event := range sync.Events {
		response.AcceptedEventIDs = append(response.AcceptedEventIDs, event.ID)
		if event.Sequence <= c.lastSequence {
			continue
		}
		c.lastSequence = event.Sequence
		c.seen = append(c.seen, event)
		if event.Kind == "started" {
			c.runningResults++
			if c.runningResults == 2 {
				c.intent.DesiredState = "stopped"
				c.intent.IntentVersion++
				c.intent.Operation.Cause = "user_stop"
			}
		}
	}
	if c.disconnect && !c.disconnected && sync.State == "running" {
		c.disconnected = true
		c.now = c.now.Add(20 * time.Second)
		return Intent{}, errors.New("offline")
	}
	if c.disconnected && c.denyResume {
		c.intent.DesiredState = "stopped"
		c.intent.IntentVersion = 2
		c.intent.Operation.Cause = "user_stop"
		response = c.intent
		for _, event := range sync.Events {
			response.AcceptedEventIDs = append(response.AcceptedEventIDs, event.ID)
		}
	}
	if sync.State == "stopped" && ((c.runningResults >= 2) || (c.disconnected && c.denyResume) || (!c.disconnect && c.claims == 1)) && len(sync.Events) > 0 {
		c.cancel()
	}
	return response, nil
}
func (c *consoleFake) Claim(_ context.Context, _ Connection, intent Intent) (Claim, error) {
	c.claims++
	return Claim{ExecutionID: "execution-" + string(rune('a'+c.claims)), Revision: intent.Revision, IntentVersion: intent.IntentVersion}, nil
}

func configForTest() connection.Configuration {
	return connection.Configuration{Version: 2, RuntimeUID: "runtime-uid", RuntimeID: "runtime", ProjectID: "console-project", Revision: 1, Name: "Example", CPUMillis: 1000, MemoryMiB: 512, Model: "test-model", Effort: "low", Components: []connection.Component{{Name: "api", Path: "api", Port: 3000, Readiness: &connection.Readiness{Kind: "http", Path: "/"}}}}
}

func TestService_DisconnectionReplaysBeforeConditionalResume(t *testing.T) {
	for _, test := range []struct {
		name             string
		stopWhileOffline bool
		wantStarts       int
	}{{"resume current intent", false, 2}, {"user stop prevents resume", true, 1}} {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			native := &nativeFake{}
			store := &storeFake{record: Connection{RuntimeUID: "runtime-uid", RuntimeID: "runtime", ProjectID: "console-project", State: "stopped", Outbox: []Event{}}}
			console := &consoleFake{intent: Intent{DesiredState: "running", Revision: 1, IntentVersion: 1, Markdown: configForTest().Markdown(), Operation: Operation{Cause: "user_start"}, LeaseSeconds: 15}, now: time.Now(), disconnect: true, denyResume: test.stopWhileOffline, cancel: cancel}
			service := Service{Projects: projectFake{}, Store: store, Console: console, Native: native, PollInterval: time.Millisecond, OperationTimeout: time.Second, Now: func() time.Time { return console.now }}
			if err := service.Run(ctx, RunRequest{}); err != nil {
				t.Fatal(err)
			}
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				t.Fatal("agent did not finish scenario")
			}
			if native.starts != test.wantStarts {
				t.Fatalf("starts=%d want %d; events=%+v", native.starts, test.wantStarts, console.seen)
			}
			foundStop, resumed := false, false
			for _, event := range console.seen {
				if event.Kind == "stopped" && event.Cause == "connection_lost" {
					foundStop = true
				}
				if event.Kind == "started" && event.Cause == "reconnect" {
					resumed = true
				}
			}
			if !foundStop {
				t.Fatalf("offline stop evidence missing: %+v", console.seen)
			}
			if resumed == test.stopWhileOffline {
				t.Fatalf("unexpected resume evidence: %+v", console.seen)
			}
		})
	}
}

func TestService_FailedStartCleansUpWithoutRetryingSameIntent(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	native := &nativeFake{failStart: true}
	store := &storeFake{record: Connection{RuntimeUID: "runtime-uid", RuntimeID: "runtime", ProjectID: "console-project", State: "stopped"}}
	console := &consoleFake{intent: Intent{DesiredState: "running", Revision: 1, IntentVersion: 1, Markdown: configForTest().Markdown(), Operation: Operation{Cause: "user_start"}, LeaseSeconds: 15}, now: time.Now(), cancel: cancel}
	service := Service{Projects: projectFake{}, Store: store, Console: console, Native: native, PollInterval: time.Millisecond, OperationTimeout: time.Second}
	if err := service.Run(ctx, RunRequest{}); err != nil {
		t.Fatal(err)
	}
	if native.starts != 1 || native.stops != 1 {
		t.Fatalf("starts=%d stops=%d", native.starts, native.stops)
	}
	for _, event := range console.seen {
		if event.Kind == "started" {
			t.Fatal("failed startup was reported running")
		}
	}
}

func TestSync_BatchesOutboxWithoutLosingRemainder(t *testing.T) {
	events := make([]Event, 250)
	for i := range events {
		events[i] = Event{ID: string(rune(i + 1)), Sequence: uint64(i + 1)}
	}
	record := Connection{Outbox: events}
	batch := syncOf(record)
	if len(batch.Events) != 100 {
		t.Fatalf("batch=%d", len(batch.Events))
	}
	ids := make([]string, 100)
	for i := range ids {
		ids[i] = events[i].ID
	}
	remaining := acknowledge(record.Outbox, ids)
	if len(remaining) != 150 || remaining[0].Sequence != 101 {
		t.Fatal("acknowledgment lost pending evidence")
	}
}

func TestService_MigrationRecoveryIsExplicit(t *testing.T) {
	native := &nativeFake{}
	service := Service{Projects: projectFake{}, Store: &storeFake{}, Native: native}
	if _, err := service.Migrate(context.Background(), MigrateRequest{Runtime: "legacy"}); err != nil {
		t.Fatal(err)
	}
	if native.stops != 0 {
		t.Fatal("ordinary migration executed cleanup")
	}
	if _, err := service.Migrate(context.Background(), MigrateRequest{Runtime: "legacy", Recover: true}); err != nil {
		t.Fatal(err)
	}
	if native.stops != 1 {
		t.Fatal("explicit migration recovery did not invoke owned cleanup")
	}
}

func TestService_MigrationTakesTheSupervisorLockAndParsesItsRuntime(t *testing.T) {
	store := &storeFake{}
	service := Service{Projects: projectFake{}, Store: store, Native: &nativeFake{}}

	if _, err := service.Migrate(context.Background(), MigrateRequest{Runtime: "not a valid id!"}); err == nil {
		t.Fatal("an unparseable runtime id was accepted")
	}

	if _, err := service.Migrate(context.Background(), MigrateRequest{Runtime: "legacy"}); err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(store.locks, "project-agent") {
		t.Fatalf("locks = %v, want the supervisor's project-agent lock", store.locks)
	}
}

func TestService_ObservedRuntimeFailureRequiresConfirmedCleanup(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	native := &nativeFake{failObserve: true}
	store := &storeFake{record: Connection{RuntimeUID: "runtime-uid", RuntimeID: "runtime", ProjectID: "console-project", State: "stopped"}}
	console := &consoleFake{intent: Intent{DesiredState: "running", Revision: 1, IntentVersion: 1, Markdown: configForTest().Markdown(), Operation: Operation{Cause: "user_start"}, LeaseSeconds: 15}, now: time.Now(), cancel: cancel}
	service := Service{Projects: projectFake{}, Store: store, Console: console, Native: native, PollInterval: time.Millisecond, OperationTimeout: time.Second}
	if err := service.Run(ctx, RunRequest{}); err != nil {
		t.Fatal(err)
	}
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		t.Fatal("observed failure was not reconciled")
	}
	failed, stopped := false, false
	for _, event := range console.seen {
		if event.Cause == "runtime_failure" {
			if event.Kind == "failed" {
				failed = true
			}
			if event.Kind == "stopped" {
				stopped = true
			}
		}
	}
	if !failed || !stopped || native.stops != 1 {
		t.Fatalf("failure/cleanup evidence missing: %+v", console.seen)
	}
}

func TestService_StopContinuesAcrossMultipleHeartbeats(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	native := &nativeFake{failStart: true, slowStop: 30 * time.Millisecond}
	store := &storeFake{record: Connection{RuntimeUID: "runtime-uid", RuntimeID: "runtime", ProjectID: "console-project", State: "stopped"}}
	console := &consoleFake{intent: Intent{DesiredState: "running", Revision: 1, IntentVersion: 1, Markdown: configForTest().Markdown(), Operation: Operation{Cause: "user_start"}, LeaseSeconds: 15}, now: time.Now(), cancel: cancel}
	// After failed startup, the same user stop intent is observed repeatedly
	// while the bounded cleanup worker is still making progress.
	wrapped := &stopConsole{consoleFake: console}
	service := Service{Projects: projectFake{}, Store: store, Console: wrapped, Native: native, PollInterval: time.Millisecond, OperationTimeout: time.Second}
	if err := service.Run(ctx, RunRequest{}); err != nil {
		t.Fatal(err)
	}
	if errors.Is(ctx.Err(), context.DeadlineExceeded) || native.stops != 1 {
		t.Fatalf("cleanup was interrupted/restarted: stops=%d err=%v", native.stops, ctx.Err())
	}
}

type stopConsole struct{ *consoleFake }

func (c *stopConsole) Sync(ctx context.Context, record Connection, snapshot Sync) (Intent, error) {
	if snapshot.State == "stopping" || snapshot.State == "unknown" {
		c.intent.DesiredState = "stopped"
		c.intent.Operation.Cause = "user_stop"
	}
	return c.consoleFake.Sync(ctx, record, snapshot)
}

func TestService_LostClaimResponseReconcilesAfterChangedIntent(t *testing.T) {
	for _, action := range []string{"stopped", "deleted", "revoked"} {
		t.Run(action, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			native := &nativeFake{}
			store := &storeFake{record: Connection{RuntimeUID: "runtime-uid", RuntimeID: "runtime", ProjectID: "console-project", State: "stopped"}}
			console := &lostClaimConsole{action: action, consoleFake: &consoleFake{intent: Intent{DesiredState: "running", Revision: 1, IntentVersion: 1, Markdown: configForTest().Markdown(), Operation: Operation{Cause: "user_start"}, LeaseSeconds: 15}, cancel: cancel}}
			service := Service{Projects: projectFake{}, Store: store, Console: console, Native: native, PollInterval: time.Millisecond, OperationTimeout: time.Second}
			if err := service.Run(ctx, RunRequest{}); err != nil {
				t.Fatal(err)
			}
			if errors.Is(ctx.Err(), context.DeadlineExceeded) || native.starts != 0 || native.stops != 1 || console.claims != 1 {
				t.Fatalf("lost claim was not reconciled: starts=%d stops=%d claims=%d err=%v", native.starts, native.stops, console.claims, ctx.Err())
			}
			confirmed := false
			for _, event := range console.seen {
				if event.Kind == "stopped" && event.ExecutionID == "claimed-but-response-lost" && event.Revision == 1 && event.IntentVersion == 1 {
					confirmed = true
				}
			}
			if !confirmed {
				t.Fatalf("claimed execution has no terminal reconciliation evidence: %+v", console.seen)
			}
		})
	}
}

type lostClaimConsole struct {
	*consoleFake
	action string
}

func (c *lostClaimConsole) Claim(context.Context, Connection, Intent) (Claim, error) {
	c.claims++
	c.intent.DesiredState = "stopped"
	c.intent.Deleted = c.action == "deleted"
	c.intent.Revoked = c.action == "revoked"
	c.intent.IntentVersion = 2
	c.intent.Operation = Operation{ID: "cleanup-request", Cause: "user_stop"}
	c.intent.UnresolvedExecution = &UnresolvedExecution{ExecutionID: "claimed-but-response-lost", IntentVersion: 1, Revision: 1, State: "starting"}
	return Claim{}, errors.New("claim committed but response was lost")
}
