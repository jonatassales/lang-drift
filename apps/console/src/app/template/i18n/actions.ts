'use server'
import { cookies, headers } from 'next/headers'

import { isLocale, localeCookie } from './config'
export async function setLocale(locale: string) {
	if (!isLocale(locale)) return { ok: false }
	const store = await cookies()
	const requestHeaders = await headers()

	store.set(localeCookie, locale, {
		path: '/',
		maxAge: 31536000,
		httpOnly: true,
		sameSite: 'lax',
		secure: requestHeaders.get('x-forwarded-proto') === 'https',
	})

	return { ok: true }
}
