import { describe, expect, it } from 'vitest'

import {
	envelope,
	envelopeSegments,
	failureScenarioIds,
	failureStatusTones,
	gateExample,
	gateFigures,
	initialWorkItemId,
	workItems,
} from './finops-content'

describe('illustrative FinOps content', () => {
	it('splits the envelope into segments that sum to the authorized amount', () => {
		const segments = envelopeSegments(envelope)
		const total =
			segments.spent +
			segments.reserved +
			segments.protected +
			segments.available
		expect(total).toBe(envelope.authorized)
		expect(segments.available).toBeGreaterThan(0)
	})

	it('keeps exactly one work item with unknown spend', () => {
		const unknown = workItems.filter((item) => item.observed === null)
		expect(unknown.map((item) => item.id)).toEqual(['invoice-export'])
	})

	it('reserves at least the estimate of every work item', () => {
		for (const item of workItems) {
			expect(item.reserved).toBeGreaterThanOrEqual(item.estimate)
		}
	})

	it('walks the gate through the work item the ledger selects first', () => {
		expect(gateExample.id).toBe(initialWorkItemId)
		expect(gateFigures(gateExample)).toEqual({
			proposed: 1800,
			reserved: 2600,
			observed: 2460,
			released: 140,
		})
	})

	it('leaves the within-envelope case to the hero and ends each failure in a warning or a danger', () => {
		expect(failureScenarioIds).not.toContain('within-envelope')
		for (const id of failureScenarioIds) {
			expect(['warning', 'danger']).toContain(failureStatusTones[id])
		}
	})
})
