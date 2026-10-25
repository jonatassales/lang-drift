import {
	getRedirectUrl,
	getRewrittenUrl,
	unstable_getResponseFromNextConfig,
} from 'next/experimental/testing/server'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import nextConfig from './next.config'

beforeEach(() => {
	vi.stubEnv('DOCS_ORIGIN', 'https://docs.example')
})

afterEach(() => vi.unstubAllEnvs())

describe('Website Docs forwarding', () => {
	it.each(['en', 'pt-br', 'zh-hant', 'ja', 'vi'])(
		'forwards the %s landing page with exactly one trailing slash',
		async (locale) => {
			const response = await unstable_getResponseFromNextConfig({
				url: `https://langdrift.md/${locale}/docs/?from=website`,
				nextConfig,
			})

			expect(getRedirectUrl(response)).toBeNull()
			expect(getRewrittenUrl(response)).toBe(
				`https://docs.example/${locale}/docs/?from=website`,
			)
		},
	)

	it.each(['en', 'pt-br', 'zh-hant', 'ja', 'vi'])(
		'preserves the %s article path, trailing slash and query',
		async (locale) => {
			const response = await unstable_getResponseFromNextConfig({
				url: `https://langdrift.md/${locale}/docs/sdk/references/?from=website`,
				nextConfig,
			})

			expect(getRedirectUrl(response)).toBeNull()
			expect(getRewrittenUrl(response)).toBe(
				`https://docs.example/${locale}/docs/sdk/references/?from=website`,
			)
		},
	)

	it('preserves a dotted path and query without an upstream slash redirect', async () => {
		const response = await unstable_getResponseFromNextConfig({
			url: 'https://langdrift.md/ja/docs/missing.part?from=website',
			nextConfig,
		})

		expect(getRedirectUrl(response)).toBeNull()
		expect(getRewrittenUrl(response)).toBe(
			'https://docs.example/ja/docs/missing.part?from=website',
		)
	})

	it.each([
		'/docs-static/_next/static/chunks/app.js',
		'/_pagefind/pagefind-entry.json',
	])('forwards the Docs asset %s without altering its path', async (path) => {
		const response = await unstable_getResponseFromNextConfig({
			url: `https://langdrift.md${path}`,
			nextConfig,
		})

		expect(getRedirectUrl(response)).toBeNull()
		expect(getRewrittenUrl(response)).toBe(`https://docs.example${path}`)
	})
})
