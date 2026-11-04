import { describe, expect, it } from 'vitest'

import { createEnvelope } from './create-envelope.compute'
import { eventEnvelopeSchema } from './event-envelope.validation'

const aggregate = {
	type: 'demo-request',
	id: '701706be-7e70-4681-9608-826d3ad548b9',
}

const now = new Date('2026-11-03T14:05:00.000Z')

describe('events/events/create-envelope.compute', () => {
	it('stamps a new id, the time and the aggregate id as the correlation id', () => {
		const envelope = createEnvelope({
			type: 'demo-request.received',
			version: 1,
			aggregate,
			payload: {
				company: 'Acme',
			},
			now,
		})

		expect(envelope).toMatchObject({
			type: 'demo-request.received',
			version: 1,
			aggregate,
			occurredAt: '2026-11-03T14:05:00.000Z',
			correlationId: '701706be-7e70-4681-9608-826d3ad548b9',
			payload: {
				company: 'Acme',
			},
		})
		expect(envelope.causationId).toBeUndefined()
		expect(eventEnvelopeSchema.safeParse(envelope).success).toBe(true)
	})

	it('gives two envelopes two different ids', () => {
		const first = createEnvelope({
			type: 'demo-request.received',
			version: 1,
			aggregate,
			payload: {},
			now,
		})

		const second = createEnvelope({
			type: 'demo-request.received',
			version: 1,
			aggregate,
			payload: {},
			now,
		})

		expect(first.id).not.toBe(second.id)
	})

	it('keeps the causation id of the event that caused this one', () => {
		const envelope = createEnvelope({
			type: 'demo-request.email-sent',
			version: 1,
			aggregate,
			payload: {},
			causationId: '0d8ce74c-4756-4bfd-8d35-c820e302898e',
			now,
		})

		expect(envelope.causationId).toBe('0d8ce74c-4756-4bfd-8d35-c820e302898e')
	})
})
