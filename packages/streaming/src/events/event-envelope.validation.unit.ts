import { describe, expect, it } from 'vitest'

import { eventEnvelopeSchema } from './event-envelope.validation.js'

const envelope = {
	id: '0d8ce74c-4756-4bfd-8d35-c820e302898e',
	type: 'demo-request.received',
	version: 1,
	aggregate: {
		type: 'demo-request',
		id: '701706be-7e70-4681-9608-826d3ad548b9',
	},
	occurredAt: '2026-11-03T14:05:00.000Z',
	correlationId: '701706be-7e70-4681-9608-826d3ad548b9',
	payload: {
		company: 'Acme',
	},
}

describe('events/events/event-envelope.validation', () => {
	it('accepts an envelope with an id, a type, a version, an aggregate, a time, a correlation id and a payload', () => {
		const result = eventEnvelopeSchema.safeParse(envelope)

		expect(result.success).toBe(true)
	})

	it('accepts a causation id when a later event names the one that caused it', () => {
		const result = eventEnvelopeSchema.safeParse({
			...envelope,
			causationId: '0d8ce74c-4756-4bfd-8d35-c820e302898e',
		})

		expect(result.success).toBe(true)
	})

	it('refuses an aggregate without an id', () => {
		const result = eventEnvelopeSchema.safeParse({
			...envelope,
			aggregate: {
				type: 'demo-request',
			},
		})

		expect(result.success).toBe(false)
	})

	it('refuses a time that is not an ISO instant', () => {
		const result = eventEnvelopeSchema.safeParse({
			...envelope,
			occurredAt: 'yesterday',
		})

		expect(result.success).toBe(false)
	})

	it('refuses a field it does not know', () => {
		const result = eventEnvelopeSchema.safeParse({
			...envelope,
			topic: 'demo-requests',
		})

		expect(result.success).toBe(false)
	})
})
