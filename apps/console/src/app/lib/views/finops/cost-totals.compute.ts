import type { RuntimeWidgetSnapshot } from '@data/types'

import type { CostTotals } from './cost-totals.types.ts'
import { runtimeCost } from './runtime-cost.compute.ts'
import { sumReported } from './sum-reported.compute.ts'

// The FinOps overview's numbers, from the runtime widgets the stub serves.
// LLM cost adds up because each execution reports its own cumulative cost.
export function costTotals(widgets: RuntimeWidgetSnapshot[]): CostTotals {
	const runtimes = widgets.map(runtimeCost)

	return {
		runtimes,
		knownLlmCostUsd: sumReported(
			runtimes.map((runtime) => runtime.knownLlmCostUsd),
		),
		unknownLlmCostCount: runtimes.reduce(
			(total, runtime) => total + runtime.unknownLlmCostCount,
			0,
		),
		executionCount: runtimes.reduce(
			(total, runtime) => total + runtime.executionCount,
			0,
		),
	}
}
