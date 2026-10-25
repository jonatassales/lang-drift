import { loadTurnstile } from './turnstile.loader'
import type { SubmissionChallenge, TurnstileApi } from './turnstile.types'

type Phase =
	| 'idle'
	| 'loading'
	| 'ready'
	| 'verifying'
	| 'verified'
	| 'consumed'
	| 'failed'
	| 'disposed'

export interface TurnstileController extends SubmissionChallenge {
	prepare: () => Promise<void>
	reset: () => void
	dispose: () => void
}

export function createTurnstileController(
	element: HTMLElement,
	sitekey: string,
): TurnstileController {
	let phase: Phase = 'idle'
	let generation = 0
	let api: TurnstileApi | undefined
	let widget: string | undefined
	let preparation: Promise<void> | undefined
	let token = ''
	let timer: ReturnType<typeof setTimeout> | undefined
	let pending:
		| { resolve: () => void; reject: (error: Error) => void }
		| undefined

	function fail() {
		token = ''
		phase = 'failed'
		clearTimeout(timer)
		pending?.reject(new Error('turnstile-challenge-failed'))
		pending = undefined
	}

	function reset() {
		generation += 1
		clearTimeout(timer)
		pending?.reject(new Error('turnstile-cancelled'))
		pending = undefined
		token = ''
		preparation = undefined
		const previousWidget = widget
		widget = undefined
		if (phase !== 'disposed') phase = 'idle'
		if (previousWidget === undefined) return

		// Cleanup must never hold the submission lock after a provider exception.
		try {
			api?.reset(previousWidget)
		} catch {
			// Removal below still runs if the provider cannot reset.
		}
		try {
			api?.remove(previousWidget)
		} catch {
			element.replaceChildren()
		}
	}

	function prepare(): Promise<void> {
		if (phase === 'disposed')
			return Promise.reject(new Error('turnstile-cancelled'))
		if (phase === 'failed') reset()
		if (phase === 'ready') return Promise.resolve()
		if (preparation) return preparation

		phase = 'loading'
		const current = generation
		preparation = (async () => {
			const loaded = await loadTurnstile()

			if (current !== generation) {
				throw new Error('turnstile-cancelled')
			}
			api = loaded
			const active = () => current === generation
			const failed = () => phase === 'failed'
			const onFailure = () => {
				if (active() && phase !== 'consumed') fail()
			}
			widget = loaded.render(element, {
				sitekey,
				action: 'early-access',
				execution: 'execute',
				'response-field': false,
				retry: 'never',
				'refresh-expired': 'never',
				callback: (response) => {
					if (!active() || phase !== 'verifying' || !pending) return
					if (!response || response.length > 2048) {
						fail()

						return
					}
					token = response
					phase = 'verified'
					clearTimeout(timer)
					pending.resolve()
					pending = undefined
				},
				'expired-callback': onFailure,
				'error-callback': onFailure,
				'unsupported-callback': onFailure,
			})
			if (failed() || !widget) {
				throw new Error('turnstile-prepare-failed')
			}
			phase = 'ready'
		})().catch(() => {
			if (current === generation && phase !== 'disposed') {
				reset()
				phase = 'failed'
			}
			throw new Error('turnstile-prepare-failed')
		})

		return preparation
	}

	async function execute() {
		await prepare()

		if (phase !== 'ready' || !api || widget === undefined) {
			throw new Error('turnstile-not-ready')
		}
		phase = 'verifying'

		return new Promise<void>((resolve, reject) => {
			pending = {
				resolve,
				reject,
			}
			timer = setTimeout(fail, 30_000)
			try {
				api?.execute(widget as string)
			} catch {
				fail()
			}
		})
	}

	function takeToken(): string {
		if (phase !== 'verified' || !token || widget === undefined || !api) {
			throw new Error('turnstile-token-unavailable')
		}
		if (api.isExpired(widget)) {
			fail()
			throw new Error('turnstile-token-expired')
		}
		const response = token
		token = ''
		phase = 'consumed'

		return response
	}

	return {
		prepare,
		execute,
		takeToken,
		reset,
		dispose: () => {
			phase = 'disposed'
			reset()
		},
	}
}
