package connected

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/jonatassales/tao/exphere/cli/src/application"
	app "github.com/jonatassales/tao/exphere/cli/src/application/connected"
	"github.com/jonatassales/tao/exphere/cli/src/domain/connection"
	"github.com/jonatassales/tao/exphere/cli/src/domain/project"
	"github.com/jonatassales/tao/exphere/cli/src/infrastructure/process"
	"github.com/jonatassales/tao/exphere/cli/src/infrastructure/storage"
)

type fakeRunner struct {
	commands    []process.Command
	failBuild   bool
	buildStderr string
}

func (r *fakeRunner) Run(_ context.Context, c process.Command) (process.Result, error) {
	r.commands = append(r.commands, c)
	if c.Name == "git" {
		if c.Args[0] == "status" {
			return process.Result{}, nil
		}
		if c.Args[0] == "ls-files" {
			return process.Result{Stdout: "build.sh\x00"}, nil
		}
		if c.Args[0] == "config" {
			return process.Result{Stdout: "git@github.com:owner/repository.git\n"}, nil
		}
		return process.Result{Stdout: strings.Repeat("a", 40)}, nil
	}
	if r.failBuild {
		if c.Stderr != nil {
			_, _ = c.Stderr.Write([]byte(r.buildStderr))
		}
		return process.Result{Stderr: r.buildStderr}, errors.New("build error")
	}
	return process.Result{}, nil
}
func ownerForTest(t *testing.T) project.Project {
	t.Helper()
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	id, _ := project.ParseID(strings.Repeat("a", 32))
	owner, err := (storage.Projects{}).Create(root, id)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "build.sh"), []byte("#!/bin/sh\nexit 1\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	return owner
}
func nativeConfig() connection.Configuration {
	return connection.Configuration{Version: 1, RuntimeUID: "runtime-uid", RuntimeID: "runtime", ProjectID: "project", Revision: 1, Name: "Example", CPUMillis: 1200, MemoryMiB: 512, Model: "test-model", Effort: "low", Components: []connection.Component{{Path: ".", Image: "api:local", Port: 8080, ReadinessPath: "/health", Build: &connection.Build{Executable: "./build.sh"}}}}
}
func TestNative_MaterializationPreservesCustomizationAndBounds(t *testing.T) {
	owner := ownerForTest(t)
	native := Native{Runner: &fakeRunner{}, Executions: storage.Executions{}}
	config := nativeConfig()
	record := app.Connection{RuntimeUID: config.RuntimeUID, RuntimeID: config.RuntimeID, Repository: app.Repository{FullName: "owner/repository"}}
	managed, hash, err := native.Materialize(context.Background(), owner, record, config, config.Markdown())
	if err != nil {
		t.Fatal(err)
	}
	if hash == "" {
		t.Fatal("missing concrete hash")
	}
	record.Managed = managed
	record.ConcreteHash = hash
	snapshot, err := native.snapshot(owner, config)
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshot.Manifests) != 2 {
		t.Fatalf("objects=%d", len(snapshot.Manifests))
	}
	path := filepath.Join(owner.Root, ".langdrift/runtimes/runtime/runtime.md")
	custom := []byte("custom authored content\n")
	if err = os.WriteFile(path, custom, 0o644); err != nil {
		t.Fatal(err)
	}
	config.Revision++
	if _, _, err = native.Materialize(context.Background(), owner, record, config, config.Markdown()); err == nil {
		t.Fatal("local edit was overwritten")
	}
	content, _ := os.ReadFile(path)
	if string(content) != string(custom) {
		t.Fatal("custom file changed")
	}
}
func TestNative_BuildFailureNeverStartsInfrastructure(t *testing.T) {
	owner := ownerForTest(t)
	runner := &fakeRunner{failBuild: true}
	native := Native{Runner: runner, Executions: storage.Executions{}}
	config := nativeConfig()
	record := app.Connection{RuntimeID: config.RuntimeID, Repository: app.Repository{FullName: "owner/repository"}}
	record.Managed, record.ConcreteHash, _ = native.Materialize(context.Background(), owner, record, config, config.Markdown())
	var kinds []string
	err := native.Start(context.Background(), owner, record, config, func(kind, detail string) error { kinds = append(kinds, kind); return nil })
	if err == nil || len(kinds) != 2 || kinds[1] != "build_failed" {
		t.Fatalf("err=%v kinds=%v", err, kinds)
	}
	// Infrastructure ports are deliberately nil: calling any would panic.
}
func TestStore_SecretsArePrivateAndNotOperationalJSON(t *testing.T) {
	owner := ownerForTest(t)
	store := Store{Runner: untrackedRunner{}}
	record := app.Connection{RuntimeUID: "runtime-uid", RuntimeID: "runtime", State: "stopped", Secret: "runtime-secret", Credential: "agent-secret"}
	if err := store.Bind(context.Background(), owner, record); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(owner.Root, ".langdrift/state/connections/runtime.json"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "secret") {
		t.Fatal("credentials leaked into state")
	}
	info, _ := os.Stat(filepath.Join(owner.Root, ".langdrift/runtimes/runtime/.env"))
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("secret mode=%v", info.Mode())
	}
	records, err := store.List(owner)
	if err != nil || len(records) != 1 || records[0].Credential != "agent-secret" {
		t.Fatalf("load: %+v %v", records, err)
	}
	if err := store.Bind(context.Background(), owner, record); err == nil {
		t.Fatal("existing .env replaced")
	}
}
func TestNative_MigrationKeepsOriginalFilesAndAppliesConsoleEnvelope(t *testing.T) {
	owner := ownerForTest(t)
	native := Native{Runner: &fakeRunner{}, Executions: storage.Executions{}}
	config := nativeConfig()
	record := app.Connection{RuntimeID: config.RuntimeID, Repository: app.Repository{FullName: "owner/repository"}}
	_, _, err := native.Materialize(context.Background(), owner, record, config, config.Markdown())
	if err != nil {
		t.Fatal(err)
	}
	// Treat the generated v2 tree as an existing independently authored runtime.
	if err = os.Remove(filepath.Join(owner.Root, ".langdrift/runtimes/runtime/runtime.md")); err != nil {
		t.Fatal(err)
	}
	original, err := os.ReadFile(filepath.Join(owner.Root, ".langdrift/runtimes/runtime/schema.json"))
	if err != nil {
		t.Fatal(err)
	}
	config.Native = &connection.Native{RuntimeID: "runtime"}
	config.CPUMillis = 2000
	managed, _, err := native.Materialize(context.Background(), owner, record, config, config.Markdown())
	if err != nil {
		t.Fatal(err)
	}
	archived, err := os.ReadFile(filepath.Join(owner.Root, ".langdrift/state/migrations/runtime/schema.json"))
	if err != nil || string(archived) != string(original) {
		t.Fatal("original descriptor not preserved")
	}
	if _, err = native.snapshot(owner, config); err != nil {
		t.Fatal(err)
	}
	record.Managed = managed
	config.CPUMillis = 3000
	config.Revision++
	if _, _, err = native.Materialize(context.Background(), owner, record, config, config.Markdown()); err != nil {
		t.Fatal(err)
	}
	migrated, err := os.ReadFile(filepath.Join(owner.Root, ".langdrift/runtimes/runtime/console-terraform/main.tf.json"))
	if err != nil || !strings.Contains(string(migrated), "3000m") {
		t.Fatal("console resource edit was ignored")
	}
}

