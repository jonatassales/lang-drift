package kubectl_test

import (
	"context"
	"errors"
	"io"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/infrastructure/kubectl"
	"github.com/exphere/ex/src/infrastructure/process"
)

// call is one recorded kubectl invocation.
type call struct {
	args     []string
	stdin    string
	deadline time.Duration
	stdout   io.Writer
	stderr   io.Writer
}

// recorder records invocations and answers them with respond, which receives
// the arguments after the explicit target flags.
type recorder struct {
	mu      sync.Mutex
	calls   []call
	respond func(arguments []string) (process.Result, error)
}

func (recorder *recorder) Run(ctx context.Context, command process.Command) (process.Result, error) {
	if command.Name != "kubectl" {
		return process.Result{}, errors.New("unexpected executable " + command.Name)
	}

	var stdin string
	if command.Stdin != nil {
		data, err := io.ReadAll(command.Stdin)
		if err != nil {
			return process.Result{}, err
		}

		stdin = string(data)
	}

	recorded := call{args: command.Args, stdin: stdin, stdout: command.Stdout, stderr: command.Stderr}
	if deadline, bounded := ctx.Deadline(); bounded {
		recorded.deadline = time.Until(deadline)
	}

	recorder.mu.Lock()
	recorder.calls = append(recorder.calls, recorded)
	recorder.mu.Unlock()

	if recorder.respond == nil {
		return process.Result{ExitCode: 0}, nil
	}

	return recorder.respond(operation(command.Args))
}

// operation strips the explicit target flags every call starts with.
func operation(arguments []string) []string {
	for index := 0; index < len(arguments); index++ {
		switch arguments[index] {
		case "--kubeconfig", "--context", "--namespace":
			index++
		case "--request-timeout=15s":
		default:
			return arguments[index:]
		}
	}

	return nil
}

func answer(stdout string) func([]string) (process.Result, error) {
	return func([]string) (process.Result, error) { return process.Result{Stdout: stdout}, nil }
}

var target = application.Target{
	ClusterAccess: application.ClusterAccess{Kubeconfig: "/workspace/.exphere/state/cluster.kubeconfig", Context: "kind-ex-0123456789ab"},
	Namespace:     "demo-local",
}

func targetFlags(namespace string) []string {
	flags := []string{"--kubeconfig", target.Kubeconfig, "--context", target.Context, "--request-timeout=15s"}
	if namespace != "" {
		flags = append(flags, "--namespace", namespace)
	}

	return flags
}

func TestClient_InvokesKubectlAgainstTheExplicitTarget(t *testing.T) {
	t.Parallel()

	runner := &recorder{respond: answer(`{"items": []}`)}
	if _, err := (kubectl.Client{Runner: runner}).List(context.Background(), target, "pods"); err != nil {
		t.Fatal(err)
	}

	if want := targetFlags("demo-local"); !slices.Equal(runner.calls[0].args[:len(want)], want) {
		t.Errorf("args = %v, want the explicit kubeconfig, context, request timeout, and namespace first", runner.calls[0].args)
	}

	if deadline := runner.calls[0].deadline; deadline <= 0 || deadline > 30*time.Second {
		t.Errorf("call deadline = %v, want the 30-second call bound", deadline)
	}

	failing := &recorder{respond: func([]string) (process.Result, error) {
		return process.Result{Stderr: "forbidden: user cannot list pods"}, errors.New("exit status 1")
	}}

	_, err := kubectl.Client{Runner: failing}.List(context.Background(), target, "pods")
	if err == nil || err.Error() != "kubectl get: exit status 1: forbidden: user cannot list pods" {
		t.Fatalf("List() error = %v, want the operation, the cause, and the tool's diagnostics", err)
	}
}
