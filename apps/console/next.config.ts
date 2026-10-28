import type { NextConfig } from 'next'
import createNextIntlPlugin from 'next-intl/plugin'

const withNextIntl = createNextIntlPlugin('./src/app/template/i18n/request.ts')

const config: NextConfig = {
	poweredByHeader: false,
	reactStrictMode: true,
	transpilePackages: ['@langdrift/web'],
	async headers() {
		return [
			{
				source: '/:path*',
				headers: [
					{
						key: 'X-Content-Type-Options',
						value: 'nosniff',
					},
					{
						key: 'Referrer-Policy',
						value: 'same-origin',
					},
					{
						key: 'X-Robots-Tag',
						value: 'noindex, nofollow',
					},
				],
			},
		]
	},
}

export default withNextIntl(config)
