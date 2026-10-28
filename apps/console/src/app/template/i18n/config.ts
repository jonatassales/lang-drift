export const locales = ['en', 'pt-br', 'zh-hant', 'ja'] as const
export type Locale = (typeof locales)[number]
export const localeCookie = 'langdrift-console-locale'
export function isLocale(value: unknown): value is Locale {
	return typeof value === 'string' && locales.some((locale) => locale === value)
}

export const languageTags = {
	en: 'en',
	'pt-br': 'pt-BR',
	'zh-hant': 'zh-Hant',
	ja: 'ja',
} as const
