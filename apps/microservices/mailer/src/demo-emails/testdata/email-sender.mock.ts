import type {
	EmailSender,
	OutgoingEmail,
	SentEmail,
} from '@app/demo-emails/email-sender.port.js'

export interface RecordedSend {
	email: OutgoingEmail
	idempotencyKey: string
}

// Records every send and fails the ones whose recipient is on its list.
export class FakeEmailSender implements EmailSender {
	readonly sends: RecordedSend[] = []
	readonly failing = new Set<string>()
	private counter = 0

	async send(email: OutgoingEmail, idempotencyKey: string): Promise<SentEmail> {
		this.sends.push({
			email,
			idempotencyKey,
		})
		if (this.failing.has(email.to)) {
			throw new Error(`Resend answered 422: invalid recipient ${email.to}`)
		}

		this.counter += 1

		return {
			id: `msg_${this.counter}`,
		}
	}
}
