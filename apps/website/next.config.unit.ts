import { unstable_getResponseFromNextConfig } from 'next/experimental/testing/server'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import nextConfig from './next.config'

// The policy production serves, directive by directive (spec 0018).
const productionPolicy = {
	'default-src': ["'self'"],
	'script-src': ["'self'", "'unsafe-inline'"],
	'style-src': ["'self'", "'unsafe-inline'"],
	'img-src': ["'self'", 'data:', 'blob:'],
	'font-src': ["'self'", 'data:'],
	'connect-src': ["'self'"],
	'worker-src': ["'self'", 'blob:'],
	'manifest-src': ["'self'"],
	'frame-src': ["'none'"],
	'object-src': ["'none'"],
	'base-uri': ["'self'"],
	'form-action': ["'self'"],
	'frame-ancestors': ["'none'"],
}

const otherHeaders = {
	'Strict-Transport-Security': 'max-age=63072000; includeSubDomains; preload',
	'X-Content-Type-Options': 'nosniff',
	'Referrer-Policy': 'strict-origin-when-cross-origin',
	'X-Frame-Options': 'DENY',
	'Permissions-Policy':
		'camera=(), microphone=(), geolocation=(), payment=(), usb=(), browsing-topics=()',
	'Cross-Origin-Opener-Policy': 'same-origin',
}

async function headersAt(route: string) {
	const response = await unstable_getResponseFromNextConfig({
		url: `https://langdrift.md${route}`,
		nextConfig,
	})

	return response.headers
}

async function policyAt(route: string) {
	const headers = await headersAt(route)
	const policy = headers.get('Content-Security-Policy') ?? ''

	return Object.fromEntries(
		policy.split('; ').map((directive) => {
			const [name, ...sources] = directive.split(' ')

			return [name, sources]
		}),
	)
}

beforeEach(() => {
	vi.stubEnv('NODE_ENV', 'production')
	vi.stubEnv('VERCEL_ENV', undefined)
})

afterEach(() => vi.unstubAllEnvs())

describe('the website security headers', () => {
	it.each(['/en/', '/api/book-a-demo/', '/robots.txt'])(
		'answers %s with every header',
		async (route) => {
			const headers = await headersAt(route)

			expect(await policyAt(route)).toEqual(productionPolicy)
			for (const [name, value] of Object.entries(otherHeaders)) {
				expect(headers.get(name)).toBe(value)
			}
		},
	)

	it('never upgrades insecure requests, so plain-HTTP loopback runs work', async () => {
		const headers = await headersAt('/en/')

		expect(headers.get('Content-Security-Policy')).not.toContain(
			'upgrade-insecure-requests',
		)
	})

	it('adds nothing on a production deployment', async () => {
		vi.stubEnv('VERCEL_ENV', 'production')

		expect(await policyAt('/en/')).toEqual(productionPolicy)
	})

	it('lets the development server evaluate and reload, and nothing more', async () => {
		vi.stubEnv('NODE_ENV', 'development')

		expect(await policyAt('/en/')).toEqual({
			...productionPolicy,
			'script-src': ["'self'", "'unsafe-inline'", "'unsafe-eval'"],
			'connect-src': ["'self'", 'ws:'],
		})
	})

	it('lets a preview deployment load the Vercel Toolbar, and nothing more', async () => {
		vi.stubEnv('VERCEL_ENV', 'preview')

		expect(await policyAt('/en/')).toEqual({
			...productionPolicy,
			'script-src': ["'self'", "'unsafe-inline'", 'https://vercel.live'],
			'style-src': ["'self'", "'unsafe-inline'", 'https://vercel.live'],
			'img-src': ["'self'", 'data:', 'blob:', 'https://vercel.live'],
			'font-src': ["'self'", 'data:', 'https://vercel.live'],
			'connect-src': ["'self'", 'https://vercel.live'],
			'frame-src': ['https://vercel.live'],
		})
	})

	it('never announces the framework', () => {
		expect(nextConfig.poweredByHeader).toBe(false)
	})
})
