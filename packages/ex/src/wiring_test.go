package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/jonatassales/tao/packages/ex/src/presentation/cli/support"
)

// project is a consumer project driven through the composed CLI. Only
// commands that stop before any native tool are run: no Docker, kind,
// kubectl, or Terraform process starts in these tests.
type project struct {
	t    *testing.T
	root string
}

func newProject(t *testing.T) project {
	t.Helper()

	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	for _, component := range []string{"services/web", "jobs/batch"} {
		if err := os.MkdirAll(filepath.Join(root, component), 0o755); err != nil {
			t.Fatal(err)
		}
	}

	return project{t: t, root: root}
}

func (p project) ex(input string, arguments ...string) (int, string, string) {
	p.t.Helper()

	var stdout, stderr bytes.Buffer
	cases, native := useCases()
	status := executeNativeCatalogs(context.Background(), arguments, support.Options{
		Version:     version,
		Directory:   p.root,
		Input:       strings.NewReader(input),
		Output:      &stdout,
		Diagnostics: &stderr,
	}, cases, native)

	return status, stdout.String(), stderr.String()
}

const guidedAnswers = `team-payments
kind-ex-0123456789ab
demo-local
services/web, jobs/batch
Deployment
web:dev
100m
500m
128Mi
256Mi
8080
1
/healthz
Job
batch:dev
100m
200m
64Mi
128Mi
1
2
1Gi
2Gi
10
`

func TestGuidedRuntimeGenerationValidates(t *testing.T) {
	t.Parallel()

	p := newProject(t)

	if status, stdout, stderr := p.ex("", "init"); status != 0 || !strings.Contains(stdout, `"status": "initialized"`) {
		t.Fatalf("ex init = %d, %s%s", status, stdout, stderr)
	}

	status, stdout, stderr := p.ex(guidedAnswers, "runtime", "init", "demo", "--guided")
	if status != 0 || !strings.Contains(stdout, `"guided": true`) || !strings.Contains(stderr, "Runtime namespace: ") {
		t.Fatalf("ex runtime init --guided = %d, %s%s", status, stdout, stderr)
	}

	status, stdout, stderr = p.ex("", "runtime", "validate", "demo", "--json")
	if status != 0 {
		t.Fatalf("ex runtime validate = %d: %s", status, stderr)
	}

	var validation struct {
		Runtime string   `json:"runtime"`
		Units   []string `json:"units"`
		Objects int      `json:"objects"`
	}
	if err := json.Unmarshal([]byte(stdout), &validation); err != nil {
		t.Fatal(err)
	}

	if validation.Runtime != "demo" || !slices.Equal(validation.Units, []string{"batch", "web"}) || validation.Objects != 3 {
		t.Fatalf("validation = %+v, want both units and three native objects", validation)
	}

	if status, _, stderr := p.ex(guidedAnswers, "runtime", "init", "demo", "--guided"); status != 1 || !strings.Contains(stderr, "already exists") {
		t.Fatalf("a second ex runtime init = %d, %s; want the existing runtime kept", status, stderr)
	}

	workload := filepath.Join(p.root, ".langdrift", "runtimes", "demo", "units", "web", "kubernetes", "workload.json")
	data, err := os.ReadFile(workload)
	if err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(workload, bytes.Replace(data, []byte(`"Never"`), []byte(`"IfNotPresent"`), 1), 0o644); err != nil {
		t.Fatal(err)
	}

	if status, _, stderr := p.ex("", "runtime", "validate", "demo"); status != 1 || !strings.Contains(stderr, "imagePullPolicy: Never") {
		t.Fatalf("ex runtime validate after an edit = %d, %s; want the policy violation reported", status, stderr)
	}
}

func TestCommandsStopBeforeInfrastructureWithoutAProject(t *testing.T) {
	t.Parallel()

	p := newProject(t)

	for _, arguments := range [][]string{
		{"runtime", "up", "demo"},
		{"runtime", "down", "demo"},
		{"cluster", "inspect"},
		{"agent", "init", "researcher"},
	} {
		if status, _, stderr := p.ex("", arguments...); status != 4 || !strings.Contains(stderr, "no LangDrift project") {
			t.Errorf("ex %v = %d, %q; want the missing project reported", arguments, status, stderr)
		}
	}
}

func TestCommandsStopBeforeInfrastructureWithoutACluster(t *testing.T) {
	t.Parallel()

	p := newProject(t)
	if status, _, stderr := p.ex("", "init"); status != 0 {
		t.Fatal(stderr)
	}

	if status, _, stderr := p.ex(guidedAnswers, "runtime", "init", "demo", "--guided"); status != 0 {
		t.Fatal(stderr)
	}

	tests := []struct {
		arguments []string
		stderr    string
	}{
		// A missing record is a missing prerequisite: exit 4 by contract.
		{arguments: []string{"cluster", "inspect"}, stderr: "no cluster is recorded"},
		{arguments: []string{"runtime", "up", "demo"}, stderr: "no cluster is recorded"},
		{arguments: []string{"runtime", "down", "demo"}, stderr: "runtime demo has no execution record"},
		{arguments: []string{"runtime", "inspect", "demo"}, stderr: "runtime demo has no execution record"},
	}

	for _, tt := range tests {
		if status, _, stderr := p.ex("", tt.arguments...); status != 4 || !strings.Contains(stderr, tt.stderr) {
			t.Errorf("ex %v = %d, %q; want %q", tt.arguments, status, stderr, tt.stderr)
		}
	}

	if _, err := os.Stat(filepath.Join(p.root, ".langdrift", "state", "runtimes")); !os.IsNotExist(err) {
		t.Errorf("a refused command recorded runtime state: %v", err)
	}
}

func TestAgentScaffoldIsPublishedInsideTheProject(t *testing.T) {
	t.Parallel()

	p := newProject(t)
	if status, _, stderr := p.ex("", "init"); status != 0 {
		t.Fatal(stderr)
	}

	status, stdout, stderr := p.ex("", "agent", "init", "researcher")
	if status != 0 || !strings.Contains(stdout, filepath.Join(p.root, ".langdrift", "agents", "researcher")) {
		t.Fatalf("ex agent init = %d, %s%s", status, stdout, stderr)
	}

	if _, err := os.Stat(filepath.Join(p.root, ".langdrift", "agents", "researcher", "Dockerfile")); err != nil {
		t.Fatalf("the agent scaffold is incomplete: %v", err)
	}
}
