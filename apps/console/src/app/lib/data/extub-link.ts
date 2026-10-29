// The rule that links the console to extub (spec 0015). Without a token the
// console serves loopback hosts only and reads an extub on plain http over
// loopback, as the local contract says. With a token, the hosted preview,
// the console serves its own host and may read an extub over https on any
// host, which the token admits. A URL never carries credentials.

const loopbackHosts = ['localhost', '127.0.0.1', '[::1]']

const localExtub = 'http://127.0.0.1:3401'

function hosted(token: string | undefined) {
	return Boolean(token?.trim())
}

export function consoleHostAllowed(
	hostname: string,
	token: string | undefined,
) {
	return hosted(token) || loopbackHosts.includes(hostname)
}

// The extub origin the console may read, or null when the configured URL
// breaks the rule. A URL that does not parse throws, as it always has, and
// the proxy reports extub as unavailable.
export function extubBase(raw: string | undefined, token: string | undefined) {
	const base = new URL(raw ?? localExtub)
	if (base.username || base.password) return null

	const onLoopback =
		base.protocol === 'http:' && loopbackHosts.includes(base.hostname)
	if (onLoopback) return base
	if (hosted(token) && base.protocol === 'https:') return base

	return null
}

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

// The request a server-rendered read hands the console's route boundary: the
// incoming host and the path, nothing else. The browser's Origin, fetch
// metadata and cookies describe its navigation, not this read, so a page
// reached from another site, or rendered over HTTPS, reads normally.
export function serverRead(host: string, path: string) {
	return new Request(new URL(`/api/console/${path}`, `http://${host}`))
}
