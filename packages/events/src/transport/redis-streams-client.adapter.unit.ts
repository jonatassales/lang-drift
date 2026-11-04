import { lastValueFrom } from 'rxjs'
import { describe, expect, it } from 'vitest'

import { RedisStreamsClient } from './redis-streams-client.adapter.js'
import { FakeStreamConnection } from './stream-connection.mock.js'

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
	payload: {},
}

describe('events/transport/redis-streams-client.adapter', () => {
	it('appends an emitted envelope under its type with the flat routing fields and a capped length', async () => {
		const connection = new FakeStreamConnection({
			addedIds: ['1762178700000-7'],
		})
		const client = new RedisStreamsClient({
			stream: 'demo-requests:events',
			maxLength: 100000,
			connect: () => connection,
		})

		const id = await lastValueFrom(
			client.emit('demo-request.received', envelope),
		)

		expect(id).toBe('1762178700000-7')
		expect(connection.calls).toEqual([
			['connect'],
			[
				'xadd',
				'demo-requests:events',
				'MAXLEN',
				'~',
				100000,
				'*',
				'type',
				'demo-request.received',
				'aggregateId',
				'701706be-7e70-4681-9608-826d3ad548b9',
				'envelope',
				JSON.stringify(envelope),
			],
		])
	})

	it('refuses to append an envelope under a type that is not its own', async () => {
		const client = new RedisStreamsClient({
			stream: 'demo-requests:events',
			connect: () => new FakeStreamConnection(),
		})

		await expect(
			lastValueFrom(client.emit('demo-request.email-sent', envelope)),
		).rejects.toThrow(/type/)
	})

	it('does not support request-response', async () => {
		const client = new RedisStreamsClient({
			stream: 'demo-requests:events',
			connect: () => new FakeStreamConnection(),
		})

		await expect(
			lastValueFrom(client.send('demo-request.received', envelope)),
		).rejects.toThrow(/request-response/)
	})

	it('rejects an emit when the connection cannot be opened and opens a new one on the next emit', async () => {
		const connections: FakeStreamConnection[] = []
		let failures = 1
		const client = new RedisStreamsClient({
			stream: 'demo-requests:events',
			connect: () => {
				const connection = new FakeStreamConnection(
					failures > 0
						? {
								connectError: new Error('connect ECONNREFUSED 127.0.0.1:6379'),
							}
						: {},
				)
				failures -= 1
				connections.push(connection)

				return connection
			},
		})

		await expect(
			lastValueFrom(client.emit('demo-request.received', envelope)),
		).rejects.toThrow(/ECONNREFUSED/)

		await lastValueFrom(client.emit('demo-request.received', envelope))

		expect(connections).toHaveLength(2)
		expect(connections[0].calls).toEqual([['connect'], ['disconnect']])
		expect(connections[1].calls[0]).toEqual(['connect'])
	})

	it('quits the connection on close and opens a new one on the next emit', async () => {
		const connections: FakeStreamConnection[] = []
		const client = new RedisStreamsClient({
			stream: 'demo-requests:events',
			connect: () => {
				const connection = new FakeStreamConnection()
				connections.push(connection)

				return connection
			},
		})

		await lastValueFrom(client.emit('demo-request.received', envelope))

		await client.close()

		await lastValueFrom(client.emit('demo-request.received', envelope))

		expect(connections).toHaveLength(2)
		expect(connections[0].calls.at(-1)).toEqual(['quit'])
	})
})
