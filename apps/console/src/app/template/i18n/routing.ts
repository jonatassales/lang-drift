import { defineRouting } from 'next-intl/routing'

import { localeCookie, locales } from './config'
export const routing = defineRouting({
	locales,
	defaultLocale: 'en',
	localePrefix: 'always',
	localeCookie: {
		name: localeCookie,
		path: '/',
		maxAge: 31536000,
		sameSite: 'lax',
	},
	alternateLinks: false,
})
