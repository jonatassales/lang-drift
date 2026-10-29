// Package stages drives the owner's coding agent for a stage's skills and
// runs a stage's committed hooks. Both are trusted processes started
// directly, never through a shell command string, with the minimal
// environment plus what the driver passes through; their output goes to a
// private log, never into evidence.
package stages

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/jonatassales/tao/packages/ex/src/application/execution"
	"github.com/jonatassales/tao/packages/ex/src/domain/harness"
	"github.com/jonatassales/tao/packages/ex/src/infrastructure/files"
	"github.com/jonatassales/tao/packages/ex/src/infrastructure/process"
)

// Runner is the direct native execution contract the driver consumes.
type Runner interface {
	Run(context.Context, process.Command) (process.Result, error)
}

// Driver runs the stages' processes inside the project.
type Driver struct{ Runner Runner }

// The placeholders a driver's arguments may carry, substituted per argument.
const (
	skillPlaceholder  = "{{skill}}"
	ticketPlaceholder = "{{ticket}}"
	runPlaceholder    = "{{run}}"
	stagePlaceholder  = "{{stage}}"
	lapPlaceholder    = "{{lap}}"
	promptPlaceholder = "{{prompt}}"
)

// summaryLimit bounds the one line a stage reports.
const summaryLimit = 120

// Run runs one skill through the coding agent: the arguments with their
// placeholders filled, the LANGDRIFT_* variables and the driver's
// pass-through in the environment, the output streamed to a private log,
// and the usage and the summary read from the output where the driver says
// they are. A process that ran reports through its exit code.
func (driver Driver) Run(ctx context.Context, config harness.Driver, call execution.StageCall) (execution.StageResult, error) {
	log, logPath, err := openLog(call.ProjectRoot, call.Runtime, call.Lap+"-"+call.Stage+"-"+call.Skill)
	if err != nil {
		return execution.StageResult{}, err
	}
	defer log.Close()

	bounded, cancel := deadline(ctx, config.TimeoutSeconds)
	defer cancel()

	output := process.SyncWriter(log)
	started := time.Now()
	result, runErr := driver.Runner.Run(bounded, process.Command{
		Name:    config.Executable,
		Args:    substitute(config.Args, call),
		Dir:     call.ProjectRoot,
		Env:     environment(call),
		PassEnv: config.PassEnv,
		Stdout:  output,
		Stderr:  output,
	})

	stageResult := execution.StageResult{ExitCode: result.ExitCode, LogPath: logPath, ElapsedMs: time.Since(started).Milliseconds()}

	if result.ExitCode < 0 {
		if errors.Is(bounded.Err(), context.DeadlineExceeded) {
			return stageResult, fmt.Errorf("%s did not finish within %d seconds; see %s", config.Executable, config.TimeoutSeconds, logPath)
		}

		return stageResult, fmt.Errorf("%s did not run: %w", config.Executable, runErr)
	}

	stageResult.InputTokens, stageResult.OutputTokens, stageResult.Summary = read(result.Stdout, config)

	return stageResult, nil
}

// Hook runs one committed script of a stage through sh, from the project
// root, with the same variables; the script must exist inside the project.
func (driver Driver) Hook(ctx context.Context, call execution.HookCall) (int, error) {
	script, err := files.Existing(call.ProjectRoot, call.Script)
	if err != nil {
		return -1, fmt.Errorf("hook %s: %w", call.Script, err)
	}

	name := strings.TrimSuffix(filepath.Base(script), filepath.Ext(script))
	log, _, err := openLog(call.ProjectRoot, call.Runtime, call.Lap+"-"+call.Stage+"-hook-"+name)
	if err != nil {
		return -1, err
	}
	defer log.Close()

	output := process.SyncWriter(log)
	result, runErr := driver.Runner.Run(ctx, process.Command{
		Name:   "sh",
		Args:   []string{script},
		Dir:    call.ProjectRoot,
		Env:    hookEnvironment(call),
		Stdout: output,
		Stderr: output,
	})
	if result.ExitCode < 0 {
		return result.ExitCode, fmt.Errorf("hook %s did not run: %w", call.Script, runErr)
	}

	return result.ExitCode, nil
}

// Prompt is what {{prompt}} expands to: the skill to run, for which ticket,
// and the one line the run reads back.
func Prompt(call execution.StageCall) string {
	return "Run the " + call.Skill + " skill for " + call.Ticket + ": stage " + call.Stage + " of run " + call.Runtime + ", lap " + call.Lap + ". End with one line that says what was done."
}

// deadline bounds a call by the driver's timeout, when it has one.
func deadline(ctx context.Context, timeoutSeconds int) (context.Context, context.CancelFunc) {
	if timeoutSeconds <= 0 {
		return context.WithCancel(ctx)
	}

	return context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
}

