import type { NextRequest, NextResponse } from 'next/server'

import type { CookieStore } from '@services/cookie.service'

export function cookieStore(
	request: NextRequest,
	response?: NextResponse,
): CookieStore {
	return {
		get(name) {
			return request.cookies.get(name)?.value
		},
		set(name, token, options) {
			if (!response) {
				throw new Error('Cookie writes require a response')
			}
			response.cookies.set(name, token, {
				httpOnly: true,
				sameSite: options.sameSite,
				path: options.path,
				secure: options.secure,
				...(options.domain ? { domain: options.domain } : {}),
				maxAge: options.maxAge,
				expires: options.expires,
			})
		},
		delete(name, options) {
			if (!response) {
				throw new Error('Cookie writes require a response')
			}
			response.cookies.set(name, '', {
				httpOnly: true,
				sameSite: options.sameSite,
				path: options.path,
				secure: options.secure,
				...(options.domain ? { domain: options.domain } : {}),
				maxAge: 0,
				expires: new Date(0),
			})
		},
	}
}
