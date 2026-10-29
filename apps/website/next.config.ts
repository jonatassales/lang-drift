import type { NextConfig } from 'next'
import createNextIntlPlugin from 'next-intl/plugin'

const withNextIntl = createNextIntlPlugin('./src/app/template/i18n/request.ts')

// One static policy keeps every page prerendered: nonces would render each
// request on demand. Next inlines its page payload as scripts, and the site
// renders no user-supplied HTML, so inline scripts are allowed. Inline styles
// cover React style props and the editor's and the orb's injected styles.
function contentSecurityPolicy(): string {
	const directives = {
		'default-src': ["'self'"],
		'script-src': [
			"'self'",
			"'unsafe-inline'",
			'https://va.vercel-scripts.com',
		],
		'style-src': ["'self'", "'unsafe-inline'"],
		'img-src': ["'self'", 'data:', 'blob:'],
		'font-src': ["'self'", 'data:'],
		'connect-src': [
			"'self'",
			'https://va.vercel-scripts.com',
			'https://vitals.vercel-insights.com',
		],
		'worker-src': ["'self'", 'blob:'],
		'manifest-src': ["'self'"],
		'frame-src': ["'none'"],
		'object-src': ["'none'"],
		'base-uri': ["'self'"],
		'form-action': ["'self'"],
		'frame-ancestors': ["'none'"],
	}

	// React rebuilds server error stacks with eval in development, and the
	// development server reloads over a websocket.
	if (process.env.NODE_ENV === 'development') {
		directives['script-src'].push("'unsafe-eval'")
		directives['connect-src'].push('ws:')
	}

	// The Vercel Toolbar on preview deployments; production never allows it.
	if (process.env.VERCEL_ENV === 'preview') {
		const toolbar = 'https://vercel.live'
		directives['script-src'].push(toolbar)
		directives['style-src'].push(toolbar)
		directives['img-src'].push(toolbar)
		directives['font-src'].push(toolbar)
		directives['connect-src'].push(toolbar)
		directives['frame-src'] = [toolbar]
	}

	return Object.entries(directives)
		.map(([directive, sources]) => [directive, ...sources].join(' '))
		.join('; ')
}

const securityHeaders = {
	'Strict-Transport-Security': 'max-age=63072000; includeSubDomains; preload',
	'X-Content-Type-Options': 'nosniff',
	'Referrer-Policy': 'strict-origin-when-cross-origin',
	'X-Frame-Options': 'DENY',
	'Permissions-Policy':
		'camera=(), microphone=(), geolocation=(), payment=(), usb=(), browsing-topics=()',
	'Cross-Origin-Opener-Policy': 'same-origin',
}

const nextConfig: NextConfig = {
	reactStrictMode: true,
	trailingSlash: true,
	poweredByHeader: false,
	allowedDevOrigins: ['127.0.0.1'],
	experimental: { serverActions: { bodySizeLimit: '4kb' } },
	transpilePackages: ['@langdrift/web'],
	// The only source of the security headers: the proxy sets none, so every
	// path, the proxy's own redirects included, carries the same set.
	async headers() {
		const headers = {
			'Content-Security-Policy': contentSecurityPolicy(),
			...securityHeaders,
		}

		return [
			{
				source: '/:path*',
				headers: Object.entries(headers).map(([key, value]) => ({
					key,
					value,
				})),
			},
		]
	},
}

export default withNextIntl(nextConfig)
