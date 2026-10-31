// The console's document policy (spec 0018). A script runs only when it
// carries this response's nonce, or when a script that does loads it
// ('strict-dynamic'). Styles keep 'unsafe-inline': React style props in the
// shared primitives and sinapsi's shadow-root style cannot carry a nonce.

export type PolicyEnvironment = 'development' | 'preview' | 'production'

// A nonce for one response: 16 random bytes, base64.
export function createNonce() {
	const bytes = crypto.getRandomValues(new Uint8Array(16))

	return btoa(String.fromCharCode(...bytes))
}

export function documentPolicy(nonce: string, environment: PolicyEnvironment) {
	// React needs eval to rebuild server error stacks, and the development
	// server pushes updates over a websocket.
	const evaluation = environment === 'development' ? ["'unsafe-eval'"] : []
	const socket = environment === 'development' ? ['ws:'] : []
	// Preview deployments load the Vercel Toolbar.
	const toolbar = environment === 'preview' ? ['https://vercel.live'] : []
	// Vercel Analytics loads its script from this origin in production and
	// from va.vercel-scripts.com in development; the hosts match the website's
	// policy (spec 0022).
	const analyticsScript = 'https://va.vercel-scripts.com'
	const analyticsVitals = 'https://vitals.vercel-insights.com'
	const directives: Record<string, string[]> = {
		'default-src': ["'self'"],
		'script-src': [
			"'self'",
			`'nonce-${nonce}'`,
			"'strict-dynamic'",
			analyticsScript,
			...evaluation,
			...toolbar,
		],
		'style-src': ["'self'", "'unsafe-inline'", ...toolbar],
		'img-src': ["'self'", 'data:', 'blob:', ...toolbar],
		'font-src': ["'self'", ...toolbar],
		'connect-src': [
			"'self'",
			analyticsScript,
			analyticsVitals,
			...socket,
			...toolbar,
		],
		'worker-src': ["'self'", 'blob:'],
		'manifest-src': ["'self'"],
		'frame-src': toolbar.length ? toolbar : ["'none'"],
		'object-src': ["'none'"],
		'base-uri': ["'self'"],
		'form-action': ["'self'"],
		'frame-ancestors': ["'none'"],
	}

	return Object.entries(directives)
		.map(([directive, sources]) => [directive, ...sources].join(' '))
		.join('; ')
}
