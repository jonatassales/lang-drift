import { describe, expect, it } from 'vitest'

import { computeEnvelope } from './compute-envelope.compute'
import { teamComputeEnvelope } from './finops-content.data'
import type { ComputeEnvelope } from './finops-content.types'

const metered = {
	ceiling: 900,
	cleanup: 'done',
	cpuHours: 120,
	gpuHours: 0,
	id: 'metered',
	idle: 100,
	idleHours: 12,
	memoryGbHours: 480,
	observed: 700,
	peakCores: 8,
	source: 'cluster/metered',
} as const

const unmetered = {
	ceiling: 1400,
	cleanup: 'unknown',
	cpuHours: null,
	gpuHours: null,
	id: 'unmetered',
	idle: 0,
	idleHours: null,
	memoryGbHours: null,
	observed: null,
	peakCores: null,
	source: 'cloud/gpu-pool',
} as const

describe('website/platform/finops/content/compute-envelope.compute', () => {
	it('splits the envelope into segments that add up to the authorized amount', () => {
		const segments = computeEnvelope(teamComputeEnvelope)

		expect(
			segments.used + segments.idle + segments.held + segments.available,
		).toBe(teamComputeEnvelope.authorized)
	})

	it('counts idle cost apart from used cost', () => {
		const envelope: ComputeEnvelope = {
			authorized: 5000,
			environments: [metered],
		}

		expect(computeEnvelope(envelope)).toEqual({
			used: 600,
			idle: 100,
			held: 0,
			available: 4300,
		})
	})

	it('holds the ceiling of an unmetered environment and uses nothing of it', () => {
		const envelope: ComputeEnvelope = {
			authorized: 5000,
			environments: [metered, unmetered],
		}

		expect(computeEnvelope(envelope)).toEqual({
			used: 600,
			idle: 100,
			held: 1400,
			available: 2900,
		})
	})

	it('leaves available at zero when spend exceeds the authorized amount', () => {
		const envelope: ComputeEnvelope = {
			authorized: 1000,
			environments: [metered, unmetered],
		}

		expect(computeEnvelope(envelope).available).toBe(0)
	})
})
