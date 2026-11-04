import type { RedisStreamsContext } from '@langdrift/events/transport'
import { Controller, Inject } from '@nestjs/common'
import { Ctx, EventPattern, Payload } from '@nestjs/microservices'

import { ArchiveService } from './archive.service.js'

// The archive consumer: every event of the stream, whatever its type.
@Controller()
export class ArchiveController {
	constructor(
		@Inject(ArchiveService) private readonly service: ArchiveService,
	) {}

	@EventPattern<string>('*')
	async archive(
		@Payload() envelope: unknown,
		@Ctx() context: RedisStreamsContext,
	): Promise<void> {
		await this.service.store(envelope, context.getEntryId())
	}
}
