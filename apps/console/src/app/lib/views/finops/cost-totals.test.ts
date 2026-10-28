import assert from 'node:assert/strict'
import test from 'node:test'

import type { Execution, RuntimeWidgetSnapshot } from '@data/types'

import { costTotals } from './cost-totals.ts'

function execution(values: Partial<Execution>): Execution {
	return {
		state: 'running',
		cpuUsed: null,
		memoryUsedMiB: null,
		...values,
	} as Execution
}

function widget(id: string, executions: Execution[]): RuntimeWidgetSnapshot {
	return {
		runtime: {
			id,
			configuration: {
				cpu: 2,
				memoryMiB: 3072,
			},
		},
		executions,
	} as RuntimeWidgetSnapshot
}

test('adds only reported LLM costs and counts the unreported ones', () => {
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

test('keeps a cost nothing reported unknown, never zero', () => {
	const totals = costTotals([widget('first', [execution({})])])
	const [runtime] = totals.runtimes

	assert.equal(totals.knownLlmCostUsd, null)
	assert.equal(runtime?.knownLlmCostUsd, null)
	assert.equal(runtime?.unknownLlmCostCount, 1)
	assert.equal(runtime?.cpuUsed, null)
	assert.equal(runtime?.memoryUsedMiB, null)
})

test('keeps the total unknown when no runtime has an execution', () => {
	const totals = costTotals([widget('first', [])])

	assert.equal(totals.knownLlmCostUsd, null)
	assert.equal(totals.executionCount, 0)
})

test('reads CPU and memory as the highest running observation, never a sum', () => {
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
			execution({
				state: 'stopped',
				cpuUsed: 4,
				memoryUsedMiB: 4096,
			}),
		]),
	]).runtimes

	assert.equal(runtime?.cpuUsed, 1)
	assert.equal(runtime?.memoryUsedMiB, 512)
})
