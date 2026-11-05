package configuration_test

import "testing"

// The simulation's fence holds the totals the plan prints:
// "✓ size       1450m of 2000m · 2304 of 3072 MiB · 5 one-replica Deployments · proposed shares".
func TestTheSimulationsFenceHoldsThePrintedTotals(t *testing.T) {
	t.Parallel()

	fence := simulationFence(t)
	cpuMillis, memoryMiB, oneReplica := 0, 0, 0

	for _, component := range fence.Components {
		cpuMillis += component.CPUMillis
		memoryMiB += component.MemoryMiB

		if component.Replicas == 1 {
			oneReplica++
		}
	}

	printed := cpuMillis == 1450 && memoryMiB == 2304 && oneReplica == 5
	if !printed || fence.CPUMillis != 2000 || fence.MemoryMiB != 3072 {
		t.Errorf("totals = %dm of %dm · %d of %d MiB · %d one-replica", cpuMillis, fence.CPUMillis, memoryMiB, fence.MemoryMiB, oneReplica)
	}
}
