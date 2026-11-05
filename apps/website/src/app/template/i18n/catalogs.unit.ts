import { createTranslator } from 'next-intl'
import { describe, expect, it } from 'vitest'

import {
	bookDemoMessages,
	ciGatesMessages,
	cliMessages,
	conversionMessages,
	finopsMessages,
	footerMessages,
	headerMessages,
	microruntimesMessages,
	platformMessages,
	privacyMessages,
	safeSeatMessages,
} from '@i18n/messages/index'
import routeState from '@i18n/messages/route-state.json'

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

function tagsOf(message: string): string[] {
	return Array.from(
		new Set(
			Array.from(message.matchAll(/<([A-Za-z]\w*)>/g), (match) => match[1]),
		),
	).sort()
}

describe('Website production translations', () => {
	it.each(['en', 'pt-BR', 'zh-Hant', 'ja', 'vi'] as const)(
		'%s preserves every key, ICU argument and rich-text tag and formats the acceptance counts',
		(locale) => {
			for (const family of [
				headerMessages,
				conversionMessages,
				footerMessages,
				platformMessages,
				privacyMessages,
				finopsMessages,
				microruntimesMessages,
				safeSeatMessages,
				cliMessages,
				ciGatesMessages,
				bookDemoMessages,
				routeState,
			]) {
				const source = leaves(family.en)
				const translated = leaves(family[locale])
				expect(Object.keys(translated).sort()).toEqual(
					Object.keys(source).sort(),
				)
				for (const [key, message] of Object.entries(translated)) {
					expect(argumentsOf(message), `${locale} ${key} arguments`).toEqual(
						argumentsOf(source[key]),
					)
					expect(tagsOf(message), `${locale} ${key} tags`).toEqual(
						tagsOf(source[key]),
					)
					const errors: unknown[] = []
					const translate = createTranslator({
						locale,
						messages: { value: message },
						onError: (error) => errors.push(error),
					})
					for (const count of [0, 1, 2, 3, 12]) {
						const values = Object.fromEntries(
							argumentsOf(message).map((name) => [name, count]),
						)
						const tags = Object.fromEntries(
							tagsOf(message).map((name) => [
								name,
								(chunks: unknown) => String(chunks),
							]),
						)
						translate.rich('value', {
							...values,
							...tags,
						})
					}
					expect(errors, `${locale} ${key}: ${message}`).toEqual([])
				}
			}
		},
	)
})
