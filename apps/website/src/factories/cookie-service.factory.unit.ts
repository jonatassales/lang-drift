import { NextRequest, NextResponse } from 'next/server'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import { cookieServiceFactory } from './cookie-service.factory'

vi.mock('server-only', () => ({}))

describe('website/src/factories/cookie-service.factory', () => {
	beforeEach(() => {
		vi.unstubAllEnvs()
	})

	afterEach(() => {
		vi.unstubAllEnvs()
	})

	it('generates secret on the fly in development when CRYPTO_SECRET is unset', () => {
		vi.stubEnv('NODE_ENV', 'development')
		vi.stubEnv('CRYPTO_SECRET', '')

		const request = new NextRequest('https://langdrift.md/en/')
		const response = NextResponse.next()

		const service = cookieServiceFactory(request, response)
		expect(service).toBeDefined()

		service.set('test-cookie', 'hello-world', { ttlSeconds: 3600 })
		const setCookie = response.headers.get('set-cookie')
		expect(setCookie).toContain('test-cookie=')
	})

	it('reuses generated development secret across multiple factory calls', () => {
		vi.stubEnv('NODE_ENV', 'development')
		vi.stubEnv('CRYPTO_SECRET', '')

		const request = new NextRequest('https://langdrift.md/en/')
		const response = NextResponse.next()

		const writer = cookieServiceFactory(request, response)
		writer.set('shared-key', 'persisted-value', { ttlSeconds: 3600 })

		const rawCookie = response.cookies.get('shared-key')?.value
		expect(rawCookie).toBeDefined()

		const readerRequest = new NextRequest('https://langdrift.md/en/', {
			headers: {
				cookie: `shared-key=${rawCookie}`,
			},
		})

		const reader = cookieServiceFactory(readerRequest)
		const readResult = reader.get(
			'shared-key',
			(val): val is string => typeof val === 'string',
		)

		expect(readResult).toEqual({
			found: true,
			value: 'persisted-value',
		})
	})

	it('uses explicit CRYPTO_SECRET when provided in environment', () => {
		const customSecret = Buffer.alloc(32, 7).toString('base64')
		vi.stubEnv('CRYPTO_SECRET', customSecret)

		const request = new NextRequest('https://langdrift.md/en/')
		const service = cookieServiceFactory(request)

		expect(service).toBeDefined()
	})

	it('throws in production when CRYPTO_SECRET is unset', () => {
		vi.stubEnv('NODE_ENV', 'production')
		vi.stubEnv('CRYPTO_SECRET', '')

		const request = new NextRequest('https://langdrift.md/en/')

		expect(() => cookieServiceFactory(request)).toThrow(
			'Missing required environment variable: CRYPTO_SECRET',
		)
	})
})
