package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ex init asks for the folder that holds the team's tickets and specs, or
// takes --artifacts, and stores it in config.json as a path relative to the
// project root. A path that leaves the project, or names no folder, is
// refused before it is stored.
func TestInitRecordsTheArtifactsFolder(test *testing.T) {
	test.Parallel()

	artifactsOf := func(test *testing.T, root string) string {
		test.Helper()

		var config struct {
			Artifacts string `json:"artifacts"`
		}
		if err := json.Unmarshal(mustRead(test, filepath.Join(root, ".langdrift", "config.json")), &config); err != nil {
			test.Fatal(err)
		}

		return config.Artifacts
	}

	withArtifacts := func(test *testing.T) project {
		test.Helper()

		consumer := newProject(test)
		if err := os.MkdirAll(filepath.Join(consumer.root, ".agents", "workflow"), 0o755); err != nil {
			test.Fatal(err)
		}

		return consumer
	}

	test.Run("stores the answer", func(test *testing.T) {
		test.Parallel()

		consumer := withArtifacts(test)

		status, stdout, stderr := consumer.ex("\nn\n.agents/workflow\n", "init", "--ask")
		if status != 0 {
			test.Fatalf("ex init --ask = %d, %s", status, stderr)
		}

		if !strings.Contains(stderr, "Artifacts folder") {
			test.Errorf("ex init did not ask for the artifacts folder:\n%s", stderr)
		}

		if got := artifactsOf(test, consumer.root); got != ".agents/workflow" {
			test.Errorf("config.json artifacts = %q, want .agents/workflow", got)
		}

		var result struct {
			Recorded []string `json:"recorded"`
		}
		if err := json.Unmarshal([]byte(stdout), &result); err != nil || len(result.Recorded) != 1 || result.Recorded[0] != "artifacts" {
			test.Errorf("result = %s, want artifacts recorded", stdout)
		}

		status, _, stderr = consumer.ex("", "init", "--ask")
		if status != 0 || strings.Contains(stderr, "Artifacts folder") {
			test.Errorf("a second ex init --ask = %d, asked again:\n%s", status, stderr)
		}
	})

	test.Run("takes the flag", func(test *testing.T) {
		test.Parallel()

		consumer := withArtifacts(test)

		if status, _, stderr := consumer.ex("", "init", "--artifacts", "./.agents/workflow/"); status != 0 {
			test.Fatalf("ex init --artifacts = %d, %s", status, stderr)
		}

		if got := artifactsOf(test, consumer.root); got != ".agents/workflow" {
			test.Errorf("config.json artifacts = %q, want .agents/workflow", got)
		}
	})

	test.Run("stores nothing when the budget is refused", func(test *testing.T) {
		test.Parallel()

		consumer := withArtifacts(test)

		status, _, stderr := consumer.ex("", "init", "--budget", "10", "--per-run", "100", "--reserve-margin", "5", "--artifacts", ".agents/workflow")
		if status != 2 {
			test.Fatalf("ex init with a per-run cap above the budget = %d, %s", status, stderr)
		}

		if got := artifactsOf(test, consumer.root); got != "" {
			test.Errorf("a refused budget left the artifacts folder %q recorded", got)
		}
	})

	test.Run("refuses a path outside the project or without a folder", func(test *testing.T) {
		test.Parallel()

		for _, path := range []string{"/tmp", "../outside", ".agents/missing", "README.md"} {
			consumer := withArtifacts(test)
			if err := os.WriteFile(filepath.Join(consumer.root, "README.md"), []byte("# Project\n"), 0o644); err != nil {
				test.Fatal(err)
			}

			status, _, stderr := consumer.ex("", "init", "--artifacts", path)
			if status != 2 || !strings.Contains(stderr, "artifacts") {
				test.Errorf("ex init --artifacts %s = %d, %s", path, status, stderr)
			}

			if _, err := os.Stat(filepath.Join(consumer.root, ".langdrift", "config.json")); err == nil && artifactsOf(test, consumer.root) != "" {
				test.Errorf("ex init stored the artifacts folder %s", path)
			}
		}
	})
}
