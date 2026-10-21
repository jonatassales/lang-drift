// Package connection defines the versioned console-to-project configuration.
package connection

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"path"
	"strings"

	"github.com/jonatassales/tao/packages/ex/src/domain/microruntime"
)

// Configuration is the concrete contract inside an authored Markdown document.
// Prose is retained as evidence but never executed as instructions.
type Configuration struct {
	// LLMBudgetUSD is a per-execution display envelope, never a billing gate.
	LLMBudgetUSD *float64    `json:"llmBudgetUsd,omitempty"`
	Version      int         `json:"version"`
	RuntimeUID   string      `json:"runtimeUid"`
	RuntimeID    string      `json:"runtimeId"`
	ProjectID    string      `json:"projectId"`
	Revision     int         `json:"revision"`
	Name         string      `json:"name"`
	CPUMillis    int         `json:"cpuMillis"`
	MemoryMiB    int         `json:"memoryMiB"`
	Model        string      `json:"model"`
	Effort       string      `json:"effort"`
	Components   []Component `json:"components"`
	Native       *Native     `json:"native,omitempty"`
}

// Component references an existing source directory and an explicit build recipe.
type Component struct {
	ID            string `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	Kind          string `json:"kind,omitempty"`
	Path          string `json:"path"`
	Image         string `json:"image"`
	Port          int    `json:"port"`
	ReadinessPath string `json:"readinessPath"`
	Build         *Build `json:"build,omitempty"`
}

// Build is a direct executable invocation, never implicit shell text.
type Build struct {
	Executable string   `json:"executable"`
	Args       []string `json:"args"`
}

// Native preserves an existing v2 descriptor and its owned resources on migration.
type Native struct {
	RuntimeID string `json:"runtimeId"`
}

// Parse requires one explicit langdrift JSON fence, rejects ambiguous or unknown
// fields, and bounds the document. Markdown never grants execution authority.
func Parse(markdown string) (Configuration, error) {
	var config Configuration
	if len(markdown) > 1<<20 {
		return config, invalid("Markdown exceeds 1 MiB")
	}
	const fence = "```langdrift\n"
	if strings.Count(markdown, fence) != 1 {
		return config, invalid("Markdown requires exactly one langdrift JSON fence")
	}
	after := strings.SplitN(markdown, fence, 2)[1]
	end := strings.Index(after, "\n```")
	if end < 0 {
		return config, invalid("configuration fence is not closed")
	}
	decoder := json.NewDecoder(strings.NewReader(after[:end]))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&config); err != nil {
		return config, fmt.Errorf("configuration block: %w", err)
	}
	if err := decoder.Decode(new(any)); !errors.Is(err, io.EOF) {
		return config, invalid("configuration block must contain one JSON object")
	}
	return config, config.Validate()
}

// Markdown renders the portable, editable migration/configuration format.
func (config Configuration) Markdown() string {
	data, _ := json.MarshalIndent(config, "", "  ")
	return "# " + strings.ReplaceAll(config.Name, "\n", " ") + "\n\nConfiguration is interpreted by the connected project agent. Prose does not execute commands.\n\n```langdrift\n" + string(data) + "\n```\n"
}

