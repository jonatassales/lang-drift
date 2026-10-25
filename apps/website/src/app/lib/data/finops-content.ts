import type { DetailTone } from '@langdrift/web/ui/detail-list'

// Illustrative FinOps content: the ids, amounts and states the FinOps page and
// its demos render. Sentences live in the finops catalog; this module holds
// what a reader does not read as prose. Amounts are simulated US dollars, and
// an observed spend of null is unknown, never zero.

// Amounts on the FinOps page are simulated US dollars, whole only.
export function dollars(locale: string) {
	return new Intl.NumberFormat(locale, {
		currency: 'USD',
		maximumFractionDigits: 0,
		style: 'currency',
	})
}

export const envelopeSegmentIds = [
	'spent',
	'reserved',
	'protected',
	'available',
] as const

export type EnvelopeSegmentId = (typeof envelopeSegmentIds)[number]

export interface Envelope {
	authorized: number
	protected: number
	reserved: number
	spent: number
}

export const envelope: Envelope = {
	authorized: 12000,
	protected: 1200,
	reserved: 2150,
	spent: 6480,
}

export function envelopeSegments(
	funding: Envelope,
): Record<EnvelopeSegmentId, number> {
	const { authorized, protected: protectedAmount, reserved, spent } = funding

	return {
		spent,
		reserved,
		protected: protectedAmount,
		available: authorized - spent - reserved - protectedAmount,
	}
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

export const workItems = [
	{
		cleanup: 'done',
		effort: 'high',
		estimate: 1800,
		evidencePath: 'receipts/2026-10-14/checkout-retries',
		id: 'checkout-retries',
		observed: 2460,
		reserved: 2600,
		retries: 2,
		tier: 'deep',
		verified: true,
	},
	{
		cleanup: 'unknown',
		effort: 'medium',
		estimate: 900,
		evidencePath: 'receipts/2026-10-21/invoice-export',
		id: 'invoice-export',
		observed: null,
		reserved: 1100,
		retries: 1,
		tier: 'standard',
		verified: false,
	},
	{
		cleanup: 'done',
		effort: 'medium',
		estimate: 2400,
		evidencePath: 'receipts/2026-10-17/tenant-migration',
		id: 'tenant-migration',
		observed: 2210,
		reserved: 2400,
		retries: 0,
		tier: 'deep',
		verified: true,
	},
	{
		cleanup: 'done',
		effort: 'low',
		estimate: 300,
		evidencePath: 'receipts/2026-10-20/onboarding-copy',
		id: 'onboarding-copy',
		observed: 240,
		reserved: 350,
		retries: 0,
		tier: 'standard',
		verified: true,
	},
] as const satisfies readonly WorkItem[]

export type WorkItemId = (typeof workItems)[number]['id']

export const initialWorkItemId: WorkItemId = 'checkout-retries'

// A work item whose spend is known: the provider's usage record arrived.
export interface ObservedWorkItem extends WorkItem {
	observed: number
}

// The work item the gate section walks through: the first in the ledger, the
// one it selects first, delivered inside the envelope.
export const gateExample = workItems[0]

export interface GateFigures {
	proposed: number
	reserved: number
	observed: number
	released: number
}

// The amounts a work item carries through the gate: the estimate the agent
// proposed, the hold the gate took before the first paid call, the spend the
// provider reported and what returned to the envelope once it was reconciled.
// Only a work item with observed spend has them; unknown is never zero.
export function gateFigures(item: ObservedWorkItem): GateFigures {
	return {
		proposed: item.estimate,
		reserved: item.reserved,
		observed: item.observed,
		released: item.reserved - item.observed,
	}
}

export const chainPillarIds = [
	'intent',
	'delivery',
	'evidence',
	'outcome',
] as const
// The two ways delivery leaves the plan. The case that lands inside the
// envelope is the hero's receipt and the gate's figures, not a scenario.
export const failureScenarioIds = ['pause-and-replan', 'unknown-spend'] as const
export const finopsFaqIds = [
	'budget',
	'spending',
	'unknown',
	'authorization',
	'availability',
] as const

export const budgetCardIds = ['envelope', 'allocation', 'reservation'] as const
export const gateStepIds = ['recommend', 'gate', 'reconcile'] as const

export type ChainPillarId = (typeof chainPillarIds)[number]
export type FailureScenarioId = (typeof failureScenarioIds)[number]
export type FinopsFaqId = (typeof finopsFaqIds)[number]
export type GateStepId = (typeof gateStepIds)[number]

// The state each failure leaves behind: a paused job needs a decision; an
// unknown charge waits for a usage record.
export const failureStatusTones: Record<FailureScenarioId, DetailTone> = {
	'pause-and-replan': 'danger',
	'unknown-spend': 'warning',
}

export interface FinopsCard {
	description: string
	eyebrow: string
	title: string
}

export interface GateStepCopy extends FinopsCard {
	figure: string
}

export interface FailureScenarioCopy extends Record<ChainPillarId, string> {
	title: string
	context: string
	status: { term: string; detail: string }
}

export interface FinopsMessages {
	skip: string
	hero: {
		eyebrow: string
		titleStart: string
		titleEnd: string
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
		items: Record<WorkItemId, { title: string; owner: string }>
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
		cards: Record<(typeof budgetCardIds)[number], FinopsCard>
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
		intendedModel: string
		items: {
			providers: { title: string; description: string }
			microruntimes: { title: string; description: string; link: string }
			evidence: { title: string; description: string }
			billing: { title: string; description: string }
		}
	}
	faq: {
		title: string
		description: string
		items: Record<FinopsFaqId, { question: string; answer: string }>
	}
}
