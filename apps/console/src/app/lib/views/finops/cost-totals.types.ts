import type { RuntimeWidgetSnapshot } from '@data/types'

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
