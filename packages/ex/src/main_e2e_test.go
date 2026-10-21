//go:build e2e

package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"
)

// The infrastructure E2E drives the built CLI against a disposable kind
// cluster: the full lifecycle, an interrupted start with explicit recovery,
// an interrupted forward, and refusals to adopt or delete what EX did not
// create. It compiles only with -tags e2e and runs only with EX_E2E=1. See
// packages/ex/README.md for prerequisites and what it may download.

const (
	runtimeID     = "e2e"
	webUnit       = "web"
	taskUnit      = "task"
	webPort       = 8080
	webReadyLine  = "ex-e2e-web-ready"
	taskDoneLine  = "ex-e2e-task-done"
	startupPause  = 10 // seconds the web unit waits before serving, so an interrupt lands mid-start
	forwardBudget = 90 * time.Second
)

// Result shapes, limited to the fields the scenario asserts.

type projectResult struct {
	ID string `json:"id"`
}

type diagnosis struct {
	ToolsAvailable  bool   `json:"toolsAvailable"`
	ContainerEngine string `json:"containerEngine"`
}

type clusterResult struct {
	Name       string `json:"name"`
	Context    string `json:"context"`
	Kubeconfig string `json:"kubeconfig"`
	UID        string `json:"uid"`
	Status     string `json:"status"`
}

type validation struct {
	Units   []string `json:"units"`
	Objects int      `json:"objects"`
}

type executionOutcome struct {
	Status       string `json:"status"`
	Settled      bool   `json:"settled"`
	RunID        string `json:"runId"`
	NamespaceUID string `json:"namespaceUid"`
	Resources    []struct {
		Kind string `json:"kind"`
		Name string `json:"name"`
		UID  string `json:"uid"`
	} `json:"resources"`
}

// recordedUID returns the UID the execution recorded for one object.
func (state executionOutcome) recordedUID(kind, name string) string {
	for _, resource := range state.Resources {
		if resource.Kind == kind && resource.Name == name {
			return resource.UID
		}
	}

	return ""
}

type inspection struct {
	State      executionOutcome `json:"state"`
	Containers []struct {
		Unit    string `json:"unit"`
		Ready   bool   `json:"ready"`
		Phase   string `json:"podPhase"`
		ImageID string `json:"observedImageId"`
	} `json:"containers"`
}

type failure struct {
	Error string `json:"error"`
}

func TestE2ELifecycle(t *testing.T) {
	h := newHarness(t)

	h.initializeProject()
	h.createCluster()
	webImage, taskImage := h.buildImages()
	h.prepareRuntime(webImage, taskImage)
	h.loadImages(webImage, taskImage)
	h.refuseForeignNamespaceOnUp()
	h.interruptAndRecoverUp()
	h.runLifecycle()
	h.stopRuntime()
	h.refuseForeignNamespaceOnClusterDelete()
	h.deleteCluster()
}

func (h *harness) namespace() string { return "ex-e2e-" + h.nonce }

// expect fails the run unless the command exited with the wanted code.
func (h *harness) expect(step string, result outcome, exit int) {
	h.t.Helper()

	if result.exit != exit {
		h.t.Fatalf("%s: exit code %d, want %d\nstdout: %s\nstderr: %s", step, result.exit, exit,
			sanitize(excerpt(result.stdout, 2048), h.workdir), sanitize(excerpt(result.stderr, 2048), h.workdir))
	}
}

// expectRefusal requires an operational refusal whose structured diagnostic,
// the last line of stderr under --json, names the reason.
func (h *harness) expectRefusal(step string, result outcome, reason string) {
	h.t.Helper()
	h.expect(step, result, 1)

	lines := strings.Split(strings.TrimSpace(result.stderr), "\n")
	report := decode[failure](h, step, outcome{stdout: lines[len(lines)-1]})

	if !strings.Contains(report.Error, reason) {
		h.t.Fatalf("%s: refused with %q, want a reason containing %q", step, sanitize(report.Error, h.workdir), reason)
	}
}

