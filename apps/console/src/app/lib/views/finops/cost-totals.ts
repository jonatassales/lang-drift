import type { RuntimeWidgetSnapshot } from '@data/types'

// One runtime's FinOps reading across its executions. A total that no
// execution measured stays null: unknown, never zero.
export interface RuntimeCost {
	snapshot: RuntimeWidgetSnapshot
	knownLlmCostUsd: number | null
	unknownLlmCostCount: number
	cpuUsed: number | null
	memoryUsedMiB: number | null
}

export interface CostTotals {
	runtimes: RuntimeCost[]
	knownLlmCostUsd: number
	unknownLlmCostCount: number
	executionCount: number
	budgetedLlmUsd: number
	unbudgetedRuntimeCount: number
}

function isMeasured(value: number | null | undefined): value is number {
	return typeof value === 'number'
}

function sumMeasured(values: (number | null | undefined)[]): number | null {
	const measured = values.filter(isMeasured)
	if (measured.length === 0) return null

	return measured.reduce((total, value) => total + value, 0)
}

function runtimeCost(snapshot: RuntimeWidgetSnapshot): RuntimeCost {
	const { executions } = snapshot
	const llmCosts = executions.map((execution) => execution.llmCostUsd)

	return {
		snapshot,
		knownLlmCostUsd: sumMeasured(llmCosts),
		unknownLlmCostCount: llmCosts.filter((cost) => !isMeasured(cost)).length,
		cpuUsed: sumMeasured(executions.map((execution) => execution.cpuUsed)),
		memoryUsedMiB: sumMeasured(
			executions.map((execution) => execution.memoryUsedMiB),
		),
	}
}

// The FinOps overview's numbers, from the runtime widgets the stub serves.
export function costTotals(widgets: RuntimeWidgetSnapshot[]): CostTotals {
	const runtimes = widgets.map(runtimeCost)
	const budgets = widgets.map(
		(widget) => widget.runtime.configuration.llmBudgetUsd,
	)

	return {
		runtimes,
		knownLlmCostUsd:
			sumMeasured(runtimes.map((runtime) => runtime.knownLlmCostUsd)) ?? 0,
		unknownLlmCostCount: runtimes.reduce(
			(count, runtime) => count + runtime.unknownLlmCostCount,
			0,
		),
		executionCount: widgets.reduce(
			(count, widget) => count + widget.executions.length,
			0,
		),
		budgetedLlmUsd: sumMeasured(budgets) ?? 0,
		unbudgetedRuntimeCount: budgets.filter((budget) => !isMeasured(budget))
			.length,
	}
}
