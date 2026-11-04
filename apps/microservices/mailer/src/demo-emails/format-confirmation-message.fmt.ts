import type {
	DemoEmailsCatalog,
	MailLocale,
} from '@app/i18n/messages/messages.types.js'

import type { EmailMessage } from './email-message.types.js'
import { formatIcuMessage } from './format-icu-message.fmt.js'

export interface ConfirmationArguments {
	company: string
}

// The confirmation to the visitor, in the locale of their request.
export function formatConfirmationMessage(
	catalog: DemoEmailsCatalog,
	locale: MailLocale,
	values: ConfirmationArguments,
): EmailMessage {
	const messages = catalog[locale].confirmation
	const format = (message: string) => formatIcuMessage(message, locale, values)

	return {
		subject: format(messages.subject),
		text: [
			format(messages.greeting),
			'',
			format(messages.thanks),
			'',
			format(messages.followUp),
			'',
			format(messages.signoff),
		].join('\n'),
	}
}
