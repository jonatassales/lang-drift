package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/domain/microruntime"
)

type fakeHookExecutor struct {
	result      application.HookResult
	err         error
	invocations []application.HookInvocation
}

func (executor *fakeHookExecutor) Run(ctx context.Context, invocation application.HookInvocation) (application.HookResult, error) {
	executor.invocations = append(executor.invocations, invocation)

	if _, bounded := ctx.Deadline(); !bounded {
		return application.HookResult{}, errors.New("hook ran without a deadline")
	}

	return executor.result, executor.err
}

func (w *world) runHook(executor *fakeHookExecutor) application.RunRuntimeHook {
	return application.RunRuntimeHook{
		Projects: w.projects, Locks: w.locks, Configuration: w.configuration,
		Scripts: w.configuration, Executor: executor, Evidence: w.evidence,
	}
}

func (w *world) hookRequest(name string) application.RunRuntimeHookRequest {
	return application.RunRuntimeHookRequest{RuntimeRequest: w.runtimeRequest(), Name: name}
}

func TestRunRuntimeHook_RunsTheDeclaredScript(t *testing.T) {
	t.Parallel()

	w := newWorld(t)
	w.configuration.hook = application.HookScript{Path: "/workspace/.exphere/runtimes/demo/hooks/seed.sh", Directory: "/workspace/.exphere/runtimes/demo"}
	executor := &fakeHookExecutor{result: application.HookResult{Stdout: "seeded", ExitCode: 0}}

	result, err := w.runHook(executor).Execute(context.Background(), w.hookRequest("seed"))
	if err != nil || result.Stdout != "seeded" {
		t.Fatalf("Execute() = %+v, %v", result, err)
	}

	if len(executor.invocations) != 1 {
		t.Fatalf("ran %d hooks, want 1", len(executor.invocations))
	}

	invocation := executor.invocations[0]
	if invocation.Script != w.configuration.hook || invocation.ProjectRoot != projectRoot || invocation.Runtime != "demo" || invocation.RunID == "" {
		t.Errorf("invocation = %+v, want the resolved script and explicit correlation", invocation)
	}

	if len(w.configuration.resolved) != 1 || w.configuration.resolved[0].String() != "hooks/seed.sh" {
		t.Errorf("resolved %v, want the declared script", w.configuration.resolved)
	}

	hooks := w.evidence.hooks
	if len(hooks) != 2 || hooks[0].Status != "running" || hooks[0].ExitCode != nil ||
		hooks[1].Status != "completed" || hooks[1].ExitCode == nil || *hooks[1].ExitCode != 0 || hooks[1].RunID != invocation.RunID {
		t.Errorf("hook evidence = %+v, want intent before execution and the outcome after", hooks)
	}

	w.requireLockBalanced(t)
}

func TestRunRuntimeHook_RecordsFailures(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		err        error
		wantStatus string
	}{
		{name: "script fails", err: errInjected, wantStatus: "failed"},
		{name: "script interrupted", err: context.Canceled, wantStatus: "interrupted"},
		{name: "script exceeds its deadline", err: context.DeadlineExceeded, wantStatus: "interrupted"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := newWorld(t)
			executor := &fakeHookExecutor{result: application.HookResult{Stderr: "boom", ExitCode: 3}, err: tt.err}

			result, err := w.runHook(executor).Execute(context.Background(), w.hookRequest("seed"))
			if !errors.Is(err, tt.err) || result.ExitCode != 3 {
				t.Fatalf("Execute() = %+v, %v; want the native result and %v", result, err, tt.err)
			}

			final := w.evidence.hooks[len(w.evidence.hooks)-1]
			if final.Status != tt.wantStatus || final.ExitCode == nil || *final.ExitCode != 3 {
				t.Fatalf("final evidence = %+v, want %s with exit code 3", final, tt.wantStatus)
			}
		})
	}
}

func TestRunRuntimeHook_Refuses(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		hook    string
		arrange func(*world)
		wantErr error
	}{
		{name: "invalid hook name", hook: "Seed", arrange: func(*world) {}, wantErr: microruntime.ErrInvalidIdentifier},
		{name: "undeclared hook", hook: "migrate", arrange: func(*world) {}},
		{name: "unresolvable script", hook: "seed", arrange: func(w *world) { w.configuration.hookErr = errInjected }, wantErr: errInjected},
		{name: "broken references", hook: "seed", arrange: func(w *world) { w.configuration.referencesErr = errInjected }, wantErr: errInjected},
		{name: "lock held elsewhere", hook: "seed", arrange: func(w *world) { w.locks.err = errInjected }, wantErr: errInjected},
		{name: "intent cannot be recorded", hook: "seed", arrange: func(w *world) { w.evidence.failHookAt = 1 }, wantErr: errInjected},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := newWorld(t)
			tt.arrange(w)
			executor := &fakeHookExecutor{}

			_, err := w.runHook(executor).Execute(context.Background(), w.hookRequest(tt.hook))
			if err == nil || (tt.wantErr != nil && !errors.Is(err, tt.wantErr)) {
				t.Fatalf("Execute() error = %v, want %v", err, tt.wantErr)
			}

			if len(executor.invocations) != 0 {
				t.Error("a refused hook ran")
			}

			w.requireLockBalanced(t)
		})
	}
}

func TestRunRuntimeHook_ReportsAnUnrecordedOutcome(t *testing.T) {
	t.Parallel()

	w := newWorld(t)
	w.evidence.failHookAt = 2
	executor := &fakeHookExecutor{result: application.HookResult{Stdout: "seeded"}}

	result, err := w.runHook(executor).Execute(context.Background(), w.hookRequest("seed"))
	if !errors.Is(err, errInjected) || result.Stdout != "seeded" {
		t.Fatalf("Execute() = %+v, %v; want the native result with the evidence failure", result, err)
	}
}
