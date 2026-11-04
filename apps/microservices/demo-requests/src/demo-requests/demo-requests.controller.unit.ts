import { demoRequestReceivedSchema } from '@langdrift/events'
import { ServiceUnavailableException } from '@nestjs/common'
import { of, throwError } from 'rxjs'
import { describe, expect, it } from 'vitest'

import { DemoRequestsController } from './demo-requests.controller.js'

interface Emitted {
	pattern: string
	data: unknown
}

class FakeProducer {
	readonly emitted: Emitted[] = []
	failing = false

	emit(pattern: string, data: unknown) {
		this.emitted.push({
			pattern,
			data,
		})
		if (this.failing) return throwError(() => new Error('Connection is closed'))

		return of('1762178700000-0')
	}
}

class FakeResponse {
	readonly headers: Record<string, string> = {}

	setHeader(name: string, value: string) {
		this.headers[name] = value
	}
}

const submission = {
	email: 'ana@acme.com',
	company: 'Acme',
	position: 'cto',
	locale: 'pt-BR',
} as const

describe('demo-requests/demo-requests/demo-requests.controller', () => {
	it('appends a received event for the submission and answers the request id with its location', async () => {
		const producer = new FakeProducer()
		const response = new FakeResponse()
		const controller = new DemoRequestsController(producer)

		const reply = await controller.create(submission, response)

		expect(producer.emitted).toHaveLength(1)
		const event = demoRequestReceivedSchema.parse(producer.emitted[0].data)
		expect(producer.emitted[0].pattern).toBe('demo-request.received')
		expect(event.aggregate.id).toBe(reply.id)
		expect(event.payload).toMatchObject(submission)
		expect(event.payload.submittedAt).toBe(event.occurredAt)
		expect(response.headers.Location).toBe(`/v1/demo-requests/${reply.id}`)
	})

	it('gives two submissions two different request ids', async () => {
		const producer = new FakeProducer()
		const controller = new DemoRequestsController(producer)

		const first = await controller.create(submission, new FakeResponse())

		const second = await controller.create(submission, new FakeResponse())

		expect(first.id).not.toBe(second.id)
	})

	it('answers 503 when the event log cannot take the event', async () => {
		const producer = new FakeProducer()
		producer.failing = true
		const controller = new DemoRequestsController(producer)

		await expect(
			controller.create(submission, new FakeResponse()),
		).rejects.toBeInstanceOf(ServiceUnavailableException)
	})
})
