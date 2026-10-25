import type { TurnstileApi } from './turnstile.types'

const scriptUrl =
	'https://challenges.cloudflare.com/turnstile/v0/api.js?render=explicit'
let loading: Promise<TurnstileApi> | undefined

// Shared readiness only. Widgets, attempts and tokens belong to each form.
export function loadTurnstile(): Promise<TurnstileApi> {
	if (window.turnstile) return Promise.resolve(window.turnstile)
	if (loading) return loading

	const pending = new Promise<TurnstileApi>((resolve, reject) => {
		const script = document.createElement('script')
		const timeout = window.setTimeout(fail, 15_000)

		function clean() {
			window.clearTimeout(timeout)
			script.onload = null
			script.onerror = null
		}

		function fail() {
			clean()
			script.remove()
			reject(new Error('turnstile-load-failed'))
		}

		script.src = scriptUrl
		script.async = true
		script.onload = () => {
			if (!window.turnstile) {
				fail()

				return
			}
			clean()
			resolve(window.turnstile)
		}
		script.onerror = fail
		document.head.append(script)
	})
	loading = pending
	void pending.catch(() => {
		if (loading === pending) loading = undefined
	})

	return pending
}
