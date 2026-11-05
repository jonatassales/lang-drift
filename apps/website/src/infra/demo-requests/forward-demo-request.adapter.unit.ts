import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import { forwardDemoRequest } from './forward-demo-request.adapter'

const submission = {
	email: 'jane@acme.io',
	company: 'Acme',
	position: 'cto',
	locale: 'en',
} as const

const token = 'a-demo-requests-token-with-thirty-two-plus-characters'

let transport: ReturnType<typeof vi.fn>
beforeEach(() => {
	transport = vi.fn(() =>
		Promise.resolve(
			Response.json(
				{
					id: '701706be-7e70-4681-9608-826d3ad548b9',
				},
				{
					status: 202,
				},
			),
		),
	)
	vi.stubGlobal('fetch', transport)
	vi.stubEnv('DEMO_REQUESTS_URL', 'http://127.0.0.1:3402')
	vi.stubEnv('DEMO_REQUESTS_TOKEN', token)
})
afterEach(() => {
	vi.unstubAllGlobals()
	vi.unstubAllEnvs()
})

describe('website/infra/demo-requests/forward-demo-request.adapter', () => {
	it('posts the submission as JSON with the bearer token, without cache or redirects, and returns the status', async () => {
		const status = await forwardDemoRequest(
			submission,
			new AbortController().signal,
		)

		expect(status).toBe(202)
		const [url, init] = transport.mock.calls[0] as [string, RequestInit]
		expect(url).toBe('http://127.0.0.1:3402/v1/demo-requests')
		expect(init).toMatchObject({
			method: 'POST',
			cache: 'no-store',
			redirect: 'error',
			headers: {
				authorization: `Bearer ${token}`,
				'content-type': 'application/json',
			},
		})
		expect(JSON.parse(String(init.body))).toEqual(submission)
		expect(init.signal).toBeInstanceOf(AbortSignal)
	})

	it('returns whatever status the service answers', async () => {
		transport.mockImplementation(() =>
			Promise.resolve(new Response(null, { status: 401 })),
		)

		await expect(
			forwardDemoRequest(submission, new AbortController().signal),
		).resolves.toBe(401)
	})

	it.each(['DEMO_REQUESTS_URL', 'DEMO_REQUESTS_TOKEN'])(
		'throws naming %s when it is not set, before calling the service',
		async (key) => {
			vi.stubEnv(key, '')

			await expect(
				forwardDemoRequest(submission, new AbortController().signal),
			).rejects.toThrow(new RegExp(key))

			expect(transport).not.toHaveBeenCalled()
		},
	)
})
