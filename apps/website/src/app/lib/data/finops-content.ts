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

export const chainPillarIds = [
	'intent',
	'delivery',
	'evidence',
	'outcome',
] as const
export const chainScenarioIds = [
	'within-envelope',
	'pause-and-replan',
	'unknown-spend',
] as const
export const finopsFaqIds = [
	'budget',
	'spending',
	'unknown',
	'authorization',
	'availability',
] as const

export const budgetCardIds = ['envelope', 'allocation', 'reservation'] as const
export const authorizationCardIds = ['recommend', 'gate', 'reconcile'] as const

export type ChainPillarId = (typeof chainPillarIds)[number]
export type ChainScenarioId = (typeof chainScenarioIds)[number]
export type FinopsFaqId = (typeof finopsFaqIds)[number]

export interface FinopsCard {
	description: string
	eyebrow: string
	title: string
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
	chain: {
		title: string
		lead: string
		selectorLabel: string
		pillars: Record<ChainPillarId, string>
		scenarios: Record<
			ChainScenarioId,
			{ tabLabel: string; context: string } & Record<ChainPillarId, string>
		>
		outro: string
		outroLink: string
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
		cards: Record<(typeof authorizationCardIds)[number], FinopsCard>
	}
	report: {
		kicker: string
		title: string
		lead: string
		weekly: string
		period: string
		rows: Record<
			'envelopeUsed' | 'largestConsumer' | 'unknownCharges' | 'attention',
			{ term: string; detail: string }
		>
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
