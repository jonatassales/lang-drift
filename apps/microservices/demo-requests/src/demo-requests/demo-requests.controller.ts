import { randomUUID } from 'node:crypto'
import type { DemoRequestSubmission } from '@langdrift/events'
import { demoRequestSubmissionSchema } from '@langdrift/events'
import {
	Body,
	Controller,
	HttpCode,
	Inject,
	Post,
	Res,
	ServiceUnavailableException,
	UseGuards,
} from '@nestjs/common'
import { lastValueFrom } from 'rxjs'

import type { EventLogProducer } from '@app/event-log/event-log-producer.port.js'
import { eventLogProducerToken } from '@app/event-log/event-log-producer.token.js'

import { BearerTokenGuard } from './bearer-token.guard.js'
import { demoRequestReceivedEnvelope } from './demo-request-received-envelope.compute.js'
import { ZodValidationPipe } from './zod-validation.pipe.js'

interface DemoRequestAccepted {
	id: string
}

// The one thing the reply needs from the HTTP response: a header.
interface ReplyHeaders {
	setHeader(name: string, value: string): unknown
}

@Controller('v1/demo-requests')
@UseGuards(BearerTokenGuard)
export class DemoRequestsController {
	constructor(
		@Inject(eventLogProducerToken)
		private readonly producer: EventLogProducer,
	) {}

	// A submission the website already validated becomes the first event of a
	// new request; the reply names the request and where to read it.
	@Post()
	@HttpCode(202)
	async create(
		@Body(new ZodValidationPipe(demoRequestSubmissionSchema))
		submission: DemoRequestSubmission,
		@Res({ passthrough: true }) response: ReplyHeaders,
	): Promise<DemoRequestAccepted> {
		const id = randomUUID()
		const envelope = demoRequestReceivedEnvelope(submission, id, new Date())

		try {
			await lastValueFrom(this.producer.emit(envelope.type, envelope))
		} catch (error) {
			const reason = error instanceof Error ? error.message : String(error)

			throw new ServiceUnavailableException(
				`The event log is unavailable: ${reason}`,
			)
		}

		response.setHeader('Location', `/v1/demo-requests/${id}`)

		return {
			id,
		}
	}
}
