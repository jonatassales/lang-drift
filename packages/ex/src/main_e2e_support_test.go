//go:build e2e

package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"
)

// Environment variables that control the opt-in infrastructure run.
const (
	// e2eEnable must be "1": the build tag alone never starts infrastructure.
	e2eEnable = "EX_E2E"

	// e2eBaseImage names a local image that provides sh, echo, sleep, and
	// BusyBox httpd. The harness never pulls it.
	e2eBaseImage = "EX_E2E_BASE_IMAGE"

	// toolLogVariable tells the counting shims where to log invocations.
	toolLogVariable = "EX_E2E_TOOL_LOG"

	defaultBaseImage = "busybox:1.36"
)

// Time bounds. Every finite CLI command receives an explicit --timeout, so
// the two-minute default never cuts a cluster operation short. Go skips
// t.Cleanup when -timeout expires, so the scenario stops cleanupBudget
// before the test deadline.
const (
	cliTimeout       = "10m"
	stepDeadline     = 12 * time.Minute
	probeDeadline    = 2 * time.Minute
	interruptGrace   = time.Minute
	readinessTimeout = 8 * time.Minute
	cleanupBudget    = 20 * time.Minute
	minimumScenario  = 30 * time.Minute
)

// nativeTools are the tools the CLI runs; each gets a counting shim.
var nativeTools = []string{"docker", "kind", "kubectl", "terraform"}

// harness owns every resource one run creates and the evidence it records.
// Only the test goroutine uses it.
type harness struct {
	t *testing.T

	// ctx ends when the scenario must stop: at the reserved deadline, or
	// when the test process is asked to stop. Cleanup does not use it.
	ctx context.Context

	nonce     string
	workdir   string // canonical; removed only when nothing is left behind
	project   string
	binary    string
	shims     string
	environ   []string          // environment of every CLI and native call
	tools     map[string]string // real executable paths, resolved before the shims
	baseImage string
	sequence  int

	// Identities learned during the run. Cleanup acts only on these.
	projectID        string
	clusterAttempted bool
	cluster          clusterResult
	upAttempted      bool
	foreignCreated   bool
	images           map[string]string // tag -> image ID recorded at build time

	report evidence
}

// outcome is one finished command.
type outcome struct {
	exit   int
	stdout string
	stderr string
}

// newHarness checks every prerequisite before any effect, prepares an
// isolated workspace, and builds the CLI from the source under test.
func newHarness(t *testing.T) *harness {
	t.Helper()

	if os.Getenv(e2eEnable) != "1" {
		t.Skipf("set %s=1 to run the infrastructure E2E; see packages/ex/README.md", e2eEnable)
	}

	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
		t.Skipf("EX supports macOS and Linux, not %s", runtime.GOOS)
	}

	h := &harness{
		t:         t,
		nonce:     randomHex(t, 4),
		tools:     map[string]string{},
		images:    map[string]string{},
		baseImage: defaultBaseImage,
	}
	h.report = newEvidence(h.nonce)

	if image := os.Getenv(e2eBaseImage); image != "" {
		h.baseImage = image
	}

	h.scenarioContext()
	h.requirePrerequisites()
	h.prepareWorkspace()
	t.Cleanup(h.cleanup)
	h.buildCLI()
	h.recordTools()

	return h
}

// scenarioContext ends the scenario before the test deadline, leaving the
// cleanup budget, and when the test process receives SIGINT or SIGTERM. The
// first signal stops the scenario and lets cleanup run; after it, signals
// regain their default effect, so a second one ends the process.
func (h *harness) scenarioContext() {
	h.t.Helper()

	signals, stopSignals := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	h.t.Cleanup(stopSignals)

	go func() {
		<-signals.Done()
		stopSignals()
	}()

	h.ctx = signals

	if deadline, bounded := h.t.Deadline(); bounded {
		stop := deadline.Add(-cleanupBudget)
		if time.Until(stop) < minimumScenario {
			h.t.Fatalf("go test -timeout leaves less than %s for the scenario after reserving %s for cleanup; use at least -timeout 70m",
				minimumScenario, cleanupBudget)
		}

		var cancel context.CancelFunc
		h.ctx, cancel = context.WithDeadline(signals, stop)
		h.t.Cleanup(cancel)
	}
}

