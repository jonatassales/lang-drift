import { localePath, type WebsiteLocale } from '@i18n/routing'

import type pages from './pages.json'

export type PageId = keyof typeof pages

// Pages whose URL is not their id, so the metadata JSON stores no route path.
const nestedPagePath: Partial<Record<PageId, string>> = {
	finops: 'platform/finops/',
	microruntimes: 'platform/microruntimes/',
}

export function pagePath(page: PageId, locale: WebsiteLocale) {
	if (page === 'home') return localePath[locale]

	return `${localePath[locale]}${nestedPagePath[page] ?? `${page}/`}`
}
