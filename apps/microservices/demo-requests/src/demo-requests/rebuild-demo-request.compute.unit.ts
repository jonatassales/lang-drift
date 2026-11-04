import { describe, expect, it } from 'vitest'

import type { DemoRequestEvent } from 'streaming'

import { rebuildDemoRequest } from './rebuild-demo-request.compute.js'

const requestId = '701706be-7e70-4681-9608-826d3ad548b9'

const received: DemoRequestEvent = {
	id: '0d8ce74c-4756-4bfd-8d35-c820e302898e',
	type: 'demo-request.received',
	version: 1,
	aggregate: {
		type: 'demo-request',
		id: requestId,
	},
	occurredAt: '2026-11-03T14:05:00.000Z',
	correlationId: requestId,
	payload: {
		email: 'ana@acme.com',
		company: 'Acme',
		position: 'cto',
		locale: 'pt-BR',
		submittedAt: '2026-11-03T14:05:00.000Z',
	},
}

const confirmationSent: DemoRequestEvent = {
	id: '3f80418a-aab6-4f07-8b70-1d7b68a03aa8',
	type: 'demo-request.email-sent',
	version: 1,
	aggregate: received.aggregate,
	occurredAt: '2026-11-03T14:05:02.000Z',
	correlationId: requestId,
	causationId: received.id,
	payload: {
		kind: 'confirmation',
		provider: 'resend',
		messageId: '49a3999c-0ce1-4ea6-ab68-afcd6dc2e794',
		to: 'ana@acme.com',
		sentAt: '2026-11-03T14:05:02.000Z',
	},
}

const confirmationFailed: DemoRequestEvent = {
	id: 'cdaa3cd1-2406-4678-9231-089e4533d4fd',
	type: 'demo-request.email-failed',
	version: 1,
	aggregate: received.aggregate,
	occurredAt: '2026-11-03T14:05:02.000Z',
	correlationId: requestId,
	causationId: received.id,
	payload: {
		kind: 'confirmation',
		provider: 'resend',
		reason: 'Resend answered 422: invalid recipient',
		failedAt: '2026-11-03T14:05:02.000Z',
	},
}

const leadNotified: DemoRequestEvent = {
	id: '7067c771-2d30-4bdf-8d4f-dcc262a7f7e5',
	type: 'demo-request.email-sent',
	version: 1,
	aggregate: received.aggregate,
	occurredAt: '2026-11-03T14:05:03.000Z',
	correlationId: requestId,
	causationId: received.id,
	payload: {
		kind: 'lead-notification',
		provider: 'resend',
		messageId: 'a808b80e-d279-4209-a1c4-555cdaa3cd12',
		to: 'team@langdrift.md',
		sentAt: '2026-11-03T14:05:03.000Z',
	},
}

describe('demo-requests/demo-requests/rebuild-demo-request.compute', () => {
	it('rebuilds nothing from no events', () => {
		expect(rebuildDemoRequest([])).toBeNull()
	})

	it('reads a request that was only received', () => {
		expect(rebuildDemoRequest([received])).toEqual({
			id: requestId,
			state: 'received',
			submission: received.payload,
			emails: [],
		})
	})

	it('reads a request as confirmed once its confirmation was sent', () => {
		const view = rebuildDemoRequest([received, confirmationSent, leadNotified])

		expect(view?.state).toBe('confirmed')
		expect(view?.emails).toEqual([
			{
				kind: 'confirmation',
				outcome: 'sent',
				provider: 'resend',
				messageId: '49a3999c-0ce1-4ea6-ab68-afcd6dc2e794',
				at: '2026-11-03T14:05:02.000Z',
			},
			{
				kind: 'lead-notification',
				outcome: 'sent',
				provider: 'resend',
				messageId: 'a808b80e-d279-4209-a1c4-555cdaa3cd12',
				at: '2026-11-03T14:05:03.000Z',
			},
		])
	})

	it('reads a request as failed when its confirmation failed', () => {
		const view = rebuildDemoRequest([received, confirmationFailed])

		expect(view?.state).toBe('failed')
		expect(view?.emails).toEqual([
			{
				kind: 'confirmation',
				outcome: 'failed',
				provider: 'resend',
				reason: 'Resend answered 422: invalid recipient',
				at: '2026-11-03T14:05:02.000Z',
			},
		])
	})

	it('reads a request as confirmed when a failed confirmation was sent later', () => {
		const view = rebuildDemoRequest([
			received,
			confirmationFailed,
			confirmationSent,
		])

		expect(view?.state).toBe('confirmed')
		expect(view?.emails).toHaveLength(2)
	})

	it('keeps the lead notification out of the state', () => {
		const view = rebuildDemoRequest([received, leadNotified])

		expect(view?.state).toBe('received')
		expect(view?.emails).toHaveLength(1)
	})
})
