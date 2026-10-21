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

// stream records a foreground operation that returns no result.
type stream[Request any] struct {
	err       error
	requests  []Request
	deadlines []time.Duration
}

func (fake *stream[Request]) Execute(ctx context.Context, request Request) error {
	fake.requests = append(fake.requests, request)
	fake.deadlines = append(fake.deadlines, remaining(ctx))

	if logs, isLogs := any(request).(app.ReadRuntimeLogsRequest); isLogs {
		fmt.Fprint(logs.Output, "2026-10-17T12:00:00Z ready\n")
	}

	return fake.err
}

// remaining is how long ctx has left, or zero when it has no deadline.
func remaining(ctx context.Context) time.Duration {
	if deadline, bounded := ctx.Deadline(); bounded {
		return time.Until(deadline)
	}

	return 0
}

type harness struct {
	projectInit      operation[app.InitializeProjectRequest, app.ProjectResult]
	doctor           operation[app.DiagnoseRequest, app.Diagnosis]
	runtimeInit      operation[app.InitializeRuntimeRequest, app.RuntimeInitResult]
	runtimeValidate  operation[app.RuntimeRequest, app.ValidationResult]
	runtimeUp        operation[app.RuntimeRequest, app.ExecutionOutcome]
	runtimeDown      operation[app.StopRuntimeRequest, app.ExecutionOutcome]
	runtimeInspect   operation[app.RuntimeRequest, app.Inspection]
	runtimeLogs      stream[app.ReadRuntimeLogsRequest]
	runtimeForward   stream[app.ForwardRuntimeRequest]
	runtimeHook      operation[app.RunRuntimeHookRequest, app.HookResult]
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
		RuntimeInit: &h.runtimeInit, RuntimeValidate: &h.runtimeValidate, RuntimeUp: &h.runtimeUp,
		RuntimeDown: &h.runtimeDown, RuntimeInspect: &h.runtimeInspect, RuntimeLogs: &h.runtimeLogs,
		RuntimeForward: &h.runtimeForward, RuntimeHook: &h.runtimeHook,
		ClusterCreate: &h.clusterCreate, ClusterDelete: &h.clusterDelete, ClusterInspect: &h.clusterInspect,
		ClusterLoadImage: &h.clusterLoadImage, AgentInit: &h.agentInit,
		ProjectAgent: &h.projectAgent,
	}
}

