// The subset of a Redis connection the transport uses, with the replies the
// stream commands give over RESP2. ioredis satisfies it; tests use a fake.
export type StreamItem = [id: string, fields: string[] | null]
export type StreamReadReply = [key: string, items: StreamItem[]][] | null

export interface StreamConnection {
	xgroup(
		subcommand: 'CREATE',
		key: string,
		group: string,
		id: '$',
		mkstream: 'MKSTREAM',
	): Promise<unknown>
	xreadgroup(
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
	): Promise<StreamReadReply>
	xautoclaim(
		key: string,
		group: string,
		consumer: string,
		minIdleMs: number,
		start: string,
		countToken: 'COUNT',
		count: number,
	): Promise<unknown[]>
	xpending(
		key: string,
		group: string,
		start: string,
		end: string,
		count: number,
	): Promise<unknown[]>
	xack(key: string, group: string, id: string): Promise<number>
	xadd(key: string, ...args: (string | number)[]): Promise<string | null>
	xrange(
		key: string,
		start: string,
		end: string,
		countToken: 'COUNT',
		count: number,
	): Promise<StreamItem[]>
	xread(
		blockToken: 'BLOCK',
		milliseconds: number,
		streamsToken: 'STREAMS',
		key: string,
		id: string,
	): Promise<StreamReadReply>
	disconnect(): void
	quit(): Promise<unknown>
}
