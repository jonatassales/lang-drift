package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"os/exec"

	app "github.com/jonatassales/tao/packages/ex/src/application"
	connectedapp "github.com/jonatassales/tao/packages/ex/src/application/connected"
	"github.com/jonatassales/tao/packages/ex/src/application/setup"
	"github.com/jonatassales/tao/packages/ex/src/domain/connection"
	"github.com/jonatassales/tao/packages/ex/src/presentation/cli/support"
)

const workingDirectory = "/workspace/project"

// fakeProjectAgent mounts the agent commands; the failure tests never reach it.
type fakeProjectAgent struct{}

func (*fakeProjectAgent) Connect(context.Context, connectedapp.ConnectRequest) (connectedapp.Result, error) {
	return connectedapp.Result{}, nil
}
func (*fakeProjectAgent) Run(context.Context, connectedapp.RunRequest) error { return nil }
func (*fakeProjectAgent) RunProject(context.Context, connectedapp.ProjectRunRequest) error {
	return nil
}
func (*fakeProjectAgent) Migrate(context.Context, connectedapp.MigrateRequest) (connectedapp.Result, error) {
	return connectedapp.Result{}, nil
}

// operation records the requests a finite use case receives and whether its
// context carried a deadline.
type operation[Request, Result any] struct {
	result    Result
	err       error
	requests  []Request
	deadlines []time.Duration
}

func (fake *operation[Request, Result]) Execute(ctx context.Context, request Request) (Result, error) {
	fake.requests = append(fake.requests, request)
	fake.deadlines = append(fake.deadlines, remaining(ctx))

	return fake.result, fake.err
}

// remaining is how long ctx has left, or zero when it has no deadline.
func remaining(ctx context.Context) time.Duration {
	if deadline, bounded := ctx.Deadline(); bounded {
		return time.Until(deadline)
	}

	return 0
}

type harness struct {
	projectInit      operation[setup.Request, setup.Result]
	doctor           operation[app.DiagnoseRequest, app.Diagnosis]
	clusterCreate    operation[app.ClusterRequest, app.ClusterResult]
	clusterDelete    operation[app.ClusterRequest, app.ClusterResult]
	clusterInspect   operation[app.ClusterRequest, app.ClusterResult]
	clusterLoadImage operation[app.LoadClusterImageRequest, app.ImageResult]
	agentInit        operation[app.InitializeAgentRequest, app.AgentResult]
	projectAgent     fakeProjectAgent
	input            string
}

func (h *harness) useCases() UseCases {
	return UseCases{
		ProjectInit: &h.projectInit, Doctor: &h.doctor,
		ClusterCreate: &h.clusterCreate, ClusterDelete: &h.clusterDelete, ClusterInspect: &h.clusterInspect,
		ClusterLoadImage: &h.clusterLoadImage, AgentInit: &h.agentInit,
		ProjectAgent: &h.projectAgent,
	}
}

func (h *harness) execute(ctx context.Context, arguments ...string) (int, string, string) {
	var stdout, stderr bytes.Buffer

	status := Execute(ctx, arguments, support.Options{
		Version:     "1.2.3",
		Directory:   workingDirectory,
		Input:       strings.NewReader(h.input),
		Output:      &stdout,
		Diagnostics: &stderr,
	}, h.useCases())

	return status, stdout.String(), stderr.String()
}

func (h *harness) run(arguments ...string) (int, string, string) {
	return h.execute(context.Background(), arguments...)
}

var discovered = app.ProjectSelection{Directory: workingDirectory}

