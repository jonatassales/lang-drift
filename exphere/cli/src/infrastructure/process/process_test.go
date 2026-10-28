package process_test

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jonatassales/tao/exphere/cli/src/infrastructure/process"
)

func shell(script string) process.Command {
	return process.Command{Name: "sh", Args: []string{"-c", script}}
}

func TestLocal_Run(t *testing.T) {
	t.Parallel()

	t.Run("captures both streams and the exit code", func(t *testing.T) {
		t.Parallel()

		result, err := process.Local{}.Run(context.Background(), shell("printf out; printf err >&2; exit 3"))

		var exitError *exec.ExitError
		if !errors.As(err, &exitError) || !strings.Contains(err.Error(), "sh failed") {
			t.Fatalf("Run() error = %v, want the native exit error", err)
		}

		if result.Stdout != "out" || result.Stderr != "err" || result.ExitCode != 3 {
			t.Fatalf("Run() = %+v, want both streams and exit code 3", result)
		}
	})

	t.Run("succeeds with exit code zero", func(t *testing.T) {
		t.Parallel()

		result, err := process.Local{}.Run(context.Background(), shell("printf ok"))
		if err != nil || result.Stdout != "ok" || result.ExitCode != 0 {
			t.Fatalf("Run() = %+v, %v", result, err)
		}
	})

	t.Run("passes directory, environment, and input", func(t *testing.T) {
		t.Parallel()

		directory, err := filepath.EvalSymlinks(t.TempDir())
		if err != nil {
			t.Fatal(err)
		}

		command := shell(`printf '%s|%s|' "$PWD" "$LANGDRIFT_TEST"; cat`)
		command.Dir, command.Env, command.Stdin = directory, []string{"LANGDRIFT_TEST=explicit"}, strings.NewReader("input")

		result, err := process.Local{}.Run(context.Background(), command)
		if err != nil || result.Stdout != directory+"|explicit|input" {
			t.Fatalf("Run() = %q, %v", result.Stdout, err)
		}
	})

	t.Run("streams output while keeping a bounded copy", func(t *testing.T) {
		t.Parallel()

		var stdout, stderr bytes.Buffer
		command := shell("printf streamed; printf diagnostics >&2")
		command.Stdout, command.Stderr = &stdout, &stderr

		result, err := process.Local{}.Run(context.Background(), command)
		if err != nil || stdout.String() != "streamed" || stderr.String() != "diagnostics" || result.Stdout != "streamed" {
			t.Fatalf("Run() = %+v, %v; streamed %q and %q", result, err, stdout.String(), stderr.String())
		}
	})

	t.Run("refuses an empty executable name", func(t *testing.T) {
		t.Parallel()

		if _, err := (process.Local{}).Run(context.Background(), process.Command{}); err == nil {
			t.Fatal("Run() accepted an empty command")
		}
	})

	t.Run("reports a missing executable", func(t *testing.T) {
		t.Parallel()

		result, err := process.Local{}.Run(context.Background(), process.Command{Name: "langdrift-test-missing-tool"})
		if err == nil || result.ExitCode != -1 {
			t.Fatalf("Run() = %+v, %v; want a failure without an exit code", result, err)
		}
	})

	t.Run("does not start after cancellation", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		marker := filepath.Join(t.TempDir(), "started")
		if _, err := (process.Local{}).Run(ctx, shell("touch "+marker)); !errors.Is(err, context.Canceled) {
			t.Fatalf("Run() error = %v, want context.Canceled", err)
		}

		if _, err := os.Stat(marker); !os.IsNotExist(err) {
			t.Fatal("a canceled command started")
		}
	})
}

func TestLocal_RunBoundsCapturedOutput(t *testing.T) {
	t.Parallel()

	const overLimit = "head -c 1048577 /dev/zero"

	t.Run("capture-only output over the limit is an error", func(t *testing.T) {
		t.Parallel()

		result, err := process.Local{}.Run(context.Background(), shell(overLimit))
		if err == nil || !strings.Contains(err.Error(), "capture limit") || len(result.Stdout) != 1<<20 {
			t.Fatalf("Run() = %d bytes, %v; want a bounded copy and an error rather than a silently truncated document", len(result.Stdout), err)
		}
	})

	t.Run("streamed output over the limit is complete in the stream", func(t *testing.T) {
		t.Parallel()

		var stream bytes.Buffer
		command := shell(overLimit)
		command.Stdout = &stream

		result, err := process.Local{}.Run(context.Background(), command)
		if err != nil || stream.Len() != 1<<20+1 || len(result.Stdout) != 1<<20 {
			t.Fatalf("Run() = %d captured, %d streamed, %v", len(result.Stdout), stream.Len(), err)
		}
	})
}

