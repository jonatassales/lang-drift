import assert from 'node:assert/strict'
import test from 'node:test'

import { parseStubFrame } from './stub-frame.ts'

const now = Date.now()
const timestamp = new Date(now).toISOString()
function frame() {
	return {
		source: 'stub',
		kind: 'runtime-widgets',
		sequence: 2,
		generatedAt: timestamp,
		widgets: [
			{
				source: 'stub',
				runtime: {
					id: 'runtime',
					productId: 'product',
					nativeRuntimeId: 'runtime',
					name: 'Seed',
					revision: 1,
					desiredState: 'running',
					components: [
						{
							path: 'api',
							name: 'API',
							kind: 'dockerfile',
						},
					],
					configuration: {
						cpu: 2,
						memoryMiB: 512,
						model: 'seed/model',
						effort: 'high',
					},
				},
				environments: [
					{
						id: 'machine',
						runtimeId: 'runtime',
						name: 'Seed machine',
						lastSeenAt: timestamp,
						connection: 'connected',
						cleanupPending: false,
					},
				],
				executions: [
					{
						id: 'execution',
						runtimeId: 'runtime',
						environmentId: 'machine',
						state: 'running',
						updatedAt: timestamp,
						observedAt: timestamp,
						cpuUsed: 0.5,
						memoryUsedMiB: 128,
						llmCostUsd: null,
						downtimeStartedAt: null,
						downtimeMilliseconds: 0,
						downtimeQuality: 'observed',
					},
				],
				startedAt: {
					execution: timestamp,
				},
				tone: 'green',
				partial: true,
			},
		],
	}
}

test('WebRTC accepts a complete fresh frame, retaining unknown LLM values', () => {
	const parsed = parseStubFrame(JSON.stringify(frame()), 1, now)
	assert.equal(parsed?.sequence, 2)
	assert.equal(parsed?.widgets[0].executions[0].llmCostUsd, null)
})

test('WebRTC accepts dense stack-tagged graphs within transport bounds', () => {
	const data = frame()
	const components = Array.from(
		{
			length: 14,
		},
		(_, index) => ({
			path: `service-${index}`,
			name: `Service ${index}`,
			kind: 'dockerfile',
			stacks: ['nextjs', 'react', 'nodejs'],
		}),
	)
	data.widgets[0].runtime.components = components
	Object.assign(data.widgets[0], {
		componentLinks: components.slice(1).flatMap((component, index) => [
			{
				from: components[index].path,
				to: component.path,
			},
			{
				from: component.path,
				to: components[0].path,
			},
		]),
	})
	assert.equal(
		parseStubFrame(JSON.stringify(data), 1, now)?.widgets[0].runtime.components
			.length,
		14,
	)
	components[0].stacks = ['nodejs', 'react', 'nextjs', 'go', 'python']
	assert.equal(parseStubFrame(JSON.stringify(data), 1, now), null)
	components[0].stacks = ['custom-stack']
	assert.ok(parseStubFrame(JSON.stringify(data), 1, now))
	components.push(...components.slice(0, 3))
	assert.equal(parseStubFrame(JSON.stringify(data), 1, now), null)
})

test('WebRTC rejects incomplete, malformed, foreign and replayed frames', () => {
	for (const mutate of [
		(data: ReturnType<typeof frame>) => {
			Reflect.deleteProperty(data.widgets[0].runtime, 'configuration')
		},
		(data: ReturnType<typeof frame>) => {
			Reflect.deleteProperty(data.widgets[0].runtime, 'components')
		},
		(data: ReturnType<typeof frame>) => {
			data.widgets[0].executions[0].observedAt = 'not-a-date'
		},
		(data: ReturnType<typeof frame>) => {
			Object.assign(data.widgets[0].executions[0], {
				cpuUsed: '50',
			})
		},
		(data: ReturnType<typeof frame>) => {
			Object.assign(data.widgets[0].runtime.configuration, {
				llmBudgetUsd: null,
			})
		},
		(data: ReturnType<typeof frame>) => {
			data.widgets[0].executions[0].environmentId = 'foreign'
		},
		(data: ReturnType<typeof frame>) => {
			data.sequence = 1
		},
		(data: ReturnType<typeof frame>) => {
			data.generatedAt = new Date(now - 6000).toISOString()
		},
		(data: ReturnType<typeof frame>) => {
			data.widgets.push(data.widgets[0])
		},
	]) {
		const value = frame()
		mutate(value)
		assert.equal(parseStubFrame(JSON.stringify(value), 1, now), null)
	}
	assert.equal(parseStubFrame('invalid', 0, now), null)
})
