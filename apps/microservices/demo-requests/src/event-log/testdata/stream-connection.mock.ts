import type {
	StreamConnection,
	StreamItem,
	StreamReadReply,
} from 'infra/streaming/transport'

export interface ReaderConnectionScript {
	ranges?: StreamItem[][]
	follows?: StreamReadReply[]
	latest?: StreamItem[]
}

// Records the read commands the reader runs and answers from a script. A
// follow past the script blocks, as Redis would, until disconnect() ends it.
export class FakeReaderConnection implements StreamConnection {
	readonly calls: (string | number)[][] = []
	private readonly ranges: StreamItem[][]
	private readonly follows: StreamReadReply[]
	private readonly latest: StreamItem[]
	private release: (() => void) | undefined
	private idleResolvers: (() => void)[] = []
	private blocking = false

	constructor(script: ReaderConnectionScript = {}) {
		this.ranges = [...(script.ranges ?? [])]
		this.follows = [...(script.follows ?? [])]
		this.latest = [...(script.latest ?? [])]
	}

	// Resolves once a follow finds the script exhausted and blocks.
	idle(): Promise<void> {
		if (this.blocking) return Promise.resolve()

		return new Promise((resolve) => {
			this.idleResolvers.push(resolve)
		})
	}

	async connect(): Promise<void> {
		this.calls.push(['connect'])
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

	async xgroup(): Promise<unknown> {
		throw new Error('the reader never creates a group')
	}

	async xreadgroup(): Promise<StreamReadReply> {
		throw new Error('the reader never reads as a group')
	}

	async xautoclaim(): Promise<unknown[]> {
		throw new Error('the reader never claims')
	}

	async xpending(): Promise<unknown[]> {
		throw new Error('the reader never inspects pending entries')
	}

	async xack(): Promise<number> {
		throw new Error('the reader never acknowledges')
	}

	async xadd(): Promise<string | null> {
		throw new Error('the reader never appends')
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