var _ application.ExecutionReader = storage.Executions{}

func TestNative_StartRechecksAuthoredMarkdownBeforeBuild(t *testing.T) {
	owner := ownerForTest(t)
	runner := &fakeRunner{}
	native := Native{Runner: runner, Executions: storage.Executions{}}
	config := nativeConfig()
	record := app.Connection{RuntimeID: config.RuntimeID, Repository: app.Repository{FullName: "owner/repository"}}
	var err error
	record.Managed, record.ConcreteHash, err = native.Materialize(context.Background(), owner, record, config, config.Markdown())
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(owner.Root, ".langdrift/runtimes/runtime/runtime.md"), []byte(config.Markdown()+"local change\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	reports := 0
	err = native.Start(context.Background(), owner, record, config, func(string, string) error { reports++; return nil })
	if err == nil || !strings.Contains(err.Error(), "conflict") || reports != 0 {
		t.Fatalf("err=%v reports=%d", err, reports)
	}
}

type inspectionFake struct {
	report application.Inspection
	err    error
}

func (f inspectionFake) Execute(context.Context, application.RuntimeRequest) (application.Inspection, error) {
	return f.report, f.err
}
func TestNative_ObserveDistinguishesLifecycleFailureFromMissingMetrics(t *testing.T) {
	healthy := application.Inspection{Execution: application.ExecutionOutcome{Status: "running", Settled: true}, Resources: []application.ResourceObservation{{Observed: "present"}}, Containers: []application.ContainerResult{{Ready: true, Phase: "Running"}}, Usage: application.UsageObservation{Status: "unavailable"}}
	native := Native{Inspect: inspectionFake{report: healthy}}
	observation, err := native.Observe(context.Background(), project.Project{}, app.Connection{})
	if err != nil || observation.State != "running" || observation.CPUUsed != nil || observation.MemoryUsedMiB != nil {
		t.Fatalf("missing metrics incorrectly changed lifecycle: %+v %v", observation, err)
	}
	healthy.Containers[0].Ready = false
	native.Inspect = inspectionFake{report: healthy}
	observation, err = native.Observe(context.Background(), project.Project{}, app.Connection{})
	if err != nil || observation.State != "unknown" {
		t.Fatalf("unready workload presented running: %+v %v", observation, err)
	}
	native.Inspect = inspectionFake{err: errors.New("cluster unavailable")}
	if _, err = native.Observe(context.Background(), project.Project{}, app.Connection{}); err == nil {
		t.Fatal("cluster loss hidden")
	}
}

func TestNative_BuildContextExcludesPlatformAndDotenvCredentials(t *testing.T) {
	owner := ownerForTest(t)
	if err := os.WriteFile(filepath.Join(owner.Root, ".env"), []byte("TOKEN=private"), 0o600); err != nil {
		t.Fatal(err)
	}
	runner := &contextRunner{}
	native := Native{Runner: runner}
	root, cleanup, err := native.buildTree(context.Background(), owner)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	if _, err := os.Stat(filepath.Join(root, "build.sh")); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, ".env")); !os.IsNotExist(err) {
		t.Fatal("dotenv entered source context")
	}
	if _, err := os.Stat(filepath.Join(root, ".langdrift")); !os.IsNotExist(err) {
		t.Fatal("platform state entered source context")
	}
}

