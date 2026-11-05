package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jonatassales/tao/lib/cli/src/presentation/cli"
	"github.com/jonatassales/tao/lib/cli/src/presentation/cli/support"
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
		{name: "explicit help names dft", args: []string{"--help"}, status: 0, stdout: "Usage:\n  dft [flags]\n  dft [command]"},
		{name: "version", args: []string{"version"}, status: 0, stdout: "dft 0.0.0-dev"},
		{name: "init initialized", args: []string{"init"}, status: 0, stdout: "initialized"},
		{name: "unexpected init flags", args: []string{"init", "--force"}, status: 2, stderr: "unknown flag"},
		{name: "unexpected init arguments", args: []string{"init", "project"}, status: 2, stderr: "requires 0 argument"},
		{name: "invalid command points to dft help", args: []string{"invalid"}, status: 2, stderr: "unknown command \"invalid\"; run 'dft --help'"},
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

// BenchmarkProcessStartup measures one complete `dft version` process, the
// latency floor of every command, after building the executable once.
func BenchmarkProcessStartup(b *testing.B) {
	goTool, err := exec.LookPath("go")
	if err != nil {
		b.Skip("go is not on PATH")
	}

	binary := filepath.Join(b.TempDir(), "dft")
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
	launcherSource, err := os.ReadFile(filepath.Join(packageDir, "bin", "dft.sh"))
	if err != nil {
		t.Fatalf("read bin/dft.sh: %v", err)
	}

	writeFakeUname := func(t *testing.T, kernel, machine string) string {
		t.Helper()
		directory := t.TempDir()
		script := fmt.Sprintf("#!/bin/sh\nif [ \"${1:-}\" = \"-s\" ]; then printf '%%s\\n' %q; else printf '%%s\\n' %q; fi\n", kernel, machine)
		unamePath := filepath.Join(directory, "uname")
		if err := os.WriteFile(unamePath, []byte(script), 0o755); err != nil {
			t.Fatalf("write fake uname: %v", err)
		}
		return directory
	}

	setupFixture := func(t *testing.T) (string, string) {
		t.Helper()
		root := t.TempDir()
		pkg := filepath.Join(root, "node_modules", "@langdrift", "dft")
		if err := os.MkdirAll(filepath.Join(pkg, "bin"), 0o755); err != nil {
			t.Fatalf("mkdir bin: %v", err)
		}
		if err := os.MkdirAll(filepath.Join(pkg, "dist"), 0o755); err != nil {
			t.Fatalf("mkdir dist: %v", err)
		}
		launcherPath := filepath.Join(pkg, "bin", "dft.sh")
		if err := os.WriteFile(launcherPath, launcherSource, 0o755); err != nil {
			t.Fatalf("write launcher: %v", err)
		}
		dotBin := filepath.Join(root, "node_modules", ".bin")
		if err := os.MkdirAll(dotBin, 0o755); err != nil {
			t.Fatalf("mkdir .bin: %v", err)
		}
		symlinkPath := filepath.Join(dotBin, "dft")
		if err := os.Symlink(launcherPath, symlinkPath); err != nil {
			t.Fatalf("symlink launcher: %v", err)
		}
		return pkg, symlinkPath
	}

	runLauncher := func(symlinkPath, fakeBin string, args ...string) (int, string, string) {
		command := exec.Command(symlinkPath, args...)
		command.Env = []string{"PATH=" + fakeBin + ":/usr/bin:/bin"}
		var stdout, stderr bytes.Buffer
		command.Stdout = &stdout
		command.Stderr = &stderr
		runErr := command.Run()
		if runErr == nil {
			return 0, stdout.String(), stderr.String()
		}
		var exitErr *exec.ExitError
		if errors.As(runErr, &exitErr) {
			return exitErr.ExitCode(), stdout.String(), stderr.String()
		}
		return -1, stdout.String(), stderr.String()
	}

	platforms := []struct {
		kernel  string
		machine string
		target  string
		status  int
	}{
		{kernel: "Darwin", machine: "arm64", target: "dft-darwin-arm64", status: 0},
		{kernel: "Darwin", machine: "x86_64", target: "dft-darwin-amd64", status: 0},
		{kernel: "Linux", machine: "x86_64", target: "dft-linux-amd64", status: 0},
		{kernel: "Linux", machine: "aarch64", target: "dft-linux-arm64", status: 0},
		{kernel: "Linux", machine: "arm64", target: "dft-linux-arm64", status: 3},
	}

	for _, platform := range platforms {
		t.Run(platform.kernel+"-"+platform.machine, func(t *testing.T) {
			pkg, symlinkPath := setupFixture(t)
			fakeBin := writeFakeUname(t, platform.kernel, platform.machine)

			targetBinary := filepath.Join(pkg, "dist", platform.target)
			stub := fmt.Sprintf("#!/bin/sh\nprintf 'target=%s args=%%s\\n' \"$*\"\nexit %d\n", platform.target, platform.status)
			if err := os.WriteFile(targetBinary, []byte(stub), 0o644); err != nil {
				t.Fatalf("write target binary: %v", err)
			}

			status, stdout, stderr := runLauncher(symlinkPath, fakeBin, "agent", "plan", "0001")
			expected := fmt.Sprintf("target=%s args=agent plan 0001\n", platform.target)
			if status != platform.status || stdout != expected || stderr != "" {
				t.Fatalf("runLauncher = (%d, %q, %q), want (%d, %q, \"\")", status, stdout, stderr, platform.status, expected)
			}
		})
	}

	t.Run("falls back to host dist/dft when platform binary is absent", func(t *testing.T) {
		pkg, symlinkPath := setupFixture(t)
		fakeBin := writeFakeUname(t, "Darwin", "arm64")

		hostBinary := filepath.Join(pkg, "dist", "dft")
		if err := os.WriteFile(hostBinary, []byte("#!/bin/sh\nprintf 'host-fallback %s\\n' \"$*\"\n"), 0o755); err != nil {
			t.Fatalf("write host binary: %v", err)
		}

		status, stdout, stderr := runLauncher(symlinkPath, fakeBin, "version")
		if status != 0 || stdout != "host-fallback version\n" || stderr != "" {
			t.Fatalf("runLauncher = (%d, %q, %q), want (0, %q, \"\")", status, stdout, stderr, "host-fallback version\n")
		}
	})

	for _, unsupported := range []struct {
		name    string
		kernel  string
		machine string
	}{
		{name: "exits 4 on unsupported kernel", kernel: "FreeBSD", machine: "arm64"},
		{name: "exits 4 on unsupported cpu architecture", kernel: "Linux", machine: "ppc64le"},
	} {
		t.Run(unsupported.name, func(t *testing.T) {
			_, symlinkPath := setupFixture(t)
			fakeBin := writeFakeUname(t, unsupported.kernel, unsupported.machine)

			status, _, stderr := runLauncher(symlinkPath, fakeBin, "version")
			if status != 4 || !strings.Contains(stderr, "macOS and Linux") {
				t.Fatalf("runLauncher = (%d, stderr=%q), want status 4 and platform diagnostic", status, stderr)
			}
		})
	}

	t.Run("exits 4 when binary is missing", func(t *testing.T) {
		_, symlinkPath := setupFixture(t)
		fakeBin := writeFakeUname(t, "Darwin", "arm64")

		status, _, stderr := runLauncher(symlinkPath, fakeBin, "version")
		if status != 4 || !strings.Contains(stderr, "binary") {
			t.Fatalf("runLauncher = (%d, stderr=%q), want status 4 and missing binary diagnostic", status, stderr)
		}
	})

	t.Run("build.sh injects package.json version and rejects unknown flags", func(t *testing.T) {
		command := exec.Command("sh", filepath.Join(packageDir, "build.sh"), "--unknown")
		var stderr bytes.Buffer
		command.Stderr = &stderr
		err := command.Run()
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) || exitErr.ExitCode() != 2 || !strings.Contains(stderr.String(), "Usage: sh build.sh [--release]") {
			t.Fatalf("unknown flag = %v (stderr=%q), want exit code 2 and usage diagnostic", err, stderr.String())
		}

		if _, lookErr := exec.LookPath("go"); lookErr != nil {
			t.Skip("go is not on PATH")
		}
		if output, buildErr := exec.Command("sh", filepath.Join(packageDir, "build.sh")).CombinedOutput(); buildErr != nil {
			t.Fatalf("build.sh: %v\n%s", buildErr, output)
		}
		versionOut, runErr := exec.Command(filepath.Join(packageDir, "bin", "dft.sh"), "version").CombinedOutput()
		if runErr != nil {
			t.Fatalf("bin/dft.sh version: %v\n%s", runErr, versionOut)
		}

		// The version is the package's own, so a release bump never edits
		// this test.
		var manifest struct {
			Version string `json:"version"`
		}

		manifestData, readErr := os.ReadFile(filepath.Join(packageDir, "package.json"))
		if readErr != nil {
			t.Fatal(readErr)
		}

		if unmarshalErr := json.Unmarshal(manifestData, &manifest); unmarshalErr != nil {
			t.Fatalf("package.json: %v", unmarshalErr)
		}

		if manifest.Version == "" {
			t.Fatal("package.json names no version")
		}

		want := "dft " + manifest.Version

		if strings.TrimSpace(string(versionOut)) != want {
			t.Fatalf("bin/dft.sh version = %q, want %q", strings.TrimSpace(string(versionOut)), want)
		}
	})
}
