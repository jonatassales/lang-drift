import conversionJson from '@i18n/messages/conversion.json'
import type { ConversionCopy } from '@i18n/messages/messages.types'
import type { WebsiteLocale } from '@i18n/routing'

import { assertCatalogShape } from './assert-catalog-shape.compute'
import { loadConversion } from './load-conversion.compute'

const conversionCatalog = assertCatalogShape('conversion', conversionJson)

export const conversionMessages: Record<WebsiteLocale, ConversionCopy> = {
	en: loadConversion(conversionCatalog.en),
	'pt-BR': loadConversion(conversionCatalog['pt-BR']),
	'zh-Hant': loadConversion(conversionCatalog['zh-Hant']),
	ja: loadConversion(conversionCatalog.ja),
	vi: loadConversion(conversionCatalog.vi),
}
