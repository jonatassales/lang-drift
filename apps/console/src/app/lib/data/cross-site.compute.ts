// A browser call another site made: its Origin is not the console's own, or
// its fetch metadata says cross-site. The console refuses it.
export function crossSite(request: Request) {
	const url = new URL(request.url)
	const host = request.headers.get('host') ?? url.host
	const own = new URL(`${url.protocol}//${host}`)
	const origin = request.headers.get('origin')
	const foreignOrigin = Boolean(origin && origin !== own.origin)
	const crossSiteFetch = request.headers.get('sec-fetch-site') === 'cross-site'

	return foreignOrigin || crossSiteFetch
}
