package configuration

import (
	"slices"
	"strconv"
	"strings"
)

// headerProblems are the problems of a fence's own fields: its version,
// identifier and revision, its bound, its suggestion and how many
// components it holds.
func headerProblems(fence Fence) []Problem {
	var problems []Problem

	if fence.Version != FenceVersion {
		problems = append(problems, problemAt("version", "must be 3; "+replan))
	}

	if CheckIdentifier(fence.RuntimeID) != nil {
		problems = append(problems, problemAt("runtimeId", "must be a lowercase DNS label"))
	}

	if fence.Revision < 1 {
		problems = append(problems, problemAt("revision", "must be at least 1"))
	}

	if fence.CPUMillis < 1 || fence.CPUMillis > maximumFenceCPUMillis {
		problems = append(problems, problemAt("cpuMillis", "must be 1 to "+strconv.Itoa(maximumFenceCPUMillis)+" millicores"))
	}

	if fence.MemoryMiB < 1 || fence.MemoryMiB > maximumFenceMemoryMiB {
		problems = append(problems, problemAt("memoryMiB", "must be 1 to "+strconv.Itoa(maximumFenceMemoryMiB)+" MiB"))
	}

	if fence.Model == "" {
		problems = append(problems, problemAt("model", "must name a model"))
	}

	if !slices.Contains(FenceEfforts, fence.Effort) {
		problems = append(problems, problemAt("effort", "must be one of "+strings.Join(FenceEfforts, ", ")))
	}

	if len(fence.Components) < 1 || len(fence.Components) > maximumFenceComponents {
		problems = append(problems, problemAt("components", "must hold 1 to "+strconv.Itoa(maximumFenceComponents)+" components"))
	}

	return problems
}
