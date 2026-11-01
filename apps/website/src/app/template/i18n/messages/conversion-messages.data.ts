import type { WebsiteLocale } from '@i18n/routing'

import conversionJson from './conversion.json'
import { loadConversion } from './load-conversion.compute'
import type { ConversionCopy } from './messages.types'

export const conversionMessages: Record<WebsiteLocale, ConversionCopy> = {
	en: loadConversion(conversionJson.en),
	'pt-BR': loadConversion(conversionJson['pt-BR']),
	'zh-Hant': loadConversion(conversionJson['zh-Hant']),
	ja: loadConversion(conversionJson.ja),
	vi: loadConversion(conversionJson.vi),
}
