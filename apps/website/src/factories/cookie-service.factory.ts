import 'server-only'

import { randomBytes } from 'node:crypto'
import type { NextRequest, NextResponse } from 'next/server'

import { type CookieDefaults, CookieService } from '@services/cookie.service'
import { CryptoService } from '@services/crypto.service'

import { cookieStore } from './cookie-store.compute'

let developmentSecret: string | undefined

export function cookieServiceFactory(
	request: NextRequest,
	response?: NextResponse,
): CookieService {
	let secret = process.env.CRYPTO_SECRET

	if (
		!secret &&
		(process.env.NODE_ENV === 'development' || process.env.NODE_ENV === 'test')
	) {
		developmentSecret ??= randomBytes(32).toString('base64')
		secret = developmentSecret
	}

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