// requireRunning fails the scenario once it must stop, so cleanup runs.
func (h *harness) requireRunning(step string) {
	h.t.Helper()

	if err := h.ctx.Err(); err != nil {
		h.t.Fatalf("%s: scenario stopped (%v); cleanup runs next", step, context.Cause(h.ctx))
	}
}

// requirePrerequisites fails closed, with instructions, before anything is
// created. It installs nothing and pulls no image.
func (h *harness) requirePrerequisites() {
	h.t.Helper()

	var missing []string

	for _, name := range append([]string{"go"}, nativeTools...) {
		path, err := exec.LookPath(name)
		if err != nil {
			missing = append(missing, name)
			continue
		}

		if path, err = filepath.Abs(path); err != nil {
			h.t.Fatal(err)
		}

		h.tools[name] = path
	}

	if len(missing) > 0 {
		h.t.Fatalf("missing prerequisites on PATH: %s; install them yourself, the harness installs nothing",
			strings.Join(missing, ", "))
	}

	ctx, cancel := context.WithTimeout(h.ctx, probeDeadline)
	defer cancel()

	if result := h.native(ctx, "docker", nil, "info", "--format", "{{.ServerVersion}}"); result.exit != 0 {
		h.t.Fatalf("the Docker daemon did not answer: %s", firstLine(result.stderr))
	}

	if result := h.native(ctx, "docker", nil, "image", "inspect", "--format", "{{.Id}}", h.baseImage); result.exit != 0 {
		h.t.Fatalf("base image %s is not present locally; pull or build it yourself, or set %s",
			h.baseImage, e2eBaseImage)
	}
}

// prepareWorkspace creates the run's private directory, the counting shims,
// and an environment whose KUBECONFIG is a disposable file, so no tool can
// read or rewrite the operator's own kubeconfig.
func (h *harness) prepareWorkspace() {
	h.t.Helper()

	directory, err := os.MkdirTemp("", "ex-e2e-")
	if err != nil {
		h.t.Fatal(err)
	}

	// The file adapters refuse symlinked ancestors, such as macOS's /var.
	if h.workdir, err = filepath.EvalSymlinks(directory); err != nil {
		h.t.Fatal(err)
	}

	h.project = filepath.Join(h.workdir, "project")
	h.binary = filepath.Join(h.workdir, "bin", "ex")
	h.shims = filepath.Join(h.workdir, "shims")

	for _, directory := range []string{h.project, filepath.Dir(h.binary), h.shims, filepath.Join(h.workdir, "tool-logs")} {
		if err := os.MkdirAll(directory, 0o700); err != nil {
			h.t.Fatal(err)
		}
	}

	// A shim logs one direct invocation by the CLI, then runs the real tool
	// with the original PATH, so calls a tool makes itself, such as kind
	// running docker, are not counted.
	for _, name := range nativeTools {
		script := fmt.Sprintf("#!/bin/sh\nprintf '%%s\\n' %s >> \"$%s\"\nPATH=%s\nexport PATH\nexec %s \"$@\"\n",
			shellQuote(name), toolLogVariable, shellQuote(os.Getenv("PATH")), shellQuote(h.tools[name]))

		if err := os.WriteFile(filepath.Join(h.shims, name), []byte(script), 0o700); err != nil {
			h.t.Fatal(err)
		}
	}

	h.environ = []string{
		"KUBECONFIG=" + filepath.Join(h.workdir, "ambient-kubeconfig"),
		"PATH=" + os.Getenv("PATH"),
		toolLogVariable + "=" + os.DevNull,
		"GOTOOLCHAIN=local",
	}

	for _, entry := range os.Environ() {
		name, _, _ := strings.Cut(entry, "=")
		if !slices.Contains([]string{"KUBECONFIG", "PATH", toolLogVariable, "GOTOOLCHAIN"}, name) {
			h.environ = append(h.environ, entry)
		}
	}
}

