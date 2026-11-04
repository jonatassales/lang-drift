import type {
	StreamConnection,
	StreamItem,
	StreamReadReply,
} from './stream-connection.types.js'

export interface StreamConnectionScript {
	reads?: StreamReadReply[]
	ranges?: StreamItem[][]
	follows?: StreamReadReply[]
	claims?: StreamItem[][]
	deliveries?: Record<string, number>
	groupError?: Error
	connectError?: Error
	readError?: Error
	latest?: StreamItem[]
	addedIds?: string[]
}

// Records every command and answers from a script. A read past the script
// blocks, as Redis would, until disconnect() releases it with no entries.
export class FakeStreamConnection implements StreamConnection {
	readonly calls: (string | number)[][] = []
	private readonly reads: StreamReadReply[]
	private readonly ranges: StreamItem[][]
	private readonly follows: StreamReadReply[]
	private readonly claims: StreamItem[][]
	private readonly deliveries: Map<string, number>
	private readonly groupError: Error | undefined
	private readonly connectError: Error | undefined
	private readError: Error | undefined
	private readonly latest: StreamItem[]
	private readonly addedIds: string[]
	private release: (() => void) | undefined
	private idleResolvers: (() => void)[] = []
	private blocking = false

	constructor(script: StreamConnectionScript = {}) {
		this.reads = [...(script.reads ?? [])]
		this.ranges = [...(script.ranges ?? [])]
		this.follows = [...(script.follows ?? [])]
		this.claims = [...(script.claims ?? [])]
		this.deliveries = new Map(Object.entries(script.deliveries ?? {}))
		this.groupError = script.groupError
		this.connectError = script.connectError
		this.readError = script.readError
		this.latest = [...(script.latest ?? [])]
		this.addedIds = [...(script.addedIds ?? [])]
	}

	// Resolves once a read finds the script exhausted: everything before it
	// has been handled and the consumer is blocking.
	idle(): Promise<void> {
		if (this.blocking) return Promise.resolve()

		return new Promise((resolve) => {
			this.idleResolvers.push(resolve)
		})
	}

	async connect(): Promise<void> {
		this.calls.push(['connect'])
		if (this.connectError) throw this.connectError
	}

	async xrevrange(
		key: string,
		end: '+',
		start: '-',
		countToken: 'COUNT',
		count: number,
	): Promise<StreamItem[]> {
		this.calls.push(['xrevrange', key, end, start, countToken, count])

		return this.latest
	}

	async xgroup(
		subcommand: 'CREATE',
		key: string,
		group: string,
		id: '$',
		mkstream: 'MKSTREAM',
	): Promise<unknown> {
		this.calls.push(['xgroup', subcommand, key, group, id, mkstream])
		if (this.groupError) throw this.groupError

		return 'OK'
	}

	async xreadgroup(
		groupToken: 'GROUP',
		group: string,
		consumer: string,
		countToken: 'COUNT',
		count: number,
		blockToken: 'BLOCK',
		milliseconds: number,
		streamsToken: 'STREAMS',
		key: string,
		id: '>',
	): Promise<StreamReadReply> {
		this.calls.push([
			'xreadgroup',
			groupToken,
			group,
			consumer,
			countToken,
			count,
			blockToken,
			milliseconds,
			streamsToken,
			key,
			id,
		])
		if (this.readError) {
			const error = this.readError
			this.readError = undefined

			throw error
		}

		const scripted = this.reads.shift()
		if (scripted !== undefined) return scripted

		this.blocking = true
		for (const resolve of this.idleResolvers) resolve()
		this.idleResolvers = []

		return new Promise((resolve) => {
			this.release = () => resolve(null)
		})
	}

	async xautoclaim(
		key: string,
		group: string,
		consumer: string,
		minIdleMs: number,
		start: string,
		countToken: 'COUNT',
		count: number,
	): Promise<unknown[]> {
		this.calls.push([
			'xautoclaim',
			key,
			group,
			consumer,
			minIdleMs,
			start,
			countToken,
			count,
		])
		const items = this.claims.shift() ?? []

		return ['0-0', items, []]
	}

	async xpending(
		key: string,
		group: string,
		start: string,
		end: string,
		count: number,
	): Promise<unknown[]> {
		this.calls.push(['xpending', key, group, start, end, count])
		const deliveries = this.deliveries.get(start)
		if (deliveries === undefined) return []

		return [[start, 'consumer-1', 61000, deliveries]]
	}

	async xack(key: string, group: string, id: string): Promise<number> {
		this.calls.push(['xack', key, group, id])

		return 1
	}

	async xadd(
		key: string,
		...args: (string | number)[]
	): Promise<string | null> {
		this.calls.push(['xadd', key, ...args])

		return this.addedIds.shift() ?? '1762178700000-0'
	}

	async xrange(
		key: string,
		start: string,
		end: string,
		countToken: 'COUNT',
		count: number,
	): Promise<StreamItem[]> {
		this.calls.push(['xrange', key, start, end, countToken, count])

		return this.ranges.shift() ?? []
	}

	async xread(
		blockToken: 'BLOCK',
		milliseconds: number,
		streamsToken: 'STREAMS',
		key: string,
		id: string,
	): Promise<StreamReadReply> {
		this.calls.push(['xread', blockToken, milliseconds, streamsToken, key, id])
		const scripted = this.follows.shift()
		if (scripted !== undefined) return scripted

		this.blocking = true
		for (const resolve of this.idleResolvers) resolve()
		this.idleResolvers = []

		return new Promise((resolve) => {
			this.release = () => resolve(null)
		})
	}

	disconnect(): void {
		this.calls.push(['disconnect'])
		this.blocking = false
		this.release?.()
		this.release = undefined
	}

	async quit(): Promise<unknown> {
		this.calls.push(['quit'])

		return 'OK'
	}
}
