package configuration_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/jonatassales/tao/lib/core/contracts/go/configuration"
)

// The simulation's fence breaks no rule: its validate line reads "0 errors".
func TestValidateAcceptsTheSimulationsFence(t *testing.T) {
	t.Parallel()

	if problems := simulationFence(t).Validate(); len(problems) != 0 {
		t.Errorf("problems = %v", problems)
	}
}

// Each rule of a version 3 fence, broken alone, is one problem naming its
// field.
func TestValidateNamesEachRule(t *testing.T) {
	t.Parallel()

	for _, scenario := range []struct {
		name      string
		breakRule func(*configuration.Fence)
		field     string
	}{
		{"another version", func(fence *configuration.Fence) { fence.Version = 2 }, "version"},
		{"an identifier that is not a DNS label", func(fence *configuration.Fence) { fence.RuntimeID = "Book_A_Demo" }, "runtimeId"},
		{"revision 0", func(fence *configuration.Fence) { fence.Revision = 0 }, "revision"},
		{"a cpu bound over the limit", func(fence *configuration.Fence) { fence.CPUMillis = 64001 }, "cpuMillis"},
		{"a memory bound over the limit", func(fence *configuration.Fence) { fence.MemoryMiB = 1048577 }, "memoryMiB"},
		{"no model", func(fence *configuration.Fence) { fence.Model = "" }, "model"},
		{"an effort outside the three", func(fence *configuration.Fence) { fence.Effort = "max" }, "effort"},
		{"no component", func(fence *configuration.Fence) { fence.Components = nil }, "components"},
		{"seventeen components", func(fence *configuration.Fence) { fence.Components = seventeen(fence.Components[0]) }, "components"},
		{"a path that climbs", func(fence *configuration.Fence) { fence.Components[1].Path = "../notify-api" }, "components[1].path"},
		{"a path inside .langdrift", func(fence *configuration.Fence) { fence.Components[1].Path = ".langdrift/agent" }, "components[1].path"},
		{"a repeated path", func(fence *configuration.Fence) { fence.Components[1].Path = "services/signup-api" }, "components[1].path"},
		{"a stack outside the set", func(fence *configuration.Fence) { fence.Components[2].Stack = "Kafka" }, "components[2].stack"},
		{"a role outside the set", func(fence *configuration.Fence) { fence.Components[2].Role = "queue" }, "components[2].role"},
		{"an invalid port", func(fence *configuration.Fence) { fence.Components[3].Port = 70000 }, "components[3].port"},
		{"a repeated port", func(fence *configuration.Fence) { fence.Components[3].Port = 8080 }, "components[3].port"},
		{"an invalid image", func(fence *configuration.Fence) { fence.Components[4].Image = "signup image" }, "components[4].image"},
		{"a repeated image", func(fence *configuration.Fence) { fence.Components[4].Image = fence.Components[0].Image }, "components[4].image"},
		{"no cpu share", func(fence *configuration.Fence) { fence.Components[0].CPUMillis = 0 }, "components[0].cpuMillis"},
		{"no memory share", func(fence *configuration.Fence) { fence.Components[0].MemoryMiB = 0 }, "components[0].memoryMiB"},
		{"two replicas", func(fence *configuration.Fence) { fence.Components[0].Replicas = 2 }, "components[0].replicas"},
		{"cpu shares over the bound", func(fence *configuration.Fence) { fence.CPUMillis = 1000 }, "components.cpuMillis"},
		{"memory shares over the bound", func(fence *configuration.Fence) { fence.MemoryMiB = 2048 }, "components.memoryMiB"},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			t.Parallel()

			fence := simulationFence(t)
			scenario.breakRule(&fence)

			problems := fence.Validate()
			if len(problems) != 1 || problems[0].Field != scenario.field || problems[0].Rule == "" {
				t.Errorf("problems = %v, want one on %s", problems, scenario.field)
			}
		})
	}
}

// Validation never stops at the first problem: a fence with two equal ports
// and shares over the bound has exactly two, which the plan counts.
func TestValidateReturnsEveryProblem(t *testing.T) {
	t.Parallel()

	fence := simulationFence(t)
	fence.Components[1].Port = 8080
	fence.CPUMillis = 1000

	problems := fence.Validate()

	fields := make([]string, 0, len(problems))
	for _, problem := range problems {
		fields = append(fields, problem.Field)
	}

	slices.Sort(fields)
	if !slices.Equal(fields, []string{"components.cpuMillis", "components[1].port"}) {
		t.Errorf("problems = %v", problems)
	}

	if !strings.Contains(problems[0].String(), problems[0].Field) {
		t.Errorf("a problem reads %q", problems[0].String())
	}
}
