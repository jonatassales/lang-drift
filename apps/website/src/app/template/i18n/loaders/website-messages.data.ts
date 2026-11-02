import headerJson from '@i18n/messages/header.json'
import type { HeaderMessages } from '@i18n/messages/messages.types'
import routeStateJson from '@i18n/messages/route-state.json'
import type { WebsiteLocale } from '@i18n/routing'

const headerMessages: Record<WebsiteLocale, HeaderMessages> = headerJson

export function getWebsiteMessages(locale: WebsiteLocale) {
	return {
		header: headerMessages[locale],
		routeState: routeStateJson[locale],
	}
}
