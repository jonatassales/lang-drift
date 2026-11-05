import { describe, expect, it } from 'vitest'

import { demoRequestEventSchema } from './demo-request-event.validation.js'

const envelope = {
	id: '0d8ce74c-4756-4bfd-8d35-c820e302898e',
	version: 1,
	aggregate: {
		type: 'demo-request',
		id: '701706be-7e70-4681-9608-826d3ad548b9',
	},
	occurredAt: '2026-11-03T14:05:00.000Z',
	correlationId: '701706be-7e70-4681-9608-826d3ad548b9',
}

const received = {
	...envelope,
	type: 'demo-request.received',
	payload: {
		email: 'ana@acme.com',
		company: 'Acme',
		position: 'cto',
		locale: 'pt-BR',
		submittedAt: '2026-11-03T14:05:00.000Z',
	},
}

const emailSent = {
	...envelope,
	id: '3f80418a-aab6-4f07-8b70-1d7b68a03aa8',
	type: 'demo-request.email-sent',
	causationId: '0d8ce74c-4756-4bfd-8d35-c820e302898e',
	payload: {
		kind: 'confirmation',
		provider: 'resend',
		messageId: '49a3999c-0ce1-4ea6-ab68-afcd6dc2e794',
		to: 'ana@acme.com',
		sentAt: '2026-11-03T14:05:02.000Z',
	},
}

const emailFailed = {
	...envelope,
	id: 'cdaa3cd1-2406-4678-9231-089e4533d4fd',
	type: 'demo-request.email-failed',
	causationId: '0d8ce74c-4756-4bfd-8d35-c820e302898e',
	payload: {
		kind: 'lead-notification',
		provider: 'resend',
		reason: 'Resend answered 422: invalid recipient',
		failedAt: '2026-11-03T14:05:02.000Z',
	},
}

describe('infra/streaming/events/demo-request-event.validation', () => {
	it('accepts a received event and narrows it by type', () => {
		const result = demoRequestEventSchema.parse(received)

		expect(result.type).toBe('demo-request.received')
		if (result.type === 'demo-request.received') {
			expect(result.payload.company).toBe('Acme')
		}
	})

	it('accepts an email-sent event with the kind, the provider and its message id', () => {
		const result = demoRequestEventSchema.parse(emailSent)

		expect(result.type).toBe('demo-request.email-sent')
		if (result.type === 'demo-request.email-sent') {
			expect(result.payload.kind).toBe('confirmation')
			expect(result.payload.messageId).toBe(
				'49a3999c-0ce1-4ea6-ab68-afcd6dc2e794',
			)
		}
	})

	it('accepts an email-failed event with the kind and the reason', () => {
		const result = demoRequestEventSchema.parse(emailFailed)

		expect(result.type).toBe('demo-request.email-failed')
		if (result.type === 'demo-request.email-failed') {
			expect(result.payload.reason).toBe(
				'Resend answered 422: invalid recipient',
			)
		}
	})

	it('refuses a type it does not know', () => {
		const result = demoRequestEventSchema.safeParse({
			...received,
			type: 'demo-request.archived',
		})

		expect(result.success).toBe(false)
	})

	it('refuses an email kind it does not know', () => {
		const result = demoRequestEventSchema.safeParse({
			...emailSent,
			payload: {
				...emailSent.payload,
				kind: 'newsletter',
			},
		})

		expect(result.success).toBe(false)
	})

	it('refuses a received event whose payload is a sent payload', () => {
		const result = demoRequestEventSchema.safeParse({
			...received,
			payload: emailSent.payload,
		})

		expect(result.success).toBe(false)
	})

	it('refuses an aggregate that is not a demo request', () => {
		const result = demoRequestEventSchema.safeParse({
			...received,
			aggregate: {
				type: 'early-access',
				id: '701706be-7e70-4681-9608-826d3ad548b9',
			},
		})

		expect(result.success).toBe(false)
	})
})
