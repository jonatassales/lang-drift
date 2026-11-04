import { describe, expect, it } from 'vitest'

import { FakeReaderConnection } from '@app/event-log/testdata/stream-connection.mock.js'

import type { FollowEvent } from './event-log-reader.service.js'
import { EventLogReaderService } from './event-log-reader.service.js'

const stream = 'demo-requests:events'
const requestA = '701706be-7e70-4681-9608-826d3ad548b9'
const requestB = 'a46d4497-4c0a-4422-9e17-187cf0983297'

function fields(aggregateId: string, type: string): string[] {
	return [
		'type',
		type,
		'aggregateId',
		aggregateId,
		'envelope',
		JSON.stringify({
			type,
			aggregate: {
				type: 'demo-request',
				id: aggregateId,
			},
		}),
	]
}

describe('demo-requests/event-log/event-log-reader.service', () => {
	it('reads one request by ranging the stream and keeping only its entries', async () => {
		const connection = new FakeReaderConnection({
			ranges: [
				[
					['1762178700000-0', fields(requestA, 'demo-request.received')],
					['1762178700001-0', fields(requestB, 'demo-request.received')],
					['1762178700002-0', fields(requestA, 'demo-request.email-sent')],
				],
			],
		})
		const reader = new EventLogReaderService({
			consumer: () => connection,
			producer: () => connection,
		})

		const entries = await reader.readAggregate(requestA)

		expect(entries.map((entry) => entry.id)).toEqual([
			'1762178700000-0',
			'1762178700002-0',
		])
		expect(entries[0].aggregateId).toBe(requestA)
		expect(connection.calls[0]).toEqual([
			'xrange',
			stream,
			'-',
			'+',
			'COUNT',
			500,
		])
	})

	it('keeps ranging from the last entry while a page comes back full', async () => {
		const fullPage: [string, string[]][] = []
		for (let index = 0; index < 500; index += 1) {
			fullPage.push([
				`1762178700000-${index}`,
				fields(requestB, 'demo-request.received'),
			])
		}
		const connection = new FakeReaderConnection({
			ranges: [
				fullPage,
				[['1762178700001-0', fields(requestA, 'demo-request.received')]],
			],
		})
		const reader = new EventLogReaderService({
			consumer: () => connection,
			producer: () => connection,
		})

		const entries = await reader.readAggregate(requestA)

		expect(entries).toHaveLength(1)
		expect(connection.calls[1]).toEqual([
			'xrange',
			stream,
			'(1762178700000-499',
			'+',
			'COUNT',
			500,
		])
	})

	it('skips an entry it cannot read', async () => {
		const connection = new FakeReaderConnection({
			ranges: [
				[
					['1762178700000-0', ['type', 'demo-request.received']],
					['1762178700001-0', fields(requestA, 'demo-request.received')],
				],
			],
		})
		const reader = new EventLogReaderService({
			consumer: () => connection,
			producer: () => connection,
		})

		const entries = await reader.readAggregate(requestA)

		expect(entries.map((entry) => entry.id)).toEqual(['1762178700001-0'])
	})

	it('follows the stream from the resolved live end on its own connection, pings on an empty block and tears down on unsubscribe', async () => {
		const connections: FakeReaderConnection[] = []
		const reader = new EventLogReaderService({
			producer: () => new FakeReaderConnection(),
			consumer: () => {
				const connection = new FakeReaderConnection({
					latest: [
						['1762178699999-0', fields(requestB, 'demo-request.received')],
					],
					follows: [
						[
							[
								stream,
								[
									[
										'1762178700000-0',
										fields(requestA, 'demo-request.received'),
									],
								],
							],
						],
						null,
					],
				})
				connections.push(connection)

				return connection
			},
		})
		const events: FollowEvent[] = []

		const subscription = reader.follow().subscribe((event) => {
			events.push(event)
		})

		await connections[0].idle()

		subscription.unsubscribe()

		expect(events).toEqual([
			{
				kind: 'entry',
				entry: {
					id: '1762178700000-0',
					type: 'demo-request.received',
					aggregateId: requestA,
					envelope: {
						type: 'demo-request.received',
						aggregate: {
							type: 'demo-request',
							id: requestA,
						},
					},
				},
			},
			{
				kind: 'ping',
			},
		])
		expect(connections[0].calls[0]).toEqual([
			'xrevrange',
			stream,
			'+',
			'-',
			'COUNT',
			1,
		])
		expect(connections[0].calls[1]).toEqual([
			'xread',
			'BLOCK',
			15000,
			'STREAMS',
			stream,
			'1762178699999-0',
		])
		expect(connections[0].calls[2]).toEqual([
			'xread',
			'BLOCK',
			15000,
			'STREAMS',
			stream,
			'1762178700000-0',
		])
		expect(connections[0].calls.at(-1)).toEqual(['disconnect'])
	})

	it('follows from a given entry when the client names one', async () => {
		const connection = new FakeReaderConnection()
		const reader = new EventLogReaderService({
			consumer: () => connection,
			producer: () => connection,
		})

		const subscription = reader.follow('1762178700000-0').subscribe()

		await connection.idle()

		subscription.unsubscribe()

		expect(connection.calls[0]).toEqual([
			'xread',
			'BLOCK',
			15000,
			'STREAMS',
			stream,
			'1762178700000-0',
		])
	})

	it('quits the shared connection on shutdown', async () => {
		const connection = new FakeReaderConnection({
			ranges: [[]],
		})
		const reader = new EventLogReaderService({
			consumer: () => connection,
			producer: () => connection,
		})

		await reader.readAggregate(requestA)

		await reader.onApplicationShutdown()

		expect(connection.calls.at(-1)).toEqual(['quit'])
	})
})
