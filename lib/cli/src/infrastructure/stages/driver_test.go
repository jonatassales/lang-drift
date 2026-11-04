package stages_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/jonatassales/tao/lib/cli/src/application/execution"
	"github.com/jonatassales/tao/lib/cli/src/domain/harness"
	"github.com/jonatassales/tao/lib/cli/src/infrastructure/process"
	"github.com/jonatassales/tao/lib/cli/src/infrastructure/stages"
)

// recorder is a runner that keeps the command and plays a scripted result,
// writing the scripted output to the streams the command set.
type recorder struct {
	command  process.Command
	deadline time.Time
	result   process.Result
	err      error
}

func (recorder *recorder) Run(ctx context.Context, command process.Command) (process.Result, error) {
	recorder.command = command
	recorder.deadline, _ = ctx.Deadline()
	if command.Stdout != nil {
		_, _ = io.WriteString(command.Stdout, recorder.result.Stdout)
	}
	if command.Stderr != nil {
		_, _ = io.WriteString(command.Stderr, recorder.result.Stderr)
	}

	return recorder.result, recorder.err
}

func claudeDriver() harness.Driver {
	return harness.Driver{
		Executable:     "claude",
		Args:           []string{"-p", "{{prompt}}", "--output-format", "json", "--append-system-prompt", "lap {{lap}} of {{run}}"},
		PassEnv:        []string{"ANTHROPIC_API_KEY", "CLAUDE_*"},
		TimeoutSeconds: 60,
		Usage:          harness.Usage{InputTokens: "usage.input_tokens", OutputTokens: "usage.output_tokens"},
		Summary:        "result",
		Unmeasured:     "pause",
	}
}

func call(root string) execution.StageCall {
	return execution.StageCall{ProjectRoot: root, Runtime: "console-microruntime", Lap: "0001", Stage: "planning", Skill: "grilling", Ticket: ".agents/workflow/tickets/0070-console-microruntime.ticket.md"}
}

func TestDriverRunsTheSkillAndReadsTheResult(t *testing.T) {
	t.Parallel()

	root := realTempDir(t)
	runner := &recorder{result: process.Result{Stdout: "{\"type\":\"result\",\"result\":\"Grilled the ticket.\\n\\n3 open questions settled\",\"usage\":{\"input_tokens\":300000,\"output_tokens\":20000},\"total_cost_usd\":0.8}\n", ExitCode: 0}}

	result, err := stages.Driver{Runner: runner}.Run(context.Background(), claudeDriver(), call(root))
	if err != nil {
		t.Fatal(err)
	}

	if result.ExitCode != 0 || result.InputTokens == nil || *result.InputTokens != 300000 || result.OutputTokens == nil || *result.OutputTokens != 20000 || result.Summary != "3 open questions settled" {
		t.Errorf("result = %+v", result)
	}

	command := runner.command
	prompt := stages.Prompt(call(root))
	wantArgs := []string{"-p", prompt, "--output-format", "json", "--append-system-prompt", "lap 0001 of console-microruntime"}
	if command.Name != "claude" || !slices.Equal(command.Args, wantArgs) || command.Dir != root {
		t.Errorf("command = %+v, want claude with the placeholders filled, from the project root", command)
	}
	if !strings.Contains(prompt, "grilling") || !strings.Contains(prompt, "0070-console-microruntime.ticket.md") {
		t.Errorf("prompt = %q", prompt)
	}
	if !slices.Equal(command.PassEnv, []string{"ANTHROPIC_API_KEY", "CLAUDE_*"}) {
		t.Errorf("passEnv = %v", command.PassEnv)
	}
	if left := time.Until(runner.deadline); left <= 50*time.Second || left > 60*time.Second {
		t.Errorf("the call had %s left, want the driver's 60-second timeout", left)
	}
	wantEnvironment := []string{"LANGDRIFT_PROJECT_ROOT=" + root, "LANGDRIFT_RUNTIME_ID=console-microruntime", "LANGDRIFT_LAP=0001", "LANGDRIFT_STAGE=planning", "LANGDRIFT_SKILL=grilling", "LANGDRIFT_TICKET=.agents/workflow/tickets/0070-console-microruntime.ticket.md"}
	if !slices.Equal(command.Env, wantEnvironment) {
		t.Errorf("environment = %v, want %v", command.Env, wantEnvironment)
	}

	if result.LogPath != ".langdrift/agent/state/runs/console-microruntime/stages/0001-planning-grilling.log" {
		t.Errorf("log path = %q", result.LogPath)
	}
	logged, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(result.LogPath)))
	if err != nil || !strings.Contains(string(logged), "3 open questions settled") {
		t.Errorf("the log holds %q, %v; want the agent's output", logged, err)
	}
}

