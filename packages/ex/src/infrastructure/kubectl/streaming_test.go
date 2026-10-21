package kubectl_test

import (
	"context"
	"slices"
	"strings"
	"testing"

	"github.com/exphere/ex/src/infrastructure/kubectl"
)

func TestClient_Logs(t *testing.T) {
	t.Parallel()

	resource := deploymentResource()

	tests := []struct {
		name   string
		follow bool
		last   string
	}{
		{name: "bounded read", last: "--limit-bytes=1048576"},
		{name: "follow until interrupted", follow: true, last: "--follow=true"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			runner := &recorder{}
			var output, diagnostics strings.Builder

			if err := (kubectl.Client{Runner: runner}).Logs(context.Background(), target, resource, tt.follow, &output, &diagnostics); err != nil {
				t.Fatal(err)
			}

			recorded := runner.calls[0]
			want := []string{
				"--kubeconfig", target.Kubeconfig, "--context", target.Context, "--namespace", "demo-local",
				"logs", "deployment/web", "--all-containers=true", "--tail=500", "--timestamps=true", tt.last,
			}

			if !slices.Equal(recorded.args, want) {
				t.Errorf("args = %v, want %v", recorded.args, want)
			}

			if recorded.deadline != 0 || recorded.stdout != &output || recorded.stderr != &diagnostics {
				t.Errorf("streams or deadline not passed through: deadline %v", recorded.deadline)
			}
		})
	}
}

func TestClient_Forward(t *testing.T) {
	t.Parallel()

	runner := &recorder{}
	var diagnostics strings.Builder

	if err := (kubectl.Client{Runner: runner}).Forward(context.Background(), target, serviceResource(), 9090, 8080, &diagnostics); err != nil {
		t.Fatal(err)
	}

	recorded := runner.calls[0]
	want := []string{
		"--kubeconfig", target.Kubeconfig, "--context", target.Context, "--namespace", "demo-local",
		"port-forward", "service/web", "9090:8080", "--address=127.0.0.1",
	}

	if !slices.Equal(recorded.args, want) || recorded.stdout != &diagnostics || recorded.stderr != &diagnostics || recorded.deadline != 0 {
		t.Fatalf("call = %+v, want a localhost-only tunnel whose streams are diagnostics", recorded)
	}
}
