import type { WebsiteLocale } from '@i18n/routing'

import bookDemoJson from './book-a-demo.json'
import conversionJson from './conversion.json'
import earlyAccessJson from './early-access.json'
import earlyAccessSentJson from './early-access-sent.json'
import finopsJson from './finops.json'
import footerJson from './footer.json'
import headerJson from './header.json'
import microruntimesJson from './microruntimes.json'
import platformJson from './platform.json'
import privacyJson from './privacy.json'
import routeStateJson from './route-state.json'
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
