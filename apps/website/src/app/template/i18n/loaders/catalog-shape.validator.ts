import bookDemoJson from '@i18n/messages/book-a-demo.json'
import ciGatesJson from '@i18n/messages/ci-gates.json'
import cliJson from '@i18n/messages/cli.json'
import conversionJson from '@i18n/messages/conversion.json'
import finopsJson from '@i18n/messages/finops.json'
import footerJson from '@i18n/messages/footer.json'
import headerJson from '@i18n/messages/header.json'
import microruntimesJson from '@i18n/messages/microruntimes.json'
import platformJson from '@i18n/messages/platform.json'
import privacyJson from '@i18n/messages/privacy.json'
import routeStateJson from '@i18n/messages/route-state.json'
import safeSeatJson from '@i18n/messages/safe-seat.json'
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
assertCatalogShape('footer', footerJson)
assertCatalogShape('platform', platformJson)
assertCatalogShape('route-state', routeStateJson)
assertCatalogShape('finops', finopsJson)
assertCatalogShape('microruntimes', microruntimesJson)
assertCatalogShape('safe-seat', safeSeatJson)
assertCatalogShape('cli', cliJson)
assertCatalogShape('ci-gates', ciGatesJson)
assertCatalogShape('book-a-demo', bookDemoJson)
