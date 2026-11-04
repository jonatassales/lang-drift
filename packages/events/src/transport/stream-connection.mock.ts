import type {
	StreamConnection,
	StreamItem,
	StreamReadReply,
} from './stream-connection.types.js'

export interface StreamConnectionScript {
	reads?: StreamReadReply[]
	claims?: StreamItem[][]
	deliveries?: Record<string, number>
	groupError?: Error
	addedIds?: string[]
}

// Records every command and answers from a script. A read past the script
// blocks, as Redis would, until disconnect() releases it with no entries.
export class FakeStreamConnection implements StreamConnection {
	readonly calls: (string | number)[][] = []
	private readonly reads: StreamReadReply[]
	private readonly claims: StreamItem[][]
	private readonly deliveries: Map<string, number>
	private readonly groupError: Error | undefined
	private readonly addedIds: string[]
	private release: (() => void) | undefined
	private idleResolvers: (() => void)[] = []

	constructor(script: StreamConnectionScript = {}) {
		this.reads = [...(script.reads ?? [])]
		this.claims = [...(script.claims ?? [])]
		this.deliveries = new Map(Object.entries(script.deliveries ?? {}))
		this.groupError = script.groupError
		this.addedIds = [...(script.addedIds ?? [])]
	}

	// Resolves once a read finds the script exhausted: everything before it
	// has been handled and the consumer is blocking.
	idle(): Promise<void> {
		return new Promise((resolve) => {
			this.idleResolvers.push(resolve)
		})
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
		const scripted = this.reads.shift()
		if (scripted !== undefined) return scripted

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

	disconnect(): void {
		this.calls.push(['disconnect'])
		this.release?.()
		this.release = undefined
	}

	async quit(): Promise<unknown> {
		this.calls.push(['quit'])

		return 'OK'
	}
}
