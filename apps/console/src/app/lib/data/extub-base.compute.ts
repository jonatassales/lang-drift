import { hosted } from './hosted.compute.ts'

const loopbackHosts = ['localhost', '127.0.0.1', '[::1]']

// The extub origin the console may read, or null when the configured URL
// breaks the rule. A URL that does not parse throws, as it always has, and
// the proxy reports extub as unavailable.
export function extubBase(raw: string, token: string | undefined) {
	const base = new URL(raw)
	if (base.username || base.password) return null

	const onLoopback =
		base.protocol === 'http:' && loopbackHosts.includes(base.hostname)
	if (onLoopback) return base
	if (hosted(token) && base.protocol === 'https:') return base

	return null
}
