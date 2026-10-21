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
// text. A child never inherits the agent's environment: it gets the minimal
// base (see baseEnvironment), the parent values PassEnv names ("DOCKER_*"
// passes a prefix), and the literal Env entries, in that order. Stdout and
// Stderr, when set, receive the streams as they are produced.
type Command struct {
	Name string
	Args []string
	Dir  string
	// Env appends explicit NAME=value entries to the minimal base.
	Env []string
	// PassEnv copies named parent variables; a trailing * passes a prefix.
	PassEnv []string
	Stdin   io.Reader
	Stdout  io.Writer
	Stderr  io.Writer
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

// baseEnvironment is what every child may see without declaring anything:
// lookup, workspace and locale basics, and proxy settings for tools that
// reach a network. Credentials are never in it; see ADR 0002.
var baseEnvironment = []string{
	"PATH",
	"HOME",
	"TMPDIR",
	"USER",
	"LOGNAME",
	"LANG",
	"LC_*",
	"HTTP_PROXY",
	"HTTPS_PROXY",
	"NO_PROXY",
	"http_proxy",
	"https_proxy",
	"no_proxy",
}

func commandEnvironment(command Command) []string {
	patterns := append(append([]string{}, baseEnvironment...), command.PassEnv...)

	passes := func(name string) bool {
		for _, pattern := range patterns {
			if prefix, wildcard := strings.CutSuffix(pattern, "*"); wildcard {
				if strings.HasPrefix(name, prefix) {
					return true
				}
			} else if name == pattern {
				return true
			}
		}

		return false
	}

	environment := make([]string, 0, len(baseEnvironment)+len(command.Env))
	for _, entry := range os.Environ() {
		name, _, _ := strings.Cut(entry, "=")
		if passes(name) {
			environment = append(environment, entry)
		}
	}

	return append(environment, command.Env...)
}

// SyncWriter serializes writes from the stdout and stderr pipe goroutines, so
// both streams can share one destination, such as a log file.
func SyncWriter(destination io.Writer) io.Writer {
	return &syncWriter{destination: destination}
}

type syncWriter struct {
	mutex       sync.Mutex
	destination io.Writer
}

func (writer *syncWriter) Write(data []byte) (int, error) {
	writer.mutex.Lock()
	defer writer.mutex.Unlock()

	return writer.destination.Write(data)
}

// Tail returns the last lines of output that fit the budget, for error
// messages that must stay bounded while naming the real cause.
func Tail(output string, budget int) string {
	trimmed := strings.TrimSpace(output)
	if len(trimmed) <= budget {
		return trimmed
	}

	cut := trimmed[len(trimmed)-budget:]
	if newline := strings.IndexByte(cut, '\n'); newline >= 0 && newline < len(cut)-1 {
		cut = cut[newline+1:]
	}

	return cut
}
