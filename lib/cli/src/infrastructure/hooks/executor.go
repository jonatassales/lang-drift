// Package hooks executes explicitly selected consumer extensions. A hook is
// trusted shell code, not a sandboxed or automatically invoked lifecycle step.
package hooks

import (
	"context"

	"github.com/jonatassales/tao/lib/cli/src/application"
	"github.com/jonatassales/tao/lib/cli/src/infrastructure/process"
)

// Runner is the direct native execution contract this adapter consumes.
// Returning must wait for the process group and preserve caller cancellation.
type Runner interface {
	Run(context.Context, process.Command) (process.Result, error)
}

// Executor runs declared scripts through the installed sh, without building
// shell command strings or retaining raw output in operation evidence.
type Executor struct{ Runner Runner }

// Run waits for the script with bounded output capture and explicit runtime
// correlation. The caller owns its deadline and the project mutation lock.
// The script sees only the minimal process base plus the LANGDRIFT_*
// correlation variables; a hook that needs more will declare it in its
// configuration when a live hook surface returns (ADR 0002).
func (executor Executor) Run(ctx context.Context, invocation application.HookInvocation) (application.HookResult, error) {
	result, err := executor.Runner.Run(ctx, process.Command{
		Name: "sh",
		Args: []string{invocation.Script.Path},
		Dir:  invocation.Script.Directory,
		Env: []string{
			"LANGDRIFT_RUNTIME_ID=" + invocation.Runtime,
			"LANGDRIFT_RUN_ID=" + invocation.RunID,
			"LANGDRIFT_PROJECT_ROOT=" + invocation.ProjectRoot,
		},
	})

	return application.HookResult{Stdout: result.Stdout, Stderr: result.Stderr, ExitCode: result.ExitCode}, err
}
