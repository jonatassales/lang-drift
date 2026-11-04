import { Module } from '@nestjs/common'

import { ArchiveModule } from '@app/archive/archive.module.js'
import { DemoRequestsModule } from '@app/demo-requests/demo-requests.module.js'
import { EventLogModule } from '@app/event-log/event-log.module.js'
import { HealthModule } from '@app/health/health.module.js'

@Module({
	imports: [HealthModule, EventLogModule, DemoRequestsModule, ArchiveModule],
})
export class AppModule {}
