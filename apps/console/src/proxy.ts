import { NextRequest, NextResponse } from 'next/server'
import createMiddleware from 'next-intl/middleware'

import { isEntityId } from '@data/entity-id'
import { isLocale, localeCookie } from '@i18n/config'
import { routing } from '@i18n/routing'

import { createNonce, documentPolicy } from './proxy/content-security-policy'

const handleLocale = createMiddleware(routing)

function policyEnvironment() {
	if (process.env.NODE_ENV === 'development') return 'development'
	if (process.env.VERCEL_ENV === 'preview') return 'preview'

	return 'production'
}

// Every segment of a console URL, a locale, a fixed name or an entity id,
// keeps to the id alphabet. A path with any other segment names no screen.
function namesNoScreen(pathname: string) {
	return !pathname.split('/').filter(Boolean).every(isEntityId)
}

// Such a path is rewritten to one no route serves, so it renders the localized
// 404 with its status before a page can stream, next-intl never decodes it,
// and no extub read is built from its segments.
function missingScreen(request: NextRequest, headers: Headers) {
	const [prefix] = request.nextUrl.pathname.split('/').filter(Boolean)
	const locale = isLocale(prefix) ? prefix : routing.defaultLocale

	return NextResponse.rewrite(new URL(`/${locale}/404`, request.url), {
		request: { headers },
	})
}

// Every document carries one policy with a nonce of its own. next-intl copies
// the headers of the request it receives into its response, so the policy
// reaches Next, which stamps the nonce on its scripts.
export default function proxy(request: NextRequest) {
	const policy = documentPolicy(createNonce(), policyEnvironment())
	const headers = new Headers(request.headers)
	headers.set('Content-Security-Policy', policy)

	const response = namesNoScreen(request.nextUrl.pathname)
		? missingScreen(request, headers)
		: handleLocale(new NextRequest(request, { headers }))
	const preference = response.cookies.get(localeCookie)
	if (preference)
		response.cookies.set({
			...preference,
			httpOnly: true,
			secure: request.nextUrl.protocol === 'https:',
		})
	response.headers.set('Content-Security-Policy', policy)
	response.headers.set('Cache-Control', 'private, no-store')

	return response
}

// Every path but the API, Next's own files and the root convention files, so
// a dotted path under a locale, and the 404, carry the policy too.
export const config = {
	matcher: [
		'/((?!(?:api|_next)(?:/|$)|(?:favicon\\.ico|icon\\.svg|apple-icon\\.png|robots\\.txt)$).*)',
	],
}
