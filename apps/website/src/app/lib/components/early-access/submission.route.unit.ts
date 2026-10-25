import { afterEach, beforeEach, expect, it, vi } from 'vitest'

vi.mock('server-only', () => ({}))

import { postEarlyAccess } from './submission.route'

const registryUrl = 'http://127.0.0.1:8080'
const verifyUrl = 'http://verify.local'

const input = {
	email: "Alias'@EXAMPLE.test",
	locale: 'en',
	source: 'landing',
	turnstileToken: 'test-token',
}
let calls: Array<{ url: string; body: Record<string, unknown> }>
let admissionStatus: number
let registrationStatus: number
let verification: unknown
beforeEach(() => {
	vi.stubEnv('WEBSITE_PREVIEW_SIGNUP_DISABLED', 'false')
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
	vi.stubEnv('VERCEL', '')
	vi.stubEnv('VERCEL_ENV', '')
	calls = []
	admissionStatus = 200
	registrationStatus = 202
	verification = {
		success: true,
		hostname: 'localhost',
		action: 'early-access',
	}
	vi.stubEnv('VERCEL', undefined)
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
it('admits before Siteverify and returns202 only after durable registry receipt', async () => {
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
	'stops before Siteverify on admission%s',
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
	expect((await postEarlyAccess(request())).status).toBe(403)
	expect(calls).toHaveLength(2)
})
it.each([429, 503, 200])(
	'does not claim success on registry%s',
	async (status) => {
		registrationStatus = status
		expect((await postEarlyAccess(request())).status).toBe(
			status === 429 ? 429 : 503,
		)
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
	expect((await postEarlyAccess(request(body))).status).toBe(400)
	expect(calls).toHaveLength(0)
})
it('enforces actual byte size and encoding before external calls', async () => {
	expect(
		(
			await postEarlyAccess(
				request(' '.repeat(4097), { 'content-length': '1' }),
			)
		).status,
	).toBe(413)
	expect(
		(await postEarlyAccess(request(undefined, { 'content-encoding': 'gzip' })))
			.status,
	).toBe(415)
	expect(calls).toHaveLength(0)
})
it('rejects alternate hostname and cross-site origin', async () => {
	expect(
		(
			await postEarlyAccess(
				request(undefined, { origin: 'https://evil.example' }),
			)
		).status,
	).toBe(403)
	expect(
		(
			await postEarlyAccess(
				request(
					undefined,
					{ host: 'alternate.example' },
					'http://alternate.example/api/early-access/',
				),
			)
		).status,
	).toBe(403)
	expect(calls).toHaveLength(0)
})
it('ignores browser-forwarded identity locally and fails closed on untrusted production', async () => {
	await postEarlyAccess(
		request(undefined, {
			'x-forwarded-for': '203.0.113.77',
			'x-vercel-forwarded-for': '203.0.113.78',
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
	expect(
		(
			await postEarlyAccess(
				request(
					undefined,
					{
						host: 'langdrift.md',
						origin: 'https://langdrift.md',
					},
					'https://langdrift.md/api/early-access/',
				),
			)
		).status,
	).toBe(503)
	expect(calls).toHaveLength(0)
})
it('requires a single provider-owned IP when running on Vercel', async () => {
	vi.stubEnv('VERCEL', '1')
	vi.stubEnv('EARLY_ACCESS_REGISTRY_URL', 'https://registry.example.test')
	vi.stubEnv('EARLY_ACCESS_PUBLIC_ORIGIN', 'https://langdrift.md')
	vi.stubEnv(
		'TURNSTILE_VERIFY_URL',
		'https://challenges.cloudflare.com/turnstile/v0/siteverify',
	)
	vi.stubEnv('TURNSTILE_SECRET_KEY', 'production-secret-fixture')
	vi.stubEnv('NEXT_PUBLIC_TURNSTILE_SITE_KEY', 'production-site-fixture')
	vi.stubEnv('TURNSTILE_EXPECTED_HOSTNAME', 'langdrift.md')
	const hosted = (headers: Record<string, string> = {}) =>
		request(
			undefined,
			{
				host: 'langdrift.md',
				origin: 'https://langdrift.md',
				...headers,
			},
			'https://langdrift.md/api/early-access/',
		)
	expect(
		(await postEarlyAccess(hosted({ 'x-forwarded-for': '203.0.113.7' })))
			.status,
	).toBe(503)
	expect(
		(
			await postEarlyAccess(
				hosted({
					'x-vercel-forwarded-for': '203.0.113.8, 203.0.113.9',
				}),
			)
		).status,
	).toBe(503)
	expect(calls).toHaveLength(0)
})
it('disabled previews produce no false success and no external calls', async () => {
	vi.stubEnv('WEBSITE_PREVIEW_SIGNUP_DISABLED', 'true')
	expect((await postEarlyAccess(request())).status).toBe(503)
	expect(calls).toHaveLength(0)
})

it.each([null, false, { success: true }])(
	'returns503 on malformed Siteverify without registering %#',
	async (body) => {
		verification = body
		const response = await postEarlyAccess(request())
		expect(response.status).toBe(503)
		expect(response.headers.get('Retry-After')).toBe('30')
		expect(calls).toHaveLength(2)
	},
)
