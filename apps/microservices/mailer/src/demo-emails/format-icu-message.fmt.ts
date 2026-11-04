import { IntlMessageFormat } from 'intl-messageformat'

import type { MailLocale } from '@app/i18n/messages/messages.types.js'

// One ICU message of the catalog, formatted for a locale with its arguments.
export function formatIcuMessage(
	message: string,
	locale: MailLocale,
	values: object,
): string {
	const formatted = new IntlMessageFormat(message, locale).format(
		Object.fromEntries(Object.entries(values)),
	)

	return String(formatted)
}