// Validate bounds resources and makes all execution choices explicit.
func (config Configuration) Validate() error {
	if config.Version != 1 || config.Revision < 1 || config.RuntimeUID == "" || config.ProjectID == "" || strings.TrimSpace(config.Name) == "" {
		return invalid("version 1, identities, revision and name are required")
	}
	if _, err := microruntime.ParseIdentifier(config.RuntimeID); err != nil {
		return err
	}
	if config.CPUMillis < 1 || config.CPUMillis > 64_000 || config.MemoryMiB < 1 || config.MemoryMiB > 1_048_576 {
		return invalid("CPU and memory must be positive bounded per-execution limits")
	}
	if config.LLMBudgetUSD != nil && (math.IsNaN(*config.LLMBudgetUSD) || *config.LLMBudgetUSD < 0.01 || *config.LLMBudgetUSD > 1_000_000) {
		return invalid("LLM budget must be positive bounded USD")
	}
	if config.Model == "" || (config.Effort != "low" && config.Effort != "medium" && config.Effort != "high") {
		return invalid("explicit agent model and low, medium, or high effort are required")
	}
	if len(config.Components) == 0 || len(config.Components) > 16 {
		return invalid("select between 1 and 16 components")
	}
	if config.Native != nil && config.Native.RuntimeID != config.RuntimeID {
		return invalid("migration must preserve the native runtime ID")
	}
	if config.CPUMillis < len(config.Components) || config.MemoryMiB < len(config.Components) {
		return invalid("resource envelope must allocate at least 1m CPU and 1 MiB per component")
	}
	names := map[string]bool{}
	for _, component := range config.Components {
		if _, err := microruntime.ParseComponentPath(component.Path); err != nil {
			return err
		}
		if microruntime.IsPlatformPath(component.Path) || names[component.Path] {
			return invalid("component sources must be outside .langdrift and have unique source paths")
		}
		names[component.Path] = true
		if !microruntime.IsImageReference(component.Image) {
			return invalid("each component requires an explicit image")
		}
		if !microruntime.ValidPort(component.Port) || !strings.HasPrefix(component.ReadinessPath, "/") {
			return invalid("each component requires a port and readiness path")
		}
		if component.Build == nil || component.Build.Executable == "" || len(component.Build.Args) > 128 || strings.ContainsAny(component.Build.Executable, "\n\r\x00") {
			return invalid("each component requires an explicit existing build recipe")
		}
		if err := component.Build.Validate(); err != nil {
			return err
		}
		encoded, _ := json.Marshal(component.Build)
		if len(encoded) > 16384 || bytes.Contains(encoded, []byte(`\u0000`)) {
			return invalid("build recipe exceeds supported limits")
		}
	}
	return nil
}

// ErrInvalidConfiguration classifies every authored-configuration rejection,
// so callers match the class with errors.Is and read the reason in the text.
var ErrInvalidConfiguration = errors.New("invalid runtime configuration")

// invalid rejects an authored configuration with its reason.
func invalid(reason string) error {
	return fmt.Errorf("%w: %s", ErrInvalidConfiguration, reason)
}

// ErrRecipeRefused marks a build recipe outside the allowlist policy, so
// callers can distinguish refusal from failure.
var ErrRecipeRefused = errors.New("build recipe refused by policy")

// The only tools a recipe may name. Anything else runs through a committed
// project script, which code review then sees. See ADR 0002.
var allowedBuildTools = map[string]bool{
	"docker": true,
	"go":     true,
	"npm":    true,
	"pnpm":   true,
	"make":   true,
}

// Validate allows exactly two recipe shapes: a committed project script
// (a relative path containing a slash, resolved inside the checkout at build
// time) or an allowlisted tool. Repository programs themselves remain trusted
// consumer code; this is not an operating-system sandbox.
func (build Build) Validate() error {
	refuse := func(reason string) error {
		return fmt.Errorf("%w: %s", ErrRecipeRefused, reason)
	}

	for _, arg := range build.Args {
		value := arg
		if _, after, ok := strings.Cut(value, "="); ok {
			value = after
		}
		if strings.HasPrefix(value, "/") || strings.Contains(value, "../") || value == ".." || strings.Contains(value, "://") {
			return refuse("build arguments must stay within the committed source context")
		}
		if strings.HasPrefix(arg, "--secret") || strings.HasPrefix(arg, "--ssh") || strings.HasPrefix(arg, "--build-context") {
			return refuse("external build credential/context mounts are unsupported")
		}
	}

	if strings.Contains(build.Executable, "/") {
		if strings.Contains(build.Executable, "\\") {
			return refuse("a committed script is named by a slash-separated relative path")
		}
		cleaned := path.Clean(build.Executable)
		if strings.HasPrefix(build.Executable, "/") || cleaned == "." || cleaned == ".." || strings.HasPrefix(cleaned, "../") {
			return refuse("a committed script must resolve inside the checkout")
		}
		return nil
	}

	if !allowedBuildTools[build.Executable] {
		return refuse("the executable is neither a committed script (./name) nor an allowed build tool (docker, go, npm, pnpm, make)")
	}
	if len(build.Args) > 0 && strings.HasPrefix(build.Args[0], "-") {
		return refuse("an allowed tool refuses a flag as its first argument")
	}
	if build.Executable == "docker" && (len(build.Args) == 0 || build.Args[0] != "build") {
		return refuse("Docker recipes must use docker build")
	}
	return nil
}
