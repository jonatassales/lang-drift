package toolchain_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/domain/microruntime"
	"github.com/exphere/ex/src/infrastructure/process"
	"github.com/exphere/ex/src/infrastructure/toolchain"
)

// fakeTools answers each tool's probe from outputs, keyed by executable.
type fakeTools struct {
	outputs map[string]string
	failing map[string]bool
}

func (tools fakeTools) Run(ctx context.Context, command process.Command) (process.Result, error) {
	if _, bounded := ctx.Deadline(); !bounded {
		return process.Result{}, errors.New("probe without a deadline")
	}

	if tools.failing[command.Name] {
		return process.Result{}, errors.New("exit status 1")
	}

	return process.Result{Stdout: tools.outputs[command.Name]}, nil
}

var installedTools = map[string]string{
	"docker":    "Docker version 27.3.1, build ce12230",
	"kind":      "kind v0.27.0 go1.23.4 linux/amd64",
	"kubectl":   `{"clientVersion": {"gitVersion": "v1.32.2"}, "kustomizeVersion": "v5.5.0"}`,
	"terraform": `{"terraform_version": "1.10.5", "platform": "linux_amd64"}`,
}

func TestToolchain_Versions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		tools fakeTools
		want  map[string]microruntime.ToolVersion
	}{
		{
			name:  "every tool answers",
			tools: fakeTools{outputs: installedTools},
			want: map[string]microruntime.ToolVersion{
				"docker": {Version: "27.3.1", Observed: true}, "kind": {Version: "v0.27.0", Observed: true},
				"kubectl": {Version: "v1.32.2", Observed: true}, "terraform": {Version: "1.10.5", Observed: true},
			},
		},
		{
			name: "failing probes and unparsable output are never guessed",
			tools: fakeTools{
				outputs: map[string]string{"docker": "Docker", "kubectl": `{"clientVersion": {}}`, "terraform": "1.10.5"},
				failing: map[string]bool{"kind": true},
			},
			want: map[string]microruntime.ToolVersion{"docker": {}, "kind": {}, "kubectl": {}, "terraform": {}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			versions := toolchain.Toolchain{Runner: tt.tools}.Versions(context.Background())

			if len(versions) != len(tt.want) {
				t.Fatalf("Versions() = %v, want %v", versions, tt.want)
			}

			for tool, want := range tt.want {
				if versions[tool] != want {
					t.Errorf("%s = %+v, want %+v", tool, versions[tool], want)
				}
			}
		})
	}
}

// TestToolchain_Check changes PATH, so it cannot run in parallel.
func TestToolchain_Check(t *testing.T) {
	bin := t.TempDir()
	for _, tool := range []string{"docker", "kind", "kubectl"} {
		if err := os.WriteFile(filepath.Join(bin, tool), []byte("#!/bin/sh\n"), 0o755); err != nil {
			t.Fatal(err)
		}
	}

	t.Setenv("PATH", bin)

	report := toolchain.Toolchain{Runner: fakeTools{outputs: installedTools, failing: map[string]bool{"kubectl": true}}}.Check(context.Background())

	want := map[string]string{
		"docker":    application.ToolAvailable,
		"kind":      application.ToolAvailable,
		"kubectl":   application.ToolUnavailable,
		"terraform": application.ToolMissing,
	}

	if len(report.Tools) != len(want) {
		t.Fatalf("Check() tools = %+v", report.Tools)
	}

	for _, tool := range report.Tools {
		if tool.Status != want[tool.Name] || tool.Install == "" {
			t.Errorf("%s = %+v, want status %s and installation guidance", tool.Name, tool, want[tool.Name])
		}

		if (tool.Status == application.ToolMissing) != (tool.Path == "") {
			t.Errorf("%s path = %q with status %s", tool.Name, tool.Path, tool.Status)
		}
	}

	if !report.ContainerEngineReachable {
		t.Error("the container engine was reported unreachable")
	}

	unreachable := toolchain.Toolchain{Runner: fakeTools{outputs: installedTools, failing: map[string]bool{"docker": true}}}.Check(context.Background())
	if unreachable.ContainerEngineReachable {
		t.Error("a failing docker info was reported reachable")
	}
}
