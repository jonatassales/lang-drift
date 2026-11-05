import { NextRequest, NextResponse } from 'next/server'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import { cookieServiceFactory } from './cookie-service.factory'

vi.mock('server-only', () => ({}))

const secret = Buffer.alloc(32, 7).toString('base64')

describe('website/factories/cookie-service.factory', () => {
	beforeEach(() => {
		vi.unstubAllEnvs()
	})

	afterEach(() => {
		vi.unstubAllEnvs()
	})

	it.each(['development', 'test', 'production'])(
		'throws in %s when CRYPTO_SECRET is unset, inventing no key',
		(environment) => {
			vi.stubEnv('NODE_ENV', environment)
			vi.stubEnv('CRYPTO_SECRET', '')

			const request = new NextRequest('https://langdrift.md/en/')

			expect(() => cookieServiceFactory(request)).toThrow(
				'Missing required environment variable: CRYPTO_SECRET',
			)
		},
	)

	it('reads back a cookie it wrote with the CRYPTO_SECRET of the environment', () => {
		vi.stubEnv('CRYPTO_SECRET', secret)

		const writerRequest = new NextRequest('https://langdrift.md/en/')
		const response = NextResponse.next()
		const writer = cookieServiceFactory(writerRequest, response)

		writer.set('shared-key', 'persisted-value', {
			ttlSeconds: 3600,
		})

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
			(value): value is string => typeof value === 'string',
		)

		expect(readResult).toEqual({
			found: true,
			value: 'persisted-value',
		})
	})
})
