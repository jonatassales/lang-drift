package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/exphere/ex/src/domain/microruntime"
	"github.com/exphere/ex/src/domain/project"
)

const hookTimeout = 5 * time.Minute

// HookScript is a verified regular .sh file and its runtime working directory.
// The file resolver refuses paths outside runtime configuration and symlinks.
type HookScript struct {
	Path      string
	Directory string
}

// HookResolver resolves the script named by an already validated descriptor.
type HookResolver interface {
	ResolveHook(owner project.Project, runtime microruntime.Identifier, script microruntime.RelativePath) (HookScript, error)
}

// HookInvocation supplies the declared script and its explicit correlation.
type HookInvocation struct {
	Script      HookScript
	ProjectRoot string
	Runtime     string
	RunID       string
}

// HookResult preserves the native process result exposed by the CLI. Its
// capitalized JSON field names are the existing public result contract.
type HookResult struct {
	Stdout   string `json:"Stdout"`
	Stderr   string `json:"Stderr"`
	ExitCode int    `json:"ExitCode"`
}

// HookExecutor runs a declared trusted script and waits for its processes.
// Hooks are not sandboxed; captured output must be bounded.
type HookExecutor interface {
	Run(ctx context.Context, invocation HookInvocation) (HookResult, error)
}

// HookOutcome is retained evidence without script output or environment data.
// ExitCode is absent while running and present once the process has returned.
type HookOutcome struct {
	Runtime  string `json:"runtime"`
	RunID    string `json:"runId"`
	Hook     string `json:"hook"`
	Status   string `json:"status"`
	ExitCode *int   `json:"exitCode,omitempty"`
}

// HookEvidence records intent before execution and its final outcome afterward.
type HookEvidence interface {
	RecordHook(owner project.Project, outcome HookOutcome) error
}

// RunRuntimeHook explicitly runs one declared extension. Hooks never run as
// a side effect of another operation, and they do not require a live cluster.
type RunRuntimeHook struct {
	Projects      ProjectOpener
	Locks         ProjectLocker
	Configuration RuntimeReferences
	Scripts       HookResolver
	Executor      HookExecutor
	Evidence      HookEvidence
}

// RunRuntimeHookRequest identifies the declared extension to execute.
type RunRuntimeHookRequest struct {
	RuntimeRequest
	Name string
}

// Execute holds the project lock, records intent, and waits at most five
// minutes or the caller's shorter deadline. Final evidence is attempted even
// after interruption, without retaining potentially sensitive script output.
func (useCase RunRuntimeHook) Execute(ctx context.Context, request RunRuntimeHookRequest) (HookResult, error) {
	runtime, err := parseRuntime(request.Runtime)
	if err != nil {
		return HookResult{}, err
	}
	name, err := microruntime.ParseIdentifier(request.Name)
	if err != nil {
		return HookResult{}, fmt.Errorf("hook name: %w", err)
	}
	owner, err := useCase.Projects.Open(request.Project)
	if err != nil {
		return HookResult{}, err
	}
	unlock, err := useCase.Locks.Lock(owner)
	if err != nil {
		return HookResult{}, err
	}
	defer unlock()

	descriptor, err := loadRuntimeReferences(useCase.Configuration, owner, runtime)
	if err != nil {
		return HookResult{}, err
	}
	relative, declared := descriptor.Hooks[name]
	if !declared {
		return HookResult{}, fmt.Errorf("hook %q is not declared", name)
	}
	script, err := useCase.Scripts.ResolveHook(owner, runtime, relative)
	if err != nil {
		return HookResult{}, err
	}

	hookContext, cancel := context.WithTimeout(ctx, hookTimeout)
	defer cancel()
	runID, err := newID()
	if err != nil {
		return HookResult{}, err
	}
	outcome := HookOutcome{Runtime: runtime.String(), RunID: runID, Hook: name.String(), Status: "running"}
	if err := useCase.Evidence.RecordHook(owner, outcome); err != nil {
		return HookResult{}, err
	}

	result, operationError := useCase.Executor.Run(hookContext, HookInvocation{
		Script: script, ProjectRoot: owner.Root, Runtime: runtime.String(), RunID: runID,
	})
	outcome.Status = "completed"
	if operationError != nil {
		outcome.Status = "failed"
		if isInterruption(operationError) {
			outcome.Status = "interrupted"
		}
	}
	outcome.ExitCode = &result.ExitCode
	if err := useCase.Evidence.RecordHook(owner, outcome); err != nil {
		return result, errors.Join(operationError, fmt.Errorf("record hook outcome: %w", err))
	}

	return result, operationError
}
