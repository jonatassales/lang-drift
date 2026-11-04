import 'reflect-metadata'

import { hostname } from 'node:os'
import { NestFactory } from '@nestjs/core'
import type { MicroserviceOptions } from '@nestjs/microservices'
import { ExpressAdapter } from '@nestjs/platform-express'

import { AppModule } from '@app/app.module.js'
import type { EventLogConnectionOpener } from '@app/event-log/event-log-connection-opener.factory.js'
import { eventLogConnectionOpenerToken } from '@app/event-log/event-log-connection-opener.token.js'
import {
	demoRequestArchiveGroup,
	demoRequestDeadLetterStream,
	demoRequestStream,
} from 'streaming'
import { RedisStreamsServer } from 'streaming/transport'

// A hybrid application: the HTTP routes, plus the archive consumer group on
// the Redis Streams transport, in one process.
async function bootstrap() {
	const port = process.env.PORT
	if (!port) throw new Error('Missing PORT')

	const app = await NestFactory.create(AppModule, new ExpressAdapter())

	const open = app.get<EventLogConnectionOpener>(eventLogConnectionOpenerToken)
	app.connectMicroservice<MicroserviceOptions>({
		strategy: new RedisStreamsServer({
			stream: demoRequestStream,
			group: demoRequestArchiveGroup,
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
