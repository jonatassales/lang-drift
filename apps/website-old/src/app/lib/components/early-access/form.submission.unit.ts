import { afterEach, describe, expect, it, vi } from 'vitest'

import { submitEarlyAccess } from './form.submission'

afterEach(() => vi.unstubAllGlobals())

describe('public submission transport', () => {
	it.each([
		[202, { ok: true }],
		[
			400,
			{
				ok: false,
				code: 'VALIDATION_ERROR',
			},
		],
		[
			403,
			{
				ok: false,
				code: 'CHALLENGE_FAILED',
			},
		],
		[
			503,
			{
				ok: false,
				code: 'UNAVAILABLE',
			},
		],
	])(
		'maps HTTP %s without exposing provider bodies',
		async (status, expected) => {
			const transport = vi
				.fn()
				.mockResolvedValue(new Response('opaque', { status }))
			vi.stubGlobal('fetch', transport)
			expect(
				await submitEarlyAccess({ email: 'synthetic@example.test' }),
			).toEqual(expected)
			expect(transport).toHaveBeenCalledWith(
				'/api/early-access/',
				expect.objectContaining({
					method: 'POST',
					credentials: 'same-origin',
					redirect: 'error',
					cache: 'no-store',
				}),
			)
		},
	)

	it.each([
		['60', 60],
		['-1', 30],
		['garbage', 30],
		['999999999', 2678400],
	])('bounds retry feedback %s', async (header, seconds) => {
		vi.stubGlobal(
			'fetch',
			vi.fn().mockResolvedValue(
				new Response(null, {
					status: 429,
					headers: { 'Retry-After': header },
				}),
			),
		)
		expect(await submitEarlyAccess({})).toEqual({
			ok: false,
			code: 'RATE_LIMITED',
			retryAfter: seconds,
		})
	})

	it('propagates navigation cancellation to the browser request', async () => {
		const controller = new AbortController()
		controller.abort()
		vi.stubGlobal(
			'fetch',
			vi.fn((_url, options) => {
				expect(options.signal.aborted).toBe(true)
				return Promise.reject(new DOMException('Aborted', 'AbortError'))
			}),
		)
		await expect(
			submitEarlyAccess({}, controller.signal),
		).rejects.toMatchObject({ name: 'AbortError' })
	})
})

it.each(['emailRequired', 'emailTooLong', 'emailInvalid'])(
	'preserves validated server field error %s',
	async (email) => {
		vi.stubGlobal(
			'fetch',
			vi
				.fn()
				.mockResolvedValue(
					Response.json({ fieldErrors: { email } }, { status: 400 }),
				),
		)
		expect(await submitEarlyAccess({})).toEqual({
			ok: false,
			code: 'VALIDATION_ERROR',
			fieldErrors: { email },
		})
	},
)