type contextRunner struct{}

func (*contextRunner) Run(context.Context, process.Command) (process.Result, error) {
	return process.Result{Stdout: "build.sh\x00.env\x00.langdrift/state/private.json\x00"}, nil
}

// untrackedRunner models an empty Git index for private storage tests.
type untrackedRunner struct{}

func (untrackedRunner) Run(context.Context, process.Command) (process.Result, error) {
	return process.Result{}, nil
}

func TestNative_ObserveRequiresCompleteSamplesAndPreservesNativeWindows(t *testing.T) {
	now := time.Date(2026, 10, 21, 12, 0, 0, 0, time.UTC)
	makeReport := func() application.Inspection {
		return application.Inspection{ObservedAt: now, Execution: application.ExecutionOutcome{Status: "running", Settled: true}, Resources: []application.ResourceObservation{{Observed: "present"}}, Containers: []application.ContainerResult{{Pod: "api-1", Container: "api", Ready: true, Phase: "Running"}, {Pod: "api-2", Container: "api", Ready: true, Phase: "Running"}}, Usage: application.UsageObservation{Status: "observed", Pods: []application.PodUsage{
			{Pod: "api-1", Timestamp: now.Add(-30 * time.Second).Format(time.RFC3339), Window: "15s", Containers: []application.ContainerUsage{{Name: "api", Usage: map[string]any{"cpu": "200m", "memory": "20Mi"}}}},
			{Pod: "api-2", Timestamp: now.Add(-25 * time.Second).Format(time.RFC3339), Window: "10s", Containers: []application.ContainerUsage{{Name: "api", Usage: map[string]any{"cpu": "300m", "memory": "30Mi"}}}},
		}}}
	}
	for _, variant := range []string{"missing pod", "missing container", "duplicate container", "invalid window", "complete"} {
		t.Run(variant, func(t *testing.T) {
			report := makeReport()
			switch variant {
			case "missing pod":
				report.Usage.Pods = report.Usage.Pods[:1]
			case "missing container":
				report.Usage.Pods[1].Containers = nil
			case "duplicate container":
				report.Usage.Pods[1].Containers = append(report.Usage.Pods[1].Containers, report.Usage.Pods[1].Containers[0])
			case "invalid window":
				report.Usage.Pods[1].Window = "unknown"
			}
			native := Native{Inspect: inspectionFake{report: report}}
			observation, err := native.Observe(context.Background(), project.Project{}, app.Connection{})
			if err != nil || observation.State != "running" {
				t.Fatalf("metrics changed workload availability: %+v %v", observation, err)
			}
			if variant != "complete" {
				if observation.CPUUsed != nil || observation.MemoryUsedMiB != nil {
					t.Fatalf("partial sample reported an aggregate: %+v", observation)
				}
				return
			}
			if observation.CPUUsed == nil || *observation.CPUUsed != 0.5 || observation.MemoryUsedMiB == nil || *observation.MemoryUsedMiB != 50 || !observation.ObservedAt.Equal(now.Add(-30*time.Second)) || observation.MeasuredFrom == nil || !observation.MeasuredFrom.Equal(now.Add(-45*time.Second)) || len(observation.MeasurementWindows) != 2 {
				t.Fatalf("native sample provenance lost: %+v", observation)
			}
			report.ObservedAt = now.Add(time.Minute)
			native.Inspect = inspectionFake{report: report}
			repeated, err := native.Observe(context.Background(), project.Project{}, app.Connection{})
			if err != nil || !repeated.ObservedAt.Equal(observation.ObservedAt) {
				t.Fatal("repeated stale sample became fresh evidence")
			}
		})
	}
}

