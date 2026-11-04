import type { Resend } from 'resend'

import type {
	EmailSender,
	OutgoingEmail,
	SentEmail,
} from './email-sender.port.js'

// Resend behind the sender port. The idempotency key makes a redelivered
// event send nothing twice; an error answer becomes a thrown reason.
export class ResendEmailSender implements EmailSender {
	constructor(private readonly resend: Resend) {}

	async send(email: OutgoingEmail, idempotencyKey: string): Promise<SentEmail> {
		const response = await this.resend.emails.send(
			{
				from: email.from,
				to: email.to,
				subject: email.subject,
				text: email.text,
			},
			{
				idempotencyKey,
			},
		)

		if (response.error) {
			throw new Error(`Resend refused the email: ${response.error.message}`)
		}

		if (!response.data?.id) {
			throw new Error('Resend answered without a message id')
		}

		return {
			id: response.data.id,
		}
	}
}
