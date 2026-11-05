import type { FaqItemCopy } from '@routes/platform/product-page/faq-item-copy.type'

import type {
	budgetTermIds,
	chainPillarIds,
	envelopeSegmentIds,
	failureScenarioIds,
	finopsFaqIds,
	gateStepIds,
	workItems,
} from './finops-content.data'

export type EnvelopeSegmentId = (typeof envelopeSegmentIds)[number]

export interface Envelope {
	authorized: number
	protected: number
	reserved: number
	spent: number
}

export type ModelTier = 'standard' | 'deep'

export type Effort = 'low' | 'medium' | 'high'

export type CleanupState = 'done' | 'pending' | 'unknown'

export interface WorkItem {
	cleanup: CleanupState
	effort: Effort
	estimate: number
	evidencePath: string
	id: string
	observed: number | null
	reserved: number
	retries: number
	tier: ModelTier
	verified: boolean
}

export type WorkItemId = (typeof workItems)[number]['id']

// A work item whose spend is known: the provider's usage record arrived.
export interface ObservedWorkItem extends WorkItem {
	observed: number
}

export interface GateFigures {
	proposed: number
	reserved: number
	observed: number
	released: number
}

export type BudgetTermId = (typeof budgetTermIds)[number]

export type ChainPillarId = (typeof chainPillarIds)[number]

export type FailureScenarioId = (typeof failureScenarioIds)[number]

export type FinopsFaqId = (typeof finopsFaqIds)[number]

export type GateStepId = (typeof gateStepIds)[number]

export interface BudgetTermCopy {
	term: string
	title: string
	description: string
}

export interface GateStepCopy {
	eyebrow: string
	title: string
	description: string
	figure: string
}

export interface FailureScenarioCopy extends Record<ChainPillarId, string> {
	title: string
	context: string
	status: { term: string; detail: string }
}

export interface WorkItemCopy {
	owner: string
	title: string
}

export interface FinopsMessages {
	skip: string
	hero: {
		taglineStart: string
		taglineEnd: string
		caption: string
		summary: string
	}
	ledger: {
		label: string
		period: string
		authorized: string
		segments: Record<EnvelopeSegmentId, string>
		columns: {
			item: string
			estimate: string
			reserved: string
			observed: string
		}
		unknown: string
		awaitingReconciliation: string
		selectHint: string
		items: Record<WorkItemId, WorkItemCopy>
	}
	failures: {
		kicker: string
		title: string
		lead: string
		pillars: Record<ChainPillarId, string>
		scenarios: Record<FailureScenarioId, FailureScenarioCopy>
	}
	receipt: {
		selected: string
		labels: {
			estimate: string
			reserved: string
			observed: string
			modelEffort: string
			retries: string
			verification: string
			cleanup: string
			evidence: string
		}
		tiers: Record<ModelTier, string>
		efforts: Record<Effort, string>
		modelValue: string
		verification: { verified: string; pending: string }
		cleanup: Record<CleanupState, string>
		retries: string
		unknownObserved: string
	}
	budget: {
		kicker: string
		title: string
		lead: string
		terms: Record<BudgetTermId, BudgetTermCopy>
	}
	authorization: {
		kicker: string
		title: string
		lead: string
		steps: Record<GateStepId, GateStepCopy>
		released: string
		source: string
	}
	sources: {
		kicker: string
		title: string
		lead: string
		intended: string
		intendedModel: string
		items: {
			providers: { title: string; description: string }
			microruntimes: { title: string; description: string; link: string }
			evidence: { title: string; description: string }
			billing: { title: string; description: string }
		}
	}
	faq: {
		items: Record<FinopsFaqId, FaqItemCopy>
	}
}
