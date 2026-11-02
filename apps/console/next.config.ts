import type { NextConfig } from 'next'
import createNextIntlPlugin from 'next-intl/plugin'

const withNextIntl = createNextIntlPlugin('./src/app/template/i18n/request.ts')

// Documents take their policy from the proxy, with a nonce per response; the
// API takes this one, so no response carries two.
const apiPolicy = "default-src 'none'; frame-ancestors 'none'"

const config: NextConfig = {
	poweredByHeader: false,
	reactStrictMode: true,
	transpilePackages: ['@langdrift/web'],
	async redirects() {
		return [
			{
				source: '/sign-in',
				destination: '/login',
				permanent: false,
			},
			{
				source: '/:locale/sign-in',
				destination: '/:locale/login',
				permanent: false,
			},
			{
				source: '/signin',
				destination: '/login',
				permanent: false,
			},
			{
				source: '/:locale/signin',
				destination: '/:locale/login',
				permanent: false,
			},
		]
	},
	async headers() {
		return [
			{
				source: '/:path*',
				headers: [
					{
						key: 'Strict-Transport-Security',
						value: 'max-age=63072000; includeSubDomains; preload',
					},
					{
						key: 'X-Content-Type-Options',
						value: 'nosniff',
					},
					{
						key: 'Referrer-Policy',
						value: 'same-origin',
					},
					{
						key: 'X-Frame-Options',
						value: 'DENY',
					},
					{
						key: 'Permissions-Policy',
						value:
							'camera=(), microphone=(), geolocation=(), payment=(), usb=(), browsing-topics=()',
					},
					{
						key: 'Cross-Origin-Opener-Policy',
						value: 'same-origin',
					},
					{
						key: 'Cross-Origin-Resource-Policy',
						value: 'same-origin',
					},
					{
						key: 'X-Robots-Tag',
						value: 'noindex, nofollow',
					},
				],
			},
			{
				source: '/api/:path*',
				headers: [
					{
						key: 'Content-Security-Policy',
						value: apiPolicy,
					},
				],
			},
		]
	},
}

export default withNextIntl(config)
