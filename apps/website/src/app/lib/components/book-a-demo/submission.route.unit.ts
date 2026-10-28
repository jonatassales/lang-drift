import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

vi.mock('server-only', () => ({}))

import { postBookDemo } from './submission.route'

const valid = {
	company: 'Acme',
	email: 'jane@acme.io',
	locale: 'en',
	position: 'cto',
}

function post(body: string) {
	return postBookDemo(
		new Request('http://localhost:3000/api/book-a-demo/', {
			body,
			headers: { 'content-type': 'application/json' },
			method: 'POST',
		}),
	)
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
})
