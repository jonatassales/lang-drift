// Package configuration is the runtime.md contract of Microruntimes: the one
// langdrift JSON fence a console authors and the connected agent executes.
//
// Its lexical rules, the Check and Is functions, hold for every Microruntimes
// value wherever it is checked: in a runtime.md fence, in the CLI's domain or
// in the server that authors the fence. They need no filesystem, cluster or
// network.
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

// Markdown renders the portable, editable migration/configuration format.
// Only a budget that is not finite keeps a configuration from encoding,
// and Validate refuses one, so a failure here is a caller's broken
// invariant.
func (config Configuration) Markdown() string {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		panic(fmt.Sprintf("configuration: %q cannot be encoded: %v", config.Name, err))
	}

	return "# " + strings.ReplaceAll(config.Name, "\n", " ") + "\n\nConfiguration is interpreted by the connected project agent. Prose does not execute commands.\n\n```langdrift\n" + string(data) + "\n```\n"
}

// Validate bounds resources and makes all execution choices explicit.
func (config Configuration) Validate() error {
	if config.Version != Version {
		return invalid(fmt.Sprintf("version %d is required; %s", Version, replan))
	}

	identityMissing := config.Revision < 1 ||
		config.RuntimeUID == "" ||
		config.ProjectID == "" ||
		strings.TrimSpace(config.Name) == ""

	if identityMissing {
		return invalid("identities, revision and name are required")
	}

	if err := CheckIdentifier(config.RuntimeID); err != nil {
		return err
	}

	resourcesOutOfBounds := config.CPUMillis < 1 ||
		config.CPUMillis > 64_000 ||
		config.MemoryMiB < 1 ||
		config.MemoryMiB > 1_048_576

	if resourcesOutOfBounds {
		return invalid("CPU and memory must be positive bounded per-execution limits")
	}

	budgetOutOfBounds := config.LLMBudgetUSD != nil &&
		(math.IsNaN(*config.LLMBudgetUSD) ||
			*config.LLMBudgetUSD < 0.01 ||
			*config.LLMBudgetUSD > 1_000_000)

	if budgetOutOfBounds {
		return invalid("LLM budget must be positive bounded USD")
	}

	effortKnown := config.Effort == "low" ||
		config.Effort == "medium" ||
		config.Effort == "high"
	choiceMissing := config.Model == "" || !effortKnown

	if choiceMissing {
		return invalid("explicit agent model and low, medium, or high effort are required")
	}

	componentCountOutOfBounds := len(config.Components) == 0 ||
		len(config.Components) > 16

	if componentCountOutOfBounds {
		return invalid("select between 1 and 16 components")
	}

	nativeChanged := config.Native != nil &&
		config.Native.RuntimeID != config.RuntimeID

	if nativeChanged {
		return invalid("migration must preserve the native runtime ID")
	}

	envelopeTooSmall := config.CPUMillis < len(config.Components) ||
		config.MemoryMiB < len(config.Components)

	if envelopeTooSmall {
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
