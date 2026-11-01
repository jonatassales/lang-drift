import type { WebsiteLocale } from '@i18n/routing'

import headerJson from './header.json'
import type { HeaderMessages } from './messages.types'
import routeStateJson from './route-state.json'

const headerMessages: Record<WebsiteLocale, HeaderMessages> = headerJson

export function getWebsiteMessages(locale: WebsiteLocale) {
	return {
		header: headerMessages[locale],
		routeState: routeStateJson[locale],
	}
}
