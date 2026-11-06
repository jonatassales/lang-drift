import type { FaqItemCopy } from '@routes/platform/product-page/faq-item-copy.type'

import type {
	chainPillarIds,
	computeSegmentIds,
	computeTermIds,
	environments,
	failureScenarioIds,
	finopsFaqIds,
	sourceIds,
} from './finops-content.data'

export type ComputeSegmentId = (typeof computeSegmentIds)[number]

export type CleanupState = 'done' | 'pending' | 'unknown'

// One environment the team's hardware ran in. Its usage and observed cost are
// null while its meter does not report: unknown, never zero.
export interface Environment {
	ceiling: number
	cleanup: CleanupState
	cpuHours: number | null
	gpuHours: number | null
	id: string
	idle: number
	idleHours: number | null
	memoryGbHours: number | null
	observed: number | null
	peakCores: number | null
	source: string
}

export interface ComputeEnvelope {
	authorized: number
	environments: readonly Environment[]
}

export type ComputeEnvelopeSegments = Record<ComputeSegmentId, number>

export type EnvironmentId = (typeof environments)[number]['id']

export type ChainPillarId = (typeof chainPillarIds)[number]

export type FailureScenarioId = (typeof failureScenarioIds)[number]

export type ComputeTermId = (typeof computeTermIds)[number]

export type SourceId = (typeof sourceIds)[number]

export type FinopsFaqId = (typeof finopsFaqIds)[number]

export interface EnvironmentCopy {
	owner: string
	title: string
}

export interface ComputeTermCopy {
	term: string
	title: string
	description: string
}

export interface FailureScenarioCopy extends Record<ChainPillarId, string> {
	context: string
	tabLabel: string
}

export interface SourceCopy {
	title: string
	description: string
}

export interface FinopsMessages {
	skip: string
	hero: {
		taglineStart: string
		taglineEnd: string
		summary: string
	}
	ledger: {
		label: string
		period: string
		authorized: string
		segments: Record<ComputeSegmentId, string>
		columns: {
			environment: string
			ceiling: string
			observed: string
		}
		unknown: string
		meterSilent: string
		usage: string
		selectHint: string
		environments: Record<EnvironmentId, EnvironmentCopy>
	}
	receipt: {
		selected: string
		labels: {
			ceiling: string
			observed: string
			cpu: string
			memory: string
			gpu: string
			peak: string
			idle: string
			cleanup: string
			source: string
		}
		cpuHours: string
		memoryHours: string
		gpuHours: string
		peakCores: string
		idleValue: string
		unknownUsage: string
		unknownObserved: string
		cleanup: Record<CleanupState, string>
	}
	failures: {
		kicker: string
		title: string
		lead: string
		selector: string
		pillars: Record<ChainPillarId, string>
		scenarios: Record<FailureScenarioId, FailureScenarioCopy>
	}
	budget: {
		kicker: string
		title: string
		lead: string
		terms: Record<ComputeTermId, ComputeTermCopy>
	}
	sources: {
		kicker: string
		title: string
		lead: string
		intended: string
		intendedModel: string
		microruntimesLink: string
		items: Record<SourceId, SourceCopy>
	}
	faq: {
		items: Record<FinopsFaqId, FaqItemCopy>
	}
}
