import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

vi.mock('server-only', () => ({}))

import { acceptedSubmission } from '@api/book-a-demo/testdata/accepted-submission.mock'
import { demoRequest } from '@api/book-a-demo/testdata/demo-request.mock'
import { postDemoRequest } from '@api/book-a-demo/testdata/post-demo-request.mock'

import { POST } from './route'

const valid = {
	company: 'Acme',
	email: 'jane@acme.io',
	locale: 'en',
	position: 'cto',
}

const token = 'a-demo-requests-token-with-thirty-two-plus-characters'

// What a browser sends with a fetch from the site's own pages.
const sameOrigin = {
	origin: 'http://localhost:3000',
	'sec-fetch-site': 'same-origin',
}

let transport: ReturnType<typeof vi.fn>
beforeEach(() => {
	transport = vi.fn(acceptedSubmission)
	vi.stubGlobal('fetch', transport)
	vi.stubEnv('DEMO_REQUESTS_URL', 'http://127.0.0.1:3402')
	vi.stubEnv('DEMO_REQUESTS_TOKEN', token)
})
afterEach(() => {
	vi.unstubAllGlobals()
	vi.unstubAllEnvs()
})

describe('website/api/book-a-demo/route', () => {
	it('refuses a body that is not a JSON object', async () => {
		const response = await postDemoRequest('[1, 2]')

		expect(response.status).toBe(400)

		expect(await response.json()).toEqual({
			code: 'VALIDATION_ERROR',
			ok: false,
		})
	})

	it('refuses malformed JSON', async () => {
		const response = await postDemoRequest('{"email": ')

		expect(response.status).toBe(400)
	})

	it('refuses a body over four kibibytes', async () => {
		const response = await postDemoRequest(
			JSON.stringify({
				...valid,
				company: 'x'.repeat(5000),
			}),
		)

		expect(response.status).toBe(413)
	})

	it('names the field that failed', async () => {
		const response = await postDemoRequest(
			JSON.stringify({
				...valid,
				email: 'jane@gmail.com',
				position: '',
			}),
		)

		expect(response.status).toBe(400)

		expect(await response.json()).toEqual({
			code: 'VALIDATION_ERROR',
			fieldErrors: {
				email: 'emailFreeMail',
				position: 'positionRequired',
			},
			ok: false,
		})

		expect(transport).not.toHaveBeenCalled()
	})

	it('refuses a submission with a key it does not know', async () => {
		const response = await postDemoRequest(
			JSON.stringify({
				...valid,
				source: 'home',
			}),
		)

		expect(response.status).toBe(400)
	})

	it('forwards an accepted submission to the demo-requests service with the bearer token and answers OK on 202', async () => {
		const response = await postDemoRequest(JSON.stringify(valid))

		expect(response.status).toBe(200)
		expect(response.headers.get('Cache-Control')).toBe('no-store')
		expect(response.headers.get('X-Content-Type-Options')).toBe('nosniff')

		expect(await response.json()).toEqual({
			code: 'OK',
			ok: true,
		})

		expect(transport).toHaveBeenCalledTimes(1)
		const [url, init] = transport.mock.calls[0] as [string, RequestInit]
		expect(url).toBe('http://127.0.0.1:3402/v1/demo-requests')
		expect(init.method).toBe('POST')
		expect(init.headers).toMatchObject({
			authorization: `Bearer ${token}`,
			'content-type': 'application/json',
		})
		expect(JSON.parse(String(init.body))).toEqual(valid)
	})

	it('answers that requests are paused when the service answers anything but 202', async () => {
		transport.mockImplementation(() =>
			Promise.resolve(
				Response.json(
					{
						message: 'boom',
					},
					{
						status: 500,
					},
				),
			),
		)

		const response = await postDemoRequest(JSON.stringify(valid))

		expect(response.status).toBe(503)
		expect(response.headers.get('Retry-After')).toBe('30')

		expect(await response.json()).toEqual({
			code: 'UNAVAILABLE',
			ok: false,
		})
	})

	it('answers that requests are paused when the service cannot be reached', async () => {
		transport.mockImplementation(() =>
			Promise.reject(new Error('fetch failed')),
		)

		const response = await postDemoRequest(JSON.stringify(valid))

		expect(response.status).toBe(503)

		expect(await response.json()).toEqual({
			code: 'UNAVAILABLE',
			ok: false,
		})
	})

	it.each(['DEMO_REQUESTS_URL', 'DEMO_REQUESTS_TOKEN'])(
		'answers that requests are paused without calling the service when %s is not set',
		async (key) => {
			vi.stubEnv(key, '')

			const response = await postDemoRequest(JSON.stringify(valid))

			expect(response.status).toBe(503)
			expect(transport).not.toHaveBeenCalled()
		},
	)

	it('processes a valid request from the site itself', async () => {
		const invalid = await postDemoRequest(
			JSON.stringify({
				...valid,
				email: 'jane@gmail.com',
			}),
			sameOrigin,
		)

		const submitted = demoRequest(JSON.stringify(valid), sameOrigin)

		const accepted = await POST(submitted)

		expect(invalid.status).toBe(400)
		expect(accepted.status).toBe(200)
		expect(submitted.bodyUsed).toBe(true)
	})

	it('accepts a request from a loopback address', async () => {
		const response = await POST(
			new Request('http://localhost:3110/api/book-a-demo/', {
				body: JSON.stringify(valid),
				headers: {
					'content-type': 'application/json',
					host: '127.0.0.1:3110',
					origin: 'http://127.0.0.1:3110',
					'sec-fetch-site': 'same-origin',
				},
				method: 'POST',
			}),
		)

		expect(response.status).toBe(200)
	})

	it.each<[string, number, Record<string, string>]>([
		['a cross-site fetch', 403, { 'sec-fetch-site': 'cross-site' }],
		['a foreign origin', 403, { origin: 'https://evil.example' }],
		['an encoded body', 415, { 'content-encoding': 'gzip' }],
		['a text body', 415, { 'content-type': 'text/plain' }],
		[
			'a form body',
			415,
			{ 'content-type': 'application/x-www-form-urlencoded' },
		],
		[
			'a declared length over four kibibytes',
			413,
			{ 'content-length': '4097' },
		],
	])('answers %s with %i', async (_refusal, status, headers) => {
		const refused = demoRequest(JSON.stringify(valid), headers)

		const response = await POST(refused)

		expect(response.status).toBe(status)
		expect(refused.bodyUsed).toBe(false)
		expect(response.headers.get('Cache-Control')).toBe('no-store')
		expect(response.headers.get('X-Content-Type-Options')).toBe('nosniff')
		expect(transport).not.toHaveBeenCalled()
	})

	it('answers any method but POST with 405', async () => {
		const refused = demoRequest(JSON.stringify(valid), sameOrigin, 'PUT')

		const response = await POST(refused)

		expect(response.status).toBe(405)
		expect(refused.bodyUsed).toBe(false)
	})
})
