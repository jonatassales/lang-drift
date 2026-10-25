// @vitest-environment jsdom
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import { createTurnstileController } from './turnstile.controller'
import type { TurnstileApi, TurnstileOptions } from './turnstile.types'

vi.mock('./turnstile.loader', () => ({ loadTurnstile: vi.fn() }))

import { loadTurnstile } from './turnstile.loader'

let options: TurnstileOptions[]
let api: TurnstileApi
beforeEach(() => {
	vi.useFakeTimers()
	vi.resetAllMocks()
	options = []
	api = {
		render: vi.fn((_element, value) => {
			options.push(value)
			return `widget-${options.length}`
		}),
		execute: vi.fn(),
		isExpired: vi.fn(() => false),
		reset: vi.fn(),
		remove: vi.fn(),
	}
	vi.mocked(loadTurnstile).mockResolvedValue(api)
})
afterEach(() => vi.useRealTimers())
const controller = () =>
	createTurnstileController(document.createElement('div'), 'site-key')

describe('Turnstile controller', () => {
	it('shares pending preparation and waits before executing', async () => {
		let loaded!: (api: TurnstileApi) => void
		vi.mocked(loadTurnstile).mockReturnValue(
			new Promise((resolve) => {
				loaded = resolve
			}),
		)
		const challenge = controller()
		const preparing = challenge.prepare()
		const execution = challenge.execute()
		expect(api.execute).not.toHaveBeenCalled()
		loaded(api)
		await preparing
		await vi.advanceTimersByTimeAsync(0)
		expect(api.render).toHaveBeenCalledTimes(1)
		expect(api.execute).toHaveBeenCalledTimes(1)
		options[0]?.callback('token')
		await execution
		expect(challenge.takeToken()).toBe('token')
		expect(() => challenge.takeToken()).toThrow()
		challenge.dispose()
	})

	it('times out verification and permits a fresh attempt', async () => {
		const challenge = controller()
		const execution = challenge.execute()
		const failed = expect(execution).rejects.toThrow()
		await vi.advanceTimersByTimeAsync(30_000)
		await failed
		challenge.reset()
		const retry = challenge.execute()
		await vi.advanceTimersByTimeAsync(0)
		options[0]?.callback('late-token')
		options[1]?.callback('fresh-token')
		await retry
		expect(challenge.takeToken()).toBe('fresh-token')
		challenge.dispose()
	})

	it('cancels pending preparation without rendering after disposal', async () => {
		let loaded!: (api: TurnstileApi) => void
		vi.mocked(loadTurnstile).mockReturnValue(
			new Promise((resolve) => {
				loaded = resolve
			}),
		)
		const challenge = controller()
		const failed = expect(challenge.execute()).rejects.toThrow()
		challenge.dispose()
		loaded(api)
		await failed
		expect(api.render).not.toHaveBeenCalled()
		await expect(challenge.prepare()).rejects.toThrow()
	})

	it.each(['render', 'execute'] as const)(
		'handles thrown %s exceptions',
		async (method) => {
			vi.mocked(api[method]).mockImplementation(() => {
				throw new Error('provider detail')
			})
			const challenge = controller()
			await expect(challenge.execute()).rejects.toThrow()
			expect(() => challenge.takeToken()).toThrow()
			challenge.dispose()
		},
	)

	it('retries failed preloading and ignores render-time token callbacks', async () => {
		vi.mocked(loadTurnstile).mockRejectedValueOnce(new Error('offline'))
		const challenge = controller()
		await expect(challenge.prepare()).rejects.toThrow()
		vi.mocked(api.render).mockImplementation((_element, callbacks) => {
			callbacks.callback('unrequested-token')
			options.push(callbacks)
			return 'widget'
		})
		await challenge.prepare()
		expect(() => challenge.takeToken()).toThrow()
		challenge.dispose()
	})

	it.each(['', 'x'.repeat(2049)])(
		'rejects malformed callback tokens (%s)',
		async (token) => {
			const challenge = controller()
			const execution = challenge.execute()
			await vi.advanceTimersByTimeAsync(0)
			const failed = expect(execution).rejects.toThrow()
			options[0]?.callback(token)
			await failed
			challenge.dispose()
		},
	)

	it('checks provider expiry immediately before consuming the token', async () => {
		const challenge = controller()
		const execution = challenge.execute()
		await vi.advanceTimersByTimeAsync(0)
		options[0]?.callback('token')
		await execution
		vi.mocked(api.isExpired).mockReturnValue(true)
		expect(() => challenge.takeToken()).toThrow()
		challenge.dispose()
	})

	it('keeps two forms independent and cleanup tolerates provider exceptions', async () => {
		const first = controller()
		const second = controller()
		const firstExecution = first.execute()
		const secondExecution = second.execute()
		await vi.advanceTimersByTimeAsync(0)
		options[1]?.callback('second-token')
		options[0]?.callback('first-token')
		await Promise.all([firstExecution, secondExecution])
		expect(first.takeToken()).toBe('first-token')
		expect(second.takeToken()).toBe('second-token')
		vi.mocked(api.reset).mockImplementation(() => {
			throw new Error('gone')
		})
		vi.mocked(api.remove).mockImplementation(() => {
			throw new Error('gone')
		})
		expect(() => {
			first.dispose()
			second.dispose()
		}).not.toThrow()
	})
})
