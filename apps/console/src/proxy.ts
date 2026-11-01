import type { NextRequest } from 'next/server'
import createMiddleware from 'next-intl/middleware'

import { localeCookie } from '@i18n/config'
import { routing } from '@i18n/routing'

import { createNonce, documentPolicy } from './proxy/content-security-policy'
import { missingScreen } from './proxy/missing-screen.compute'
import { namesNoScreen } from './proxy/names-no-screen.compute'
import { policyEnvironment } from './proxy/policy-environment.compute'

const handleLocale = createMiddleware(routing)

// Every document carries one policy with a nonce of its own. The policy is
// set on the request itself, which keeps Next's routing context; next-intl
// copies the request's headers into its response, so the policy reaches
// Next, which stamps the nonce on its scripts.
export default function proxy(request: NextRequest) {
	const policy = documentPolicy(createNonce(), policyEnvironment())
	request.headers.set('Content-Security-Policy', policy)

	const response = namesNoScreen(request.nextUrl.pathname)
		? missingScreen(request)
		: handleLocale(request)
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
