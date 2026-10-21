package hooks_test

import (
	"context"
	"errors"
	"slices"
	"testing"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/infrastructure/hooks"
	"github.com/exphere/ex/src/infrastructure/process"
)

type recorder struct {
	command process.Command
	result  process.Result
	err     error
}

func (recorder *recorder) Run(_ context.Context, command process.Command) (process.Result, error) {
	recorder.command = command

	return recorder.result, recorder.err
}

func TestExecutor_Run(t *testing.T) {
	t.Parallel()

	invocation := application.HookInvocation{
		Script:      application.HookScript{Path: "/workspace/.exphere/runtimes/demo/hooks/seed.sh", Directory: "/workspace/.exphere/runtimes/demo"},
		ProjectRoot: "/workspace",
		Runtime:     "demo",
		RunID:       "run-1",
	}

	runner := &recorder{result: process.Result{Stdout: "seeded", Stderr: "warning", ExitCode: 0}}

	result, err := hooks.Executor{Runner: runner}.Run(context.Background(), invocation)
	if err != nil || result != (application.HookResult{Stdout: "seeded", Stderr: "warning", ExitCode: 0}) {
		t.Fatalf("Run() = %+v, %v", result, err)
	}

	command := runner.command
	if command.Name != "sh" || !slices.Equal(command.Args, []string{invocation.Script.Path}) || command.Dir != invocation.Script.Directory {
		t.Errorf("command = %+v, want sh running the script from the runtime directory", command)
	}

	wantEnvironment := []string{"EXPHERE_RUNTIME_ID=demo", "EXPHERE_RUN_ID=run-1", "EXPHERE_PROJECT_ROOT=/workspace"}
	if !slices.Equal(command.Env, wantEnvironment) {
		t.Errorf("environment = %v, want %v", command.Env, wantEnvironment)
	}

	failing := &recorder{result: process.Result{Stderr: "boom", ExitCode: 3}, err: errors.New("exit status 3")}

	result, err = hooks.Executor{Runner: failing}.Run(context.Background(), invocation)
	if err == nil || result.ExitCode != 3 || result.Stderr != "boom" {
		t.Fatalf("Run() = %+v, %v; want the native result and failure", result, err)
	}
}
