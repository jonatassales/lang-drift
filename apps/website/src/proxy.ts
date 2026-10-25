import type { NextRequest } from 'next/server'

import { LocaleService } from './proxy/locale'

const localeService = new LocaleService()

export default function proxy(request: NextRequest) {
	return localeService.handle(request)
}

export const config = {
	matcher: [
		'/((?!api|_next|_vercel|docs-static|_pagefind|opengraph-image|metadata|.*\\..*).*)',
		'/([eE][nN]|[pP][tT]-[bB][rR]|[zZ][hH]-[hH][aA][nN][tT]|[jJ][aA]|[vV][iI]|[zZ][hH]-[cC][nN]|[zZ][hH]-[hH][aA][nN][sS]|[zZ][hH]-[sS][gG])/:path*',
	],
}
