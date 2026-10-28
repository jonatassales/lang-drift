import { createTranslator } from 'next-intl'
import { describe, expect, it } from 'vitest'

import {
	bookDemoMessages,
	conversionMessages,
	demoMessages,
	earlyAccessMessages,
	earlyAccessSentMessages,
	executiveReviewMessages,
	finopsMessages,
	footerMessages,
	headerMessages,
	homeMessages,
	integrationsMessages,
	microruntimesMessages,
	platformMessages,
	privacyMessages,
	visionLoopMessages,
} from './messages'
import routeState from './messages/route-state.json'
import type { WebsiteLocale } from './routing'

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
				homeMessages,
				headerMessages,
				demoMessages,
				conversionMessages,
				earlyAccessMessages,
				earlyAccessSentMessages,
				footerMessages,
				integrationsMessages,
				platformMessages,
				privacyMessages,
				visionLoopMessages,
				executiveReviewMessages,
				finopsMessages,
				microruntimesMessages,
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
					const t = createTranslator({
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
						t.rich('value', {
							...values,
							...tags,
						})
					}
					expect(errors, `${locale} ${key}: ${message}`).toEqual([])
				}
			}
		},
	)
	it.each([
		['en', '0 points'],
		['pt-BR', '0 pontos'],
		['zh-Hant', '0 點'],
		['ja', '0 ポイント'],
		['vi', '0 điểm'],
	] as const)('formats the %s empty count without errors', (locale, text) => {
		const errors: unknown[] = []
		const t = createTranslator({
			locale,
			messages: homeMessages[locale as WebsiteLocale],
			onError: (error) => errors.push(error),
		})
		expect(t('intent.points', { count: 0 })).toBe(text)
		expect(errors).toEqual([])
	})
})
