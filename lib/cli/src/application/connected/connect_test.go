package connected

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/jonatassales/tao/lib/cli/src/domain/connection"
	"github.com/jonatassales/tao/lib/cli/src/domain/project"
)

// bindingConsole returns a configurable binding for Connect tests.
type bindingConsole struct {
	consoleFake
	binding Binding
	err     error
}

func (c *bindingConsole) Bind(_ context.Context, _, _, environment string) (Binding, error) {
	if c.binding.EnvironmentID == "" {
		c.binding.EnvironmentID = environment
	}
	return c.binding, c.err
}

// emptyStore lists no existing connections.
type emptyStore struct{ storeFake }

func (*emptyStore) List(project.Project) ([]Connection, error) { return nil, nil }

func TestService_Connect(t *testing.T) {
	healthy := Binding{Credential: "agent-credential", RuntimeUID: "runtime-uid", RuntimeID: "runtime", ProjectID: "console-project"}

	t.Run("requires the secret", func(t *testing.T) {
		service := Service{Projects: projectFake{}}
		if _, err := service.Connect(context.Background(), ConnectRequest{}); err == nil || !strings.Contains(err.Error(), "RUNTIME_SECRET") {
			t.Fatalf("err = %v", err)
		}
	})

	t.Run("binds and stores a stopped connection", func(t *testing.T) {
		store := &emptyStore{}
		service := Service{Projects: projectFake{}, Store: store, Console: &bindingConsole{binding: healthy}, Native: &nativeFake{}}
		result, err := service.Connect(context.Background(), ConnectRequest{URL: "https://console.example", Secret: "runtime-secret"})
		if err != nil {
			t.Fatal(err)
		}
		if result.Runtime != "runtime" || result.Environment == "" {
			t.Fatalf("result = %+v", result)
		}
		// The fake persists through the JSON contract, which excludes
		// credentials by design; identity and state prove the bind.
		record := store.record
		if record.State != connection.StateStopped || record.RuntimeUID != "runtime-uid" || record.EnvironmentID != result.Environment {
			t.Fatalf("record = %+v", record)
		}
	})

	t.Run("refuses an incomplete binding", func(t *testing.T) {
		service := Service{Projects: projectFake{}, Store: &emptyStore{}, Console: &bindingConsole{binding: Binding{RuntimeID: "runtime"}}, Native: &nativeFake{}}
		if _, err := service.Connect(context.Background(), ConnectRequest{URL: "https://console.example", Secret: "s"}); err == nil || !strings.Contains(err.Error(), "incomplete binding") {
			t.Fatalf("err = %v", err)
		}
	})

	t.Run("refuses a binding error", func(t *testing.T) {
		service := Service{Projects: projectFake{}, Store: &emptyStore{}, Console: &bindingConsole{err: errors.New("offline")}, Native: &nativeFake{}}
		if _, err := service.Connect(context.Background(), ConnectRequest{URL: "https://console.example", Secret: "s"}); err == nil {
			t.Fatal("binding failure was swallowed")
		}
	})

	t.Run("refuses a runtime that is already connected", func(t *testing.T) {
		store := &storeFake{record: Connection{RuntimeUID: "runtime-uid"}}
		service := Service{Projects: projectFake{}, Store: store, Console: &bindingConsole{binding: healthy}, Native: &nativeFake{}}
		if _, err := service.Connect(context.Background(), ConnectRequest{URL: "https://console.example", Secret: "s"}); err == nil || !strings.Contains(err.Error(), "already connected") {
			t.Fatalf("err = %v", err)
		}
	})
}

func TestConnection_StringHidesCredentials(t *testing.T) {
	record := Connection{RuntimeID: "runtime", EnvironmentID: "env", State: connection.StateRunning, Secret: "private", Credential: "private"}
	text := record.String()
	if strings.Contains(text, "private") {
		t.Fatalf("String() leaks credentials: %q", text)
	}
	if !strings.Contains(text, "runtime") || !strings.Contains(text, "running") {
		t.Fatalf("String() = %q", text)
	}
}
