package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/jonatassales/tao/lib/cli/src/presentation/cli/support"
)

// ex init asks whether to generate CI Gates, or takes --ci-gates, and on a yes
// writes the one file it may write outside .langdrift/: a workflow pinned to
// the ex version that generated it. An existing workflow is kept byte for
// byte, and an unreleased build, which npx could not resolve, writes none.
func TestInitGeneratesTheCIGatesWorkflow(test *testing.T) {
	test.Parallel()

	golden, err := os.ReadFile(filepath.Join("testdata", "ci-gates", "langdrift-gates.yml"))
	if err != nil {
		test.Fatal(err)
	}

	released := func(test *testing.T, consumer project, input string, arguments ...string) (int, string, string) {
		test.Helper()

		var stdout, stderr bytes.Buffer
		cases, native := useCases()
		status := executeNativeCatalogs(context.Background(), arguments, support.Options{
			Version:     "0.1.0",
			Directory:   consumer.root,
			Input:       strings.NewReader(input),
			Output:      &stdout,
			Diagnostics: &stderr,
		}, cases, native)

		return status, stdout.String(), stderr.String()
	}

	workflowOf := func(consumer project) string {
		return filepath.Join(consumer.root, ".github", "workflows", "langdrift-gates.yml")
	}

	test.Run("writes the workflow on a yes", func(test *testing.T) {
		test.Parallel()

		consumer := newProject(test)

		status, stdout, stderr := released(test, consumer, "\nn\n\ny\n", "init", "--ask")
		if status != 0 {
			test.Fatalf("ex init --ask = %d, %s", status, stderr)
		}

		if !strings.Contains(stderr, "Generate the CI Gates workflow") {
			test.Errorf("ex init did not ask for CI Gates:\n%s", stderr)
		}

		written, err := os.ReadFile(workflowOf(consumer))
		if err != nil {
			test.Fatal(err)
		}

		if string(written) != string(golden) {
			test.Errorf("the workflow differs from the golden copy:\n%s", written)
		}

		var result struct {
			Recorded []string `json:"recorded"`
		}
		if err := json.Unmarshal([]byte(stdout), &result); err != nil || !slices.Equal(result.Recorded, []string{"ci-gates"}) {
			test.Errorf("result = %s, want ci-gates recorded", stdout)
		}
	})

	test.Run("writes nothing outside .langdrift without a yes", func(test *testing.T) {
		test.Parallel()

		consumer := newProject(test)

		if status, _, stderr := released(test, consumer, "\nn\n\n\n", "init", "--ask"); status != 0 {
			test.Fatalf("ex init --ask = %d, %s", status, stderr)
		}

		if _, err := os.Stat(filepath.Join(consumer.root, ".github")); err == nil {
			test.Error("ex init wrote .github without a yes")
		}
	})

	test.Run("keeps an existing workflow", func(test *testing.T) {
		test.Parallel()

		consumer := newProject(test)
		edited := []byte("# The team's own gates\nname: gates\n")

		if err := os.MkdirAll(filepath.Dir(workflowOf(consumer)), 0o755); err != nil {
			test.Fatal(err)
		}

		if err := os.WriteFile(workflowOf(consumer), edited, 0o644); err != nil {
			test.Fatal(err)
		}

		status, stdout, stderr := released(test, consumer, "", "init", "--ci-gates")
		if status != 0 {
			test.Fatalf("ex init --ci-gates = %d, %s", status, stderr)
		}

		if kept, _ := os.ReadFile(workflowOf(consumer)); string(kept) != string(edited) {
			test.Errorf("ex init replaced the existing workflow with:\n%s", kept)
		}

		if !strings.Contains(stdout, `"kept": [`) || !strings.Contains(stdout, `"ci-gates"`) {
			test.Errorf("result = %s, want ci-gates kept", stdout)
		}
	})

	test.Run("refuses an unreleased build", func(test *testing.T) {
		test.Parallel()

		consumer := newProject(test)

		status, _, stderr := consumer.ex("", "init", "--ci-gates")
		if status != 2 || !strings.Contains(stderr, "unreleased") {
			test.Errorf("ex init --ci-gates on 0.0.0-dev = %d, %s", status, stderr)
		}

		if _, err := os.Stat(workflowOf(consumer)); err == nil {
			test.Error("an unreleased build wrote the workflow")
		}
	})

	test.Run("the golden workflow follows the workflow conventions", func(test *testing.T) {
		test.Parallel()

		text := string(golden)
		pinned := regexp.MustCompile(`(?m)^\s+- uses: [a-z-]+/[a-z-]+@[0-9a-f]{40} # v[0-9]+\.[0-9]+\.[0-9]+$`)

		if uses := strings.Count(text, "uses:"); len(pinned.FindAllString(text, -1)) != uses || uses != 2 {
			test.Errorf("every action must be pinned to a SHA with its version; %d of %d are", len(pinned.FindAllString(text, -1)), uses)
		}

		if !strings.Contains(text, "\npermissions:\n  contents: read\n\n") {
			test.Error("the workflow must grant contents: read at the top level, and nothing more")
		}

		for _, line := range strings.Split(text, "\n") {
			if strings.Contains(line, "run:") && strings.Contains(line, "${{") {
				test.Errorf("an expression inside run: %s", line)
			}
		}
	})
}
