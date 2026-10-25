import { describe, expect, it } from 'vitest'

import { envelope, envelopeSegments, workItems } from './finops-content'

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
})
