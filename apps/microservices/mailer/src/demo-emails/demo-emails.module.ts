import { Module } from '@nestjs/common'

import { EventLogModule } from '@app/event-log/event-log.module.js'

import { DemoEmailsController } from './demo-emails.controller.js'
import { createEmailSender } from './email-sender.factory.js'
import { emailSenderToken } from './email-sender.token.js'

@Module({
	imports: [EventLogModule],
	controllers: [DemoEmailsController],
	providers: [
		{
			provide: emailSenderToken,
			useFactory: createEmailSender,
		},
	],
})
export class DemoEmailsModule {}
