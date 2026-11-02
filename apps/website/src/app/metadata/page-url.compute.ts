import type { WebsiteLocale } from '@i18n/routing'

import { type PageId, pagePath } from './page-path.compute'
import site from './site.json'

export function pageUrl(page: PageId, locale: WebsiteLocale) {
	return new URL(pagePath(page, locale), site.url).toString()
}
