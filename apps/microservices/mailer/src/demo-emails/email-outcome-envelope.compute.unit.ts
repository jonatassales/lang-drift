import { describe, expect, it } from 'vitest'

import type { DemoRequestReceived } from 'infra/streaming'
import { demoRequestEventSchema } from 'infra/streaming'

import { emailOutcomeEnvelope } from './email-outcome-envelope.compute.js'

const received: DemoRequestReceived = {
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
		email: 'ana@acme.com',
		company: 'Acme',
		position: 'cto',
		locale: 'pt-BR',
		submittedAt: '2026-11-03T14:05:00.000Z',
	},
}

const now = new Date('2026-11-03T14:05:02.000Z')

describe('mailer/demo-emails/email-outcome-envelope.compute', () => {
	it('records a sent email as an event caused by the received one', () => {
		const envelope = emailOutcomeEnvelope(received, {
			kind: 'confirmation',
			outcome: 'sent',
			messageId: '49a3999c-0ce1-4ea6-ab68-afcd6dc2e794',
			to: 'ana@acme.com',
			now,
		})

		expect(demoRequestEventSchema.safeParse(envelope).success).toBe(true)
		expect(envelope).toMatchObject({
			type: 'demo-request.email-sent',
			version: 1,
			aggregate: received.aggregate,
			occurredAt: '2026-11-03T14:05:02.000Z',
			correlationId: received.aggregate.id,
			causationId: received.id,
			payload: {
				kind: 'confirmation',
				provider: 'resend',
				messageId: '49a3999c-0ce1-4ea6-ab68-afcd6dc2e794',
				to: 'ana@acme.com',
				sentAt: '2026-11-03T14:05:02.000Z',
			},
		})
	})

	it('records a failed email with the reason', () => {
		const envelope = emailOutcomeEnvelope(received, {
			kind: 'lead-notification',
			outcome: 'failed',
			reason: 'Resend answered 422: invalid recipient',
			now,
		})

		expect(demoRequestEventSchema.safeParse(envelope).success).toBe(true)
		expect(envelope).toMatchObject({
			type: 'demo-request.email-failed',
			causationId: received.id,
			payload: {
				kind: 'lead-notification',
				provider: 'resend',
				reason: 'Resend answered 422: invalid recipient',
				failedAt: '2026-11-03T14:05:02.000Z',
			},
		})
	})
})
