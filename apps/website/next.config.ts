import type { NextConfig } from 'next'
import createNextIntlPlugin from 'next-intl/plugin'

const withNextIntl = createNextIntlPlugin('./src/app/template/i18n/request.ts')

const nextConfig: NextConfig = {
	reactStrictMode: true,
	trailingSlash: true,
	redirects: async () => [
		{
			source: '/docs/:path*',
			destination: '/en/docs/:path*/',
			permanent: true,
		},
	],
	async rewrites() {
		try {
			const configured = process.env.DOCS_ORIGIN
			if (!configured) throw new Error('DOCS_ORIGIN is required')
			const origin = new URL(configured)
			const local =
				process.env.LANGDRIFT_LOCAL_BUILD === '1' ||
				process.env.NODE_ENV === 'development'
			const localOrigin =
				local && ['localhost', '127.0.0.1'].includes(origin.hostname)
			if (
				origin.protocol !== 'https:' &&
				!(localOrigin && origin.protocol === 'http:')
			)
				throw new Error('DOCS_ORIGIN must use HTTPS')
			if (
				origin.username ||
				origin.password ||
				origin.pathname !== '/' ||
				origin.search ||
				origin.hash
			) {
				throw new Error(
					'DOCS_ORIGIN must be an origin without credentials or a path',
				)
			}
			const redirectsBackToWebsite =
				['langdrift.md', 'www.langdrift.md', 'docs.langdrift.md'].includes(
					origin.hostname,
				) ||
				origin.host === process.env.VERCEL_URL ||
				origin.host === process.env.VERCEL_PROJECT_PRODUCTION_URL
			if (redirectsBackToWebsite) {
				throw new Error('DOCS_ORIGIN must be a dedicated upstream deployment')
			}
			return {
				beforeFiles: [
					{
						source: '/:locale(en|pt-br|zh-hant|ja|vi)/docs/:path(.*)',
						destination: `${origin.origin}/:locale/docs/:path`,
					},
					{
						source: '/docs-static/:path*',
						destination: `${origin.origin}/docs-static/:path*`,
					},
					{
						source: '/_pagefind/:path*',
						destination: `${origin.origin}/_pagefind/:path*`,
					},
				],
			}
		} catch {
			throw new Error(
				'Configure a valid server-side DOCS_ORIGIN before serving the Docs zone',
			)
		}
	},
	// Preview URLs are not authorized hostnames on the production Turnstile widget.
	env: {
		WEBSITE_PREVIEW_SIGNUP_DISABLED:
			process.env.VERCEL_ENV === 'preview' ||
			process.env.WEBSITE_PREVIEW_SIGNUP_DISABLED !== 'false'
				? 'true'
				: 'false',
	},
	allowedDevOrigins: ['127.0.0.1'],
	experimental: { serverActions: { bodySizeLimit: '4kb' } },
	transpilePackages: ['@pack/react'],
}

export default withNextIntl(nextConfig)
