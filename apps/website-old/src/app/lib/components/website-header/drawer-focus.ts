import type { WebsiteLocale } from '@i18n/routing'

// Only set by a browser drawer interaction; survives the locale route remount
// until focus returns to the replacement navigation trigger.
export const drawerFocus: { locale: WebsiteLocale | undefined } = {
	locale: undefined,
}
