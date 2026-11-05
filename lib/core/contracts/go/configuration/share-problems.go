package configuration

// shareProblems are the components' shares summed against the bound: a
// draft never proposes more than the design doc allows.
func shareProblems(fence Fence) []Problem {
	cpuMillis, memoryMiB := 0, 0

	for _, component := range fence.Components {
		cpuMillis += component.CPUMillis
		memoryMiB += component.MemoryMiB
	}

	var problems []Problem

	if cpuMillis > fence.CPUMillis {
		problems = append(problems, problemAt("components.cpuMillis", "the shares exceed the bound"))
	}

	if memoryMiB > fence.MemoryMiB {
		problems = append(problems, problemAt("components.memoryMiB", "the shares exceed the bound"))
	}

	return problems
}
