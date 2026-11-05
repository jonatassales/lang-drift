import { describe, expect, it } from 'vitest'

import type { RedisStreamsContext } from './redis-streams.context.js'
import { RedisStreamsServer } from './redis-streams.strategy.js'
import type { StreamConnectionScript } from './stream-connection.mock.js'
import { FakeStreamConnection } from './stream-connection.mock.js'

const stream = 'demo-requests:events'
const group = 'mailer'
const consumer = 'mailer-1'
const deadLetterStream = 'demo-requests:events:dead-letter'

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

const fields = [
	'type',
	'demo-request.received',
	'aggregateId',
	'701706be-7e70-4681-9608-826d3ad548b9',
	'envelope',
	JSON.stringify(envelope),
]

interface Harness {
	server: RedisStreamsServer
	blocking: FakeStreamConnection
	commands: FakeStreamConnection
}

function harness(
	blockingScript: StreamConnectionScript = {},
	commandsScript: StreamConnectionScript = {},
): Harness {
	const blocking = new FakeStreamConnection(blockingScript)
	const commands = new FakeStreamConnection(commandsScript)
	const connections = [blocking, commands]
	const server = new RedisStreamsServer({
		stream,
		group,
		consumer,
		deadLetterStream,
		blockMs: 5000,
		batchSize: 16,
		claimMinIdleMs: 60000,
		claimEveryMs: 30000,
		maxDeliveries: 5,
		recoveryDelayMs: 0,
		connect: () => {
			const connection = connections.shift()
			if (!connection) throw new Error('only two connections are opened')

			return connection
		},
	})

	return {
		server,
		blocking,
		commands,
	}
}

