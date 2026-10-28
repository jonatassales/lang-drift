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
