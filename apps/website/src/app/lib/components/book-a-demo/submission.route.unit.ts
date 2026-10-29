import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

vi.mock('server-only', () => ({}))

import { postBookDemo } from './submission.route'

const valid = {
	company: 'Acme',
	email: 'jane@acme.io',
	locale: 'en',
	position: 'cto',
}

// What a browser sends with a fetch from the site's own pages.
const sameOrigin = {
	origin: 'http://localhost:3000',
	'sec-fetch-site': 'same-origin',
}

function request(
	body: string,
	headers: Record<string, string> = {},
	method = 'POST',
) {
	return new Request('http://localhost:3000/api/book-a-demo/', {
		body,
		headers: {
			'content-type': 'application/json',
			...headers,
		},
		method,
	})
}

function post(body: string, headers?: Record<string, string>) {
	return postBookDemo(request(body, headers))
}

let transport: ReturnType<typeof vi.fn>
beforeEach(() => {
	transport = vi.fn()
	vi.stubGlobal('fetch', transport)
})
afterEach(() => {
	vi.unstubAllGlobals()
})

describe('a demo request at the route', () => {
	it('refuses a body that is not a JSON object', async () => {
		const response = await post('[1, 2]')

		expect(response.status).toBe(400)
		expect(await response.json()).toEqual({
			ok: false,
			code: 'VALIDATION_ERROR',
		})
	})

	it('refuses malformed JSON', async () => {
		const response = await post('{"email": ')

		expect(response.status).toBe(400)
	})

	it('refuses a body over four kibibytes', async () => {
		const response = await post(
			JSON.stringify({
				...valid,
				company: 'x'.repeat(5000),
			}),
		)

		expect(response.status).toBe(413)
	})

	it('names the field that failed', async () => {
		const response = await post(
			JSON.stringify({
				...valid,
				email: 'jane@gmail.com',
				position: '',
			}),
		)

		expect(response.status).toBe(400)
		expect(await response.json()).toEqual({
			ok: false,
			code: 'VALIDATION_ERROR',
			fieldErrors: {
				email: 'emailFreeMail',
				position: 'positionRequired',
			},
		})
	})

	it('refuses a submission with a key it does not know', async () => {
		const response = await post(
			JSON.stringify({
				...valid,
				source: 'home',
			}),
		)

		expect(response.status).toBe(400)
	})

	it('answers that requests are paused for a valid submission, storing and sending nothing', async () => {
		const response = await post(JSON.stringify(valid))

		expect(response.status).toBe(503)
		expect(response.headers.get('Retry-After')).toBe('30')
		expect(response.headers.get('Cache-Control')).toBe('no-store')
		expect(response.headers.get('X-Content-Type-Options')).toBe('nosniff')
		expect(await response.json()).toEqual({
			ok: false,
			code: 'UNAVAILABLE',
		})
		expect(transport).not.toHaveBeenCalled()
	})

	it('keeps field errors and the pause for a request from the site itself', async () => {
		const invalid = await post(
			JSON.stringify({
				...valid,
				email: 'jane@gmail.com',
			}),
			sameOrigin,
		)
		const submitted = request(JSON.stringify(valid), sameOrigin)
		const paused = await postBookDemo(submitted)

		expect(invalid.status).toBe(400)
		expect(paused.status).toBe(503)
		expect(submitted.bodyUsed).toBe(true)
	})

	it('keeps the pause for a page served from a loopback address', async () => {
		// Next names every loopback host localhost in the request URL, while the
		// browser names the address it loaded the page from.
		const response = await postBookDemo(
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

		expect(response.status).toBe(503)
	})
})

describe('a demo request refused before its body is read', () => {
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
		const refused = request(JSON.stringify(valid), headers)
		const response = await postBookDemo(refused)

		expect(response.status).toBe(status)
		expect(refused.bodyUsed).toBe(false)
		expect(response.headers.get('Cache-Control')).toBe('no-store')
		expect(response.headers.get('X-Content-Type-Options')).toBe('nosniff')
	})

	it('answers any method but POST with 405', async () => {
		const refused = request(JSON.stringify(valid), sameOrigin, 'PUT')
		const response = await postBookDemo(refused)

		expect(response.status).toBe(405)
		expect(refused.bodyUsed).toBe(false)
	})
})
