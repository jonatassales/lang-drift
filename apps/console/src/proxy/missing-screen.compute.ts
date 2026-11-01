import { type NextRequest, NextResponse } from 'next/server'

import { isLocale } from '@i18n/config'
import { routing } from '@i18n/routing'

export function missingScreen(request: NextRequest) {
	const [prefix] = request.nextUrl.pathname.split('/').filter(Boolean)
	const locale = isLocale(prefix) ? prefix : routing.defaultLocale

	return NextResponse.rewrite(new URL(`/${locale}/404`, request.url), {
		request: { headers: request.headers },
	})
}
