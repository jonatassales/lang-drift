import type { RuntimeWidgetSnapshot } from '@data/types'

import type { RuntimeCost } from './cost-totals.types.ts'
import { isReported } from './is-reported.compute.ts'
import { peakRunning } from './peak-running.compute.ts'
import { sumReported } from './sum-reported.compute.ts'

export function runtimeCost(snapshot: RuntimeWidgetSnapshot): RuntimeCost {
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
