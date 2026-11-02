import conversionJson from '@i18n/messages/conversion.json'
import type { ConversionCopy } from '@i18n/messages/messages.types'
import type { WebsiteLocale } from '@i18n/routing'

import { loadConversion } from './load-conversion.compute'

export const conversionMessages: Record<WebsiteLocale, ConversionCopy> = {
	en: loadConversion(conversionJson.en),
	'pt-BR': loadConversion(conversionJson['pt-BR']),
	'zh-Hant': loadConversion(conversionJson['zh-Hant']),
	ja: loadConversion(conversionJson.ja),
	vi: loadConversion(conversionJson.vi),
}
