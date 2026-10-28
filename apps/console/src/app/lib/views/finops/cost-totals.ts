import type { Execution, RuntimeWidgetSnapshot } from '@data/types'

// One runtime's FinOps reading. A value no execution reported stays null:
// unknown, never zero.
export interface RuntimeCost {
	snapshot: RuntimeWidgetSnapshot
	knownLlmCostUsd: number | null
	unknownLlmCostCount: number
	executionCount: number
	cpuUsed: number | null
	memoryUsedMiB: number | null
}

export interface CostTotals {
	runtimes: RuntimeCost[]
	knownLlmCostUsd: number | null
	unknownLlmCostCount: number
	executionCount: number
}

function isReported(value: number | null | undefined): value is number {
	return typeof value === 'number'
}

function sumReported(values: (number | null | undefined)[]): number | null {
	const reported = values.filter(isReported)
	if (reported.length === 0) return null

	return reported.reduce((total, value) => total + value, 0)
}

// The highest observation among running executions, as the runtime widget
// reads it: never a sum across machines.
function peakRunning(
	executions: Execution[],
	key: 'cpuUsed' | 'memoryUsedMiB',
): number | null {
	const observed = executions
		.filter((execution) => execution.state === 'running')
		.map((execution) => execution[key])
		.filter(isReported)
	if (observed.length === 0) return null

	return Math.max(...observed)
}

function runtimeCost(snapshot: RuntimeWidgetSnapshot): RuntimeCost {
	const { executions } = snapshot
	const llmCosts = executions.map((execution) => execution.llmCostUsd)

	return {
		snapshot,
		knownLlmCostUsd: sumReported(llmCosts),
		unknownLlmCostCount: llmCosts.filter((cost) => !isReported(cost)).length,
		executionCount: executions.length,
		cpuUsed: peakRunning(executions, 'cpuUsed'),
		memoryUsedMiB: peakRunning(executions, 'memoryUsedMiB'),
	}
}

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
			(count, runtime) => count + runtime.unknownLlmCostCount,
			0,
		),
		executionCount: runtimes.reduce(
			(count, runtime) => count + runtime.executionCount,
			0,
		),
	}
}
