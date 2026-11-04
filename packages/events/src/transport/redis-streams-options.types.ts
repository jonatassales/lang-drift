import type { StreamConnection } from './stream-connection.types.js'

// What a consumer needs: its stream, its group, its own name, where to park
// an entry it gives up on, and how to open a connection. Two connections are
// opened: one blocks on reads, the other runs every other command.
export interface RedisStreamsServerOptions {
	stream: string
	group: string
	consumer: string
	deadLetterStream: string
	connect: () => StreamConnection
	blockMs?: number
	batchSize?: number
	claimMinIdleMs?: number
	claimEveryMs?: number
	maxDeliveries?: number
}

// What a producer needs: its stream, how to open a connection, and the
// approximate length the stream is trimmed to on every append.
export interface RedisStreamsClientOptions {
	stream: string
	connect: () => StreamConnection
	maxLength?: number
}
