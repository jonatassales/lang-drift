import { Resend } from 'resend'

import type { EmailSender } from './email-sender.port.js'
import { ResendEmailSender } from './resend-email-sender.adapter.js'

// Builds the Resend sender from the key, read here where it is needed.
export function createEmailSender(): EmailSender {
	const apiKey = process.env.RESEND_API_KEY
	if (!apiKey) throw new Error('Missing RESEND_API_KEY')

	return new ResendEmailSender(new Resend(apiKey))
}
