import 'reflect-metadata'

import { NestFactory } from '@nestjs/core'
import { ExpressAdapter } from '@nestjs/platform-express'

import { AppModule } from '@app/app.module.js'

async function bootstrap() {
	const port = process.env.PORT
	if (!port) throw new Error('Missing PORT')

	const app = await NestFactory.create(AppModule, new ExpressAdapter())

	app.enableShutdownHooks()

	await app.listen(Number(port))
}

await bootstrap()
