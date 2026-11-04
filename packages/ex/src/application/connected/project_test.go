package connected

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jonatassales/tao/packages/ex/src/domain/connection"
	"github.com/jonatassales/tao/packages/ex/src/domain/project"
)

type projectConsoleFake struct {
	consoleFake
	calls int
	binds int
}

func (console *projectConsoleFake) ProjectSync(context.Context, string, string) (ProjectMembership, error) {
	console.calls++
	result := ProjectMembership{ProjectID: "project"}
	if console.calls > 1 {
		result.RuntimeIDs = []string{"runtime-uid"}
	}
	return result, nil
}
func (console *projectConsoleFake) ProjectBind(_ context.Context, _, _, runtime, environment string) (Binding, error) {
	console.binds++
	return Binding{EnvironmentID: environment, Credential: "synthetic-recovery", RuntimeUID: runtime, RuntimeID: "runtime", ProjectID: "project", Repository: Repository{Kind: "local"}}, nil
}

type projectStoreFake struct {
	storeFake
	bindings []Connection
}

func (store *projectStoreFake) List(project.Project) ([]Connection, error) { return nil, nil }
func (store *projectStoreFake) Bind(_ context.Context, _ project.Project, record Connection) error {
	store.bindings = append(store.bindings, record)
	return nil
}

type projectNativeFake struct {
	nativeFake
	once       sync.Once
	configured chan struct{}
}

func (native *projectNativeFake) Materialize(context.Context, project.Project, Connection, connection.Configuration, string) (map[string]string, string, error) {
	native.once.Do(func() { close(native.configured) })
	return map[string]string{}, "", nil
}

func TestProjectAgentDiscoversLaterRuntimeWithoutStartingIt(t *testing.T) {
	config := connection.Configuration{Version: 2, Name: "api", ProjectID: "project", RuntimeUID: "runtime-uid", RuntimeID: "runtime", Revision: 1, CPUMillis: 1000, MemoryMiB: 256, Model: "example/model", Effort: "medium", Components: []connection.Component{{Path: ".", Name: "api", Port: 3000, Readiness: &connection.Readiness{Kind: "http", Path: "/"}}}}
	console := &projectConsoleFake{consoleFake: consoleFake{intent: Intent{DesiredState: "stopped", Revision: 1, Markdown: config.Markdown(), LeaseSeconds: 15}}}
	store := &projectStoreFake{}
	native := &projectNativeFake{configured: make(chan struct{})}
	service := Service{Projects: projectFake{}, Store: store, Console: console, ProjectConsole: console, Native: native, PollInterval: time.Millisecond}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- service.RunProject(ctx, ProjectRunRequest{URL: "http://127.0.0.1:3001", Secret: "synthetic-project"})
	}()
	select {
	case <-native.configured:
	case err := <-done:
		cancel()
		t.Fatalf("agent exited before configuration: %v", err)
	case <-time.After(2 * time.Second):
		cancel()
		<-done
		t.Fatal("agent did not discover later runtime")
	}
	cancel()
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	if len(store.bindings) != 1 || console.binds != 1 {
		t.Fatalf("bindings=%d calls=%d", len(store.bindings), console.binds)
	}
	if native.starts != 0 || console.claims != 0 {
		t.Fatal("project connection started a microruntime")
	}
}

type recoveryProjectConsole struct {
	consoleFake
	stopped chan struct{}
	once    sync.Once
}

func (*recoveryProjectConsole) ProjectSync(context.Context, string, string) (ProjectMembership, error) {
	return ProjectMembership{}, ErrAccessRevoked
}
func (*recoveryProjectConsole) ProjectBind(context.Context, string, string, string, string) (Binding, error) {
	panic("revoked project must not bind")
}
func (console *recoveryProjectConsole) Sync(_ context.Context, _ Connection, sync Sync) (Intent, error) {
	result := Intent{Revoked: false, DesiredState: "running", Revision: 1, IntentVersion: 1}
	for _, event := range sync.Events {
		result.AcceptedEventIDs = append(result.AcceptedEventIDs, event.ID)
	}
	if sync.State == "stopped" && len(sync.Events) > 0 {
		console.once.Do(func() { close(console.stopped) })
	}
	return result, nil
}
func TestRevokedProjectRestartsForCleanupAndOutboxRecoveryOnly(t *testing.T) {
	console := &recoveryProjectConsole{stopped: make(chan struct{})}
	store := &storeFake{record: Connection{URL: "http://127.0.0.1:3001", RuntimeUID: "uid", RuntimeID: "runtime", ProjectID: "project", State: "running", Revision: 1, Credential: "synthetic-recovery", Outbox: []Event{{ID: "pending", Sequence: 1, Kind: "started"}}, Sequence: 1}}
	native := &nativeFake{}
	service := Service{Projects: projectFake{}, Store: store, Console: console, ProjectConsole: console, Native: native, PollInterval: time.Millisecond}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	done := make(chan error, 1)
	go func() {
		done <- service.RunProject(ctx, ProjectRunRequest{URL: "http://127.0.0.1:3001", Secret: "revoked"})
	}()
	select {
	case <-console.stopped:
	case err := <-done:
		t.Fatalf("recovery exited before acknowledgement: %v", err)
	case <-ctx.Done():
		<-done
		t.Fatal("recovery did not complete")
	}
	cancel()
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	if native.starts != 0 || native.stops == 0 || console.claims != 0 {
		t.Fatal("revoked recovery started work or skipped cleanup")
	}
	if len(store.record.Outbox) != 0 || store.record.State != "stopped" {
		t.Fatal("recovery did not persist acknowledgement")
	}
}

func TestProjectRevocationFencesLegacyStoppedWorker(t *testing.T) {
	var revoked atomic.Bool
	revoked.Store(true)
	base := &consoleFake{intent: Intent{DesiredState: "running", Revision: 1}}
	gated := projectRecoveryConsole{Console: base, revoked: &revoked}
	intent, err := gated.Sync(context.Background(), Connection{State: "stopped"}, Sync{})
	if err != nil || !intent.Revoked {
		t.Fatal("legacy intent bypassed project revocation")
	}
	if _, err := gated.Claim(context.Background(), Connection{}, intent); err != ErrAccessRevoked || base.claims != 0 {
		t.Fatal("legacy claim bypassed project revocation")
	}
}
