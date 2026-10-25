import type { NextConfig } from 'next'
import createNextIntlPlugin from 'next-intl/plugin'

const withNextIntl = createNextIntlPlugin('./src/app/template/i18n/request.ts')

const nextConfig: NextConfig = {
	reactStrictMode: true,
	trailingSlash: true,
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
