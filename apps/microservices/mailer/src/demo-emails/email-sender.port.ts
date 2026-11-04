export interface OutgoingEmail {
	from: string
	to: string
	subject: string
	text: string
}

export interface SentEmail {
	id: string
}

// What the mailer needs from an email provider: one idempotent send that
// answers with the provider's message id.
export interface EmailSender {
	send(email: OutgoingEmail, idempotencyKey: string): Promise<SentEmail>
}
