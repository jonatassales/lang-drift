import type {
	DemoEmailsCatalog,
	MailLocale,
} from '@app/i18n/messages/messages.types.js'
import type { DemoRequestReceived } from 'infra/streaming'

import type { EmailMessage } from './email-message.types.js'
import { formatIcuMessage } from './format-icu-message.fmt.js'

// The notification to the team, in the visitor's locale so whoever answers
// knows which language the visitor expects.
export function formatLeadNotificationMessage(
	catalog: DemoEmailsCatalog,
	locale: MailLocale,
	submission: DemoRequestReceived['payload'],
): EmailMessage {
	const messages = catalog[locale].leadNotification
	const format = (message: string) =>
		formatIcuMessage(message, locale, submission)

	return {
		subject: format(messages.subject),
		text: [
			format(messages.intro),
			'',
			format(messages.company),
			format(messages.contact),
			format(messages.position),
			format(messages.locale),
			format(messages.submittedAt),
		].join('\n'),
	}
}
