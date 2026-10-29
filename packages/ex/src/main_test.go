package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/jonatassales/tao/packages/ex/src/presentation/cli"
	"github.com/jonatassales/tao/packages/ex/src/presentation/cli/support"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// mainUseCases is the production composition, without the native lifecycles.
func mainUseCases() cli.UseCases {
	cases, _ := useCases()
	return cases
}

func TestCommand(t *testing.T) {
	tests := []struct {
		name   string
		args   []string
		status int
		stdout string
		stderr string
	}{
		{name: "help on no args", status: 0, stdout: "LangDrift CLI"},
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
			status := cli.Execute(context.Background(), tt.args, support.Options{Version: version, Directory: t.TempDir(), Output: &out, Diagnostics: &err}, mainUseCases())
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
			if status := cli.Execute(context.Background(), []string{command}, support.Options{Version: version, Output: failingWriter{}, Diagnostics: &diagnostics}, mainUseCases()); status != 1 {
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
			status := cli.Execute(context.Background(), []string{"runtime", verb, "demo"}, support.Options{Output: &out, Diagnostics: &diagnostics}, mainUseCases())
			if status != 1 || !strings.Contains(diagnostics.String(), "moved to the console") {
				t.Fatalf("retired command = %d %s", status, diagnostics.String())
			}
		})
	}
}

