import { randomUUID } from 'node:crypto'
import type { MessageEvent } from '@nestjs/common'
import {
	Body,
	Controller,
	Get,
	Headers,
	HttpCode,
	Inject,
	NotFoundException,
	Param,
	Post,
	Query,
	Res,
	ServiceUnavailableException,
	Sse,
	UseGuards,
} from '@nestjs/common'
import type { Observable } from 'rxjs'
import { lastValueFrom, map } from 'rxjs'

import type { EventLogProducer } from '@app/event-log/event-log-producer.port.js'
import { eventLogProducerToken } from '@app/event-log/event-log-producer.token.js'
import { EventLogReaderService } from '@app/event-log/event-log-reader.service.js'
import type { DemoRequestSubmission } from 'infra/streaming'
import { demoRequestSubmissionSchema } from 'infra/streaming'

import { BearerTokenGuard } from './bearer-token.guard.js'
import { demoRequestEventsOf } from './demo-request-events-of.compute.js'
import { demoRequestReceivedEnvelope } from './demo-request-received-envelope.compute.js'
import type { DemoRequestView } from './demo-request-view.types.js'
import { rebuildDemoRequest } from './rebuild-demo-request.compute.js'
import { toSseEvent } from './to-sse-event.map.js'
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
		@Inject(EventLogReaderService)
		private readonly reader: EventLogReaderService,
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

	// The log as it grows. A client that reconnects sends the last id it saw
	// and misses nothing; `from` replays from any entry, `0-0` from the start.
	// Declared before the id route so `events` is never read as an id.
	@Sse('events')
	events(
		@Headers('last-event-id') lastEventId: string | undefined,
		@Query('from') from: string | undefined,
	): Observable<MessageEvent> {
		return this.reader.follow(from ?? lastEventId).pipe(map(toSseEvent))
	}

	@Get(':id')
	async read(@Param('id') id: string): Promise<DemoRequestView> {
		const entries = await this.reader.readAggregate(id)

		const view = rebuildDemoRequest(demoRequestEventsOf(entries))
		if (!view) throw new NotFoundException(`No demo request ${id}`)

		return view
	}
}
