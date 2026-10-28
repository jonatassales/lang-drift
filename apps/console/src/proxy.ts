import type { NextRequest } from 'next/server'
import createMiddleware from 'next-intl/middleware'

import { localeCookie } from '@i18n/config'
import { routing } from '@i18n/routing'

const handleLocale = createMiddleware(routing)
export default function proxy(request: NextRequest) {
	const response = handleLocale(request)
	const preference = response.cookies.get(localeCookie)
	if (preference)
		response.cookies.set({
			...preference,
			httpOnly: true,
			secure: request.nextUrl.protocol === 'https:',
		})
	response.headers.set('Cache-Control', 'private, no-store')

	return response
}
export const config = { matcher: ['/((?!api|_next|.*\\..*).*)'] }
