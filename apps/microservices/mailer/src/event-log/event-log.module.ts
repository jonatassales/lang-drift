import type { RedisStreamsClient } from '@langdrift/events/transport'
import type { OnApplicationShutdown } from '@nestjs/common'
import { Inject, Module } from '@nestjs/common'

import { createEventLogConnectionOpener } from './event-log-connection-opener.factory.js'
import { eventLogConnectionOpenerToken } from './event-log-connection-opener.token.js'
import { createEventLogProducer } from './event-log-producer.factory.js'
import { eventLogProducerToken } from './event-log-producer.token.js'

@Module({
	providers: [
		{
			provide: eventLogConnectionOpenerToken,
			useFactory: createEventLogConnectionOpener,
		},
		{
			provide: eventLogProducerToken,
			useFactory: createEventLogProducer,
			inject: [eventLogConnectionOpenerToken],
		},
	],
	exports: [eventLogConnectionOpenerToken, eventLogProducerToken],
})
export class EventLogModule implements OnApplicationShutdown {
	constructor(
		@Inject(eventLogProducerToken)
		private readonly producer: RedisStreamsClient,
	) {}

	async onApplicationShutdown(): Promise<void> {
		await this.producer.close()
	}
}
