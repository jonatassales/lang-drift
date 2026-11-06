import type {
	ComputeEnvelope,
	Environment,
	EnvironmentId,
} from './finops-content.types'

// FinOps content: the ids, amounts and usage the FinOps page and its demos
// render. Sentences live in the finops catalog; this module holds what a
// reader does not read as prose. Amounts are US dollars, and usage or an
// observed cost of null is unknown, never zero.

export const computeSegmentIds = ['used', 'idle', 'held', 'available'] as const

export const environments = [
	{
		ceiling: 420,
		cleanup: 'done',
		cpuHours: 96,
		gpuHours: 0,
		id: 'checkout-microruntime',
		idle: 18,
		idleHours: 3,
		memoryGbHours: 384,
		observed: 312,
		peakCores: 4,
		source: 'microruntimes/checkout-retries',
	},
	{
		ceiling: 1600,
		cleanup: 'done',
		cpuHours: 1480,
		gpuHours: 0,
		id: 'ci-runners',
		idle: 210,
		idleHours: 41,
		memoryGbHours: 2960,
		observed: 1240,
		peakCores: 16,
		source: 'ci/runner-minutes',
	},
	{
		ceiling: 1800,
		cleanup: 'pending',
		cpuHours: 2210,
		gpuHours: 0,
		id: 'preview-cluster',
		idle: 640,
		idleHours: 96,
		memoryGbHours: 8840,
		observed: 1530,
		peakCores: 24,
		source: 'cluster/preview',
	},
	{
		ceiling: 1400,
		cleanup: 'unknown',
		cpuHours: null,
		gpuHours: null,
		id: 'embedding-gpu-job',
		idle: 0,
		idleHours: null,
		memoryGbHours: null,
		observed: null,
		peakCores: null,
		source: 'cloud/gpu-pool',
	},
] as const satisfies readonly Environment[]

export const teamComputeEnvelope: ComputeEnvelope = {
	authorized: 8000,
	environments,
}

export const initialEnvironmentId: EnvironmentId = 'checkout-microruntime'

export const chainPillarIds = [
	'intent',
	'delivery',
	'evidence',
	'outcome',
] as const
// The two ways hardware leaves the plan: an environment that outlives its
// run, and a meter that stops reporting.
export const failureScenarioIds = [
	'orphaned-environment',
	'silent-meter',
] as const
export const computeTermIds = ['envelope', 'quota', 'cleanup'] as const
export const sourceIds = [
	'microruntimes',
	'cluster',
	'runners',
	'billing',
] as const
export const finopsFaqIds = [
	'scope',
	'llmops',
	'unknown',
	'idle',
	'license',
	'availability',
] as const
