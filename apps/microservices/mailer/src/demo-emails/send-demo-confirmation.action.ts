import type { DemoEmailsCatalog } from '@app/i18n/messages/messages.types.js'
import type { DemoRequestReceived } from 'infra/streaming'

import type { DeliveredEmail } from './delivered-email.types.js'
import type { EmailSender } from './email-sender.port.js'
import { formatConfirmationMessage } from './format-confirmation-message.fmt.js'
import { senderAddress } from './sender-address.data.js'

// Tells the visitor their request arrived and the team will be in touch,
// once per received event however often the event is delivered.
export async function sendDemoConfirmation(
	sender: EmailSender,
	received: DemoRequestReceived,
	catalog: DemoEmailsCatalog,
): Promise<DeliveredEmail> {
	const message = formatConfirmationMessage(catalog, received.payload.locale, {
		company: received.payload.company,
	})

	const sent = await sender.send(
		{
			from: senderAddress,
			to: received.payload.email,
			subject: message.subject,
			text: message.text,
		},
		`${received.id}/confirmation`,
	)

	return {
		id: sent.id,
		to: received.payload.email,
	}
}