// buildCLI compiles the test-only native lifecycle fixture. Connected-console
// execution remains separately covered by fake component tests.
func (h *harness) buildCLI() {
	h.t.Helper()

	ctx, cancel := context.WithTimeout(h.ctx, stepDeadline)
	defer cancel()

	command := exec.CommandContext(ctx, h.tools["go"], "test", "-c", "-tags", "e2e", "-o", h.binary, ".")
	command.Env = h.environ

	if output, err := command.CombinedOutput(); err != nil {
		h.t.Fatalf("build the CLI: %v\n%s", err, output)
	}
}

// cliEnvironment routes one CLI invocation through the counting shims, with
// its own invocation log, so concurrent commands are counted separately.
func (h *harness) cliEnvironment() ([]string, string) {
	h.sequence++
	log := filepath.Join(h.workdir, "tool-logs", strconv.Itoa(h.sequence)+".log")

	// exec uses the last value of a duplicated variable.
	return append(slices.Clone(h.environ),
		"PATH="+h.shims+string(os.PathListSeparator)+os.Getenv("PATH"),
		toolLogVariable+"="+log,
	), log
}

// cliCommand prepares one CLI invocation. Canceling its context interrupts
// the CLI, so it settles its records and stops its own tool processes, and
// kills it only if it is still running after the grace period.
func (h *harness) cliCommand(ctx context.Context, arguments []string) (*exec.Cmd, string) {
	environment, log := h.cliEnvironment()

	command := exec.CommandContext(ctx, h.binary, append([]string{"-test.run=^TestNativeLifecycleFixture$", "--"}, arguments...)...)
	command.Dir = h.project
	command.Env = append(environment, "LANGDRIFT_NATIVE_FIXTURE=1")
	command.Cancel = func() error { return command.Process.Signal(os.Interrupt) }
	command.WaitDelay = interruptGrace

	return command, log
}

// ex runs one finite CLI command and records it as a scenario step.
func (h *harness) ex(step, stdin string, arguments ...string) outcome {
	h.t.Helper()
	h.requireRunning(step)

	ctx, cancel := context.WithTimeout(h.ctx, stepDeadline)
	defer cancel()

	result := h.runCLI(ctx, &h.report.Steps, step, stdin, arguments)

	h.requireRunning(step)

	if ctx.Err() != nil {
		h.t.Fatalf("%s: no result within %s", step, stepDeadline)
	}

	return result
}

// query runs one finite CLI command that the scenario polls; polls are not
// recorded as steps.
func (h *harness) query(arguments ...string) outcome {
	h.t.Helper()

	ctx, cancel := context.WithTimeout(h.ctx, probeDeadline)
	defer cancel()

	var discarded []stepRecord

	return h.runCLI(ctx, &discarded, "poll", "", arguments)
}

func (h *harness) runCLI(ctx context.Context, target *[]stepRecord, step, stdin string, arguments []string) outcome {
	command, log := h.cliCommand(ctx, arguments)
	command.Stdin = strings.NewReader(stdin)

	var stdout, stderr bytes.Buffer
	command.Stdout, command.Stderr = &stdout, &stderr

	started := time.Now()
	err := command.Run()
	result := outcome{exit: exitCodeOf(err), stdout: stdout.String(), stderr: stderr.String()}

	h.record(target, step, arguments, result, time.Since(started), countInvocations(log))

	return result
}

// background is a CLI process the scenario interrupts.
type background struct {
	step      string
	arguments []string
	log       string
	command   *exec.Cmd
	stdout    lockedBuffer
	stderr    lockedBuffer
	started   time.Time
	done      chan struct{}
	err       error
}

// start runs a CLI command that keeps running until the scenario stops it.
func (h *harness) start(step string, arguments ...string) *background {
	h.t.Helper()
	h.requireRunning(step)

	command, log := h.cliCommand(context.Background(), arguments)
	process := &background{step: step, arguments: arguments, log: log, command: command, done: make(chan struct{})}
	command.Stdout, command.Stderr = &process.stdout, &process.stderr
	process.started = time.Now()

	if err := command.Start(); err != nil {
		h.t.Fatalf("%s: %v", step, err)
	}

	go func() {
		process.err = command.Wait()
		close(process.done)
	}()

	// Registered after the harness cleanup, so it runs first: a failed
	// scenario still stops the process before resources are removed.
	h.t.Cleanup(process.stop)

	return process
}

