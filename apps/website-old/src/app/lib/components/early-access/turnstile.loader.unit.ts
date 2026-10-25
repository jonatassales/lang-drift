// @vitest-environment jsdom
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import type { TurnstileApi } from './turnstile.types'

const api: TurnstileApi = {
	render: vi.fn(),
	execute: vi.fn(),
	reset: vi.fn(),
	remove: vi.fn(),
	isExpired: vi.fn(),
}
beforeEach(() => {
	vi.resetModules()
	vi.useFakeTimers()
	delete window.turnstile
})
afterEach(() => {
	document.head.innerHTML = ''
	delete window.turnstile
	vi.useRealTimers()
})

describe('lazy Turnstile script loader', () => {
	it('loads once on demand and shares readiness only', async () => {
		const { loadTurnstile } = await import('./turnstile.loader')
		expect(document.querySelector('script')).toBeNull()
		const first = loadTurnstile()
		expect(loadTurnstile()).toBe(first)
		expect(document.querySelectorAll('script')).toHaveLength(1)
		window.turnstile = api
		document.querySelector('script')?.dispatchEvent(new Event('load'))
		expect(await first).toBe(api)
		expect(await loadTurnstile()).toBe(api)
	})

	it.each(['error', 'timeout', 'missing-api'])(
		'recovers after %s without caching a rejected load',
		async (failure) => {
			const { loadTurnstile } = await import('./turnstile.loader')
			const rejected = expect(loadTurnstile()).rejects.toThrow(
				'turnstile-load-failed',
			)
			if (failure === 'timeout') await vi.advanceTimersByTimeAsync(15_000)
			else
				document
					.querySelector('script')
					?.dispatchEvent(new Event(failure === 'error' ? 'error' : 'load'))
			await rejected
			expect(document.querySelector('script')).toBeNull()
			const retry = loadTurnstile()
			window.turnstile = api
			document.querySelector('script')?.dispatchEvent(new Event('load'))
			expect(await retry).toBe(api)
		},
	)
})
