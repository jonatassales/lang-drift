package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// dft gates report turns a pull request's commit range into the CI Gates
// report with no project: it reads git alone, and every FinOps reading stays
// unknown because CI has no committed FinOps source. It changes environment
// variables, so it never runs in parallel.
func TestGatesReport(test *testing.T) {
	root, err := filepath.EvalSymlinks(test.TempDir())
	if err != nil {
		test.Fatal(err)
	}

	git := func(arguments ...string) string {
		test.Helper()

		command := exec.Command("git", append([]string{"-C", root}, arguments...)...)
		command.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=Gates",
			"GIT_AUTHOR_EMAIL=gates@example.com",
			"GIT_COMMITTER_NAME=Gates",
			"GIT_COMMITTER_EMAIL=gates@example.com",
		)

		output, err := command.CombinedOutput()
		if err != nil {
			test.Fatalf("git %v: %v\n%s", arguments, err, output)
		}

		return strings.TrimSpace(string(output))
	}

	git("init", "--quiet")
	git("commit", "--quiet", "--allow-empty", "--message", "chore: start")
	base := git("rev-parse", "HEAD")
	git("commit", "--quiet", "--allow-empty", "--message", "feat(dft): ask on init")
	git("commit", "--quiet", "--allow-empty", "--message", "fix(dft): keep a | in the subject")
	head := git("rev-parse", "HEAD")
	first := git("rev-parse", "--short", "HEAD~1")
	second := git("rev-parse", "--short", "HEAD")

	consumer := project{
		t:    test,
		root: root,
	}

	test.Run("prints the range and leaves FinOps unknown", func(test *testing.T) {
		test.Setenv("BASE_SHA", base)
		test.Setenv("HEAD_SHA", head)
		test.Setenv("GITHUB_STEP_SUMMARY", "")

		status, stdout, stderr := consumer.dft("", "gates", "report")
		if status != 0 {
			test.Fatalf("dft gates report = %d, %s", status, stderr)
		}

		for _, want := range []string{
			"## LangDrift CI Gates",
			"`" + base + ".." + head + "`: 2 commits",
			"| `" + first + "` | feat(dft): ask on init |",
			"| `" + second + "` | fix(dft): keep a \\| in the subject |",
			"| Budget | unknown |",
			"| Spend | unknown |",
			"| Ledger | unknown |",
			"no committed FinOps source",
		} {
			if !strings.Contains(stdout, want) {
				test.Errorf("the report lacks %q:\n%s", want, stdout)
			}
		}

		if strings.Contains(stdout, "$0") {
			test.Errorf("the report reads an unknown spend as zero:\n%s", stdout)
		}

		if _, err := os.Stat(filepath.Join(root, ".langdrift")); err == nil {
			test.Error("dft gates report created .langdrift")
		}
	})

	test.Run("appends to the job summary", func(test *testing.T) {
		summary := filepath.Join(test.TempDir(), "summary.md")
		if err := os.WriteFile(summary, []byte("earlier step\n"), 0o644); err != nil {
			test.Fatal(err)
		}

		test.Setenv("BASE_SHA", base)
		test.Setenv("HEAD_SHA", head)
		test.Setenv("GITHUB_STEP_SUMMARY", summary)

		status, stdout, stderr := consumer.dft("", "gates", "report")
		if status != 0 || stdout != "" {
			test.Fatalf("dft gates report = %d, stdout %q, %s", status, stdout, stderr)
		}

		written, err := os.ReadFile(summary)
		if err != nil {
			test.Fatal(err)
		}

		if !strings.HasPrefix(string(written), "earlier step\n## LangDrift CI Gates") {
			test.Errorf("the job summary = %q, want the report appended", written)
		}
	})

	test.Run("names a missing SHA", func(test *testing.T) {
		test.Setenv("BASE_SHA", "")
		test.Setenv("HEAD_SHA", head)

		status, _, stderr := consumer.dft("", "gates", "report")
		if status != 1 || !strings.Contains(stderr, "BASE_SHA") {
			test.Errorf("dft gates report without BASE_SHA = %d, %s", status, stderr)
		}
	})

	test.Run("fails when git cannot read the range", func(test *testing.T) {
		test.Setenv("BASE_SHA", strings.Repeat("0", 40))
		test.Setenv("HEAD_SHA", head)
		test.Setenv("GITHUB_STEP_SUMMARY", "")

		status, stdout, stderr := consumer.dft("", "gates", "report")
		if status != 1 || stdout != "" || !strings.Contains(stderr, "commit range") {
			test.Errorf("dft gates report on an unknown SHA = %d, stdout %q, %s", status, stdout, stderr)
		}
	})
}