// interrupt sends the same signal as Ctrl-C and waits for the CLI to exit.
func (h *harness) interrupt(process *background) outcome {
	h.t.Helper()

	if err := process.command.Process.Signal(os.Interrupt); err != nil && !errors.Is(err, os.ErrProcessDone) {
		h.t.Fatalf("%s: interrupt: %v", process.step, err)
	}

	select {
	case <-process.done:
	case <-time.After(interruptGrace):
		process.stop()
		h.t.Fatalf("%s: still running %s after the interrupt", process.step, interruptGrace)
	}

	result := outcome{exit: exitCodeOf(process.err), stdout: process.stdout.String(), stderr: process.stderr.String()}
	h.record(&h.report.Steps, process.step, process.arguments, result, time.Since(process.started), countInvocations(process.log))

	return result
}

// exited reports whether the process ended on its own.
func (process *background) exited() bool {
	select {
	case <-process.done:
		return true
	default:
		return false
	}
}

// stop interrupts the CLI so it settles its records and stops the tools it
// started in their own process groups, and kills it only after the grace
// period.
func (process *background) stop() {
	if process.exited() {
		return
	}

	_ = process.command.Process.Signal(os.Interrupt)

	select {
	case <-process.done:
	case <-time.After(interruptGrace):
		_ = process.command.Process.Kill()
		<-process.done
	}
}

// finiteCommand builds a finite command against this run's project, with
// compact results and an explicit deadline.
func (h *harness) finiteCommand(arguments ...string) []string {
	return append(arguments, "--project", h.project, "--json", "--timeout", cliTimeout)
}

// streamingCommand builds a runtime logs or forward command, which takes no
// --json.
func (h *harness) streamingCommand(verb string, flags ...string) []string {
	return append(append([]string{"runtime", verb, runtimeID}, flags...), "--project", h.project)
}

// native runs one native tool directly, outside the counting shims.
func (h *harness) native(ctx context.Context, tool string, stdin []byte, arguments ...string) outcome {
	command := exec.CommandContext(ctx, h.tools[tool], arguments...)
	command.Env = h.environ
	command.WaitDelay = 10 * time.Second

	if stdin != nil {
		command.Stdin = bytes.NewReader(stdin)
	}

	var stdout, stderr bytes.Buffer
	command.Stdout, command.Stderr = &stdout, &stderr

	err := command.Run()

	return outcome{exit: exitCodeOf(err), stdout: stdout.String(), stderr: stderr.String()}
}

// kubectl reaches the run's own cluster through its private credentials. It
// changes only objects the run creates itself.
func (h *harness) kubectl(ctx context.Context, stdin []byte, arguments ...string) outcome {
	target := []string{"--kubeconfig", h.cluster.Kubeconfig, "--context", h.cluster.Context, "--request-timeout=30s"}

	return h.native(ctx, "kubectl", stdin, append(target, arguments...)...)
}

// observe runs one read-only kubectl query and fails the run if it fails.
func (h *harness) observe(arguments ...string) string {
	h.t.Helper()

	ctx, cancel := context.WithTimeout(h.ctx, probeDeadline)
	defer cancel()

	result := h.kubectl(ctx, nil, arguments...)
	if result.exit != 0 {
		h.requireRunning("kubectl " + arguments[0])
		h.t.Fatalf("kubectl %s: %s", strings.Join(arguments, " "), sanitize(firstLine(result.stderr), h.workdir))
	}

	return strings.TrimSpace(result.stdout)
}

// uidOf observes one object; an absent object has no UID.
func (h *harness) uidOf(arguments ...string) string {
	h.t.Helper()

	return h.observe(append(arguments, "--ignore-not-found=true", "-o", "jsonpath={.metadata.uid}")...)
}

// namespaceUID observes a namespace; an absent namespace has no UID.
func (h *harness) namespaceUID(name string) string {
	h.t.Helper()

	return h.uidOf("get", "namespace", name)
}

// waitFor polls condition until it holds, the timeout elapses, or the
// scenario must stop.
func (h *harness) waitFor(description string, timeout time.Duration, condition func() bool) {
	h.t.Helper()

	deadline := time.Now().Add(timeout)

	for !condition() {
		h.requireRunning("wait for " + description)

		if time.Now().After(deadline) {
			h.t.Fatalf("timed out after %s waiting for %s", timeout, description)
		}

		time.Sleep(500 * time.Millisecond)
	}
}

