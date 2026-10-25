import { afterEach, beforeEach, expect, it, vi } from 'vitest'

vi.mock('server-only', () => ({}))

import { verifyEarlyAccessChallenge } from './turnstile.validation'

beforeEach(() => {
	vi.stubEnv('EARLY_ACCESS_MODE', 'test')
	vi.stubEnv('TURNSTILE_SECRET_KEY', '1x0000000000000000000000000000000AA')
	vi.stubEnv('TURNSTILE_EXPECTED_HOSTNAME', 'localhost')
	vi.stubEnv('TURNSTILE_VERIFY_URL', 'http://verify.local')
	vi.stubEnv('NEXT_PUBLIC_TURNSTILE_SITE_KEY', '1x00000000000000000000BB')
	vi.stubEnv('VERCEL', '')
	vi.stubEnv('VERCEL_ENV', '')
})
afterEach(() => {
	vi.unstubAllGlobals()
	vi.unstubAllEnvs()
})
it('never contacts providers for invalid token boundaries', async () => {
	const transport = vi.fn()
	vi.stubGlobal('fetch', transport)
	expect(await verifyEarlyAccessChallenge('')).toBe('rejected')
	expect(await verifyEarlyAccessChallenge('x'.repeat(2049))).toBe('rejected')
	expect(transport).not.toHaveBeenCalled()
})
it.each([
	[
		{
			success: true,
			hostname: 'localhost',
			action: 'test',
		},
		'verified',
	],
	[
		{
			success: true,
			hostname: 'wrong',
			action: 'test',
		},
		'rejected',
	],
	[
		{
			success: true,
			hostname: 'localhost',
			action: 'other',
		},
		'rejected',
	],
	[
		{
			success: false,
			'error-codes': ['timeout-or-duplicate'],
		},
		'rejected',
	],
	[null, 'unavailable'],
	[{ success: true }, 'unavailable'],
])(
	'distinguishes verdict from broken provider contract %#',
	async (body, result) => {
		vi.stubGlobal('fetch', vi.fn().mockResolvedValue(Response.json(body)))
		expect(await verifyEarlyAccessChallenge('token')).toBe(result)
	},
)
it.each([
	new Response('bad', { status: 500 }),
	new Response('{'),
	new Response('x'.repeat(4097)),
])('bounds provider failures %#', async (response) => {
	vi.stubGlobal('fetch', vi.fn().mockResolvedValue(response))
	expect(await verifyEarlyAccessChallenge('token')).toBe('unavailable')
})
it('classifies timeout and transport failure as unavailable', async () => {
	vi.stubGlobal(
		'fetch',
		vi.fn().mockRejectedValue(new DOMException('timeout', 'TimeoutError')),
	)
	expect(await verifyEarlyAccessChallenge('token')).toBe('unavailable')
})
