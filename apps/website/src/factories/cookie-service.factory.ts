import 'server-only'
import type { NextRequest, NextResponse } from 'next/server'

import { type CookieDefaults, CookieService } from '@services/cookie.service'
import { CryptoService } from '@services/crypto.service'

import { cookieStore } from './cookie-store.compute'

// The cookie service the locale proxy encrypts its cookies with. Without
// CRYPTO_SECRET it fails in every environment rather than invent a key.
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
