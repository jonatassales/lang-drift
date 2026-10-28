import assert from 'node:assert/strict'
import test from 'node:test'

import type { Execution, RuntimeWidgetSnapshot } from '@data/types'

import { costTotals } from './cost-totals.ts'

function execution(values: Partial<Execution>): Execution {
	return {
		cpuUsed: null,
		memoryUsedMiB: null,
		...values,
	} as Execution
}

function widget(
	id: string,
	executions: Execution[],
	llmBudgetUsd?: number,
): RuntimeWidgetSnapshot {
	return {
		runtime: {
			id,
			configuration: {
				cpu: 2,
				memoryMiB: 3072,
				llmBudgetUsd,
			},
		},
		executions,
	} as RuntimeWidgetSnapshot
}

test('adds only measured LLM costs and counts the unknown ones', () => {
	const totals = costTotals([
		widget('first', [
			execution({ llmCostUsd: 2.5 }),
			execution({ llmCostUsd: null }),
		]),
		widget('second', [execution({}), execution({ llmCostUsd: 0 })]),
	])

	assert.equal(totals.knownLlmCostUsd, 2.5)
	assert.equal(totals.unknownLlmCostCount, 2)
	assert.equal(totals.executionCount, 4)
})

test('keeps a runtime with no measurement unknown, never zero', () => {
	const [runtime] = costTotals([widget('first', [execution({})])]).runtimes

	assert.equal(runtime?.knownLlmCostUsd, null)
	assert.equal(runtime?.cpuUsed, null)
	assert.equal(runtime?.memoryUsedMiB, null)
})

test('sums measured CPU and memory across a runtime executions', () => {
	const [runtime] = costTotals([
		widget('first', [
			execution({
				cpuUsed: 0.5,
				memoryUsedMiB: 512,
			}),
			execution({
				cpuUsed: 1,
				memoryUsedMiB: null,
			}),
		]),
	]).runtimes

	assert.equal(runtime?.cpuUsed, 1.5)
	assert.equal(runtime?.memoryUsedMiB, 512)
})

test('adds the configured LLM budgets and counts runtimes without one', () => {
	const totals = costTotals([
		widget('first', [], 10),
		widget('second', [], 5),
		widget('third', []),
	])

	assert.equal(totals.budgetedLlmUsd, 15)
	assert.equal(totals.unbudgetedRuntimeCount, 1)
})