describe('infra/streaming/transport/redis-streams.strategy', () => {
	it('creates the consumer group with MKSTREAM and tolerates one that exists', async () => {
		const { server, commands, blocking } = harness(
			{},
			{
				groupError: new Error('BUSYGROUP Consumer Group name already exists'),
			},
		)
		let ready = false

		await server.listen(() => {
			ready = true
		})

		await blocking.idle()

		await server.close()

		expect(ready).toBe(true)
		expect(commands.calls[0]).toEqual([
			'xgroup',
			'CREATE',
			stream,
			group,
			'$',
			'MKSTREAM',
		])
	})

	it('hands an entry to the handler registered for its type and acknowledges it once the handler completes', async () => {
		const { server, commands, blocking } = harness({
			reads: [[[stream, [['1762178700000-0', fields]]]]],
		})
		const handled: [unknown, RedisStreamsContext][] = []
		server.addHandler(
			'demo-request.received',
			async (data: unknown, context: RedisStreamsContext) => {
				handled.push([data, context])
			},
			true,
		)

		await server.listen(() => undefined)

		await blocking.idle()

		await server.close()

		expect(handled).toHaveLength(1)
		expect(handled[0][0]).toEqual(envelope)
		expect(handled[0][1].getEntryId()).toBe('1762178700000-0')
		expect(handled[0][1].getStream()).toBe(stream)
		expect(handled[0][1].getGroup()).toBe(group)
		expect(handled[0][1].getDeliveryCount()).toBe(1)
		expect(commands.calls).toContainEqual([
			'xack',
			stream,
			group,
			'1762178700000-0',
		])
	})

	it('reads with the group, the consumer, the batch size and the block time, from new entries only', async () => {
		const { server, blocking } = harness()

		await server.listen(() => undefined)

		await blocking.idle()

		await server.close()

		expect(blocking.calls).toContainEqual([
			'xreadgroup',
			'GROUP',
			group,
			consumer,
			'COUNT',
			16,
			'BLOCK',
			5000,
			'STREAMS',
			stream,
			'>',
		])
	})

	it('breaks the blocking read on close and quits the command connection', async () => {
		const { server, blocking, commands } = harness()

		await server.listen(() => undefined)

		await blocking.idle()

		await server.close()

		expect(blocking.calls.at(-1)).toEqual(['disconnect'])
		expect(commands.calls.at(-1)).toEqual(['quit'])
	})

	it('logs a failed read, ensures the group again and keeps consuming', async () => {
		const { server, blocking, commands } = harness({
			readError: new Error(
				"NOGROUP No such key 'demo-requests:events' or consumer group 'mailer'",
			),
			reads: [[[stream, [['1762178700000-8', fields]]]]],
		})
		const handled: unknown[] = []
		server.addHandler(
			'demo-request.received',
			async (data: unknown) => {
				handled.push(data)
			},
			true,
		)

		await server.listen(() => undefined)

		await blocking.idle()

		await server.close()

		expect(handled).toEqual([envelope])
		expect(commands.calls.filter((call) => call[0] === 'xgroup')).toHaveLength(
			2,
		)
	})

	it('acknowledges and skips an entry no handler of the group wants', async () => {
		const { server, blocking, commands } = harness({
			reads: [[[stream, [['1762178700000-1', fields]]]]],
		})

		await server.listen(() => undefined)

		await blocking.idle()

		await server.close()

		expect(commands.calls).toContainEqual([
			'xack',
			stream,
			group,
			'1762178700000-1',
		])
		expect(commands.calls.some((call) => call[0] === 'xadd')).toBe(false)
	})

	it('hands an entry of any type to the wildcard handler', async () => {
		const { server, blocking, commands } = harness({
			reads: [[[stream, [['1762178700000-2', fields]]]]],
		})
		const handled: unknown[] = []
		server.addHandler(
			'*',
			async (data: unknown) => {
				handled.push(data)
			},
			true,
		)

		await server.listen(() => undefined)

		await blocking.idle()

		await server.close()

		expect(handled).toEqual([envelope])
		expect(commands.calls).toContainEqual([
			'xack',
			stream,
			group,
			'1762178700000-2',
		])
	})

	it('leaves an entry pending when its handler throws', async () => {
		const { server, blocking, commands } = harness({
			reads: [[[stream, [['1762178700000-3', fields]]]]],
		})
		server.addHandler(
			'demo-request.received',
			async () => {
				throw new Error('Resend is down')
			},
			true,
		)

		await server.listen(() => undefined)

		await blocking.idle()

		await server.close()

		expect(commands.calls.some((call) => call[0] === 'xack')).toBe(false)
	})

	it('hands a claimed entry under the maximum deliveries to its handler with its delivery count', async () => {
		const { server, blocking, commands } = harness(
			{},
			{
				claims: [[['1762178700000-4', fields]]],
				deliveries: {
					'1762178700000-4': 2,
				},
			},
		)
		const counts: number[] = []
		server.addHandler(
			'demo-request.received',
			async (_data: unknown, context: RedisStreamsContext) => {
				counts.push(context.getDeliveryCount())
			},
			true,
		)

		await server.listen(() => undefined)

		await blocking.idle()

		await server.close()

		expect(counts).toEqual([2])
		expect(commands.calls).toContainEqual([
			'xautoclaim',
			stream,
			group,
			consumer,
			60000,
			'0-0',
			'COUNT',
			16,
		])
		expect(commands.calls).toContainEqual([
			'xack',
			stream,
			group,
			'1762178700000-4',
		])
	})

	it('parks a claimed entry past the maximum deliveries on the dead-letter stream and acknowledges it', async () => {
		const { server, blocking, commands } = harness(
			{},
			{
				claims: [[['1762178700000-5', fields]]],
				deliveries: {
					'1762178700000-5': 6,
				},
			},
		)
		const handled: unknown[] = []
		server.addHandler(
			'demo-request.received',
			async (data: unknown) => {
				handled.push(data)
			},
			true,
		)

		await server.listen(() => undefined)

		await blocking.idle()

		await server.close()

		expect(handled).toEqual([])
		const parked = commands.calls.find((call) => call[0] === 'xadd')
		expect(parked?.slice(0, 2)).toEqual(['xadd', deadLetterStream])
		expect(parked).toContain('delivered 6 times')
		expect(parked).toContain('1762178700000-5')
		expect(commands.calls).toContainEqual([
			'xack',
			stream,
			group,
			'1762178700000-5',
		])
	})

	it('parks an entry whose envelope is not JSON with the reason and acknowledges it', async () => {
		const { server, blocking, commands } = harness({
			reads: [
				[
					[
						stream,
						[
							[
								'1762178700000-6',
								[
									'type',
									'demo-request.received',
									'aggregateId',
									'701706be-7e70-4681-9608-826d3ad548b9',
									'envelope',
									'{not json',
								],
							],
						],
					],
				],
			],
		})
		const handled: unknown[] = []
		server.addHandler(
			'demo-request.received',
			async (data: unknown) => {
				handled.push(data)
			},
			true,
		)

		await server.listen(() => undefined)

		await blocking.idle()

		await server.close()

		expect(handled).toEqual([])
		const parked = commands.calls.find((call) => call[0] === 'xadd')
		expect(parked?.slice(0, 2)).toEqual(['xadd', deadLetterStream])
		expect(
			String(parked?.find((value) => String(value).includes('JSON'))),
		).toMatch(/JSON/)
		expect(commands.calls).toContainEqual([
			'xack',
			stream,
			group,
			'1762178700000-6',
		])
	})
})
