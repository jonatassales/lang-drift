// Package delivery observes what the client's harness delivered, through the
// repository's own tools: git for the branch and the push, the forge's CLI
// for the pull request and its checks. It changes nothing, and anything it
// cannot observe stays unknown with the reason.
package delivery

import (
	"context"
	"encoding/json"
	"errors"
	"os/exec"
	"strconv"
	"strings"

	"github.com/jonatassales/tao/packages/ex/src/application/execution"
	"github.com/jonatassales/tao/packages/ex/src/domain/run"
	"github.com/jonatassales/tao/packages/ex/src/infrastructure/process"
)

// Runner is the direct native execution contract the observer consumes.
type Runner interface {
	Run(context.Context, process.Command) (process.Result, error)
}

// Observer runs git and gh inside the project.
type Observer struct{ Runner Runner }

// Branch is the checked-out branch, its short commit and whether the upstream
// has it.
func (observer Observer) Branch(ctx context.Context, root string) execution.Branch {
	name, err := observer.run(ctx, root, "git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return execution.Branch{Reason: "not a git repository"}
	}

	commit, err := observer.run(ctx, root, "git", "rev-parse", "--short", "HEAD")
	if err != nil {
		return execution.Branch{Name: name, Reason: "no commit"}
	}

	branch := execution.Branch{Name: name, Commit: commit}
	if _, err := observer.run(ctx, root, "git", "rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{upstream}"); err != nil {
		branch.Reason = "no upstream"

		return branch
	}

	ahead, err := observer.run(ctx, root, "git", "rev-list", "--count", "@{upstream}..HEAD")
	if err != nil || ahead != "0" {
		branch.Reason = ahead + " commits ahead of the upstream"

		return branch
	}

	branch.Upstream = true

	return branch
}

// PullRequest is the pull request of a branch, as gh sees it.
func (observer Observer) PullRequest(ctx context.Context, root, branch string) (*run.PullRequest, string) {
	output, err := observer.run(ctx, root, "gh", "pr", "view", branch, "--json", "number,title,state")
	if errors.Is(err, exec.ErrNotFound) {
		return nil, "gh is not installed"
	}
	if err != nil {
		return nil, "no pull request for " + branch
	}

	var view struct {
		Number int    `json:"number"`
		Title  string `json:"title"`
		State  string `json:"state"`
	}
	if err := json.Unmarshal([]byte(output), &view); err != nil || view.Number == 0 {
		return nil, "gh answered with no pull request"
	}

	return &run.PullRequest{Number: view.Number, Title: view.Title, State: strings.ToLower(view.State)}, ""
}

// Checks are the pull request's checks, as gh buckets them. gh exits with a
// code of its own while checks are pending or failed, so its output counts
// whenever it has some.
func (observer Observer) Checks(ctx context.Context, root string, number int) ([]run.CheckResult, string) {
	result, err := observer.Runner.Run(ctx, process.Command{Name: "gh", Args: []string{"pr", "checks", strconv.Itoa(number), "--json", "name,bucket"}, Dir: root})
	if errors.Is(err, exec.ErrNotFound) {
		return nil, "gh is not installed"
	}

	var buckets []struct {
		Name   string `json:"name"`
		Bucket string `json:"bucket"`
	}
	if decodeErr := json.Unmarshal([]byte(strings.TrimSpace(result.Stdout)), &buckets); decodeErr != nil {
		if err != nil {
			return nil, "gh could not list the checks"
		}

		return nil, "gh answered with no checks"
	}
	if len(buckets) == 0 {
		return nil, "no checks"
	}

	checks := make([]run.CheckResult, 0, len(buckets))
	for _, bucket := range buckets {
		checks = append(checks, run.CheckResult{Name: bucket.Name, Result: resultOf(bucket.Bucket)})
	}

	return checks, ""
}

func resultOf(bucket string) string {
	switch bucket {
	case "pass":
		return execution.CheckPassed
	case "fail":
		return execution.CheckFailed
	case "pending":
		return execution.CheckPending
	case "skipping":
		return "skipped"
	case "cancel":
		return "cancelled"
	}

	return bucket
}

// run runs one tool inside the project and returns what it printed, trimmed.
func (observer Observer) run(ctx context.Context, root, tool string, args ...string) (string, error) {
	result, err := observer.Runner.Run(ctx, process.Command{Name: tool, Args: args, Dir: root})
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result.Stdout), nil
}
