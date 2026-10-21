package main

import (
	"bytes"
	"context"
	"errors"
	"github.com/exphere/ex/src/presentation/cli"
	"github.com/exphere/ex/src/presentation/cli/support"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCommand(t *testing.T) {
	tests := []struct {
		name   string
		args   []string
		status int
		stdout string
		stderr string
	}{
		{name: "help on no args", status: 0, stdout: "Exphere CLI"},
		{name: "explicit help", args: []string{"--help"}, status: 0, stdout: "Usage:"},
		{name: "version", args: []string{"version"}, status: 0, stdout: "ex 0.0.0-dev"},
		{name: "init initialized", args: []string{"init"}, status: 0, stdout: "initialized"},
		{name: "unexpected init flags", args: []string{"init", "--force"}, status: 2, stderr: "unknown flag"},
		{name: "unexpected init arguments", args: []string{"init", "project"}, status: 2, stderr: "requires 0 argument"},
		{name: "invalid command", args: []string{"invalid"}, status: 2, stderr: "unknown command"},
		{name: "unexpected version arguments", args: []string{"version", "extra"}, status: 2, stderr: "requires 0 argument"},
		{name: "unexpected help arguments", args: []string{"help", "extra"}, status: 2, stderr: "unknown help command"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out, err bytes.Buffer
			status := cli.Execute(context.Background(), tt.args, support.Options{Version: version, Directory: t.TempDir(), Output: &out, Diagnostics: &err}, useCases())
			if status != tt.status || !strings.Contains(out.String(), tt.stdout) || !strings.Contains(err.String(), tt.stderr) {
				t.Fatalf("run(%q) = (%d, %q, %q), want status %d and substrings %q, %q", tt.args, status, out.String(), err.String(), tt.status, tt.stdout, tt.stderr)
			}
			if tt.status == 0 && err.Len() != 0 {
				t.Errorf("unexpected stderr: %s", err.String())
			}
			if tt.status != 0 && out.Len() != 0 {
				t.Errorf("unexpected stdout: %s", out.String())
			}
		})
	}
}

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) {
	return 0, errors.New("output unavailable")
}

func TestOutputFailure(t *testing.T) {
	for _, command := range []string{"help", "version"} {
		t.Run(command, func(t *testing.T) {
			var diagnostics bytes.Buffer
			if status := cli.Execute(context.Background(), []string{command}, support.Options{Version: version, Output: failingWriter{}, Diagnostics: &diagnostics}, useCases()); status != 1 {
				t.Fatalf("failed output returned status %d, want 1", status)
			}
			if !strings.Contains(diagnostics.String(), "output unavailable") {
				t.Fatalf("missing output failure diagnostic: %q", diagnostics.String())
			}
		})
	}
}

// BenchmarkProcessStartup measures one complete `ex version` process, the
// latency floor of every command, after building the executable once.
func BenchmarkProcessStartup(b *testing.B) {
	goTool, err := exec.LookPath("go")
	if err != nil {
		b.Skip("go is not on PATH")
	}

	binary := filepath.Join(b.TempDir(), "ex")
	if output, err := exec.Command(goTool, "build", "-o", binary, ".").CombinedOutput(); err != nil {
		b.Fatalf("build: %v\n%s", err, output)
	}

	b.ResetTimer()

	for range b.N {
		if err := exec.Command(binary, "version").Run(); err != nil {
			b.Fatal(err)
		}
	}
}

func TestRuntimeCommandsAreRetired(t *testing.T) {
	for _, verb := range []string{"init", "validate", "up", "down", "inspect", "logs", "forward", "hook"} {
		t.Run(verb, func(t *testing.T) {
			var out, diagnostics bytes.Buffer
			status := cli.Execute(context.Background(), []string{"runtime", verb, "demo"}, support.Options{Output: &out, Diagnostics: &diagnostics}, useCases())
			if status != 1 || !strings.Contains(diagnostics.String(), "moved to the console") {
				t.Fatalf("retired command = %d %s", status, diagnostics.String())
			}
		})
	}
}