// buildImage builds a disposable image from the local base image without
// pulling, labels it with the run nonce, and records its ID for cleanup.
func (h *harness) buildImage(name string, preparation string, command []string) string {
	h.t.Helper()
	h.requireRunning("build image " + name)

	tag := fmt.Sprintf("ex-e2e-%s:%s", name, h.nonce)
	buildContext := filepath.Join(h.workdir, "images", name)

	if err := os.MkdirAll(buildContext, 0o700); err != nil {
		h.t.Fatal(err)
	}

	commandJSON, err := json.Marshal(command)
	if err != nil {
		h.t.Fatal(err)
	}

	dockerfile := fmt.Sprintf("FROM %s\n%sLABEL io.langdrift.e2e.run=%q\nCMD %s\n", h.baseImage, preparation, h.nonce, commandJSON)
	if err := os.WriteFile(filepath.Join(buildContext, "Dockerfile"), []byte(dockerfile), 0o600); err != nil {
		h.t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(h.ctx, stepDeadline)
	defer cancel()

	started := time.Now()
	result := h.native(ctx, "docker", nil, "build", "--pull=false", "-t", tag, buildContext)
	h.record(&h.report.Steps, "build image "+name, []string{"docker", "build", "--pull=false", "-t", tag, "<context>"},
		result, time.Since(started), nil)

	if result.exit != 0 {
		h.t.Fatalf("build image %s: %s", tag, sanitize(result.stderr, h.workdir))
	}

	id := h.imageID(ctx, tag)
	if id.exit != 0 || strings.TrimSpace(id.stdout) == "" {
		h.t.Fatalf("inspect image %s: %s", tag, sanitize(firstLine(id.stderr), h.workdir))
	}

	h.images[tag] = strings.TrimSpace(id.stdout)

	return tag
}

// imageID lists the ID a tag names. An absent tag lists nothing and still
// succeeds, so a failure means the daemon could not answer.
func (h *harness) imageID(ctx context.Context, tag string) outcome {
	return h.native(ctx, "docker", nil, "image", "ls", "--no-trunc", "--quiet", "--filter", "reference="+tag)
}

// freePort finds a localhost port that is free now. Another process could
// take it before the forward binds, which then fails visibly.
func freePort(t *testing.T) int {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	return listener.Addr().(*net.TCPAddr).Port
}

func exitCodeOf(err error) int {
	var exit *exec.ExitError

	switch {
	case err == nil:
		return 0
	case errors.As(err, &exit) && exit.ExitCode() >= 0:
		return exit.ExitCode()
	default:
		return -1
	}
}

// decode reads one JSON result from a command's stdout. A failure reports
// the output sanitized, because test logs may be kept.
func decode[T any](h *harness, step string, result outcome) T {
	h.t.Helper()

	var value T
	if err := json.Unmarshal([]byte(result.stdout), &value); err != nil {
		h.t.Fatalf("%s: decode result: %v\nstdout: %s", step, err, sanitize(excerpt(result.stdout, 1024), h.workdir))
	}

	return value
}

func randomHex(t *testing.T, size int) string {
	t.Helper()

	data := make([]byte, size)
	if _, err := rand.Read(data); err != nil {
		t.Fatal(err)
	}

	return hex.EncodeToString(data)
}

func shellQuote(text string) string {
	return "'" + strings.ReplaceAll(text, "'", `'\''`) + "'"
}

func firstLine(text string) string {
	line, _, _ := strings.Cut(strings.TrimSpace(text), "\n")
	return line
}

// lockedBuffer collects a background process's output while it runs.
type lockedBuffer struct {
	mutex  sync.Mutex
	buffer bytes.Buffer
}

func (buffer *lockedBuffer) Write(data []byte) (int, error) {
	buffer.mutex.Lock()
	defer buffer.mutex.Unlock()

	return buffer.buffer.Write(data)
}

func (buffer *lockedBuffer) String() string {
	buffer.mutex.Lock()
	defer buffer.mutex.Unlock()

	return buffer.buffer.String()
}