func (h *harness) initializeProject() {
	result := h.ex("init", "", h.finiteCommand("init")...)
	h.expect("init", result, 0)
	h.projectID = decode[projectResult](h, "init", result).ID

	doctor := h.ex("doctor", "", h.finiteCommand("doctor")...)
	h.expect("doctor", doctor, 0)

	if report := decode[diagnosis](h, "doctor", doctor); !report.ToolsAvailable || report.ContainerEngine != "reachable" {
		h.t.Fatalf("doctor reported %+v; every tool and the container engine must be available", report)
	}
}

func (h *harness) createCluster() {
	h.clusterAttempted = true

	result := h.ex("cluster create", "", h.finiteCommand("cluster", "create")...)
	h.expect("cluster create", result, 0)
	h.cluster = decode[clusterResult](h, "cluster create", result)

	if h.cluster.Status != "ready" || h.cluster.UID == "" || h.cluster.Name != h.expectedClusterName() {
		h.t.Fatalf("cluster create reported %+v; want a ready cluster named after the project", h.cluster)
	}

	h.requireClusterUnchanged("cluster inspect")

	ctx, cancel := context.WithTimeout(h.ctx, probeDeadline)
	defer cancel()

	if server := h.kubectl(ctx, nil, "version", "-o", "json"); server.exit == 0 {
		h.report.Tools["kubernetes"] = excerpt(sanitize(strings.TrimSpace(server.stdout), h.workdir), 1024)
	}
}

// requireClusterUnchanged requires the recorded cluster to be the created
// one and still ready.
func (h *harness) requireClusterUnchanged(step string) {
	h.t.Helper()

	inspected := h.ex(step, "", h.finiteCommand("cluster", "inspect")...)
	h.expect(step, inspected, 0)

	if observed := decode[clusterResult](h, step, inspected); observed.UID != h.cluster.UID ||
		observed.Name != h.cluster.Name || observed.Status != "ready" {
		h.t.Fatalf("%s reported %+v, want the created cluster %s, ready", step, observed, h.cluster.UID)
	}
}

// buildImages builds the web image, which logs, pauses, then serves
// /healthz, and the task image, which logs and exits.
func (h *harness) buildImages() (string, string) {
	webImage := h.buildImage(webUnit,
		"RUN mkdir -p /www && echo ok > /www/healthz\n",
		[]string{"sh", "-c", fmt.Sprintf("echo %s; sleep %d; exec httpd -f -v -p %d -h /www", webReadyLine, startupPause, webPort)})
	taskImage := h.buildImage(taskUnit, "", []string{"sh", "-c", "echo " + taskDoneLine})

	return webImage, taskImage
}

// prepareRuntime generates a runtime with a Deployment and Service unit and
// a Job unit through the guided prompts, then validates it.
func (h *harness) prepareRuntime(webImage, taskImage string) {
	for _, unit := range []string{webUnit, taskUnit} {
		if err := os.MkdirAll(filepath.Join(h.project, "components", unit), 0o700); err != nil {
			h.t.Fatal(err)
		}
	}

	answers := strings.Join([]string{
		"langdrift-e2e",
		h.cluster.Context,
		h.namespace(),
		"components/" + webUnit + ", components/" + taskUnit,
		"Deployment", webImage, "50m", "200m", "32Mi", "64Mi", strconv.Itoa(webPort), "1", "/healthz",
		"Job", taskImage, "50m", "200m", "32Mi", "64Mi",
		"1", "2", "512Mi", "1Gi", "10", // requests.cpu, limits.cpu, requests.memory, limits.memory, pods
	}, "\n") + "\n"

	initialized := h.ex("runtime init --guided", answers, h.finiteCommand("runtime", "init", runtimeID, "--guided")...)
	h.expect("runtime init --guided", initialized, 0)

	validated := h.ex("runtime validate", "", h.finiteCommand("runtime", "validate", runtimeID)...)
	h.expect("runtime validate", validated, 0)

	if report := decode[validation](h, "runtime validate", validated); len(report.Units) != 2 || report.Objects != 3 {
		h.t.Fatalf("runtime validate reported %+v; want two units and three objects", report)
	}
}

func (h *harness) loadImages(images ...string) {
	for _, image := range images {
		result := h.ex("cluster load-image "+image, "", h.finiteCommand("cluster", "load-image", image)...)
		h.expect("cluster load-image", result, 0)
	}
}

