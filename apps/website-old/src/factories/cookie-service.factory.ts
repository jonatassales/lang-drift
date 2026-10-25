import 'server-only'

import type { NextRequest, NextResponse } from 'next/server'

import {
	type CookieDefaults,
	type CookieDeleteAttributes,
	CookieService,
	type CookieStore,
	type CookieWriteAttributes,
} from '@services/cookie.service'
import { CryptoService } from '@services/crypto.service'

export function cookieServiceFactory(
	request: NextRequest,
	response?: NextResponse,
): CookieService {
	const secret = process.env.CRYPTO_SECRET
	if (!secret) {
		throw new Error('Missing required environment variable: CRYPTO_SECRET')
	}

	const defaults: CookieDefaults = {
		httpOnly: true,
		sameSite: 'lax',
		path: '/',
		secure: process.env.NODE_ENV === 'production',
	}

	return new CookieService(
		new CryptoService(secret),
		cookieStore(request, response),
		defaults,
	)
}

function cookieStore(
	request: NextRequest,
	response?: NextResponse,
): CookieStore {
	return {
		get(name) {
			return request.cookies.get(name)?.value
		},
		set(name, token, options) {
			const target = requireResponse(response)
			target.cookies.set(name, token, {
				...sharedOptions(options),
				maxAge: options.maxAge,
				expires: options.expires,
			})
		},
		delete(name, options) {
			const target = requireResponse(response)
			target.cookies.set(name, '', {
				...sharedOptions(options),
				maxAge: 0,
				expires: new Date(0),
			})
		},
	}
}

function requireResponse(response: NextResponse | undefined): NextResponse {
	if (!response) {
		throw new Error('Cookie writes require a response')
	}

	return response
}

function sharedOptions(
	options: CookieWriteAttributes | CookieDeleteAttributes,
) {
	return {
		httpOnly: true as const,
		sameSite: options.sameSite,
		path: options.path,
		secure: options.secure,
		...(options.domain ? { domain: options.domain } : {}),
	}
}
