import type {
	CustomTransportStrategy,
	MessageHandler,
} from '@nestjs/microservices'
import { Server } from '@nestjs/microservices'
import { finalize, isObservable, lastValueFrom } from 'rxjs'

import { parseStreamEntry } from './parse-stream-entry.compute.js'
import { RedisStreamsContext } from './redis-streams.context.js'
import type { RedisStreamsServerOptions } from './redis-streams-options.types.js'
import { shouldDeadLetter } from './should-dead-letter.compute.js'
import type { StreamConnection, StreamItem } from './stream-connection.types.js'
import type { StreamEntry } from './stream-entry.types.js'

const defaults = {
	blockMs: 5000,
	batchSize: 16,
	claimMinIdleMs: 60000,
	claimEveryMs: 30000,
	maxDeliveries: 5,
}

const wildcardPattern = '*'
const groupExists = 'BUSYGROUP'

interface Delivery {
	entry: StreamEntry
	deliveries: number
}

// Consumes one stream as a member of one consumer group. An entry is handed
// to the handler registered for its type (or to the wildcard handler) and
// acknowledged once that handler completes; a throw leaves it pending for the
// next claim, and an entry claimed too many times is parked on the
// dead-letter stream.
export class RedisStreamsServer
	extends Server
	implements CustomTransportStrategy
{
	override readonly transportId = Symbol('redis-streams')
	override readonly propagatesEventHandlerErrors = true
	private readonly blocking: StreamConnection
	private readonly commands: StreamConnection
	private readonly blockMs: number
	private readonly batchSize: number
	private readonly claimMinIdleMs: number
	private readonly claimEveryMs: number
	private readonly maxDeliveries: number
	private stopping = false
	private loop: Promise<void> | undefined
	private lastClaimAt: number | undefined

	constructor(private readonly options: RedisStreamsServerOptions) {
		super()
		this.blocking = options.connect()
		this.commands = options.connect()
		this.blockMs = options.blockMs ?? defaults.blockMs
		this.batchSize = options.batchSize ?? defaults.batchSize
		this.claimMinIdleMs = options.claimMinIdleMs ?? defaults.claimMinIdleMs
		this.claimEveryMs = options.claimEveryMs ?? defaults.claimEveryMs
		this.maxDeliveries = options.maxDeliveries ?? defaults.maxDeliveries
	}

	async listen(callback: (...optionalParams: unknown[]) => unknown) {
		await this.createGroup()

		this.loop = this.consume()
		callback()
	}

	async close() {
		this.stopping = true
		this.blocking.disconnect()

		await this.loop

		await this.commands.quit()
	}

	on() {
		return undefined
	}

	unwrap<Connection>(): Connection {
		return this.commands as Connection
	}

	private async createGroup() {
		try {
			await this.commands.xgroup(
				'CREATE',
				this.options.stream,
				this.options.group,
				'$',
				'MKSTREAM',
			)
		} catch (error) {
			const message = error instanceof Error ? error.message : String(error)
			if (!message.startsWith(groupExists)) throw error
		}
	}

	private async consume() {
		while (!this.stopping) {
			await this.claimWhenDue()

			const items = await this.read()

			for (const item of items) {
				await this.deliver(item, 1)
			}
		}
	}

	// An entry that cannot be read is parked as it is, with the reason, so the
	// group moves on and the raw fields stay inspectable.
	private async deliver(item: StreamItem, deliveries: number) {
		let entry: StreamEntry
		try {
			entry = parseStreamEntry(item[0], item[1])
		} catch (error) {
			const reason = error instanceof Error ? error.message : String(error)

			await this.parkRaw(item, reason)

			return
		}

		await this.handle({
			entry,
			deliveries,
		})
	}

	private async read(): Promise<StreamItem[]> {
		try {
			const reply = await this.blocking.xreadgroup(
				'GROUP',
				this.options.group,
				this.options.consumer,
				'COUNT',
				this.batchSize,
				'BLOCK',
				this.blockMs,
				'STREAMS',
				this.options.stream,
				'>',
			)

			return reply?.[0]?.[1] ?? []
		} catch (error) {
			if (this.stopping) return []

			throw error
		}
	}

	private async claimWhenDue() {
		const now = Date.now()
		const due =
			this.lastClaimAt === undefined ||
			now - this.lastClaimAt >= this.claimEveryMs
		if (!due) return

		this.lastClaimAt = now

		const reply = await this.commands.xautoclaim(
			this.options.stream,
			this.options.group,
			this.options.consumer,
			this.claimMinIdleMs,
			'0-0',
			'COUNT',
			this.batchSize,
		)

		const items = (reply[1] ?? []) as StreamItem[]
		for (const item of items) {
			const deliveries = await this.deliveriesOf(item[0])

			await this.deliver(item, deliveries)
		}
	}

	private async deliveriesOf(id: string): Promise<number> {
		const reply = await this.commands.xpending(
			this.options.stream,
			this.options.group,
			id,
			id,
			1,
		)

		const pending = reply[0] as [string, string, number, number] | undefined

		return pending?.[3] ?? 1
	}

	private async handle(delivery: Delivery) {
		const { entry, deliveries } = delivery

		if (shouldDeadLetter(deliveries, this.maxDeliveries)) {
			await this.park(entry, `delivered ${deliveries} times`)

			return
		}

		const handler =
			this.getHandlerByPattern(entry.type) ??
			this.getHandlerByPattern(wildcardPattern)
		if (!handler) {
			await this.acknowledge(entry.id)

			return
		}

		const context = new RedisStreamsContext([
			this.options.stream,
			entry.id,
			this.options.group,
			deliveries,
		])

		try {
			await this.dispatch(handler, entry, context)
		} catch (error) {
			const message = error instanceof Error ? error.message : String(error)
			this.logger.error(
				`Entry ${entry.id} (${entry.type}) stays pending: ${message}`,
			)

			return
		}

		await this.acknowledge(entry.id)
	}

	private dispatch(
		handler: MessageHandler,
		entry: StreamEntry,
		context: RedisStreamsContext,
	): Promise<void> {
		return new Promise((resolve, reject) => {
			this.runWithProcessingHooks(context, async (runEndHook) => {
				try {
					const result = await handler(entry.envelope, context)

					if (isObservable(result)) {
						await lastValueFrom(result.pipe(finalize(runEndHook)), {
							defaultValue: undefined,
						})
					}

					runEndHook()
					resolve()
				} catch (error) {
					runEndHook()
					reject(error)
				}
			})
		})
	}

	private async park(entry: StreamEntry, reason: string) {
		await this.commands.xadd(
			this.options.deadLetterStream,
			'*',
			'type',
			entry.type,
			'aggregateId',
			entry.aggregateId,
			'envelope',
			JSON.stringify(entry.envelope),
			'reason',
			reason,
			'source',
			entry.id,
		)

		await this.acknowledge(entry.id)
	}

	private async parkRaw(item: StreamItem, reason: string) {
		await this.commands.xadd(
			this.options.deadLetterStream,
			'*',
			'fields',
			JSON.stringify(item[1]),
			'reason',
			reason,
			'source',
			item[0],
		)

		await this.acknowledge(item[0])
	}

	private async acknowledge(id: string) {
		await this.commands.xack(this.options.stream, this.options.group, id)
	}
}
