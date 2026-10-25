import { describe, expect, it } from 'vitest'

import {
	formatDescriptor,
	lifecycleStageIds,
	runtimeDescriptor,
} from './microruntimes-content'

describe('illustrative Microruntimes content', () => {
	it('keeps the six lifecycle stages in order', () => {
		expect(lifecycleStageIds).toEqual([
			'configure',
			'validate',
			'prepare',
			'execute',
			'inspect',
			'cleanup',
		])
	})

	it('keeps the runtime file within the limits the ex agent accepts', () => {
		const { components, cpuMillis, memoryMiB } = runtimeDescriptor
		expect(components.length).toBeGreaterThanOrEqual(1)
		expect(components.length).toBeLessThanOrEqual(16)
		expect(cpuMillis).toBeLessThanOrEqual(64_000)
		expect(cpuMillis).toBeGreaterThanOrEqual(components.length)
		expect(memoryMiB).toBeLessThanOrEqual(1_048_576)
		expect(['low', 'medium', 'high']).toContain(runtimeDescriptor.effort)
	})

	it('uses only committed scripts or allowlisted build tools', () => {
		const allowed = ['docker', 'go', 'npm', 'pnpm', 'make']
		for (const { build } of runtimeDescriptor.components) {
			const committedScript = build.executable.startsWith('./')
			expect(committedScript || allowed.includes(build.executable)).toBe(true)
		}
	})

	it('formats the descriptor as parseable JSON with short arrays inline', () => {
		const text = formatDescriptor(runtimeDescriptor)
		expect(JSON.parse(text)).toEqual(runtimeDescriptor)
		expect(text).toContain('"args": ["build", "-t", "storefront:local", "."]')
		expect(text).toContain('"args": []')
	})
})
