// Package process runs native tools directly, never through a shell, with
// bounded output capture and cancellation that reaches the whole process group.
package process

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"
)

const (
	// captureLimit bounds how much of each output stream is retained.
	captureLimit = 1 << 20

	// pipeDrainDelay bounds how long Run waits for pipes after the process exits.
	pipeDrainDelay = 2 * time.Second
)

// Command describes one direct executable invocation. Args are never shell
// text. Env adds explicit overrides to the inherited environment. Stdout and
// Stderr, when set, receive the streams as they are produced.
type Command struct {
	Name string
	Args []string
	Dir  string
	Env  []string
	// UnsetEnv removes inherited credentials before explicit Env overrides.
	UnsetEnv []string
	Stdin    io.Reader
	Stdout   io.Writer
	Stderr   io.Writer
}

// Result records a completed process. A negative ExitCode means it did not
// exit normally.
type Result struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

// Local runs commands on the macOS or Linux host. Its zero value is ready.
type Local struct{}

// Run waits for the process and its pipes. Cancellation kills the owned
// process group, including subprocesses that Terraform or kubectl start.
// Captured-only output that overflows is an error, so callers never parse a
// truncated document as complete.
func (Local) Run(ctx context.Context, command Command) (Result, error) {
	result := Result{ExitCode: -1}

	if command.Name == "" {
		return result, errors.New("executable name is required")
	}

	if err := ctx.Err(); err != nil {
		return result, err
	}

	stdout := &capture{stream: command.Stdout}
	stderr := &capture{stream: command.Stderr}

	executable := exec.CommandContext(ctx, command.Name, command.Args...)
	executable.Dir = command.Dir
	executable.Env = commandEnvironment(command)
	executable.Stdin, executable.Stdout, executable.Stderr = command.Stdin, stdout, stderr
	executable.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	executable.WaitDelay = pipeDrainDelay
	executable.Cancel = func() error { return killGroup(executable.Process) }

	runErr := executable.Run()

	// A short-lived tool must not leave children holding its pipes.
	_ = killGroup(executable.Process)

	result.Stdout, result.Stderr = stdout.String(), stderr.String()
	if executable.ProcessState != nil {
		result.ExitCode = executable.ProcessState.ExitCode()
	}

	switch {
	case ctx.Err() != nil:
		return result, ctx.Err()

	case runErr != nil:
		return result, fmt.Errorf("%s failed: %w", command.Name, runErr)

	case (stdout.overflow && command.Stdout == nil) || (stderr.overflow && command.Stderr == nil):
		return result, fmt.Errorf("%s output exceeded the %d-byte capture limit", command.Name, captureLimit)
	}

	return result, nil
}

// Explain wraps a tool failure with what was being done and the tool's own
// diagnostics, when it printed any.
func Explain(action string, result Result, err error) error {
	if diagnostics := strings.TrimSpace(result.Stderr); diagnostics != "" {
		return fmt.Errorf("%s: %w: %s", action, err, diagnostics)
	}

	return fmt.Errorf("%s: %w", action, err)
}

func killGroup(process *os.Process) error {
	if process == nil {
		return nil
	}

	err := syscall.Kill(-process.Pid, syscall.SIGKILL)
	if errors.Is(err, syscall.ESRCH) {
		return os.ErrProcessDone
	}

	return err
}

// capture retains a bounded copy of a stream and optionally forwards it.
type capture struct {
	mutex    sync.Mutex
	data     bytes.Buffer
	overflow bool
	stream   io.Writer
}

func (capture *capture) Write(data []byte) (int, error) {
	capture.mutex.Lock()
	defer capture.mutex.Unlock()

	remaining := captureLimit - capture.data.Len()
	if remaining > 0 {
		_, _ = capture.data.Write(data[:min(remaining, len(data))])
	}

	if len(data) > remaining {
		capture.overflow = true
	}

	if capture.stream != nil {
		return capture.stream.Write(data)
	}

	return len(data), nil
}

func (capture *capture) String() string {
	capture.mutex.Lock()
	defer capture.mutex.Unlock()

	return capture.data.String()
}

func commandEnvironment(command Command) []string {
	excluded := map[string]bool{}
	for _, name := range command.UnsetEnv {
		excluded[name] = true
	}
	environment := make([]string, 0, len(os.Environ())+len(command.Env))
	for _, entry := range os.Environ() {
		name, _, _ := strings.Cut(entry, "=")
		if !excluded[name] {
			environment = append(environment, entry)
		}
	}
	return append(environment, command.Env...)
}
