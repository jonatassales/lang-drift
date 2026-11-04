import type {
	DemoRequestEmailKind,
	DemoRequestReceived,
} from '@langdrift/events'
import {
	demoRequestEventNames,
	demoRequestReceivedSchema,
} from '@langdrift/events'
import { Controller, Inject, Logger } from '@nestjs/common'
import { EventPattern, Payload } from '@nestjs/microservices'
import { lastValueFrom } from 'rxjs'

import type { EventLogProducer } from '@app/event-log/event-log-producer.port.js'
import { eventLogProducerToken } from '@app/event-log/event-log-producer.token.js'
import { demoEmailsMessages } from '@app/i18n/messages/index.js'

import type { DeliveredEmail } from './delivered-email.types.js'
import type { EmailOutcome } from './email-outcome-envelope.compute.js'
import { emailOutcomeEnvelope } from './email-outcome-envelope.compute.js'
import type { EmailSender } from './email-sender.port.js'
import { emailSenderToken } from './email-sender.token.js'
import { sendDemoConfirmation } from './send-demo-confirmation.action.js'
import { sendDemoLeadNotification } from './send-demo-lead-notification.action.js'

// The mailer consumer: a received request gets its confirmation and its
// lead notification, and each outcome goes back on the log as an event.
@Controller()
export class DemoEmailsController {
	private readonly logger = new Logger(DemoEmailsController.name)

	constructor(
		@Inject(emailSenderToken) private readonly sender: EmailSender,
		@Inject(eventLogProducerToken) private readonly producer: EventLogProducer,
	) {}

	@EventPattern<string>(demoRequestEventNames.received)
	async onReceived(@Payload() envelope: unknown): Promise<void> {
		const received = demoRequestReceivedSchema.parse(envelope)

		await this.deliver(received, 'confirmation', () =>
			sendDemoConfirmation(this.sender, received, demoEmailsMessages),
		)

		await this.deliver(received, 'lead-notification', () =>
			sendDemoLeadNotification(this.sender, received, demoEmailsMessages),
		)
	}

	// A failed send is a recorded outcome, never a lost event; only a failure
	// to record it leaves the entry pending for another delivery.
	private async deliver(
		received: DemoRequestReceived,
		kind: DemoRequestEmailKind,
		send: () => Promise<DeliveredEmail>,
	): Promise<void> {
		let outcome: EmailOutcome
		try {
			const delivered = await send()

			outcome = {
				kind,
				outcome: 'sent',
				messageId: delivered.id,
				to: delivered.to,
				now: new Date(),
			}
		} catch (error) {
			const reason = error instanceof Error ? error.message : String(error)
			this.logger.warn(
				`${kind} for request ${received.aggregate.id} failed: ${reason}`,
			)
			outcome = {
				kind,
				outcome: 'failed',
				reason,
				now: new Date(),
			}
		}

		const event = emailOutcomeEnvelope(received, outcome)

		await lastValueFrom(this.producer.emit(event.type, event))
	}
}
