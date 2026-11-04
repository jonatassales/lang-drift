package delivery_test

import (
	"context"
	"errors"
	"os/exec"
	"strings"
	"testing"

	"github.com/jonatassales/tao/lib/cli/src/infrastructure/delivery"
	"github.com/jonatassales/tao/lib/cli/src/infrastructure/process"
)

// scripted answers each command line with a canned result.
type scripted struct {
	outputs  map[string]string
	failures map[string]error
	commands []string
	dir      string
}

func (scripted *scripted) Run(_ context.Context, command process.Command) (process.Result, error) {
	key := command.Name + " " + strings.Join(command.Args, " ")
	scripted.commands = append(scripted.commands, key)
	scripted.dir = command.Dir
	if err, failed := scripted.failures[key]; failed {
		return process.Result{Stdout: scripted.outputs[key], ExitCode: 1}, err
	}

	return process.Result{Stdout: scripted.outputs[key]}, nil
}

var notInstalled = &exec.Error{Name: "gh", Err: exec.ErrNotFound}

func TestBranchSaysWhetherTheCommitReachedTheUpstream(t *testing.T) {
	t.Parallel()

	runner := &scripted{outputs: map[string]string{
		"git rev-parse --abbrev-ref HEAD":                             "feat/0070-console-microruntime\n",
		"git rev-parse --short HEAD":                                  "3f9c2e1\n",
		"git rev-parse --abbrev-ref --symbolic-full-name @{upstream}": "origin/feat/0070-console-microruntime\n",
		"git rev-list --count @{upstream}..HEAD":                      "0\n",
	}, failures: map[string]error{}}

	branch := delivery.Observer{Runner: runner}.Branch(context.Background(), "/workspace")
	if branch.Name != "feat/0070-console-microruntime" || branch.Commit != "3f9c2e1" || !branch.Upstream || branch.Reason != "" || runner.dir != "/workspace" {
		t.Errorf("Branch() = %+v in %s", branch, runner.dir)
	}

	runner.outputs["git rev-list --count @{upstream}..HEAD"] = "2\n"
	if branch := (delivery.Observer{Runner: runner}).Branch(context.Background(), "/workspace"); branch.Upstream || branch.Reason != "2 commits ahead of the upstream" {
		t.Errorf("a branch ahead of its upstream = %+v", branch)
	}

	runner.failures["git rev-list --count @{upstream}..HEAD"] = errors.New("git failed: exit status 128")
	if branch := (delivery.Observer{Runner: runner}).Branch(context.Background(), "/workspace"); branch.Upstream || branch.Reason != "the upstream could not be compared" {
		t.Errorf("an upstream that cannot be compared = %+v", branch)
	}
	delete(runner.failures, "git rev-list --count @{upstream}..HEAD")

	runner.failures["git rev-parse --abbrev-ref --symbolic-full-name @{upstream}"] = errors.New("git failed: exit status 128")
	if branch := (delivery.Observer{Runner: runner}).Branch(context.Background(), "/workspace"); branch.Upstream || branch.Reason != "no upstream" || branch.Name == "" {
		t.Errorf("a branch without an upstream = %+v", branch)
	}

	runner.failures["git rev-parse --abbrev-ref HEAD"] = errors.New("git failed: exit status 128")
	if branch := (delivery.Observer{Runner: runner}).Branch(context.Background(), "/workspace"); branch.Name != "" || branch.Reason != "not a git repository" {
		t.Errorf("no repository = %+v", branch)
	}
}

func TestPullRequestReadsTheForge(t *testing.T) {
	t.Parallel()

	runner := &scripted{outputs: map[string]string{
		"gh pr view feat/0070-console-microruntime --json number,title,state": `{"number":19,"title":"feat(console): run the console preview as a microruntime","state":"OPEN"}` + "\n",
	}, failures: map[string]error{}}

	pullRequest, reason := delivery.Observer{Runner: runner}.PullRequest(context.Background(), "/workspace", "feat/0070-console-microruntime")
	if pullRequest == nil || pullRequest.Number != 19 || pullRequest.State != "open" || !strings.HasPrefix(pullRequest.Title, "feat(console)") || reason != "" {
		t.Errorf("PullRequest() = %+v, %q", pullRequest, reason)
	}

	runner.failures["gh pr view feat/0070-console-microruntime --json number,title,state"] = errors.New("gh failed: exit status 1")
	if pullRequest, reason := (delivery.Observer{Runner: runner}).PullRequest(context.Background(), "/workspace", "feat/0070-console-microruntime"); pullRequest != nil || reason != "no pull request for feat/0070-console-microruntime" {
		t.Errorf("no pull request = %+v, %q", pullRequest, reason)
	}

	runner.failures["gh pr view feat/0070-console-microruntime --json number,title,state"] = notInstalled
	if pullRequest, reason := (delivery.Observer{Runner: runner}).PullRequest(context.Background(), "/workspace", "feat/0070-console-microruntime"); pullRequest != nil || reason != "gh is not installed" {
		t.Errorf("no gh = %+v, %q", pullRequest, reason)
	}
}

func TestChecksReadTheBucketsEvenWhileGhExitsForPendingOnes(t *testing.T) {
	t.Parallel()

	runner := &scripted{outputs: map[string]string{
		"gh pr checks 19 --json name,bucket": `[{"name":"build","bucket":"pass"},{"name":"check","bucket":"fail"},{"name":"test","bucket":"pending"},{"name":"guardrails","bucket":"skipping"}]` + "\n",
	}, failures: map[string]error{"gh pr checks 19 --json name,bucket": errors.New("gh failed: exit status 8")}}

	checks, reason := delivery.Observer{Runner: runner}.Checks(context.Background(), "/workspace", 19)
	if reason != "" || len(checks) != 4 || checks[0].Result != "passed" || checks[1].Result != "failed" || checks[2].Result != "pending" || checks[3].Result != "skipped" || checks[3].Name != "guardrails" {
		t.Errorf("Checks() = %+v, %q", checks, reason)
	}

	runner.outputs["gh pr checks 19 --json name,bucket"] = "[]\n"
	delete(runner.failures, "gh pr checks 19 --json name,bucket")
	if checks, reason := (delivery.Observer{Runner: runner}).Checks(context.Background(), "/workspace", 19); checks != nil || reason != "no checks" {
		t.Errorf("no checks = %+v, %q", checks, reason)
	}

	runner.failures["gh pr checks 19 --json name,bucket"] = notInstalled
	if checks, reason := (delivery.Observer{Runner: runner}).Checks(context.Background(), "/workspace", 19); checks != nil || reason != "gh is not installed" {
		t.Errorf("no gh = %+v, %q", checks, reason)
	}
}
