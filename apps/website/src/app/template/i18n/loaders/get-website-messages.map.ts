import routeStateJson from '@i18n/messages/route-state.json'
import type { WebsiteLocale } from '@i18n/routing'

import { assertCatalogShape } from './assert-catalog-shape.compute'
import { headerMessages } from './catalogs.data'

const routeStateMessages = assertCatalogShape('route-state', routeStateJson)

export function getWebsiteMessages(locale: WebsiteLocale) {
	return {
		header: headerMessages[locale],
		routeState: routeStateMessages[locale],
	}
}
