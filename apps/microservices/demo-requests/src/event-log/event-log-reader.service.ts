import type { OnApplicationShutdown } from '@nestjs/common'
import { Inject, Injectable, Logger } from '@nestjs/common'
import { Observable } from 'rxjs'

import { demoRequestStream } from 'streaming'
import type {
	StreamConnection,
	StreamEntry,
	StreamItem,
} from 'streaming/transport'
import { parseStreamEntry } from 'streaming/transport'

import { aggregateIdOf } from './aggregate-id-of.compute.js'
import type { EventLogConnectionOpener } from './event-log-connection-opener.factory.js'
import { eventLogConnectionOpenerToken } from './event-log-connection-opener.token.js'

export interface FollowEntry {
	kind: 'entry'
	entry: StreamEntry
}

export interface FollowPing {
	kind: 'ping'
}

export type FollowEvent = FollowEntry | FollowPing

const pageSize = 500
const blockMs = 15000
const streamStart = '0-0'

// Reads the event log without a consumer group: a request's entries by
// ranging the stream, or the stream as it grows, one blocking connection per
// follower, released when the follower leaves.
@Injectable()
export class EventLogReaderService implements OnApplicationShutdown {
	private readonly logger = new Logger(EventLogReaderService.name)
	private shared: StreamConnection | undefined

	constructor(
		@Inject(eventLogConnectionOpenerToken)
		private readonly open: EventLogConnectionOpener,
	) {}

	async readAggregate(aggregateId: string): Promise<StreamEntry[]> {
		const connection = this.sharedConnection()
		const entries: StreamEntry[] = []
		let start = '-'

		while (true) {
			const page = await connection.xrange(
				demoRequestStream,
				start,
				'+',
				'COUNT',
				pageSize,
			)

			for (const item of page) {
				if (aggregateIdOf(item) !== aggregateId) continue

				const entry = this.parse(item)
				if (entry) entries.push(entry)
			}

			if (page.length < pageSize) return entries

			start = `(${page[page.length - 1][0]}`
		}
	}

	follow(fromId?: string): Observable<FollowEvent> {
		return new Observable((subscriber) => {
			const connection = this.open.consumer()
			let stopped = false
			let last = fromId

			// The live end is resolved to a concrete id once, so nothing appended
			// between two blocking reads is skipped.
			const run = async () => {
				if (last === undefined) {
					try {
						last = await this.latestEntryId(connection)
					} catch (error) {
						if (!stopped) subscriber.error(error)

						return
					}
				}

				while (!stopped) {
					let reply: Awaited<ReturnType<StreamConnection['xread']>>
					try {
						reply = await connection.xread(
							'BLOCK',
							blockMs,
							'STREAMS',
							demoRequestStream,
							last,
						)
					} catch (error) {
						// Unsubscribing disconnects the blocking read, which rejects it: that
						// rejection is the teardown, not a failure to report.
						if (!stopped) subscriber.error(error)

						return
					}

					if (stopped) return

					const items = reply?.[0]?.[1] ?? []
					if (items.length === 0) {
						subscriber.next({
							kind: 'ping',
						})
						continue
					}

					for (const item of items) {
						last = item[0]
						const entry = this.parse(item)
						if (entry) {
							subscriber.next({
								kind: 'entry',
								entry,
							})
						}
					}
				}
			}

			void run()

			return () => {
				stopped = true
				connection.disconnect()
			}
		})
	}

	async onApplicationShutdown(): Promise<void> {
		const connection = this.shared
		if (!connection) return

		this.shared = undefined

		await connection.quit()
	}

	// The id of the newest entry, or the start when the stream is empty.
	private async latestEntryId(connection: StreamConnection): Promise<string> {
		const newest = await connection.xrevrange(
			demoRequestStream,
			'+',
			'-',
			'COUNT',
			1,
		)

		return newest[0]?.[0] ?? streamStart
	}

	private sharedConnection(): StreamConnection {
		if (!this.shared) this.shared = this.open.consumer()

		return this.shared
	}

	// An entry that cannot be read is skipped here; the consumer groups park
	// it on the dead-letter stream with the reason.
	private parse(item: StreamItem): StreamEntry | null {
		try {
			return parseStreamEntry(item[0], item[1])
		} catch (error) {
			const reason = error instanceof Error ? error.message : String(error)
			this.logger.warn(`Entry ${item[0]} skipped: ${reason}`)

			return null
		}
	}
}
