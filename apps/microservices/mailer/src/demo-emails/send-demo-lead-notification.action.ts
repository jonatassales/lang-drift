import type { DemoEmailsCatalog } from '@app/i18n/messages/messages.types.js'
import type { DemoRequestReceived } from 'streaming'

import type { DeliveredEmail } from './delivered-email.types.js'
import type { EmailSender } from './email-sender.port.js'
import { formatLeadNotificationMessage } from './format-lead-notification-message.fmt.js'
import { senderAddress } from './sender-address.data.js'

// Tells the team about the lead, at the address the team configured, read
// here where it is needed.
export async function sendDemoLeadNotification(
	sender: EmailSender,
	received: DemoRequestReceived,
	catalog: DemoEmailsCatalog,
): Promise<DeliveredEmail> {
	const to = process.env.DEMO_LEAD_EMAIL
	if (!to) throw new Error('Missing DEMO_LEAD_EMAIL')

	const message = formatLeadNotificationMessage(
		catalog,
		received.payload.locale,
		received.payload,
	)

	const sent = await sender.send(
		{
			from: senderAddress,
			to,
			subject: message.subject,
			text: message.text,
		},
		`${received.id}/lead-notification`,
	)

	return {
		id: sent.id,
		to,
	}
}
