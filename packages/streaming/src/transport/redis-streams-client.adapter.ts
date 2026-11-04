import type { ReadPacket, WritePacket } from '@nestjs/microservices'
import { ClientProxy } from '@nestjs/microservices'

import type { RedisStreamsClientOptions } from './redis-streams-options.types.js'
import type { StreamConnection } from './stream-connection.types.js'
import { streamFields } from './stream-fields.map.js'

const defaultMaxLength = 100000

// Appends events to a stream through Nest's client contract: `emit` resolves
// with the entry id. Request-response has no place on an event log.
export class RedisStreamsClient extends ClientProxy {
	private connection: StreamConnection | undefined

	constructor(private readonly options: RedisStreamsClientOptions) {
		super()
	}

	// The connection is opened and connected once; an opening that fails is
	// forgotten, so the next emit tries again instead of reusing a dead one.
	async connect(): Promise<StreamConnection> {
		if (this.connection) return this.connection

		const connection = this.options.connect()
		try {
			await connection.connect()
		} catch (error) {
			connection.disconnect()

			throw error
		}

		this.connection = connection

		return connection
	}

	async close(): Promise<void> {
		const connection = this.connection
		if (!connection) return

		this.connection = undefined

		await connection.quit()
	}

	unwrap<Connection>(): Connection {
		return this.connection as Connection
	}

	protected publish(
		packet: ReadPacket,
		callback: (packet: WritePacket) => void,
	): () => void {
		const pattern = this.normalizePattern(packet.pattern)
		callback({
			err: new Error(
				`request-response is not supported by the Redis Streams transport (pattern ${pattern})`,
			),
		})

		return () => undefined
	}

	protected async dispatchEvent<Result = string>(
		packet: ReadPacket,
	): Promise<Result> {
		const type = this.normalizePattern(packet.pattern)
		const fields = streamFields(type, packet.data)
		const maxLength = this.options.maxLength ?? defaultMaxLength

		const connection = await this.connect()

		const id = await connection.xadd(
			this.options.stream,
			'MAXLEN',
			'~',
			maxLength,
			'*',
			...fields,
		)

		if (!id)
			throw new Error(`Redis appended no entry to ${this.options.stream}`)

		return id as Result
	}
}
