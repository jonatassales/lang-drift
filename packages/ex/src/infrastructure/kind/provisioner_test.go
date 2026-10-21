package kind_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/exphere/ex/src/infrastructure/kind"
	"github.com/exphere/ex/src/infrastructure/process"
)

// recorder records each invocation as its executable followed by its
// arguments, and answers with respond.
type recorder struct {
	calls   [][]string
	respond func(invocation []string) (process.Result, error)
}

func (recorder *recorder) Run(_ context.Context, command process.Command) (process.Result, error) {
	invocation := append([]string{command.Name}, command.Args...)
	recorder.calls = append(recorder.calls, invocation)

	if recorder.respond == nil {
		return process.Result{}, nil
	}

	return recorder.respond(invocation)
}

func fails(stderr string) func([]string) (process.Result, error) {
	return func([]string) (process.Result, error) {
		return process.Result{Stderr: stderr, ExitCode: 1}, errors.New("exit status 1")
	}
}

func TestProvisioner_Clusters(t *testing.T) {
	t.Parallel()

	runner := &recorder{respond: func([]string) (process.Result, error) { return process.Result{Stdout: "ex-a\nex-b\n"}, nil }}

	clusters, err := kind.Provisioner{Runner: runner}.Clusters(context.Background())
	if err != nil || !slices.Equal(clusters, []string{"ex-a", "ex-b"}) || !slices.Equal(runner.calls[0], []string{"kind", "get", "clusters"}) {
		t.Fatalf("Clusters() = %v, %v with calls %v", clusters, err, runner.calls)
	}

	if _, err := (kind.Provisioner{Runner: &recorder{respond: fails("docker is not running")}}).Clusters(context.Background()); err == nil ||
		!strings.Contains(err.Error(), "list kind clusters") || !strings.Contains(err.Error(), "docker is not running") {
		t.Fatalf("Clusters() error = %v, want the action and the tool's diagnostics", err)
	}
}

func TestProvisioner_Create(t *testing.T) {
	t.Parallel()

	kubeconfig := filepath.Join(t.TempDir(), "cluster.kubeconfig")
	if err := os.WriteFile(kubeconfig, []byte("apiVersion: v1\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	runner := &recorder{}

	contextName, err := kind.Provisioner{Runner: runner}.Create(context.Background(), "ex-0123456789ab", kubeconfig)
	if err != nil || contextName != "kind-ex-0123456789ab" {
		t.Fatalf("Create() = %q, %v", contextName, err)
	}

	want := []string{"kind", "create", "cluster", "--name", "ex-0123456789ab", "--kubeconfig", kubeconfig, "--wait", "120s"}
	if !slices.Equal(runner.calls[0], want) {
		t.Errorf("call = %v, want %v", runner.calls[0], want)
	}

	if info, err := os.Stat(kubeconfig); err != nil || info.Mode().Perm() != 0o600 {
		t.Errorf("kubeconfig mode = %v, %v; want private credentials", info.Mode().Perm(), err)
	}

	if _, err := (kind.Provisioner{Runner: &recorder{respond: fails("node failed")}}).Create(context.Background(), "ex-0123456789ab", kubeconfig); err == nil ||
		!strings.Contains(err.Error(), "create kind cluster") {
		t.Fatalf("Create() error = %v, want the failure explained", err)
	}
}

func TestProvisioner_Delete(t *testing.T) {
	t.Parallel()

	runner := &recorder{}
	if err := (kind.Provisioner{Runner: runner}).Delete(context.Background(), "ex-0123456789ab"); err != nil {
		t.Fatal(err)
	}

	if want := []string{"kind", "delete", "cluster", "--name", "ex-0123456789ab"}; !slices.Equal(runner.calls[0], want) {
		t.Fatalf("call = %v, want %v", runner.calls[0], want)
	}
}

func TestProvisioner_LoadImage(t *testing.T) {
	t.Parallel()

	runner := &recorder{}
	if err := (kind.Provisioner{Runner: runner}).LoadImage(context.Background(), "ex-a", "web:dev"); err != nil {
		t.Fatal(err)
	}

	if want := []string{"kind", "load", "docker-image", "web:dev", "--name", "ex-a"}; !slices.Equal(runner.calls[0], want) {
		t.Fatalf("call = %v, want %v", runner.calls[0], want)
	}

	refused := &recorder{}
	if err := (kind.Provisioner{Runner: refused}).LoadImage(context.Background(), "ex-a", "--help"); err == nil || len(refused.calls) != 0 {
		t.Fatal("LoadImage() passed an option-like image to kind")
	}
}

func TestProvisioner_RequireImages(t *testing.T) {
	t.Parallel()

	nodes := func(invocation []string) (process.Result, error) {
		if invocation[0] == "kind" {
			return process.Result{Stdout: "ex-a-control-plane\nex-a-worker\n"}, nil
		}

		return process.Result{}, nil
	}

	runner := &recorder{respond: nodes}
	if err := (kind.Provisioner{Runner: runner}).RequireImages(context.Background(), "ex-a", []string{"batch:dev", "web:dev"}); err != nil {
		t.Fatal(err)
	}

	want := [][]string{
		{"kind", "get", "nodes", "--name", "ex-a"},
		{"docker", "exec", "ex-a-control-plane", "crictl", "inspecti", "batch:dev"},
		{"docker", "exec", "ex-a-worker", "crictl", "inspecti", "batch:dev"},
		{"docker", "exec", "ex-a-control-plane", "crictl", "inspecti", "web:dev"},
		{"docker", "exec", "ex-a-worker", "crictl", "inspecti", "web:dev"},
	}
	if !slices.EqualFunc(runner.calls, want, slices.Equal) {
		t.Fatalf("calls = %v, want every image probed on every node", runner.calls)
	}

	missing := &recorder{respond: func(invocation []string) (process.Result, error) {
		if invocation[0] == "docker" && invocation[2] == "ex-a-worker" {
			return process.Result{Stderr: "image not found"}, errors.New("exit status 1")
		}

		return nodes(invocation)
	}}

	err := kind.Provisioner{Runner: missing}.RequireImages(context.Background(), "ex-a", []string{"web:dev"})
	if err == nil || !strings.Contains(err.Error(), `image "web:dev" is not loaded on node ex-a-worker`) {
		t.Fatalf("RequireImages() error = %v, want the missing image and node named", err)
	}

	noNodes := &recorder{respond: func([]string) (process.Result, error) { return process.Result{Stdout: "\n"}, nil }}
	if err := (kind.Provisioner{Runner: noNodes}).RequireImages(context.Background(), "ex-a", []string{"web:dev"}); err == nil {
		t.Fatal("RequireImages() accepted a cluster without nodes")
	}

	option := &recorder{respond: nodes}
	if err := (kind.Provisioner{Runner: option}).RequireImages(context.Background(), "ex-a", []string{"-v"}); err == nil || len(option.calls) != 1 {
		t.Fatal("RequireImages() probed an option-like image")
	}
}
