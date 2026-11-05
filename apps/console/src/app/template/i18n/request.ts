import { getRequestConfig } from 'next-intl/server'

import en from '@i18n/messages/en.json'
import ja from '@i18n/messages/ja.json'
import pt from '@i18n/messages/pt-br.json'
import zh from '@i18n/messages/zh-hant.json'

import { isLocale } from './config'

const catalogs = {
	en,
	'pt-br': pt,
	'zh-hant': zh,
	ja,
}

export default getRequestConfig(async ({ requestLocale }) => {
	const requested = await requestLocale

	const locale = isLocale(requested) ? requested : 'en'

	return {
		locale,
		timeZone: 'UTC',
		messages: catalogs[locale],
	}
})
