import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

vi.mock('server-only', () => ({}))

import { POST, postEarlyAccess } from './route'

const registryUrl = 'http://127.0.0.1:8080'
const verifyUrl = 'http://verify.local'

const input = {
	email: "Alias'@EXAMPLE.test",
	locale: 'en',
	source: 'landing',
	turnstileToken: 'test-token',
}

describe('apps/website/src/app/api/early-access/route.ts', () => {
	describe('paused route handler', () => {
		afterEach(() => {
			vi.unstubAllGlobals()
		})

		it('answers unavailable without calling a provider while sign-up is paused', async () => {
			const transport = vi.fn()
			vi.stubGlobal('fetch', transport)

			const response = await POST()

			expect(response.status).toBe(503)
			expect(response.headers.get('Retry-After')).toBe('30')
			expect(response.headers.get('Cache-Control')).toBe('no-store')
			expect(await response.json()).toEqual({
				ok: false,
				code: 'UNAVAILABLE',
			})
			expect(transport).not.toHaveBeenCalled()
		})
	})

	describe('submission logic', () => {
		let calls: Array<{ url: string; body: Record<string, unknown> }>
		let admissionStatus: number
		let registrationStatus: number
		let verification: unknown

		beforeEach(() => {
			vi.stubEnv('EARLY_ACCESS_MODE', 'test')
			vi.stubEnv('EARLY_ACCESS_REGISTRY_URL', registryUrl)
			vi.stubEnv(
				'EARLY_ACCESS_REGISTRY_KEY',
				'synthetic-key-0000000000000000000000',
			)
			vi.stubEnv('EARLY_ACCESS_PUBLIC_ORIGIN', 'http://localhost:3000')
			vi.stubEnv('TURNSTILE_SECRET_KEY', 'synthetic-verification-secret')
			vi.stubEnv('TURNSTILE_EXPECTED_HOSTNAME', 'localhost')
			vi.stubEnv('TURNSTILE_VERIFY_URL', verifyUrl)
			vi.stubEnv('NEXT_PUBLIC_TURNSTILE_SITE_KEY', 'site-key-fixture')
			calls = []
			admissionStatus = 200
			registrationStatus = 202
			verification = {
				success: true,
				hostname: 'localhost',
				action: 'early-access',
			}
			vi.stubGlobal(
				'fetch',
				vi.fn(async (url: string, options: RequestInit) => {
					calls.push({
						url,
						body: JSON.parse(String(options.body)),
					})
					if (url.endsWith('/v1/admissions'))
						return Response.json(
							{ grant: 'a'.repeat(32) },
							{
								status: admissionStatus,
								headers: { 'Retry-After': '60' },
							},
						)
					if (url === verifyUrl) return Response.json(verification)
					if (url.endsWith('/v1/registrations'))
						return Response.json(
							{ status: 'accepted' },
							{
								status: registrationStatus,
								headers: { 'Retry-After': '60' },
							},
						)
					throw new Error('Unexpected transport')
				}),
			)
		})

		afterEach(() => {
			vi.unstubAllGlobals()
			vi.unstubAllEnvs()
		})

		function request(
			body = JSON.stringify(input),
			headers: Record<string, string> = {},
			url = 'http://localhost:3000/api/early-access/',
		) {
			return new Request(url, {
				method: 'POST',
				headers: {
					host: 'localhost:3000',
					origin: 'http://localhost:3000',
					'content-type': 'application/json',
					...headers,
				},
				body,
			})
		}

		it('admits before Siteverify and returns 202 only after durable registry receipt', async () => {
			const response = await postEarlyAccess(request())

			expect(response.status).toBe(202)
			expect(calls.map((call) => call.url)).toEqual([
				'http://127.0.0.1:8080/v1/admissions',
				'http://verify.local',
				'http://127.0.0.1:8080/v1/registrations',
			])
			expect(calls[0]?.body).toMatchObject({
				email: "Alias'@example.test",
				origin: '127.0.0.1',
			})
			expect(Object.keys(calls[2]?.body ?? {}).sort()).toEqual([
				'email',
				'grant',
				'locale',
				'request',
				'source',
			])
			expect(calls[1]?.body).not.toHaveProperty('email')
		})

		it.each([429, 503])(
			'stops before Siteverify on admission %s',
			async (status) => {
				admissionStatus = status
				const response = await postEarlyAccess(request())

				expect(response.status).toBe(status)
				expect(response.headers.get('Retry-After')).toBe('60')
				expect(calls).toHaveLength(1)
			},
		)

		it.each([
			{ success: false },
			{
				success: true,
				hostname: 'wrong',
				action: 'early-access',
			},
			{
				success: true,
				hostname: 'localhost',
				action: 'wrong',
			},
		])('never registers when verification fails %#', async (body) => {
			verification = body
			const response = await postEarlyAccess(request())

			expect(response.status).toBe(403)
			expect(calls).toHaveLength(2)
		})

		it.each([429, 503, 200])(
			'does not claim success on registry %s',
			async (status) => {
				registrationStatus = status
				const response = await postEarlyAccess(request())

				expect(response.status).toBe(status === 429 ? 429 : 503)
			},
		)

		it.each([
			'{"email":"a@example.test","email":"b@example.test"}',
			JSON.stringify({
				...input,
				Email: 'other@example.test',
			}),
			JSON.stringify({
				...input,
				email: ['a@example.test'],
			}),
			JSON.stringify({
				...input,
				email: 'a@example.test\n',
			}),
			JSON.stringify({
				...input,
				turnstileToken: 'x'.repeat(2049),
			}),
		])('rejects malformed input without external calls %#', async (body) => {
			const response = await postEarlyAccess(request(body))

			expect(response.status).toBe(400)
			expect(calls).toHaveLength(0)
		})

		it('enforces actual byte size and encoding before external calls', async () => {
			const responseOversized = await postEarlyAccess(
				request(' '.repeat(4097), { 'content-length': '1' }),
			)

			expect(responseOversized.status).toBe(413)

			const responseEncoded = await postEarlyAccess(
				request(undefined, { 'content-encoding': 'gzip' }),
			)

			expect(responseEncoded.status).toBe(415)
			expect(calls).toHaveLength(0)
		})

		it('refuses a non-JSON type and an oversized declared length before external calls', async () => {
			const responseType = await postEarlyAccess(
				request(undefined, { 'content-type': 'text/plain' }),
			)

			expect(responseType.status).toBe(415)

			const responseLength = await postEarlyAccess(
				request(undefined, { 'content-length': '4097' }),
			)

			expect(responseLength.status).toBe(413)
			expect(calls).toHaveLength(0)
		})

		it('rejects alternate hostname and cross-site origin', async () => {
			const responseOrigin = await postEarlyAccess(
				request(undefined, { origin: 'https://evil.example' }),
			)

			expect(responseOrigin.status).toBe(403)

			const responseHost = await postEarlyAccess(
				request(
					undefined,
					{ host: 'alternate.example' },
					'http://alternate.example/api/early-access/',
				),
			)

			expect(responseHost.status).toBe(403)
			expect(calls).toHaveLength(0)
		})

		it('ignores browser-forwarded identity locally and fails closed on untrusted production', async () => {
			await postEarlyAccess(
				request(undefined, {
					'x-forwarded-for': '203.0.113.77',
				}),
			)

			expect(calls[0]?.body.origin).toBe('127.0.0.1')
			calls = []
			vi.stubEnv('EARLY_ACCESS_MODE', '')
			vi.stubEnv('EARLY_ACCESS_REGISTRY_URL', 'https://registry.example.test')
			vi.stubEnv('EARLY_ACCESS_PUBLIC_ORIGIN', 'https://langdrift.md')
			vi.stubEnv(
				'TURNSTILE_VERIFY_URL',
				'https://challenges.cloudflare.com/turnstile/v0/siteverify',
			)
			vi.stubEnv('TURNSTILE_SECRET_KEY', 'production-secret-fixture')
			vi.stubEnv('NEXT_PUBLIC_TURNSTILE_SITE_KEY', 'production-site-fixture')
			vi.stubEnv('TURNSTILE_EXPECTED_HOSTNAME', 'langdrift.md')

			const response = await postEarlyAccess(
				request(
					undefined,
					{
						host: 'langdrift.md',
						origin: 'https://langdrift.md',
						'x-forwarded-for': '203.0.113.7',
					},
					'https://langdrift.md/api/early-access/',
				),
			)

			expect(response.status).toBe(503)
			expect(calls).toHaveLength(0)
		})

		it.each([null, false, { success: true }])(
			'returns 503 on malformed Siteverify without registering %#',
			async (body) => {
				verification = body
				const response = await postEarlyAccess(request())

				expect(response.status).toBe(503)
				expect(response.headers.get('Retry-After')).toBe('30')
				expect(calls).toHaveLength(2)
			},
		)
	})
})
