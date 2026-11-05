import { IntlMessageFormat } from 'intl-messageformat'
import { describe, expect, it } from 'vitest'

import { demoEmailsMessages } from '@app/i18n/messages/index.js'
import { demoRequestLocales } from 'infra/streaming'

function leaves(value: unknown, prefix = ''): Record<string, string> {
	if (typeof value === 'string') return { [prefix]: value }
	if (!value || typeof value !== 'object') return {}

	return Object.fromEntries(
		Object.entries(value).flatMap(([key, child]) =>
			Object.entries(leaves(child, prefix ? `${prefix}.${key}` : key)),
		),
	)
}

function argumentsOf(message: string): string[] {
	return Array.from(
		new Set(
			Array.from(
				message.matchAll(/\{([A-Za-z]\w*)\s*(?:,|\})/g),
				(match) => match[1],
			),
		),
	).sort()
}

const sampleArguments = {
	company: 'Acme',
	email: 'ana@acme.com',
	position: 'cto',
	locale: 'pt-BR',
	submittedAt: '2026-11-03T14:05:00.000Z',
}

describe('mailer/i18n/catalogs', () => {
	it.each(demoRequestLocales)(
		'%s carries every key and ICU argument of en, and every message formats for the counts 0, 1, 2, 3 and 12',
		(locale) => {
			const reference = leaves(demoEmailsMessages.en)
			const translated = leaves(demoEmailsMessages[locale])

			expect(Object.keys(translated).sort()).toEqual(
				Object.keys(reference).sort(),
			)

			for (const [key, message] of Object.entries(translated)) {
				expect(argumentsOf(message), key).toEqual(argumentsOf(reference[key]))

				for (const count of [0, 1, 2, 3, 12]) {
					const formatted = new IntlMessageFormat(message, locale).format({
						...sampleArguments,
						count,
					})

					expect(typeof formatted, key).toBe('string')
					expect(formatted, key).not.toBe('')
				}
			}
		},
	)
})