func TestExecute_MapsEveryCommandToItsRequest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		arguments []string
		request   func(*harness) any
		want      any
	}{
		{[]string{"init"}, func(h *harness) any { return h.projectInit.requests }, []setup.Request{{
			Directory: workingDirectory,
			By:        "owner",
		}}},
		{[]string{"doctor"}, func(h *harness) any { return h.doctor.requests }, []app.DiagnoseRequest{{Project: discovered}}},
		{[]string{"cluster", "create"}, func(h *harness) any { return h.clusterCreate.requests }, []app.ClusterRequest{{Project: discovered}}},
		{[]string{"cluster", "delete"}, func(h *harness) any { return h.clusterDelete.requests }, []app.ClusterRequest{{Project: discovered}}},
		{[]string{"cluster", "inspect"}, func(h *harness) any { return h.clusterInspect.requests }, []app.ClusterRequest{{Project: discovered}}},
		{
			[]string{"cluster", "load-image", "web:dev"},
			func(h *harness) any { return h.clusterLoadImage.requests },
			[]app.LoadClusterImageRequest{{Project: discovered, Image: "web:dev"}},
		},
		{
			[]string{"agent", "init", "researcher"},
			func(h *harness) any { return h.agentInit.requests },
			[]app.InitializeAgentRequest{{Project: discovered, Name: "researcher"}},
		},
		{
			[]string{"--project", "/elsewhere/project", "cluster", "create"},
			func(h *harness) any { return h.clusterCreate.requests },
			[]app.ClusterRequest{{Project: app.ProjectSelection{Directory: "/elsewhere/project", Explicit: true}}},
		},
	}

	for _, tt := range tests {
		t.Run(strings.Join(tt.arguments, " "), func(t *testing.T) {
			t.Parallel()

			h := &harness{}
			status, stdout, stderr := h.run(tt.arguments...)

			if status != 0 || stderr != "" {
				t.Fatalf("status %d, stderr %q", status, stderr)
			}

			if got := tt.request(h); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("requests = %+v, want %+v", got, tt.want)
			}

			// doctor and cluster inspect render for people by default; every
			// other command emits one JSON document.
			if tt.arguments[0] == "doctor" || (tt.arguments[0] == "cluster" && tt.arguments[1] == "inspect") {
				if stdout == "" || strings.Contains(stdout, "{") {
					t.Errorf("stdout = %q, want human text", stdout)
				}
			} else if !json.Valid([]byte(stdout)) {
				t.Errorf("stdout = %q, want one JSON result", stdout)
			}
		})
	}
}

func TestExecute_ResolvesAnExplicitRelativeProject(t *testing.T) {
	t.Parallel()

	h := &harness{}
	if status, _, stderr := h.run("--project", "relative/project", "cluster", "inspect"); status != 0 {
		t.Fatalf("status %d: %s", status, stderr)
	}

	want, err := filepath.Abs("relative/project")
	if err != nil {
		t.Fatal(err)
	}

	if selection := h.clusterInspect.requests[0].Project; selection.Directory != want || !selection.Explicit {
		t.Fatalf("selection = %+v, want the absolute explicit project %s", selection, want)
	}
}

func TestExecute_RendersResults(t *testing.T) {
	t.Parallel()

	h := &harness{}
	h.clusterCreate.result = app.ClusterResult{Name: "ex-0123456789ab", Status: "ready"}

	_, indented, _ := h.run("cluster", "create")
	if !strings.Contains(indented, "\n  \"name\": \"ex-0123456789ab\",\n") || !strings.HasSuffix(indented, "}\n") {
		t.Errorf("default output = %q, want indented JSON", indented)
	}

	_, compact, _ := h.run("--json", "cluster", "create")
	if strings.Count(compact, "\n") != 1 || !strings.Contains(compact, `"name":"ex-0123456789ab"`) {
		t.Errorf("--json output = %q, want one compact line", compact)
	}
}