func TestLauncherDispatchAndFallback(t *testing.T) {
	sourceDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	packageDir := filepath.Dir(sourceDir)
	launcherSource, err := os.ReadFile(filepath.Join(packageDir, "bin", "ex"))
	if err != nil {
		t.Fatalf("read bin/ex: %v", err)
	}

	writeFakeUname := func(t *testing.T, directory, kernel, machine string) {
		t.Helper()
		script := fmt.Sprintf("#!/bin/sh\nif [ \"${1:-}\" = \"-s\" ]; then printf '%%s\\n' %q; else printf '%%s\\n' %q; fi\n", kernel, machine)
		unamePath := filepath.Join(directory, "uname")
		if err := os.WriteFile(unamePath, []byte(script), 0o755); err != nil {
			t.Fatalf("write fake uname: %v", err)
		}
	}

	setupFixture := func(t *testing.T) (string, string) {
		t.Helper()
		root := t.TempDir()
		pkg := filepath.Join(root, "node_modules", "@langdrift", "ex")
		if err := os.MkdirAll(filepath.Join(pkg, "bin"), 0o755); err != nil {
			t.Fatalf("mkdir bin: %v", err)
		}
		if err := os.MkdirAll(filepath.Join(pkg, "dist"), 0o755); err != nil {
			t.Fatalf("mkdir dist: %v", err)
		}
		launcherPath := filepath.Join(pkg, "bin", "ex")
		if err := os.WriteFile(launcherPath, launcherSource, 0o755); err != nil {
			t.Fatalf("write launcher: %v", err)
		}
		dotBin := filepath.Join(root, "node_modules", ".bin")
		if err := os.MkdirAll(dotBin, 0o755); err != nil {
			t.Fatalf("mkdir .bin: %v", err)
		}
		symlinkPath := filepath.Join(dotBin, "ex")
		if err := os.Symlink(launcherPath, symlinkPath); err != nil {
			t.Fatalf("symlink launcher: %v", err)
		}
		return pkg, symlinkPath
	}

	platforms := []struct {
		kernel  string
		machine string
		target  string
	}{
		{kernel: "Darwin", machine: "arm64", target: "ex-darwin-arm64"},
		{kernel: "Darwin", machine: "x86_64", target: "ex-darwin-amd64"},
		{kernel: "Linux", machine: "x86_64", target: "ex-linux-amd64"},
		{kernel: "Linux", machine: "aarch64", target: "ex-linux-arm64"},
		{kernel: "Linux", machine: "arm64", target: "ex-linux-arm64"},
	}

	for _, platform := range platforms {
		t.Run(platform.kernel+"-"+platform.machine, func(t *testing.T) {
			pkg, symlinkPath := setupFixture(t)
			fakeBin := t.TempDir()
			writeFakeUname(t, fakeBin, platform.kernel, platform.machine)

			targetBinary := filepath.Join(pkg, "dist", platform.target)
			stub := fmt.Sprintf("#!/bin/sh\nprintf 'target=%s args=%%s\\n' \"$*\"\nexit 0\n", platform.target)
			if err := os.WriteFile(targetBinary, []byte(stub), 0o644); err != nil {
				t.Fatalf("write target binary: %v", err)
			}

			command := exec.Command(symlinkPath, "agent", "plan", "0001")
			command.Env = []string{"PATH=" + fakeBin + ":/usr/bin:/bin"}
			output, err := command.CombinedOutput()
			if err != nil {
				t.Fatalf("run launcher: %v\n%s", err, output)
			}
			expected := fmt.Sprintf("target=%s args=agent plan 0001\n", platform.target)
			if string(output) != expected {
				t.Fatalf("output = %q, want %q", string(output), expected)
			}
		})
	}

	t.Run("falls back to host dist/ex when platform binary is absent", func(t *testing.T) {
		pkg, symlinkPath := setupFixture(t)
		fakeBin := t.TempDir()
		writeFakeUname(t, fakeBin, "Darwin", "arm64")

		hostBinary := filepath.Join(pkg, "dist", "ex")
		if err := os.WriteFile(hostBinary, []byte("#!/bin/sh\nprintf 'host-fallback %s\\n' \"$*\"\n"), 0o755); err != nil {
			t.Fatalf("write host binary: %v", err)
		}

		command := exec.Command(symlinkPath, "version")
		command.Env = []string{"PATH=" + fakeBin + ":/usr/bin:/bin"}
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run fallback: %v\n%s", err, output)
		}
		if string(output) != "host-fallback version\n" {
			t.Fatalf("output = %q, want %q", string(output), "host-fallback version\n")
		}
	})

	t.Run("exits 4 on unsupported platform", func(t *testing.T) {
		_, symlinkPath := setupFixture(t)
		fakeBin := t.TempDir()
		writeFakeUname(t, fakeBin, "FreeBSD", "riscv64")

		command := exec.Command(symlinkPath, "version")
		command.Env = []string{"PATH=" + fakeBin + ":/usr/bin:/bin"}
		var stderr bytes.Buffer
		command.Stderr = &stderr
		err := command.Run()
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) || exitErr.ExitCode() != 4 {
			t.Fatalf("err = %v, want exit code 4 (stderr: %s)", err, stderr.String())
		}
		if !strings.Contains(stderr.String(), "macOS and Linux") {
			t.Fatalf("stderr = %q, want platform diagnostic", stderr.String())
		}
	})

	t.Run("exits 4 when binary is missing", func(t *testing.T) {
		_, symlinkPath := setupFixture(t)
		fakeBin := t.TempDir()
		writeFakeUname(t, fakeBin, "Darwin", "arm64")

		command := exec.Command(symlinkPath, "version")
		command.Env = []string{"PATH=" + fakeBin + ":/usr/bin:/bin"}
		var stderr bytes.Buffer
		command.Stderr = &stderr
		err := command.Run()
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) || exitErr.ExitCode() != 4 {
			t.Fatalf("err = %v, want exit code 4 (stderr: %s)", err, stderr.String())
		}
		if !strings.Contains(stderr.String(), "binary") {
			t.Fatalf("stderr = %q, want missing binary diagnostic", stderr.String())
		}
	})

	t.Run("build.sh rejects unknown flags with exit code 2", func(t *testing.T) {
		command := exec.Command("sh", filepath.Join(packageDir, "build.sh"), "--unknown")
		var stderr bytes.Buffer
		command.Stderr = &stderr
		err := command.Run()
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) || exitErr.ExitCode() != 2 {
			t.Fatalf("err = %v, want exit code 2 (stderr: %s)", err, stderr.String())
		}
		if !strings.Contains(stderr.String(), "Usage: sh build.sh [--release]") {
			t.Fatalf("stderr = %q, want usage diagnostic", stderr.String())
		}
	})
}
