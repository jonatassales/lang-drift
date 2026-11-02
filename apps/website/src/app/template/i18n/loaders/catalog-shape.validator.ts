import bookDemoJson from '@i18n/messages/book-a-demo.json'
import conversionJson from '@i18n/messages/conversion.json'
import earlyAccessJson from '@i18n/messages/early-access.json'
import earlyAccessSentJson from '@i18n/messages/early-access-sent.json'
import finopsJson from '@i18n/messages/finops.json'
import footerJson from '@i18n/messages/footer.json'
import headerJson from '@i18n/messages/header.json'
import microruntimesJson from '@i18n/messages/microruntimes.json'
import platformJson from '@i18n/messages/platform.json'
import privacyJson from '@i18n/messages/privacy.json'
import routeStateJson from '@i18n/messages/route-state.json'
import type { WebsiteLocale } from '@i18n/routing'

import { shape } from './shape.compute'

const locales = ['en', 'pt-BR', 'zh-Hant', 'ja', 'vi'] as const

export function assertCatalogShape(
	name: string,
	catalog: Record<WebsiteLocale, unknown>,
) {
	const base = shape(catalog.en)
	for (const locale of locales) {
		if (shape(catalog[locale]) !== base) {
			throw new Error(
				`${name} locale ${locale} does not match the catalog shape`,
			)
		}
	}
}

assertCatalogShape('header', headerJson)
assertCatalogShape('privacy', privacyJson)
assertCatalogShape('conversion', conversionJson)
assertCatalogShape('early-access', earlyAccessJson)
assertCatalogShape('early-access-sent', earlyAccessSentJson)
assertCatalogShape('footer', footerJson)
assertCatalogShape('platform', platformJson)
assertCatalogShape('route-state', routeStateJson)
assertCatalogShape('finops', finopsJson)
assertCatalogShape('microruntimes', microruntimesJson)
assertCatalogShape('book-a-demo', bookDemoJson)