func TestExecute_ClassifiesFailures(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		arguments  []string
		failure    error
		wantStatus int
		wantStderr string
	}{
		{name: "retired runtime tree refuses with guidance", arguments: []string{"runtime", "up", "demo"}, wantStatus: 1, wantStderr: "runtime operations moved to the console"},
		{name: "command group without a subcommand", arguments: []string{"cluster"}, wantStatus: 2, wantStderr: "ex cluster requires a subcommand"},
		{name: "unknown subcommand", arguments: []string{"cluster", "resize"}, wantStatus: 2, wantStderr: "unknown command"},
		{name: "unknown root command", arguments: []string{"deploy"}, wantStatus: 2, wantStderr: `unknown command "deploy"`},
		{name: "missing argument", arguments: []string{"agent", "init"}, wantStatus: 2, wantStderr: "ex agent init requires 1 argument(s)"},
		{name: "extra argument", arguments: []string{"cluster", "create", "now"}, wantStatus: 2, wantStderr: "requires 0 argument(s)"},
		{name: "unknown flag", arguments: []string{"cluster", "create", "--force"}, wantStatus: 2, wantStderr: "unknown flag: --force"},
		{name: "zero timeout", arguments: []string{"--timeout", "0s", "cluster", "create"}, wantStatus: 2, wantStderr: "--timeout must be positive"},
		{name: "invalid timeout", arguments: []string{"--timeout", "soon", "cluster", "create"}, wantStatus: 2, wantStderr: "invalid argument"},
		{name: "unknown help topic", arguments: []string{"help", "deploy"}, wantStatus: 2, wantStderr: "unknown help command"},
		{name: "agent connect without URL", arguments: []string{"agent", "connect"}, wantStatus: 2, wantStderr: "ex agent connect requires 1 argument(s)"},
		{name: "agent migrate without runtime", arguments: []string{"agent", "migrate"}, wantStatus: 2, wantStderr: "ex agent migrate requires 1 argument(s)"},
		{name: "operational failure", arguments: []string{"cluster", "create"}, failure: errors.New("cluster is not ready"), wantStatus: 1, wantStderr: "ex: cluster is not ready\n"},
		{name: "adoption refused", arguments: []string{"cluster", "create"}, failure: fmt.Errorf("namespace exists: %w", app.ErrAdoptionRefused), wantStatus: 3, wantStderr: "refusing to adopt"},
		{name: "target mismatch refused", arguments: []string{"cluster", "create"}, failure: fmt.Errorf("guard: %w", app.ErrTargetMismatch), wantStatus: 3, wantStderr: "target differs"},
		{name: "build recipe refused", arguments: []string{"cluster", "create"}, failure: fmt.Errorf("component api: %w", connection.ErrRecipeRefused), wantStatus: 3, wantStderr: "refused by policy"},
		{name: "console access revoked", arguments: []string{"cluster", "create"}, failure: fmt.Errorf("sync: %w", connectedapp.ErrAccessRevoked), wantStatus: 3, wantStderr: "revoked"},
		{name: "replacement refused", arguments: []string{"cluster", "create"}, failure: fmt.Errorf("record: %w", app.ErrAlreadyExists), wantStatus: 3, wantStderr: "already exists"},
		{name: "missing prerequisite record", arguments: []string{"cluster", "create"}, failure: app.NotFound("no cluster for this project"), wantStatus: 4, wantStderr: "no cluster"},
		{name: "missing tool", arguments: []string{"cluster", "create"}, failure: fmt.Errorf("kind failed: %w", exec.ErrNotFound), wantStatus: 4, wantStderr: "executable file not found"},
		{name: "interruption", arguments: []string{"cluster", "create"}, failure: fmt.Errorf("await: %w", context.Canceled), wantStatus: 130},
		{name: "deadline", arguments: []string{"cluster", "create"}, failure: context.DeadlineExceeded, wantStatus: 1, wantStderr: "deadline exceeded"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := &harness{}
			h.clusterCreate.err = tt.failure

			status, stdout, stderr := h.run(tt.arguments...)

			if status != tt.wantStatus || !strings.Contains(stderr, tt.wantStderr) {
				t.Fatalf("status %d, stderr %q; want %d and %q", status, stderr, tt.wantStatus, tt.wantStderr)
			}

			if stdout != "" {
				t.Errorf("a failure wrote a result: %q", stdout)
			}
		})
	}
}

func TestExecute_ReportsStructuredDiagnostics(t *testing.T) {
	t.Parallel()

	h := &harness{}
	h.clusterCreate.err = errors.New("cluster is not ready")

	status, _, stderr := h.run("--json", "cluster", "create")

	var diagnostic struct {
		Error    string `json:"error"`
		ExitCode int    `json:"exitCode"`
	}

	if err := json.Unmarshal([]byte(stderr), &diagnostic); err != nil || status != 1 ||
		diagnostic.Error != "cluster is not ready" || diagnostic.ExitCode != 1 {
		t.Fatalf("status %d, stderr %q; want a JSON diagnostic with exit code 1", status, stderr)
	}
}

func TestExecute_ReportsAnInterruptedApplication(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	h := &harness{}
	h.clusterCreate.err = errors.New("stopped")

	if status, _, _ := h.execute(ctx, "cluster", "create"); status != 130 {
		t.Fatalf("status = %d, want 130 once the application context is canceled", status)
	}
}

func TestExecute_TimeoutBoundsOnlyFiniteOperations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		arguments []string
		deadline  func(*harness) time.Duration
		bounded   bool
	}{
		{name: "finite operation", arguments: []string{"--timeout", "5s", "cluster", "create"}, deadline: func(h *harness) time.Duration { return h.clusterCreate.deadlines[0] }, bounded: true},
		{name: "doctor", arguments: []string{"--timeout", "5s", "doctor"}, deadline: func(h *harness) time.Duration { return h.doctor.deadlines[0] }, bounded: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := &harness{}
			if status, _, stderr := h.run(tt.arguments...); status != 0 {
				t.Fatalf("status %d: %s", status, stderr)
			}

			deadline := tt.deadline(h)
			if tt.bounded != (deadline > 0) || deadline > 5*time.Second {
				t.Fatalf("deadline = %v, want bounded %t by --timeout", deadline, tt.bounded)
			}
		})
	}
}

