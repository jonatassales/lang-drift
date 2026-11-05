import type { DetailTone } from 'web/ui/detail-list'

import type {
	Envelope,
	FailureScenarioId,
	WorkItem,
	WorkItemId,
} from './finops-content.types'

// FinOps content: the ids, amounts and states the FinOps page and
// its demos render. Sentences live in the finops catalog; this module holds
// what a reader does not read as prose. Amounts are simulated US dollars, and
// an observed spend of null is unknown, never zero.

export const envelopeSegmentIds = [
	'spent',
	'reserved',
	'protected',
	'available',
] as const

export const envelope: Envelope = {
	authorized: 12000,
	protected: 1200,
	reserved: 2150,
	spent: 6480,
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

export const initialWorkItemId: WorkItemId = 'checkout-retries'

// The work item the gate section walks through: the first in the ledger, the
// one it selects first, delivered inside the envelope.
export const gateExample = workItems[0]

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
	'license',
	'availability',
] as const

export const budgetTermIds = ['envelope', 'allocation', 'reservation'] as const
export const gateStepIds = ['recommend', 'gate', 'reconcile'] as const
export const sourceIds = [
	'providers',
	'microruntimes',
	'evidence',
	'billing',
] as const

// The state each failure leaves behind: a paused job needs a decision; an
// unknown charge waits for a usage record.
export const failureStatusTones: Record<FailureScenarioId, DetailTone> = {
	'pause-and-replan': 'danger',
	'unknown-spend': 'warning',
}