// createForeignNamespace creates a namespace EX did not create, labeled with
// the run nonce so that cleanup finds exactly what this run made.
func (h *harness) createForeignNamespace(name string) string {
	h.t.Helper()

	manifest, err := json.Marshal(map[string]any{
		"apiVersion": "v1",
		"kind":       "Namespace",
		"metadata":   map[string]any{"name": name, "labels": map[string]string{foreignLabelKey: h.nonce}},
	})
	if err != nil {
		h.t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(h.ctx, probeDeadline)
	defer cancel()

	h.foreignCreated = true

	if result := h.kubectl(ctx, manifest, "create", "-f", "-"); result.exit != 0 {
		h.t.Fatalf("create foreign namespace %s: %s", name, sanitize(firstLine(result.stderr), h.workdir))
	}

	uid := h.namespaceUID(name)
	if uid == "" {
		h.t.Fatalf("foreign namespace %s was not observed after creation", name)
	}

	return uid
}

// removeForeignNamespace deletes the run's labeled foreign namespaces and
// waits until name is gone.
func (h *harness) removeForeignNamespace(name string) {
	h.t.Helper()

	ctx, cancel := context.WithTimeout(h.ctx, readinessTimeout)
	defer cancel()

	if removed := h.deleteForeignNamespaces(ctx); removed.exit != 0 {
		h.t.Fatalf("remove foreign namespace %s: %s", name, sanitize(firstLine(removed.stderr), h.workdir))
	}

	h.waitNamespaceGone(name)
}

// waitNamespaceGone waits until a namespace no longer exists.
func (h *harness) waitNamespaceGone(name string) {
	h.t.Helper()
	h.waitFor("namespace "+name+" to be deleted", readinessTimeout, func() bool {
		return h.namespaceUID(name) == ""
	})
}

// refuseForeignNamespaceOnUp requires up to refuse a namespace EX did not
// create, to leave it untouched, and to record nothing.
func (h *harness) refuseForeignNamespaceOnUp() {
	foreignUID := h.createForeignNamespace(h.namespace())

	h.upAttempted = true
	result := h.ex("runtime up with a foreign namespace", "", h.finiteCommand("runtime", "up", runtimeID)...)
	h.expectRefusal("runtime up with a foreign namespace", result, "refusing to adopt")

	if uid := h.namespaceUID(h.namespace()); uid != foreignUID {
		h.t.Fatalf("foreign namespace UID is %q after the refused up, want %q", uid, foreignUID)
	}

	h.expectRefusal("runtime inspect after the refused up",
		h.ex("runtime inspect after the refused up", "", h.finiteCommand("runtime", "inspect", runtimeID)...),
		"no execution record")

	h.requireClusterUnchanged("cluster inspect after the refused up")
	h.removeForeignNamespace(h.namespace())
}

// interruptAndRecoverUp interrupts up while it waits for its workloads,
// then requires the recorded state, the refusals that protect it, and
// explicit recovery to remove everything.
func (h *harness) interruptAndRecoverUp() {
	process := h.start("runtime up, interrupted", h.finiteCommand("runtime", "up", runtimeID)...)

	// Up creates units in name order and each unit's files in lexical order,
	// so the web Deployment (workload.json) is the last object created and
	// recorded. Once inspect reports its recorded UID, up is waiting for
	// readiness, which the web unit's startup pause holds back. Interrupting
	// a create in flight is a different, timing-dependent case this
	// scenario does not target.
	h.waitFor("up to record the web Deployment", readinessTimeout, func() bool {
		if process.exited() {
			h.t.Fatalf("runtime up ended before it could be interrupted:\n%s", sanitize(process.stderr.String(), h.workdir))
		}

		polled := h.query(h.finiteCommand("runtime", "inspect", runtimeID)...)
		if polled.exit != 0 {
			return false
		}

		var report inspection

		return json.Unmarshal([]byte(polled.stdout), &report) == nil &&
			report.State.recordedUID("Deployment", webUnit) != ""
	})

	h.expect("runtime up, interrupted", h.interrupt(process), 130)

	interrupted := h.inspectState("runtime inspect after the interrupt")
	if interrupted.Status != "interrupted" || !interrupted.Settled {
		h.t.Fatalf("after the interrupt the runtime is %q (settled %t), want a settled interruption",
			interrupted.Status, interrupted.Settled)
	}

	deploymentUID := h.uidOf("get", "deployment", webUnit, "--namespace", h.namespace())
	if deploymentUID == "" || deploymentUID != interrupted.recordedUID("Deployment", webUnit) {
		h.t.Fatalf("live web Deployment UID %q differs from the recorded one", deploymentUID)
	}

	h.expectRefusal("runtime up after the interrupt",
		h.ex("runtime up after the interrupt", "", h.finiteCommand("runtime", "up", runtimeID)...), "run down before a new up")

	h.expectRefusal("runtime down without --recover",
		h.ex("runtime down without --recover", "", h.finiteCommand("runtime", "down", runtimeID)...), "explicit recovery")

	// Neither refusal may change the record or the live objects.
	unchanged := h.inspectState("runtime inspect after the refusals")
	if unchanged.Status != interrupted.Status || unchanged.RunID != interrupted.RunID ||
		unchanged.NamespaceUID != interrupted.NamespaceUID ||
		h.uidOf("get", "deployment", webUnit, "--namespace", h.namespace()) != deploymentUID {
		h.t.Fatalf("the refusals changed the runtime: %+v, before %+v", unchanged, interrupted)
	}

	recovered := h.ex("runtime down --recover", "", h.finiteCommand("runtime", "down", runtimeID, "--recover")...)
	h.expect("runtime down --recover", recovered, 0)

	if state := decode[executionOutcome](h, "runtime down --recover", recovered); state.Status != "down" {
		h.t.Fatalf("recovery left the runtime %q, want down", state.Status)
	}

	h.waitNamespaceGone(h.namespace())
}

// inspectState inspects the runtime and returns its recorded execution.
func (h *harness) inspectState(step string) executionOutcome {
	h.t.Helper()

	inspected := h.ex(step, "", h.finiteCommand("runtime", "inspect", runtimeID)...)
	h.expect(step, inspected, 0)

	return decode[inspection](h, step, inspected).State
}

// runLifecycle starts the runtime and exercises every observation.
func (h *harness) runLifecycle() {
	started := h.ex("runtime up", "", h.finiteCommand("runtime", "up", runtimeID)...)
	h.expect("runtime up", started, 0)

	state := decode[executionOutcome](h, "runtime up", started)
	if state.Status != "running" || !state.Settled || state.NamespaceUID == "" || len(state.Resources) == 0 {
		h.t.Fatalf("runtime up reported %+v; want a settled running runtime with recorded identities", state)
	}

	for _, resource := range state.Resources {
		if resource.UID == "" {
			h.t.Fatalf("runtime up recorded %s/%s without a UID", resource.Kind, resource.Name)
		}
	}

	if uid := h.namespaceUID(h.namespace()); uid != state.NamespaceUID {
		h.t.Fatalf("live namespace UID %q differs from the recorded %q", uid, state.NamespaceUID)
	}

	h.inspectRunning()
	h.readLogs(webUnit, webReadyLine)
	h.readLogs(taskUnit, taskDoneLine)
	h.forward()
}

func (h *harness) inspectRunning() {
	inspected := h.ex("runtime inspect", "", h.finiteCommand("runtime", "inspect", runtimeID)...)
	h.expect("runtime inspect", inspected, 0)

	report := decode[inspection](h, "runtime inspect", inspected)
	if report.State.Status != "running" {
		h.t.Fatalf("runtime inspect reported %q, want running", report.State.Status)
	}

	var webReady, taskSucceeded bool

	for _, container := range report.Containers {
		webReady = webReady || (container.Unit == webUnit && container.Ready && container.ImageID != "")
		taskSucceeded = taskSucceeded || (container.Unit == taskUnit && container.Phase == "Succeeded")
	}

	if !webReady || !taskSucceeded {
		h.t.Fatalf("runtime inspect containers %+v; want a ready web container and a succeeded task", report.Containers)
	}
}

func (h *harness) readLogs(unit, wanted string) {
	step := "runtime logs --unit " + unit
	result := h.ex(step, "", h.streamingCommand("logs", "--unit", unit, "--timeout", cliTimeout)...)
	h.expect(step, result, 0)

	if !strings.Contains(result.stdout, wanted) {
		h.t.Fatalf("%s did not contain %q:\n%s", step, wanted, sanitize(excerpt(result.stdout, 2048), h.workdir))
	}
}

// forward opens a tunnel, answers HTTP through it while another observation
// runs, and requires an interrupt to end it with exit code 130 and the
// local port released.
func (h *harness) forward() {
	port := freePort(h.t)
	address := net.JoinHostPort("127.0.0.1", strconv.Itoa(port))
	process := h.start("runtime forward", h.streamingCommand("forward", "--unit", webUnit, "--port", strconv.Itoa(port))...)

	client := http.Client{Timeout: 5 * time.Second}

	h.waitFor("the forward to answer HTTP", forwardBudget, func() bool {
		if process.exited() {
			h.t.Fatalf("runtime forward ended on its own:\n%s", sanitize(process.stderr.String(), h.workdir))
		}

		response, err := client.Get("http://" + address + "/healthz")
		if err != nil {
			return false
		}
		defer response.Body.Close()

		body, err := io.ReadAll(io.LimitReader(response.Body, 1024))

		return err == nil && response.StatusCode == http.StatusOK && strings.TrimSpace(string(body)) == "ok"
	})

	// Inspection takes no project lock, so it coexists with the tunnel.
	h.expect("runtime inspect during forward",
		h.ex("runtime inspect during forward", "", h.finiteCommand("runtime", "inspect", runtimeID)...), 0)

	h.expect("runtime forward", h.interrupt(process), 130)

	h.waitFor("the forwarded port to be released", 10*time.Second, func() bool {
		connection, err := net.DialTimeout("tcp", address, time.Second)
		if err != nil {
			return true // nothing accepts connections any more
		}

		connection.Close()

		return false
	})
}

func (h *harness) stopRuntime() {
	stopped := h.ex("runtime down", "", h.finiteCommand("runtime", "down", runtimeID)...)
	h.expect("runtime down", stopped, 0)

	if state := decode[executionOutcome](h, "runtime down", stopped); state.Status != "down" {
		h.t.Fatalf("runtime down left the runtime %q, want down", state.Status)
	}

	h.waitNamespaceGone(h.namespace())
}

// refuseForeignNamespaceOnClusterDelete requires cluster deletion to refuse
// while the cluster holds a namespace EX did not create, and to change
// nothing.
func (h *harness) refuseForeignNamespaceOnClusterDelete() {
	name := "ex-e2e-foreign-" + h.nonce
	foreignUID := h.createForeignNamespace(name)

	result := h.ex("cluster delete with a foreign namespace", "", h.finiteCommand("cluster", "delete")...)
	h.expectRefusal("cluster delete with a foreign namespace", result, "cluster deletion is unsafe")

	if uid := h.namespaceUID(name); uid != foreignUID {
		h.t.Fatalf("foreign namespace UID is %q after the refused deletion, want %q", uid, foreignUID)
	}

	h.requireClusterUnchanged("cluster inspect after the refused deletion")
	h.removeForeignNamespace(name)
	h.foreignCreated = false
}

func (h *harness) deleteCluster() {
	result := h.ex("cluster delete", "", h.finiteCommand("cluster", "delete")...)
	h.expect("cluster delete", result, 0)

	if deleted := decode[clusterResult](h, "cluster delete", result); deleted.Status != "deleted" {
		h.t.Fatalf("cluster delete reported %q, want deleted", deleted.Status)
	}

	ctx, cancel := context.WithTimeout(h.ctx, probeDeadline)
	defer cancel()

	listed := h.native(ctx, "kind", nil, "get", "clusters")
	if listed.exit != 0 || slices.Contains(strings.Fields(listed.stdout), h.cluster.Name) {
		h.t.Fatalf("kind still lists %s after deletion, or listing failed: %s", h.cluster.Name, firstLine(listed.stderr))
	}

	if _, err := os.Lstat(h.cluster.Kubeconfig); !errors.Is(err, os.ErrNotExist) {
		h.t.Fatalf("the private kubeconfig remains after deletion: %v", err)
	}

	// Nothing remains for cleanup to remove except the images.
	h.upAttempted, h.clusterAttempted = false, false
}