func TestExecute_PrintsTheVersionWithoutAProject(t *testing.T) {
	t.Parallel()

	for _, arguments := range [][]string{{"version"}, {"--version"}, {"-V"}} {
		h := &harness{}
		if status, stdout, _ := h.run(arguments...); status != 0 || stdout != "ex 1.2.3\n" {
			t.Errorf("%v = %d, %q", arguments, status, stdout)
		}
	}

	for _, arguments := range [][]string{{"version", "--json"}, {"--json", "--version"}} {
		h := &harness{}
		if status, stdout, _ := h.run(arguments...); status != 0 || stdout != "{\"version\":\"1.2.3\"}\n" {
			t.Errorf("%v = %d, %q", arguments, status, stdout)
		}
	}
}

func BenchmarkExecute(b *testing.B) {
	for _, arguments := range [][]string{{"version"}, {"--help"}} {
		b.Run(strings.Join(arguments, " "), func(b *testing.B) {
			options := support.Options{Version: "bench", Directory: workingDirectory, Output: io.Discard, Diagnostics: io.Discard}

			b.ReportAllocs()

			for range b.N {
				if status := Execute(context.Background(), arguments, options, UseCases{}); status != 0 {
					b.Fatalf("Execute(%q) = %d, want 0", arguments, status)
				}
			}
		})
	}
}

func TestExecute_RendersDoctorForPeople(t *testing.T) {
	t.Parallel()

	h := &harness{}
	h.doctor.result = app.Diagnosis{
		Tools: []app.ToolStatus{
			{Name: "docker", Status: app.ToolAvailable, Version: "27.0", Path: "/usr/local/bin/docker"},
			{Name: "kind", Status: app.ToolMissing, Install: "https://kind.sigs.k8s.io"},
		},
		ContainerEngine: "reachable",
		ClusterNote:     "no verified project cluster: no cluster is recorded for this project",
		Installation:    "manual; no tools were installed",
	}

	status, stdout, _ := h.run("doctor")
	if status != 0 {
		t.Fatalf("status %d", status)
	}
	for _, want := range []string{
		"ok       docker 27.0 (/usr/local/bin/docker)",
		"missing  kind — install: https://kind.sigs.k8s.io",
		"Container engine: reachable",
		"Cluster: none — no verified project cluster: no cluster is recorded for this project",
		"Installation: manual; no tools were installed",
	} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("doctor output %q lacks %q", stdout, want)
		}
	}
	if strings.Contains(stdout, "{") {
		t.Fatalf("doctor output still renders JSON: %q", stdout)
	}

	h = &harness{}
	h.doctor.result = app.Diagnosis{Installation: "manual; no tools were installed"}
	if _, stdout, _ = h.run("doctor", "--json"); !strings.HasPrefix(stdout, `{"tools":null`) {
		t.Fatalf("doctor --json = %q, want the compact document", stdout)
	}
}

func TestExecute_RendersClusterInspectionForPeople(t *testing.T) {
	t.Parallel()

	h := &harness{}
	h.clusterInspect.result = app.ClusterResult{
		Name:             "kind-demo",
		Context:          "kind-demo",
		Status:           "ready",
		UID:              "uid-123",
		Kubeconfig:       "/private/kubeconfig",
		NamespaceUIDs:    map[string]string{"demo-local": "ns-1"},
		InventoryVersion: 7,
		ResourceTypes:    []string{"pods", "services"},
	}

	status, stdout, _ := h.run("cluster", "inspect")
	if status != 0 {
		t.Fatalf("status %d", status)
	}
	for _, want := range []string{
		"Cluster kind-demo",
		"Context:    kind-demo",
		"Status:     ready",
		"UID:        uid-123",
		"Kubeconfig: /private/kubeconfig",
		"Namespaces: 1",
		"Inventory:  version 7, 2 resource types",
	} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("inspect output %q lacks %q", stdout, want)
		}
	}
}

func TestExecute_ErrorsCarryAHint(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		failure  error
		wantHint string
	}{
		{"refusal points at inspection", fmt.Errorf("x: %w", app.ErrAdoptionRefused), "refused to change what it does not own"},
		{"missing prerequisite points at doctor", app.NotFound("no cluster"), "run ex doctor"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := &harness{}
			h.clusterInspect.err = tt.failure
			_, _, stderr := h.run("cluster", "inspect")
			if !strings.Contains(stderr, "hint: ") || !strings.Contains(stderr, tt.wantHint) {
				t.Fatalf("stderr %q lacks hint %q", stderr, tt.wantHint)
			}
		})
	}

	h := &harness{}
	h.clusterInspect.err = app.NotFound("no cluster")
	_, _, stderr := h.run("--json", "cluster", "inspect")
	if strings.Contains(stderr, "hint") {
		t.Fatalf("JSON diagnostics must stay structured: %q", stderr)
	}
}
