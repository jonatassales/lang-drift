import { getRequestConfig } from 'next-intl/server'

import { isLocale } from './config'
import en from './messages/en.json'
import ja from './messages/ja.json'
import pt from './messages/pt-br.json'
import zh from './messages/zh-hant.json'

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
