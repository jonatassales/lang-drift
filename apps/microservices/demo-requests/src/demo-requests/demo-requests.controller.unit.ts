import { randomUUID } from 'node:crypto'
import { NotFoundException, ServiceUnavailableException } from '@nestjs/common'
import { lastValueFrom, of, throwError, toArray } from 'rxjs'
import { describe, expect, it } from 'vitest'

import type {
	EventLogReaderService,
	FollowEvent,
} from '@app/event-log/event-log-reader.service.js'
import { demoRequestReceivedSchema } from 'infra/streaming'
import type { StreamEntry } from 'infra/streaming/transport'

import { DemoRequestsController } from './demo-requests.controller.js'

interface Emitted {
	pattern: string
	data: unknown
}

const requestId = '701706be-7e70-4681-9608-826d3ad548b9'

function entryOf(id: string, type: string, payload: unknown): StreamEntry {
	return {
		id,
		type,
		aggregateId: requestId,
		envelope: {
			id: randomUUID(),
			type,
			version: 1,
			aggregate: {
				type: 'demo-request',
				id: requestId,
			},
			occurredAt: '2026-11-03T14:05:00.000Z',
			correlationId: requestId,
			payload,
		},
	}
}

class FakeReader {
	readonly entries: StreamEntry[] = []
	readonly followedFrom: (string | undefined)[] = []
	followEvents: FollowEvent[] = []

	async readAggregate(aggregateId: string): Promise<StreamEntry[]> {
		return this.entries.filter((entry) => entry.aggregateId === aggregateId)
	}

	follow(fromId?: string) {
		this.followedFrom.push(fromId)

		return of(...this.followEvents)
	}
}

function readerOf(fake: FakeReader): EventLogReaderService {
	return fake as unknown as EventLogReaderService
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
		const controller = new DemoRequestsController(
			producer,
			readerOf(new FakeReader()),
		)

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
		const controller = new DemoRequestsController(
			producer,
			readerOf(new FakeReader()),
		)

		const first = await controller.create(submission, new FakeResponse())

		const second = await controller.create(submission, new FakeResponse())

		expect(first.id).not.toBe(second.id)
	})

	it('answers 503 when the event log cannot take the event', async () => {
		const producer = new FakeProducer()
		producer.failing = true
		const controller = new DemoRequestsController(
			producer,
			readerOf(new FakeReader()),
		)

		await expect(
			controller.create(submission, new FakeResponse()),
		).rejects.toBeInstanceOf(ServiceUnavailableException)
	})

	it('reads a request by folding its events into its state', async () => {
		const reader = new FakeReader()
		reader.entries.push(
			entryOf('1762178700000-0', 'demo-request.received', {
				...submission,
				submittedAt: '2026-11-03T14:05:00.000Z',
			}),
			entryOf('1762178700002-0', 'demo-request.email-sent', {
				kind: 'confirmation',
				provider: 'resend',
				messageId: '49a3999c-0ce1-4ea6-ab68-afcd6dc2e794',
				to: 'ana@acme.com',
				sentAt: '2026-11-03T14:05:02.000Z',
			}),
		)
		const controller = new DemoRequestsController(
			new FakeProducer(),
			readerOf(reader),
		)

		const view = await controller.read(requestId)

		expect(view.state).toBe('confirmed')
		expect(view.submission.company).toBe('Acme')
		expect(view.emails).toHaveLength(1)
	})

	it('answers 404 for a request with no events', async () => {
		const controller = new DemoRequestsController(
			new FakeProducer(),
			readerOf(new FakeReader()),
		)

		await expect(controller.read(requestId)).rejects.toBeInstanceOf(
			NotFoundException,
		)
	})

	it('skips an entry whose envelope is not a demo-request event when reading', async () => {
		const reader = new FakeReader()
		reader.entries.push(
			entryOf('1762178700000-0', 'demo-request.received', {
				...submission,
				submittedAt: '2026-11-03T14:05:00.000Z',
			}),
			entryOf('1762178700001-0', 'demo-request.archived', {}),
		)
		const controller = new DemoRequestsController(
			new FakeProducer(),
			readerOf(reader),
		)

		const view = await controller.read(requestId)

		expect(view.state).toBe('received')
	})

	it('streams the log as server-sent events from the last event the client saw, else from the live end', async () => {
		const reader = new FakeReader()
		reader.followEvents = [
			{
				kind: 'entry',
				entry: entryOf('1762178700000-0', 'demo-request.received', {}),
			},
			{
				kind: 'ping',
			},
		]
		const controller = new DemoRequestsController(
			new FakeProducer(),
			readerOf(reader),
		)

		const events = await lastValueFrom(
			controller.events('1762178699999-0', undefined).pipe(toArray()),
		)

		await lastValueFrom(controller.events(undefined, undefined).pipe(toArray()))

		await lastValueFrom(controller.events(undefined, '0-0').pipe(toArray()))

		expect(events.map((event) => event.type)).toEqual([
			'demo-request.received',
			'ping',
		])
		expect(events[0].id).toBe('1762178700000-0')
		expect(reader.followedFrom).toEqual(['1762178699999-0', undefined, '0-0'])
	})
})
