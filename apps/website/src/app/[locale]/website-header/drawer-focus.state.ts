import type { WebsiteLocale } from '@i18n/routing'

interface DrawerFocus {
	locale: WebsiteLocale | undefined
}

// Only set by a browser drawer interaction; survives the locale route remount
// until focus returns to the replacement navigation trigger.
export const drawerFocus: DrawerFocus = {
	locale: undefined,
}
