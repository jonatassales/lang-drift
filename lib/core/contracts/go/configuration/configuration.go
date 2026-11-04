// Package configuration is the runtime.md contract of Microruntimes: the one
// langdrift JSON fence a console authors and the connected agent executes.
package configuration

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
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

// Native preserves an existing v2 descriptor and its owned resources on migration.
type Native struct {
	RuntimeID string `json:"runtimeId"`
}

// Markdown renders the portable, editable migration/configuration format.
func (config Configuration) Markdown() string {
	data, _ := json.MarshalIndent(config, "", "  ")

	return "# " + strings.ReplaceAll(config.Name, "\n", " ") + "\n\nConfiguration is interpreted by the connected project agent. Prose does not execute commands.\n\n```langdrift\n" + string(data) + "\n```\n"
}

// Validate bounds resources and makes all execution choices explicit.
func (config Configuration) Validate() error {
	if config.Version != Version {
		return invalid(fmt.Sprintf("version %d is required; %s", Version, replan))
	}

	if config.Revision < 1 || config.RuntimeUID == "" || config.ProjectID == "" || strings.TrimSpace(config.Name) == "" {
		return invalid("identities, revision and name are required")
	}

	if err := CheckIdentifier(config.RuntimeID); err != nil {
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

	paths := map[string]bool{}
	names := map[string]bool{}
	for _, component := range config.Components {
		if err := component.Validate(); err != nil {
			return err
		}

		if paths[component.Path] {
			return invalid("components must have unique source paths")
		}

		if names[component.Name] {
			return invalid("components must have unique names")
		}

		paths[component.Path] = true
		names[component.Name] = true
	}

	return nil
}
