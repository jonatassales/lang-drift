package connected

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/jonatassales/tao/lib/cli/src/domain/connection"
	"github.com/jonatassales/tao/lib/cli/src/domain/project"
)

// emptyStore lists no existing connections.
type emptyStore struct{ storeFake }

func (*emptyStore) List(project.Project) ([]Connection, error) { return nil, nil }

// judgeCheckFake answers the connection check as the console's judge route
// would, and remembers what it was asked.
type judgeCheckFake struct {
	err    error
	url    string
	secret string
}

func (check *judgeCheckFake) Check(_ context.Context, consoleURL string, secret string) error {
	check.url, check.secret = consoleURL, secret

	return check.err
}

// judgeConnectionsFake keeps the recorded console in memory.
type judgeConnectionsFake struct {
	url    string
	secret string
}

func (connections *judgeConnectionsFake) Record(_ context.Context, _ project.Project, consoleURL string, secret string) error {
	connections.url, connections.secret = consoleURL, secret

	return nil
}

func TestService_Connect(t *testing.T) {
	t.Run("requires the secret", func(t *testing.T) {
		service := Service{Projects: projectFake{}}

		_, err := service.Connect(context.Background(), ConnectRequest{})
		if err == nil || !strings.Contains(err.Error(), "RUNTIME_SECRET") {
			t.Fatalf("err = %v", err)
		}
	})

	t.Run("records the console once the judge route accepts the secret", func(t *testing.T) {
		check := &judgeCheckFake{}
		connections := &judgeConnectionsFake{}
		service := Service{
			Projects:         projectFake{},
			Judge:            check,
			JudgeConnections: connections,
		}

		result, err := service.Connect(context.Background(), ConnectRequest{
			URL:    "https://console.example",
			Secret: "runtime-secret",
		})
		if err != nil {
			t.Fatal(err)
		}

		if check.url != "https://console.example" || check.secret != "runtime-secret" {
			t.Fatalf("checked %q with %q", check.url, check.secret)
		}
		if connections.url != "https://console.example" || connections.secret != "runtime-secret" {
			t.Fatalf("recorded %q with %q", connections.url, connections.secret)
		}
		if !strings.Contains(result.Status, "https://console.example") || strings.Contains(result.Status, "runtime-secret") {
			t.Fatalf("status = %q", result.Status)
		}
	})

	t.Run("records nothing when the judge route refuses the secret", func(t *testing.T) {
		connections := &judgeConnectionsFake{}
		service := Service{
			Projects:         projectFake{},
			Judge:            &judgeCheckFake{err: errors.New("the console at https://console.example refused RUNTIME_SECRET")},
			JudgeConnections: connections,
		}

		_, err := service.Connect(context.Background(), ConnectRequest{
			URL:    "https://console.example",
			Secret: "wrong",
		})
		if err == nil || !strings.Contains(err.Error(), "https://console.example") {
			t.Fatalf("err = %v", err)
		}
		if connections.url != "" {
			t.Fatalf("a refused connection was recorded: %q", connections.url)
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
