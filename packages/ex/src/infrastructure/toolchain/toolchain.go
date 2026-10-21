// Package toolchain probes the consumer-installed native tools. It never
// installs, updates, or configures anything.
package toolchain

import (
	"context"
	"encoding/json"
	"os/exec"
	"regexp"
	"time"

	"github.com/exphere/ex/src/application"
	"github.com/exphere/ex/src/domain/microruntime"
	"github.com/exphere/ex/src/infrastructure/process"
)

// probeTimeout bounds each version probe.
const probeTimeout = 10 * time.Second

var versionPattern = regexp.MustCompile(`\bv?[0-9]+\.[0-9]+\.[0-9]+(?:[-+][0-9A-Za-z.-]+)?\b`)

// tool describes how to find out a native tool's version. When versionField
// is set, the probe prints JSON and the version is read from that path.
type tool struct {
	name         string
	install      string
	arguments    []string
	versionField []string
}

// tools are the native tools EX relies on.
var tools = []tool{
	{
		name:      "docker",
		install:   "https://docs.docker.com/get-started/get-docker/",
		arguments: []string{"--version"},
	},
	{
		name:      "kind",
		install:   "https://kind.sigs.k8s.io/docs/user/quick-start/",
		arguments: []string{"version"},
	},
	{
		name:         "kubectl",
		install:      "https://kubernetes.io/docs/tasks/tools/",
		arguments:    []string{"version", "--client", "-o", "json"},
		versionField: []string{"clientVersion", "gitVersion"},
	},
	{
		name:         "terraform",
		install:      "https://developer.hashicorp.com/terraform/install",
		arguments:    []string{"version", "-json"},
		versionField: []string{"terraform_version"},
	},
}

// commandRunner runs a bounded native tool probe and returns its captured
// output. Cancellation must stop the invocation before Run returns.
type commandRunner interface {
	Run(context.Context, process.Command) (process.Result, error)
}

// Toolchain implements the application's tool version and prerequisite ports.
type Toolchain struct {
	Runner commandRunner
}

// Versions observes each tool's version for evidence. A tool that cannot be
// probed is reported as unobserved, never guessed.
func (toolchain Toolchain) Versions(ctx context.Context) map[string]microruntime.ToolVersion {
	versions := make(map[string]microruntime.ToolVersion, len(tools))

	for _, tool := range tools {
		version, observed := toolchain.version(ctx, tool)
		versions[tool.name] = microruntime.ToolVersion{Version: version, Observed: observed}
	}

	return versions
}

// Check reports where each tool is installed, whether it answers, and
// whether the container engine is reachable.
func (toolchain Toolchain) Check(ctx context.Context) application.PrerequisiteReport {
	var report application.PrerequisiteReport

	for _, tool := range tools {
		status := application.ToolStatus{Name: tool.name, Install: tool.install, Status: application.ToolMissing}

		if path, err := exec.LookPath(tool.name); err == nil {
			status.Path, status.Status = path, application.ToolUnavailable

			if version, observed := toolchain.version(ctx, tool); observed {
				status.Version, status.Status = version, application.ToolAvailable
			}
		}

		report.Tools = append(report.Tools, status)
	}

	report.ContainerEngineReachable = toolchain.engineReachable(ctx)

	return report
}

func (toolchain Toolchain) version(ctx context.Context, tool tool) (string, bool) {
	probeContext, cancel := context.WithTimeout(ctx, probeTimeout)
	defer cancel()

	result, err := toolchain.Runner.Run(probeContext, process.Command{Name: tool.name, Args: tool.arguments})
	if err != nil {
		return "", false
	}

	text := result.Stdout
	if len(tool.versionField) > 0 {
		text = jsonField(result.Stdout, tool.versionField)
	}

	version := versionPattern.FindString(text)

	return version, version != ""
}

func (toolchain Toolchain) engineReachable(ctx context.Context) bool {
	probeContext, cancel := context.WithTimeout(ctx, probeTimeout)
	defer cancel()

	_, err := toolchain.Runner.Run(probeContext, process.Command{
		Name: "docker",
		Args: []string{"info", "--format", "{{.ServerVersion}}"},
	})

	return err == nil
}

// jsonField reads a string at path in a JSON document, or returns "".
func jsonField(document string, path []string) string {
	var value any
	if json.Unmarshal([]byte(document), &value) != nil {
		return ""
	}

	for _, key := range path {
		object, _ := value.(map[string]any)
		value = object[key]
	}

	text, _ := value.(string)

	return text
}