func TestStore_NearCapacityReopensAndRejectsOversizeReplacement(t *testing.T) {
	owner := ownerForTest(t)
	store := Store{Runner: untrackedRunner{}}
	record := app.Connection{RuntimeUID: "runtime-uid", RuntimeID: "runtime", State: "stopped", Secret: "runtime-secret", Credential: "agent-secret", Markdown: strings.Repeat("# Configuration\n", 2048)}
	if err := store.Bind(context.Background(), owner, record); err != nil {
		t.Fatal(err)
	}
	event := app.Event{ID: strings.Repeat("a", 32), Kind: "failed", Cause: "connection_lost", ExecutionID: strings.Repeat("b", 36), Revision: 5, IntentVersion: 10, OccurredAt: time.Now().UTC(), Detail: strings.Repeat("operation_failed ", 24), ConcreteHash: strings.Repeat("c", 64), MarkdownHash: strings.Repeat("d", 64)}
	low, high := 0, 10000
	for low < high {
		count := (low + high + 1) / 2
		record.Outbox = make([]app.Event, count)
		for i := range record.Outbox {
			record.Outbox[i] = event
			record.Outbox[i].Sequence = uint64(i + 1)
		}
		data, err := app.MarshalConnection(record)
		if err == nil && len(data) <= app.MaxConnectionStateBytes-(64<<10) {
			low = count
		} else {
			high = count - 1
		}
	}
	record.Outbox = record.Outbox[:low]
	for _, kind := range []string{"stop_started", "stopped"} {
		event.Kind = kind
		record.Outbox = append(record.Outbox, event)
	}
	if err := store.Save(owner, record); err != nil {
		t.Fatal(err)
	}
	reopened, err := store.List(owner)
	if err != nil || len(reopened) != 1 || len(reopened[0].Outbox) != len(record.Outbox) || reopened[0].Outbox[len(record.Outbox)-1].Kind != "stopped" {
		t.Fatalf("near-capacity cleanup ledger could not be reopened: %v", err)
	}
	record.Markdown = strings.Repeat("x", app.MaxConnectionStateBytes)
	if err := store.Save(owner, record); err == nil {
		t.Fatal("oversize state replaced a readable ledger")
	}
	if _, err := store.List(owner); err != nil {
		t.Fatalf("rejected replacement damaged prior state: %v", err)
	}
}