func (h *harness) execute(ctx context.Context, arguments ...string) (int, string, string) {
	var stdout, stderr bytes.Buffer

	status := executeCatalogs(ctx, arguments, support.Options{
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

	runtime := app.RuntimeRequest{Project: discovered, Runtime: "demo"}

	tests := []struct {
		arguments []string
		request   func(*harness) any
		want      any
	}{
		{[]string{"init"}, func(h *harness) any { return h.projectInit.requests }, []app.InitializeProjectRequest{{Directory: workingDirectory}}},
		{[]string{"doctor"}, func(h *harness) any { return h.doctor.requests }, []app.DiagnoseRequest{{Project: discovered}}},
		{
			[]string{"runtime", "init", "demo"},
			func(h *harness) any { return h.runtimeInit.requests },
			[]app.InitializeRuntimeRequest{{Project: discovered, Runtime: "demo"}},
		},
		{[]string{"runtime", "validate", "demo"}, func(h *harness) any { return h.runtimeValidate.requests }, []app.RuntimeRequest{runtime}},
		{[]string{"runtime", "up", "demo"}, func(h *harness) any { return h.runtimeUp.requests }, []app.RuntimeRequest{runtime}},
		{[]string{"runtime", "down", "demo"}, func(h *harness) any { return h.runtimeDown.requests }, []app.StopRuntimeRequest{{RuntimeRequest: runtime}}},
		{
			[]string{"runtime", "down", "demo", "--recover"},
			func(h *harness) any { return h.runtimeDown.requests },
			[]app.StopRuntimeRequest{{RuntimeRequest: runtime, Recover: true}},
		},
		{[]string{"runtime", "inspect", "demo"}, func(h *harness) any { return h.runtimeInspect.requests }, []app.RuntimeRequest{runtime}},
		{
			[]string{"runtime", "hook", "demo", "seed"},
			func(h *harness) any { return h.runtimeHook.requests },
			[]app.RunRuntimeHookRequest{{RuntimeRequest: runtime, Name: "seed"}},
		},
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
			[]string{"--project", "/elsewhere/project", "runtime", "up", "demo"},
			func(h *harness) any { return h.runtimeUp.requests },
			[]app.RuntimeRequest{{Project: app.ProjectSelection{Directory: "/elsewhere/project", Explicit: true}, Runtime: "demo"}},
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
	h.runtimeInspect.result = app.Inspection{ObservationID: "obs-1"}

	_, indented, _ := h.run("runtime", "inspect", "demo")
	if !strings.Contains(indented, "\"observationId\": \"obs-1\"") || !strings.HasSuffix(indented, "}\n") {
		t.Errorf("default output = %q, want indented JSON", indented)
	}

	_, compact, _ := h.run("--json", "runtime", "inspect", "demo")
	if strings.Count(compact, "\n") != 1 || !strings.Contains(compact, `"observationId":"obs-1"`) {
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
		{name: "command group without a subcommand", arguments: []string{"runtime"}, wantStatus: 2, wantStderr: "ex runtime requires a subcommand"},
		{name: "unknown subcommand", arguments: []string{"cluster", "resize"}, wantStatus: 2, wantStderr: "unknown command"},
		{name: "unknown root command", arguments: []string{"deploy"}, wantStatus: 2, wantStderr: `unknown command "deploy"`},
		{name: "missing argument", arguments: []string{"runtime", "up"}, wantStatus: 2, wantStderr: "ex runtime up requires 1 argument(s)"},
		{name: "extra argument", arguments: []string{"cluster", "create", "now"}, wantStatus: 2, wantStderr: "requires 0 argument(s)"},
		{name: "hook without a name", arguments: []string{"runtime", "hook", "demo"}, wantStatus: 2, wantStderr: "requires 2 argument(s)"},
		{name: "unknown flag", arguments: []string{"runtime", "up", "demo", "--force"}, wantStatus: 2, wantStderr: "unknown flag: --force"},
		{name: "zero timeout", arguments: []string{"--timeout", "0s", "runtime", "up", "demo"}, wantStatus: 2, wantStderr: "--timeout must be positive"},
		{name: "invalid timeout", arguments: []string{"--timeout", "soon", "runtime", "up", "demo"}, wantStatus: 2, wantStderr: "invalid argument"},
		{name: "logs without a unit", arguments: []string{"runtime", "logs", "demo"}, wantStatus: 2, wantStderr: "--unit is required"},
		{name: "logs as JSON", arguments: []string{"runtime", "logs", "demo", "--unit", "web", "--json"}, wantStatus: 2, wantStderr: "--json is unsupported"},
		{name: "forward without a unit", arguments: []string{"runtime", "forward", "demo"}, wantStatus: 2, wantStderr: "--unit is required"},
		{name: "unknown help topic", arguments: []string{"help", "deploy"}, wantStatus: 2, wantStderr: "unknown help command"},
		{name: "agent connect without URL", arguments: []string{"agent", "connect"}, wantStatus: 2, wantStderr: "ex agent connect requires 1 argument(s)"},
		{name: "agent migrate without runtime", arguments: []string{"agent", "migrate"}, wantStatus: 2, wantStderr: "ex agent migrate requires 1 argument(s)"},
		{name: "operational failure", arguments: []string{"runtime", "up", "demo"}, failure: errors.New("cluster is not ready"), wantStatus: 1, wantStderr: "ex: cluster is not ready\n"},
		{name: "adoption refused", arguments: []string{"runtime", "up", "demo"}, failure: fmt.Errorf("namespace exists: %w", app.ErrAdoptionRefused), wantStatus: 3, wantStderr: "refusing to adopt"},
		{name: "target mismatch refused", arguments: []string{"runtime", "up", "demo"}, failure: fmt.Errorf("guard: %w", app.ErrTargetMismatch), wantStatus: 3, wantStderr: "target differs"},
		{name: "build recipe refused", arguments: []string{"runtime", "up", "demo"}, failure: fmt.Errorf("component api: %w", connection.ErrRecipeRefused), wantStatus: 3, wantStderr: "refused by policy"},
		{name: "console access revoked", arguments: []string{"runtime", "up", "demo"}, failure: fmt.Errorf("sync: %w", connectedapp.ErrAccessRevoked), wantStatus: 3, wantStderr: "revoked"},
		{name: "replacement refused", arguments: []string{"runtime", "up", "demo"}, failure: fmt.Errorf("record: %w", app.ErrAlreadyExists), wantStatus: 3, wantStderr: "already exists"},
		{name: "missing prerequisite record", arguments: []string{"runtime", "up", "demo"}, failure: app.NotFound("no cluster for this project"), wantStatus: 4, wantStderr: "no cluster"},
		{name: "missing tool", arguments: []string{"runtime", "up", "demo"}, failure: fmt.Errorf("kind failed: %w", exec.ErrNotFound), wantStatus: 4, wantStderr: "executable file not found"},
		{name: "interruption", arguments: []string{"runtime", "up", "demo"}, failure: fmt.Errorf("await: %w", context.Canceled), wantStatus: 130},
		{name: "deadline", arguments: []string{"runtime", "up", "demo"}, failure: context.DeadlineExceeded, wantStatus: 1, wantStderr: "deadline exceeded"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := &harness{}
			h.runtimeUp.err = tt.failure

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
	h.runtimeUp.err = errors.New("cluster is not ready")

	status, _, stderr := h.run("--json", "runtime", "up", "demo")

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
	h.runtimeUp.err = errors.New("stopped")

	if status, _, _ := h.execute(ctx, "runtime", "up", "demo"); status != 130 {
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
		{name: "finite operation", arguments: []string{"--timeout", "5s", "runtime", "up", "demo"}, deadline: func(h *harness) time.Duration { return h.runtimeUp.deadlines[0] }, bounded: true},
		{name: "unguided runtime init", arguments: []string{"--timeout", "5s", "runtime", "init", "demo"}, deadline: func(h *harness) time.Duration { return h.runtimeInit.deadlines[0] }, bounded: true},
		{name: "bounded log read", arguments: []string{"--timeout", "5s", "runtime", "logs", "demo", "--unit", "web"}, deadline: func(h *harness) time.Duration { return h.runtimeLogs.deadlines[0] }, bounded: true},
		{name: "followed logs", arguments: []string{"--timeout", "5s", "runtime", "logs", "demo", "--unit", "web", "--follow"}, deadline: func(h *harness) time.Duration { return h.runtimeLogs.deadlines[0] }},
		{name: "port forwarding", arguments: []string{"--timeout", "5s", "runtime", "forward", "demo", "--unit", "web"}, deadline: func(h *harness) time.Duration { return h.runtimeForward.deadlines[0] }},
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

func TestExecute_StreamsLogsAsNativeText(t *testing.T) {
	t.Parallel()

	h := &harness{}
	status, stdout, _ := h.run("runtime", "logs", "demo", "--unit", "web", "--follow")

	if status != 0 || stdout != "2026-10-17T12:00:00Z ready\n" {
		t.Fatalf("status %d, stdout %q; want the native log text", status, stdout)
	}

	request := h.runtimeLogs.requests[0]
	if request.Unit != "web" || !request.Follow || request.Runtime != "demo" || request.Diagnostics == nil {
		t.Fatalf("request = %+v", request)
	}
}

func TestExecute_SelectsForwardingPorts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		arguments []string
		local     int
		remote    int
	}{
		{arguments: []string{"runtime", "forward", "demo", "--unit", "web"}, local: 8080},
		{arguments: []string{"runtime", "forward", "demo", "--unit", "web", "--port", "9090", "--remote-port", "8081"}, local: 9090, remote: 8081},
	}

	for _, tt := range tests {
		h := &harness{}
		if status, _, stderr := h.run(tt.arguments...); status != 0 {
			t.Fatalf("status %d: %s", status, stderr)
		}

		if request := h.runtimeForward.requests[0]; request.LocalPort != tt.local || request.RemotePort != tt.remote || request.Unit != "web" {
			t.Errorf("%v: request = %+v, want ports %d:%d", tt.arguments, request, tt.local, tt.remote)
		}
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

// guidedInit runs the guide it receives, as runtime initialization does, and
// records what the guide collected.
type guidedInit struct {
	input    app.BlueprintInput
	deadline time.Duration
}

func (collector *guidedInit) Execute(ctx context.Context, request app.InitializeRuntimeRequest) (app.RuntimeInitResult, error) {
	collector.deadline = remaining(ctx)

	if request.Guide == nil {
		return app.RuntimeInitResult{}, errors.New("no guide")
	}

	input, err := request.Guide.Blueprint(ctx)
	collector.input = input

	return app.RuntimeInitResult{Runtime: request.Runtime, Guided: true}, err
}

func runGuided(answers string, arguments ...string) (*guidedInit, int, string, string) {
	collector := &guidedInit{}
	h := &harness{input: answers}
	useCases := h.useCases()
	useCases.RuntimeInit = collector

	var stdout, stderr bytes.Buffer
	status := executeCatalogs(context.Background(), arguments, support.Options{
		Directory: workingDirectory, Input: strings.NewReader(answers), Output: &stdout, Diagnostics: &stderr,
	}, useCases)

	return collector, status, stdout.String(), stderr.String()
}

func TestExecute_GuidedRuntimeInitAsksOnDiagnostics(t *testing.T) {
	t.Parallel()

	answers := strings.Join([]string{
		"team-payments", "kind-ex-0123456789ab", "demo-local", "services/web, jobs/batch,",
		"Deployment", "web:dev", "100m", "500m", "128Mi", "256Mi", "8080", "2", "/healthz",
		"Job", "batch:dev", "100m", "200m", "64Mi", "128Mi",
		"1", "2", "1Gi", "2Gi", "10",
	}, "\n") + "\n"

	collector, status, stdout, stderr := runGuided(answers, "--timeout", "1s", "runtime", "init", "demo", "--guided")
	if status != 0 {
		t.Fatalf("status %d: %s", status, stderr)
	}

	if collector.deadline != 0 {
		t.Errorf("guided questions had deadline %v, want --timeout not to bound them", collector.deadline)
	}

	input := collector.input
	if input.Owner != "team-payments" || input.Namespace != "demo-local" || len(input.Workloads) != 2 ||
		input.Workloads[0].Component != "services/web" || input.Workloads[0].Port != 8080 || input.Workloads[0].Replicas != 2 ||
		input.Workloads[1].Kind != "Job" || input.Workloads[1].Port != 0 || input.Quotas["pods"] != "10" || len(input.Quotas) != 5 {
		t.Errorf("collected input = %+v", input)
	}

	if !strings.Contains(stderr, "Owner reference: ") || !strings.Contains(stderr, "Configure jobs/batch\n") {
		t.Errorf("stderr = %q, want every question on diagnostics", stderr)
	}

	if !json.Valid([]byte(stdout)) || strings.Contains(stdout, "Owner") {
		t.Errorf("stdout = %q, want only the JSON result", stdout)
	}
}

func TestExecute_GuidedRuntimeInitRefusesInvalidAnswers(t *testing.T) {
	t.Parallel()

	prefix := "team-payments\nkind-ex\ndemo-local\nservices/web\n"

	tests := []struct {
		name       string
		answers    string
		wantStderr string
	}{
		{name: "empty answer", answers: "team-payments\n\n", wantStderr: "Kubernetes context of the project cluster is required"},
		{name: "input ends early", answers: "team-payments\n", wantStderr: "is required; guidance canceled without writing files"},
		{name: "unsupported workload kind", answers: prefix + "StatefulSet\n", wantStderr: "workload kind must be Deployment or Job"},
		{
			name:       "port that is not a number",
			answers:    prefix + "Deployment\nweb:dev\n100m\n500m\n128Mi\n256Mi\neighty\n",
			wantStderr: "Application port must be a whole number",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, status, stdout, stderr := runGuided(tt.answers, "runtime", "init", "demo", "--guided")

			if status != 2 || !strings.Contains(stderr, tt.wantStderr) || stdout != "" {
				t.Fatalf("status %d, stderr %q; want usage error %q", status, stderr, tt.wantStderr)
			}
		})
	}
}

// BenchmarkExecute measures building the command tree and dispatching one
// command in process. Process startup is outside this measurement.
func BenchmarkExecute(b *testing.B) {
	for _, arguments := range [][]string{{"version"}, {"--help"}} {
		b.Run(strings.Join(arguments, " "), func(b *testing.B) {
			options := support.Options{Version: "bench", Directory: workingDirectory, Output: io.Discard, Diagnostics: io.Discard}

			b.ReportAllocs()

			for range b.N {
				if status := executeCatalogs(context.Background(), arguments, options, UseCases{}); status != 0 {
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
