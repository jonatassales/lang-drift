import type { NextConfig } from 'next'
import createNextIntlPlugin from 'next-intl/plugin'

const withNextIntl = createNextIntlPlugin('./src/app/template/i18n/request.ts')

const nextConfig: NextConfig = {
	reactStrictMode: true,
	trailingSlash: true,
	allowedDevOrigins: ['127.0.0.1'],
	experimental: { serverActions: { bodySizeLimit: '4kb' } },
	transpilePackages: ['@langdrift/web'],
}

export default withNextIntl(nextConfig)