func TestStore_RefusesDeletedButGitTrackedCredentialPath(t *testing.T) {
	owner := ownerForTest(t)
	runner := process.Local{}
	run := func(args ...string) process.Result {
		t.Helper()
		result, err := runner.Run(context.Background(), process.Command{Name: "git", Args: args, Dir: owner.Root})
		if err != nil {
			t.Fatal(err)
		}
		return result
	}
	run("init", "-q")
	relative := ".langdrift/runtimes/runtime/.env"
	absolute := filepath.Join(owner.Root, relative)
	if err := os.MkdirAll(filepath.Dir(absolute), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(absolute, []byte("previous nonsecret example\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	run("add", "-f", "--", relative)
	if err := os.Remove(absolute); err != nil {
		t.Fatal(err)
	}
	store := Store{Runner: runner}
	err := store.Bind(context.Background(), owner, app.Connection{RuntimeUID: "runtime-uid", RuntimeID: "runtime", Credential: "agent-secret", Secret: "runtime-secret"})
	if err == nil || !strings.Contains(err.Error(), "tracked by Git") {
		t.Fatalf("tracked credential path was accepted: %v", err)
	}
	if _, err := os.Stat(absolute); !os.IsNotExist(err) {
		t.Fatal("private credentials were written into a tracked path")
	}
	if run("ls-files", "--", relative).Stdout != relative+"\n" {
		t.Fatal("consumer Git index was changed")
	}
}

func TestNative_CommittedScriptRecipeResolvesInsideTheCheckout(t *testing.T) {
	owner := ownerForTest(t)
	// failBuild stops Start right after the build phase, before the nil
	// infrastructure ports; the resolved command is recorded either way.
	runner := &fakeRunner{failBuild: true}
	native := Native{Runner: runner, Executions: storage.Executions{}}
	config := nativeConfig()
	record := app.Connection{RuntimeID: config.RuntimeID, Repository: app.Repository{FullName: "owner/repository"}}
	record.Managed, record.ConcreteHash, _ = native.Materialize(context.Background(), owner, record, config, config.Markdown())

	if err := native.Start(context.Background(), owner, record, config, func(string, string) error { return nil }); err == nil {
		t.Fatal("start succeeded with a failing build")
	}

	var build *process.Command
	for index, command := range runner.commands {
		if command.Name != "git" {
			build = &runner.commands[index]

			break
		}
	}
	if build == nil {
		t.Fatal("no build command ran")
	}
	if !filepath.IsAbs(build.Name) || !strings.HasSuffix(build.Name, "/build.sh") {
		t.Fatalf("script ran as %q, want the resolved absolute path", build.Name)
	}
	if !slices.Equal(build.PassEnv, []string{"DOCKER_*"}) || len(build.Env) != 0 {
		t.Fatalf("recipe environment PassEnv=%v Env=%v, want only the container engine pass-through", build.PassEnv, build.Env)
	}
	if rel, err := filepath.Rel(owner.Root, build.Name); err != nil || strings.HasPrefix(rel, "..") {
		t.Fatalf("script %q resolved outside the checkout", build.Name)
	}
}

func TestNative_MissingCommittedScriptFailsBeforeRunning(t *testing.T) {
	owner := ownerForTest(t)
	runner := &fakeRunner{}
	native := Native{Runner: runner, Executions: storage.Executions{}}
	config := nativeConfig()
	config.Components[0].Build = &connection.Build{Executable: "./absent.sh"}
	record := app.Connection{RuntimeID: config.RuntimeID, Repository: app.Repository{FullName: "owner/repository"}}
	record.Managed, record.ConcreteHash, _ = native.Materialize(context.Background(), owner, record, config, config.Markdown())

	err := native.Start(context.Background(), owner, record, config, func(string, string) error { return nil })
	if err == nil || !strings.Contains(err.Error(), "committed build script") {
		t.Fatalf("err=%v, want a committed-script resolution failure", err)
	}
	for _, command := range runner.commands {
		if command.Name != "git" {
			t.Fatalf("a non-git command ran: %v", command.Name)
		}
	}
}

func TestNative_FailedBuildExplainsItself(t *testing.T) {
	owner := ownerForTest(t)
	runner := &fakeRunner{failBuild: true, buildStderr: "compile exploded\nmissing header"}
	native := Native{Runner: runner, Executions: storage.Executions{}}
	config := nativeConfig()
	record := app.Connection{RuntimeID: config.RuntimeID, Repository: app.Repository{FullName: "owner/repository"}}
	record.Managed, record.ConcreteHash, _ = native.Materialize(context.Background(), owner, record, config, config.Markdown())

	err := native.Start(context.Background(), owner, record, config, func(string, string) error { return nil })
	if err == nil {
		t.Fatal("build failure reported success")
	}
	for _, want := range []string{"./build.sh", "missing header", "build.log"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("error %q does not name %q", err, want)
		}
	}

	log, readErr := os.ReadFile(filepath.Join(owner.Root, ".langdrift", "state", "runtimes", config.RuntimeID, "build.log"))
	if readErr != nil {
		t.Fatalf("build log: %v", readErr)
	}
	if !strings.Contains(string(log), "compile exploded") {
		t.Fatalf("build log %q lacks the tool output", log)
	}
}

func TestRelocateMigration_RefusesUnexpectedShapes(t *testing.T) {
	for name, schema := range map[string]string{
		"terraform not an object": `{"terraform": "text", "units": {}}`,
		"unit not an object":      `{"terraform": {"directory": "d"}, "units": {"api": 7}}`,
		"manifests not text":      `{"terraform": {"directory": "d"}, "units": {"api": {"kubernetes": {"manifests": 7}}}}`,
		"missing kubernetes":      `{"terraform": {"directory": "d"}, "units": {"api": {}}}`,
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := relocateMigration(map[string][]byte{"schema.json": []byte(schema)}); err == nil {
				t.Fatal("unexpected shape accepted")
			}
		})
	}
}
