// Package connection defines the versioned console-to-project configuration.
package connection

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"strings"

	"github.com/exphere/ex/src/domain/microruntime"
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

// Parse requires one explicit exphere JSON fence, rejects ambiguous or unknown
// fields, and bounds the document. Markdown never grants execution authority.
func Parse(markdown string) (Configuration, error) {
	var config Configuration
	if len(markdown) > 1<<20 {
		return config, errors.New("Markdown exceeds 1 MiB")
	}
	const fence = "```exphere\n"
	if strings.Count(markdown, fence) != 1 {
		return config, errors.New("Markdown requires exactly one exphere JSON fence")
	}
	after := strings.SplitN(markdown, fence, 2)[1]
	end := strings.Index(after, "\n```")
	if end < 0 {
		return config, errors.New("configuration fence is not closed")
	}
	decoder := json.NewDecoder(strings.NewReader(after[:end]))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&config); err != nil {
		return config, fmt.Errorf("configuration block: %w", err)
	}
	if err := decoder.Decode(new(any)); !errors.Is(err, io.EOF) {
		return config, errors.New("configuration block must contain one JSON object")
	}
	return config, config.Validate()
}

// Markdown renders the portable, editable migration/configuration format.
func (config Configuration) Markdown() string {
	data, _ := json.MarshalIndent(config, "", "  ")
	return "# " + strings.ReplaceAll(config.Name, "\n", " ") + "\n\nConfiguration is interpreted by the connected project agent. Prose does not execute commands.\n\n```exphere\n" + string(data) + "\n```\n"
}

// Validate bounds resources and makes all execution choices explicit.
func (config Configuration) Validate() error {
	if config.Version != 1 || config.Revision < 1 || config.RuntimeUID == "" || config.ProjectID == "" || strings.TrimSpace(config.Name) == "" {
		return errors.New("version 1, identities, revision and name are required")
	}
	if _, err := microruntime.ParseIdentifier(config.RuntimeID); err != nil {
		return err
	}
	if config.CPUMillis < 1 || config.CPUMillis > 64_000 || config.MemoryMiB < 1 || config.MemoryMiB > 1_048_576 {
		return errors.New("CPU and memory must be positive bounded per-execution limits")
	}
	if config.LLMBudgetUSD != nil && (math.IsNaN(*config.LLMBudgetUSD) || *config.LLMBudgetUSD < 0.01 || *config.LLMBudgetUSD > 1_000_000) {
		return errors.New("LLM budget must be positive bounded USD")
	}
	if config.Model == "" || (config.Effort != "low" && config.Effort != "medium" && config.Effort != "high") {
		return errors.New("explicit agent model and low, medium, or high effort are required")
	}
	if len(config.Components) == 0 || len(config.Components) > 16 {
		return errors.New("select between 1 and 16 components")
	}
	if config.Native != nil && config.Native.RuntimeID != config.RuntimeID {
		return errors.New("migration must preserve the native runtime ID")
	}
	if config.CPUMillis < len(config.Components) || config.MemoryMiB < len(config.Components) {
		return errors.New("resource envelope must allocate at least 1m CPU and 1 MiB per component")
	}
	names := map[string]bool{}
	for _, component := range config.Components {
		if _, err := microruntime.ParseComponentPath(component.Path); err != nil {
			return err
		}
		if strings.HasPrefix(component.Path, ".exphere/") || names[component.Path] {
			return errors.New("component sources must be outside .exphere and have unique source paths")
		}
		names[component.Path] = true
		if !microruntime.IsImageReference(component.Image) {
			return errors.New("each component requires an explicit image")
		}
		if component.Port < 1 || component.Port > 65535 || !strings.HasPrefix(component.ReadinessPath, "/") {
			return errors.New("each component requires a port and readiness path")
		}
		if component.Build == nil || component.Build.Executable == "" || len(component.Build.Args) > 128 || strings.ContainsAny(component.Build.Executable, "\n\r\x00") {
			return errors.New("each component requires an explicit existing build recipe")
		}
		if err := component.Build.Validate(); err != nil {
			return err
		}
		encoded, _ := json.Marshal(component.Build)
		if len(encoded) > 16384 || bytes.Contains(encoded, []byte(`\u0000`)) {
			return errors.New("build recipe exceeds supported limits")
		}
	}
	return nil
}

// Validate refuses command-line context escapes. Repository programs themselves
// remain trusted consumer code; this is not an operating-system sandbox.
func (build Build) Validate() error {
	if strings.ContainsAny(build.Executable, "/\\") {
		return errors.New("build executable must be a tool name resolved from PATH")
	}
	for _, arg := range build.Args {
		value := arg
		if _, after, ok := strings.Cut(value, "="); ok {
			value = after
		}
		if strings.HasPrefix(value, "/") || strings.Contains(value, "../") || value == ".." || strings.Contains(value, "://") {
			return errors.New("build arguments must stay within the committed source context")
		}
		if strings.HasPrefix(arg, "--secret") || strings.HasPrefix(arg, "--ssh") || strings.HasPrefix(arg, "--build-context") {
			return errors.New("external build credential/context mounts are unsupported")
		}
	}
	if (build.Executable == "sh" || build.Executable == "bash" || build.Executable == "zsh") && len(build.Args) > 0 && strings.HasPrefix(build.Args[0], "-") {
		return errors.New("shell build recipes must name an existing project script")
	}
	if build.Executable == "docker" && (len(build.Args) == 0 || build.Args[0] != "build") {
		return errors.New("Docker recipes must use docker build")
	}
	return nil
}
