import { addressedOrigin } from './addressed-origin.compute'

// A browser names the page behind a request in Origin and Sec-Fetch-Site.
// Book a demo has no configured public origin, so a page that names one must
// share the origin the request was addressed to. Browsers send Origin with
// every cross-origin POST, so a request without one is not a cross-origin
// browser request.
export function isForeignRequest(request: Request): boolean {
	if (request.headers.get('sec-fetch-site') === 'cross-site') return true
	const origin = request.headers.get('origin')
	if (origin === null) return false

	return origin !== addressedOrigin(request)
}
