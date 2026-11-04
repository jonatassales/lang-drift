import { Module } from '@nestjs/common'

import { DemoEmailsModule } from '@app/demo-emails/demo-emails.module.js'
import { EventLogModule } from '@app/event-log/event-log.module.js'
import { HealthModule } from '@app/health/health.module.js'

@Module({
	imports: [HealthModule, EventLogModule, DemoEmailsModule],
})
export class AppModule {}
