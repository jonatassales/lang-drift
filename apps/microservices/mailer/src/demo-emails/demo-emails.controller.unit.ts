import { demoRequestEventSchema } from '@langdrift/events'
import { of } from 'rxjs'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import { FakeEmailSender } from '@app/demo-emails/testdata/email-sender.mock.js'

import { DemoEmailsController } from './demo-emails.controller.js'

interface Emitted {
	pattern: string
	data: unknown
}

class FakeProducer {
	readonly emitted: Emitted[] = []

	emit(pattern: string, data: unknown) {
		this.emitted.push({
			pattern,
			data,
		})

		return of('1762178700000-9')
	}
}

const received = {
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

describe('mailer/demo-emails/demo-emails.controller', () => {
	beforeEach(() => {
		vi.stubEnv('DEMO_LEAD_EMAIL', 'team@langdrift.md')
	})

	afterEach(() => {
		vi.unstubAllEnvs()
	})

	it('sends the confirmation to the visitor and the notification to the team, each idempotent per event, and records both as sent', async () => {
		const sender = new FakeEmailSender()
		const producer = new FakeProducer()
		const controller = new DemoEmailsController(sender, producer)

		await controller.onReceived(received)

		expect(sender.sends.map((send) => send.email.to)).toEqual([
			'ana@acme.com',
			'team@langdrift.md',
		])
		expect(sender.sends.map((send) => send.idempotencyKey)).toEqual([
			`${received.id}/confirmation`,
			`${received.id}/lead-notification`,
		])
		expect(sender.sends[0].email.from).toBe('LangDrift <no-reply@langdrift.md>')
		expect(sender.sends[0].email.subject).toBe(
			'Recebemos sua solicitação de demo — LangDrift',
		)
		expect(sender.sends[1].email.subject).toBe(
			'[LangDrift] Nova solicitação de demo: Acme',
		)

		const events = producer.emitted.map((emitted) =>
			demoRequestEventSchema.parse(emitted.data),
		)
		expect(producer.emitted.map((emitted) => emitted.pattern)).toEqual([
			'demo-request.email-sent',
			'demo-request.email-sent',
		])
		expect(
			events.map((event) =>
				'kind' in event.payload ? event.payload.kind : null,
			),
		).toEqual(['confirmation', 'lead-notification'])
		expect(events.every((event) => event.causationId === received.id)).toBe(
			true,
		)
	})

	it('records a failed send with the reason and still sends the other email', async () => {
		const sender = new FakeEmailSender()
		sender.failing.add('ana@acme.com')
		const producer = new FakeProducer()
		const controller = new DemoEmailsController(sender, producer)

		await controller.onReceived(received)

		expect(sender.sends).toHaveLength(2)
		expect(producer.emitted.map((emitted) => emitted.pattern)).toEqual([
			'demo-request.email-failed',
			'demo-request.email-sent',
		])
		const failed = demoRequestEventSchema.parse(producer.emitted[0].data)
		expect(failed.type).toBe('demo-request.email-failed')
		if (failed.type === 'demo-request.email-failed') {
			expect(failed.payload.reason).toMatch(/422/)
		}
	})

	it('reuses the same idempotency keys when the event is delivered again', async () => {
		const sender = new FakeEmailSender()
		const controller = new DemoEmailsController(sender, new FakeProducer())

		await controller.onReceived(received)

		await controller.onReceived(received)

		expect(new Set(sender.sends.map((send) => send.idempotencyKey)).size).toBe(
			2,
		)
	})

	it('refuses an envelope that is not a received event, so the entry stays pending', async () => {
		const sender = new FakeEmailSender()
		const producer = new FakeProducer()
		const controller = new DemoEmailsController(sender, producer)

		await expect(
			controller.onReceived({
				...received,
				type: 'demo-request.email-sent',
			}),
		).rejects.toThrow()
		expect(sender.sends).toEqual([])
		expect(producer.emitted).toEqual([])
	})

	it('fails closed when DEMO_LEAD_EMAIL is not set', async () => {
		vi.stubEnv('DEMO_LEAD_EMAIL', '')
		const sender = new FakeEmailSender()
		const producer = new FakeProducer()
		const controller = new DemoEmailsController(sender, producer)

		await controller.onReceived(received)

		expect(sender.sends).toHaveLength(1)
		expect(producer.emitted.map((emitted) => emitted.pattern)).toEqual([
			'demo-request.email-sent',
			'demo-request.email-failed',
		])
		const failed = demoRequestEventSchema.parse(producer.emitted[1].data)
		if (failed.type === 'demo-request.email-failed') {
			expect(failed.payload.reason).toMatch(/DEMO_LEAD_EMAIL/)
		}
	})
})