// openLog opens the private log of one call, replacing an earlier one.
func openLog(root, runtime, name string) (*os.File, string, error) {
	relative := filepath.Join(".langdrift", "agent", "state", "runs", runtime, "stages", name+".log")

	path, err := files.Within(root, relative)
	if err != nil {
		return nil, "", err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil, "", err
	}

	log, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return nil, "", err
	}

	return log, filepath.ToSlash(relative), nil
}

// substitute fills the placeholders of each argument; an argument is never
// split, so a value with spaces stays one argument.
func substitute(args []string, call execution.StageCall) []string {
	replacer := strings.NewReplacer(
		skillPlaceholder, call.Skill,
		ticketPlaceholder, call.Ticket,
		runPlaceholder, call.Runtime,
		stagePlaceholder, call.Stage,
		lapPlaceholder, call.Lap,
		promptPlaceholder, Prompt(call),
	)

	substituted := make([]string, 0, len(args))
	for _, arg := range args {
		substituted = append(substituted, replacer.Replace(arg))
	}

	return substituted
}

func environment(call execution.StageCall) []string {
	variables := []string{
		"LANGDRIFT_PROJECT_ROOT=" + call.ProjectRoot,
		"LANGDRIFT_RUNTIME_ID=" + call.Runtime,
		"LANGDRIFT_LAP=" + call.Lap,
		"LANGDRIFT_STAGE=" + call.Stage,
		"LANGDRIFT_SKILL=" + call.Skill,
		"LANGDRIFT_TICKET=" + call.Ticket,
	}
	if call.Execution != "" {
		variables = append(variables, "LANGDRIFT_EXECUTION_ID="+call.Execution)
	}

	return variables
}

func hookEnvironment(call execution.HookCall) []string {
	return []string{
		"LANGDRIFT_PROJECT_ROOT=" + call.ProjectRoot,
		"LANGDRIFT_RUNTIME_ID=" + call.Runtime,
		"LANGDRIFT_LAP=" + call.Lap,
		"LANGDRIFT_STAGE=" + call.Stage,
		"LANGDRIFT_HOOK=" + call.Script,
	}
}

// read takes the usage and the summary out of the output. A JSON object is
// read where the driver says the fields are; any other output reports no
// usage and its last line as the summary.
func read(output string, config harness.Driver) (*int64, *int64, string) {
	document, isObject := decode(output)
	if !isObject {
		return nil, nil, line(lastLine(output))
	}

	return integerAt(document, config.Usage.InputTokens), integerAt(document, config.Usage.OutputTokens), line(lastLine(stringAt(document, config.Summary)))
}

func decode(output string) (map[string]any, bool) {
	decoder := json.NewDecoder(strings.NewReader(output))
	decoder.UseNumber()

	var document map[string]any
	if err := decoder.Decode(&document); err != nil || document == nil {
		return nil, false
	}

	return document, true
}

// valueAt walks a dotted path into the document; an empty path reads nothing.
func valueAt(document map[string]any, path string) (any, bool) {
	if path == "" {
		return nil, false
	}

	var current any = document
	for _, segment := range strings.Split(path, ".") {
		object, isObject := current.(map[string]any)
		if !isObject {
			return nil, false
		}

		next, found := object[segment]
		if !found {
			return nil, false
		}
		current = next
	}

	return current, true
}

func integerAt(document map[string]any, path string) *int64 {
	value, found := valueAt(document, path)
	if !found {
		return nil
	}

	switch typed := value.(type) {
	case json.Number:
		if whole, err := typed.Int64(); err == nil {
			return &whole
		}
		if real, err := typed.Float64(); err == nil {
			whole := int64(real)

			return &whole
		}
	case string:
		if whole, err := strconv.ParseInt(strings.TrimSpace(typed), 10, 64); err == nil {
			return &whole
		}
	}

	return nil
}

func stringAt(document map[string]any, path string) string {
	value, found := valueAt(document, path)
	if !found {
		return ""
	}

	switch typed := value.(type) {
	case string:
		return typed
	case json.Number:
		return typed.String()
	}

	return ""
}

// line bounds one line of text.
func line(text string) string {
	first, _, _ := strings.Cut(strings.TrimSpace(text), "\n")
	first = strings.TrimSpace(first)
	if utf8.RuneCountInString(first) <= summaryLimit {
		return first
	}

	runes := []rune(first)

	return string(runes[:summaryLimit-1]) + "…"
}

// lastLine is the last line of the output that says anything.
func lastLine(output string) string {
	lines := strings.Split(output, "\n")
	for index := len(lines) - 1; index >= 0; index-- {
		if trimmed := strings.TrimSpace(lines[index]); trimmed != "" {
			return trimmed
		}
	}

	return ""
}
