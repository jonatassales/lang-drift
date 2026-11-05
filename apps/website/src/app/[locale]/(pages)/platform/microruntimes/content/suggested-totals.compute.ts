import type { PlanComponent, Resources } from './microruntimes-content'

// The totals the plan proposes: the sum of the components' shares.
export function suggestedTotals(
	components: readonly PlanComponent[],
): Resources {
	return components.reduce<Resources>(
		(totals, component) => ({
			cpuMillis: totals.cpuMillis + component.cpuMillis,
			memoryMiB: totals.memoryMiB + component.memoryMiB,
		}),
		{
			cpuMillis: 0,
			memoryMiB: 0,
		},
	)
}
