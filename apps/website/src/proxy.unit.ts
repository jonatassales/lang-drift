import { unstable_doesMiddlewareMatch } from 'next/experimental/testing/server'
import { NextRequest } from 'next/server'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import { pageMetadata } from '@app/metadata/page-metadata'
import pages from '@app/metadata/pages.json'
import { routing } from '@i18n/routing'

import proxy, { config } from './proxy'

vi.mock('server-only', () => ({}))

const pageIds = Object.keys(pages) as Array<keyof typeof pages>

const localizedPages = pageIds.flatMap((page) =>
	routing.locales.map((locale) => {
		const canonical = pageMetadata(page, locale).alternates?.canonical

		return {
			locale,
			pathname: new URL(String(canonical)).pathname,
		}
	}),
)

beforeEach(() => {
	vi.stubEnv('CRYPTO_SECRET', Buffer.alloc(32, 1).toString('base64'))
})

afterEach(() => vi.unstubAllEnvs())

describe('Website locale routing', () => {
	it.each(['en', 'pt-br', 'zh-hant', 'ja', 'vi'])(
		'serves the explicit %s prefix without redirecting to itself',
		(prefix) => {
			const url = `https://langdrift.test/${prefix}/`
			const response = proxy(new NextRequest(url))
			expect(response.headers.get('location')).not.toBe(url)
			expect(response.status).toBe(200)
		},
	)

	it.each(['/ja/missing.mdx', '/pt-BR/missing.html', '/zh-CN/missing.file'])(
		'routes dotted page %s through locale negotiation',
		(url) => {
			expect(
				unstable_doesMiddlewareMatch({
					config,
					nextConfig: {},
					url,
				}),
			).toBe(true)
		},
	)

	// Next serves an image file of the [locale] segment once, at /-/<file>.
	it.each(['/-/opengraph-image.png', '/-/twitter-image.png'])(
		'serves social image %s without a cookie or no-store',
		(url) => {
			expect(
				unstable_doesMiddlewareMatch({
					config,
					nextConfig: {},
					url,
				}),
			).toBe(false)
		},
	)

	it.each(['/api/health', '/_next/static/app.js', '/favicon.ico'])(
		'leaves infrastructure request %s outside locale negotiation',
		(url) => {
			expect(
				unstable_doesMiddlewareMatch({
					config,
					nextConfig: {},
					url,
				}),
			).toBe(false)
		},
	)

	it.each(localizedPages)(
		'leaves hreflang for $pathname to the HTML and the sitemap',
		({ pathname }) => {
			const response = proxy(new NextRequest(`https://langdrift.md${pathname}`))

			expect(response.headers.get('link')).toBeNull()
		},
	)

	it('canonicalizes casing without discarding the search string', () => {
		const response = proxy(
			new NextRequest('https://langdrift.test/pt-BR/pricing/?from=menu'),
		)
		expect(response.status).toBe(308)
		expect(response.headers.get('location')).toBe(
			'https://langdrift.test/pt-br/pricing/?from=menu',
		)
	})

	it('uses English and an availability notice for Simplified Chinese', () => {
		const response = proxy(
			new NextRequest('https://langdrift.test/zh-CN/missing.file?from=menu'),
		)
		expect(response.headers.get('location')).toBe(
			'https://langdrift.test/en/missing.file?from=menu&availability=zh-CN',
		)
	})

	it('prefers the saved encrypted cookie to the browser language and the URL to both', () => {
		const first = proxy(new NextRequest('https://langdrift.test/vi/'))
		const cookie = first.cookies.get('langdrift-website-locale')
		expect(cookie?.httpOnly).toBe(true)
		const headers = {
			cookie: `langdrift-website-locale=${cookie?.value}`,
			'accept-language': 'ja',
		}
		const saved = proxy(
			new NextRequest('https://langdrift.test/pricing/?from=menu', { headers }),
		)
		expect(saved.headers.get('location')).toBe(
			'https://langdrift.test/vi/pricing/?from=menu',
		)
		const explicit = proxy(
			new NextRequest('https://langdrift.test/pt-br/', { headers }),
		)
		expect(explicit.headers.get('x-middleware-rewrite')).toContain('/pt-BR')
	})

	it('uses a supported browser preference before the English default', () => {
		const response = proxy(
			new NextRequest('https://langdrift.test/', {
				headers: { 'accept-language': 'vi-VN,ja;q=0.8' },
			}),
		)
		expect(response.headers.get('location')).toBe('https://langdrift.test/vi/')
	})

	it('never silently substitutes Traditional Chinese for a Simplified Chinese browser preference', () => {
		const response = proxy(
			new NextRequest('https://langdrift.test/pricing/?from=browser', {
				headers: { 'accept-language': 'zh-CN,zh;q=0.9,en;q=0.8' },
			}),
		)
		expect(response.headers.get('location')).toBe(
			'https://langdrift.test/en/pricing/?from=browser&availability=zh-CN',
		)
	})
})