// TestLocal_RunCancelsTheProcessGroup observes a grandchild that would
// create a file shortly after cancellation if it were still alive.
func TestLocal_RunCancelsTheProcessGroup(t *testing.T) {
	t.Parallel()

	marker := filepath.Join(t.TempDir(), "survived")
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	started := time.Now()

	_, err := process.Local{}.Run(ctx, shell("(sleep 0.5; touch "+marker+") & wait"))
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Run() error = %v, want context.DeadlineExceeded", err)
	}

	if elapsed := time.Since(started); elapsed > 3*time.Second {
		t.Fatalf("Run() returned after %v, want it to stop the group promptly", elapsed)
	}

	time.Sleep(time.Second)

	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatal("a process in the canceled group kept running")
	}
}

func TestExplain(t *testing.T) {
	t.Parallel()

	cause := errors.New("exit status 1")

	withDiagnostics := process.Explain("kubectl get", process.Result{Stderr: "  forbidden\n"}, cause)
	if !errors.Is(withDiagnostics, cause) || withDiagnostics.Error() != "kubectl get: exit status 1: forbidden" {
		t.Errorf("Explain() = %q", withDiagnostics)
	}

	without := process.Explain("kind load image", process.Result{}, cause)
	if !errors.Is(without, cause) || without.Error() != "kind load image: exit status 1" {
		t.Errorf("Explain() = %q", without)
	}
}

// BenchmarkLocal_Run measures starting and waiting for one short process, the
// floor every native tool invocation pays before the tool's own work.
func BenchmarkLocal_Run(b *testing.B) {
	runner := process.Local{}

	b.ReportAllocs()

	for range b.N {
		if _, err := runner.Run(context.Background(), process.Command{Name: "true"}); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkLocal_RunCancellation measures how long Run takes to return after
// its context is canceled while a grandchild holds the output pipes.
func BenchmarkLocal_RunCancellation(b *testing.B) {
	runner := process.Local{}

	var latency time.Duration

	for range b.N {
		ctx, cancel := context.WithCancel(context.Background())
		started := &firstWrite{signal: make(chan struct{})}

		command := shell("sleep 30 & echo started; wait")
		command.Stdout = started

		done := make(chan error, 1)
		go func() {
			_, err := runner.Run(ctx, command)
			done <- err
		}()

		<-started.signal
		canceled := time.Now()
		cancel()

		if err := <-done; !errors.Is(err, context.Canceled) {
			b.Fatalf("Run() = %v, want cancellation", err)
		}

		latency += time.Since(canceled)
	}

	b.ReportMetric(float64(latency.Nanoseconds())/float64(b.N), "cancel-ns/op")
}

// firstWrite signals once, when the process first writes.
type firstWrite struct {
	once   sync.Once
	signal chan struct{}
}

func (writer *firstWrite) Write(data []byte) (int, error) {
	writer.once.Do(func() { close(writer.signal) })

	return len(data), nil
}

func TestLocal_ChildrenNeverInheritUndeclaredVariables(t *testing.T) {
	t.Setenv("PROJECT_SECRET", "private-value")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "private-value")
	t.Setenv("LANGDRIFT_TEST_PASS", "declared-value")

	command := shell(`printf '%s|%s|%s|%s' "${PROJECT_SECRET-unset}" "${AWS_SECRET_ACCESS_KEY-unset}" "${LANGDRIFT_TEST_PASS-unset}" "${HOME-missing}"`)
	command.PassEnv = []string{"LANGDRIFT_TEST_*"}

	result, err := (process.Local{}).Run(context.Background(), command)
	if err != nil {
		t.Fatal(err)
	}
	fields := strings.Split(result.Stdout, "|")
	if fields[0] != "unset" || fields[1] != "unset" {
		t.Fatalf("a credential was inherited: %q", result.Stdout)
	}
	if fields[2] != "declared-value" {
		t.Fatalf("declared pass-through missing: %q", result.Stdout)
	}
	if fields[3] == "missing" || fields[3] == "" {
		t.Fatal("the minimal base must keep HOME")
	}
}

func TestLocal_ExplicitAdditionsJoinTheMinimalBase(t *testing.T) {
	command := shell(`printf '%s|%s' "$LANGDRIFT_TEST" "${PATH:+present}"`)
	command.Env = []string{"LANGDRIFT_TEST=explicit"}

	result, err := (process.Local{}).Run(context.Background(), command)
	if err != nil || result.Stdout != "explicit|present" {
		t.Fatalf("Run() = %q, %v", result.Stdout, err)
	}
}