func TestDriverReportsWhatItCannotRead(t *testing.T) {
	t.Parallel()

	root := realTempDir(t)

	plain := &recorder{result: process.Result{Stdout: "working...\nDone: 2 findings fixed\n\n", ExitCode: 0}}
	result, err := stages.Driver{Runner: plain}.Run(context.Background(), claudeDriver(), call(root))
	if err != nil || result.InputTokens != nil || result.OutputTokens != nil || result.Summary != "Done: 2 findings fixed" {
		t.Errorf("plain output gave %+v, %v; want no usage and the last line", result, err)
	}

	failing := &recorder{result: process.Result{Stdout: "{\"result\":\"it broke\",\"usage\":{\"input_tokens\":\"12\"}}", ExitCode: 1}, err: errors.New("claude failed: exit status 1")}
	result, err = stages.Driver{Runner: failing}.Run(context.Background(), claudeDriver(), call(root))
	if err != nil || result.ExitCode != 1 || result.InputTokens == nil || *result.InputTokens != 12 || result.OutputTokens != nil || result.Summary != "it broke" {
		t.Errorf("a failing skill gave %+v, %v; want its exit code and what it reported", result, err)
	}

	missing := stages.Driver{Runner: &recorder{result: process.Result{ExitCode: -1}, err: fmt.Errorf("claude failed: %w", &exec.Error{Name: "claude", Err: exec.ErrNotFound})}}
	if _, err := missing.Run(context.Background(), claudeDriver(), call(root)); !errors.Is(err, execution.ErrNotStarted) || !strings.Contains(err.Error(), "did not run") {
		t.Errorf("a missing executable gave %v, want ErrNotStarted", err)
	}

	killed := stages.Driver{Runner: &recorder{result: process.Result{ExitCode: -1}, err: errors.New("signal: killed")}}
	if _, err := killed.Run(context.Background(), claudeDriver(), call(root)); err == nil || errors.Is(err, execution.ErrNotStarted) || !strings.Contains(err.Error(), "unknown") {
		t.Errorf("a killed agent gave %v, want its use unknown", err)
	}
}

func TestDriverRunsACommittedHook(t *testing.T) {
	t.Parallel()

	root := realTempDir(t)
	if err := os.MkdirAll(filepath.Join(root, "platform", "hooks"), 0o755); err != nil {
		t.Fatal(err)
	}
	script := filepath.Join(root, "platform", "hooks", "lint.sh")
	if err := os.WriteFile(script, []byte("#!/bin/sh\necho linted\n"), 0o755); err != nil {
		t.Fatal(err)
	}

	runner := &recorder{result: process.Result{Stdout: "linted\n", ExitCode: 0}}
	hook := execution.HookCall{ProjectRoot: root, Runtime: "console-microruntime", Lap: "0001", Stage: "implement", Script: "platform/hooks/lint.sh"}

	exitCode, err := stages.Driver{Runner: runner}.Hook(context.Background(), hook)
	if err != nil || exitCode != 0 {
		t.Fatalf("Hook() = %d, %v", exitCode, err)
	}

	command := runner.command
	if command.Name != "sh" || !slices.Equal(command.Args, []string{script}) || command.Dir != root {
		t.Errorf("command = %+v, want sh running the script from the project root", command)
	}
	if !slices.Contains(command.Env, "LANGDRIFT_HOOK=platform/hooks/lint.sh") || !slices.Contains(command.Env, "LANGDRIFT_STAGE=implement") {
		t.Errorf("environment = %v", command.Env)
	}

	driver := stages.Driver{Runner: runner}
	if _, err := driver.Hook(context.Background(), execution.HookCall{ProjectRoot: root, Runtime: "console-microruntime", Lap: "0001", Stage: "implement", Script: "platform/hooks/missing.sh"}); err == nil {
		t.Error("a hook that does not exist ran")
	}

	outside := filepath.Join(filepath.Dir(root), "outside.sh")
	if _, err := driver.Hook(context.Background(), execution.HookCall{ProjectRoot: root, Runtime: "console-microruntime", Lap: "0001", Stage: "implement", Script: outside}); err == nil {
		t.Error("a hook outside the project ran")
	}
}

// realTempDir is a temporary directory with its symlinks resolved, since
// the driver refuses symlinked paths.
func realTempDir(t *testing.T) string {
	t.Helper()

	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	return root
}
