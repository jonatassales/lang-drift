import type { EarlyAccessMode } from './early-access-mode.compute'

export function assertCredentialFreeOrigin(
	endpoint: URL,
	mode: EarlyAccessMode,
): void {
	if (
		endpoint.username ||
		endpoint.password ||
		endpoint.search ||
		endpoint.hash ||
		endpoint.pathname !== '/' ||
		!['http:', 'https:'].includes(endpoint.protocol)
	) {
		throw new Error('Early Access endpoints must be credential-free origins')
	}
	const requiresHttps = mode === 'production' && endpoint.protocol !== 'https:'
	if (requiresHttps) {
		throw new Error('Production Early Access endpoints require HTTPS')
	}
	if (
		mode !== 'production' &&
		!['localhost', '127.0.0.1', 'email-store'].includes(endpoint.hostname)
	) {
		throw new Error('Development Early Access endpoints must be local')
	}
}
