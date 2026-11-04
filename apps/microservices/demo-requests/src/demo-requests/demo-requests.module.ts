import { Module } from '@nestjs/common'

import { EventLogModule } from '@app/event-log/event-log.module.js'

import { BearerTokenGuard } from './bearer-token.guard.js'
import { DemoRequestsController } from './demo-requests.controller.js'

@Module({
	imports: [EventLogModule],
	controllers: [DemoRequestsController],
	providers: [BearerTokenGuard],
})
export class DemoRequestsModule {}
