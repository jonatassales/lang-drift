import 'reflect-metadata'

import { hostname } from 'node:os'
import {
	demoRequestDeadLetterStream,
	demoRequestMailerGroup,
	demoRequestStream,
} from '@langdrift/events'
import { RedisStreamsServer } from '@langdrift/events/transport'
import { NestFactory } from '@nestjs/core'
import type { MicroserviceOptions } from '@nestjs/microservices'
import { ExpressAdapter } from '@nestjs/platform-express'

import { AppModule } from '@app/app.module.js'
import type { EventLogConnectionOpener } from '@app/event-log/event-log-connection-opener.factory.js'
import { eventLogConnectionOpenerToken } from '@app/event-log/event-log-connection-opener.token.js'

// A hybrid application: the health route, plus the mailer consumer group on
// the Redis Streams transport, in one process.
async function bootstrap() {
	const port = process.env.PORT
	if (!port) throw new Error('Missing PORT')

	const app = await NestFactory.create(AppModule, new ExpressAdapter())

	const open = app.get<EventLogConnectionOpener>(eventLogConnectionOpenerToken)
	app.connectMicroservice<MicroserviceOptions>({
		strategy: new RedisStreamsServer({
			stream: demoRequestStream,
			group: demoRequestMailerGroup,
			consumer: hostname(),
			deadLetterStream: demoRequestDeadLetterStream,
			connect: () => open.consumer(),
		}),
	})
	app.enableShutdownHooks()

	await app.startAllMicroservices()

	await app.listen(Number(port))
}

await bootstrap()
