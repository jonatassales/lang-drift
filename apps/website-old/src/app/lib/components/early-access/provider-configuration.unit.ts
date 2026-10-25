import { afterEach, beforeEach, expect, it, vi } from 'vitest'

vi.mock('server-only', () => ({}))

import { postEarlyAccess } from './submission.route'
import { verifyEarlyAccessChallenge } from './turnstile.validation'

const input = {
	email: 'a@example.test',
	locale: 'en',
	source: 'landing',
	turnstileToken: 'test-token',
}

beforeEach(() => {
	vi.resetModules()
	vi.stubEnv('NODE_ENV', 'test')
	vi.stubEnv('EARLY_ACCESS_MODE', 'test')
	vi.stubEnv('TURNSTILE_SECRET_KEY', '1x0000000000000000000000000000000AA')
	vi.stubEnv('NEXT_PUBLIC_TURNSTILE_SITE_KEY', '1x00000000000000000000BB')
	vi.stubEnv('TURNSTILE_EXPECTED_HOSTNAME', 'localhost')
	vi.stubEnv('TURNSTILE_VERIFY_URL', 'http://localhost:8082/siteverify')
	vi.stubEnv('EARLY_ACCESS_REGISTRY_URL', 'http://127.0.0.1:8080')
	vi.stubEnv('EARLY_ACCESS_REGISTRY_KEY', 'synthetic-registry-key-000000000000')
	vi.stubEnv('EARLY_ACCESS_PUBLIC_ORIGIN', 'http://localhost:3000')
})
afterEach(() => {
	vi.unstubAllEnvs()
	vi.unstubAllGlobals()
})

function request() {
	return new Request('http://localhost:3000/api/early-access/', {
		method: 'POST',
		headers: {
			host: 'localhost:3000',
			origin: 'http://localhost:3000',
			'content-type': 'application/json',
		},
		body: JSON.stringify(input),
	})
}

function production() {
	vi.stubEnv('NODE_ENV', 'production')
	vi.stubEnv('EARLY_ACCESS_MODE', '')
	vi.stubEnv('TURNSTILE_SECRET_KEY', 'production-secret-fixture')
	vi.stubEnv('NEXT_PUBLIC_TURNSTILE_SITE_KEY', 'production-site-fixture')
	vi.stubEnv(
		'TURNSTILE_VERIFY_URL',
		'https://challenges.cloudflare.com/turnstile/v0/siteverify',
	)
	vi.stubEnv('TURNSTILE_EXPECTED_HOSTNAME', 'langdrift.md')
	vi.stubEnv('EARLY_ACCESS_REGISTRY_URL', 'https://registry.example.test')
	vi.stubEnv('EARLY_ACCESS_PUBLIC_ORIGIN', 'https://langdrift.md')
}

it('uses registry configuration and needs no Resend or legacy ingress credentials', async () => {
	for (const name of [
		'RESEND_API_KEY',
		'RESEND_FROM_EMAIL',
		'RESEND_REPLY_TO',
		'EMAIL_SERVICE_URL',
		'EMAIL_SERVICE_API_KEY',
	])
		vi.stubEnv(name, '')
	const calls: string[] = []
	vi.stubGlobal(
		'fetch',
		vi.fn(async (url: string) => {
			calls.push(url)
			if (url.endsWith('/v1/admissions'))
				return Response.json({ grant: 'a'.repeat(32) })
			return Response.json({
				success: true,
				hostname: 'localhost',
				action: 'test',
			})
		}),
	)
	await postEarlyAccess(request())
	expect(calls[0]).toBe('http://127.0.0.1:8080/v1/admissions')
})

it.each([
	'EARLY_ACCESS_REGISTRY_URL',
	'EARLY_ACCESS_REGISTRY_KEY',
	'EARLY_ACCESS_PUBLIC_ORIGIN',
])('returns unavailable when %s is unset', async (name) => {
	vi.stubEnv(name, '')
	const transport = vi.fn()
	vi.stubGlobal('fetch', transport)
	expect((await postEarlyAccess(request())).status).toBe(503)
	expect(transport).not.toHaveBeenCalled()
})

it('returns unavailable when TURNSTILE_VERIFY_URL is unset', async () => {
	vi.stubEnv('TURNSTILE_VERIFY_URL', '')
	const transport = vi.fn()
	vi.stubGlobal('fetch', transport)
	expect(await verifyEarlyAccessChallenge('token')).toBe('unavailable')
	expect(transport).not.toHaveBeenCalled()
})

it('rejects official test keys in production', async () => {
	production()
	vi.stubEnv('TURNSTILE_SECRET_KEY', '1x0000000000000000000000000000000AA')
	const transport = vi.fn()
	vi.stubGlobal('fetch', transport)
	expect(await verifyEarlyAccessChallenge('token')).toBe('unavailable')
	expect(transport).not.toHaveBeenCalled()
})

it.each([
	'http://challenges.cloudflare.com/turnstile/v0/siteverify',
	'https://elsewhere.example/siteverify',
	'https://challenges.cloudflare.com/turnstile/v0/siteverify?x=1',
	'https://user:pass@challenges.cloudflare.com/turnstile/v0/siteverify',
])('rejects unofficial production Siteverify %s', async (url) => {
	production()
	vi.stubEnv('TURNSTILE_VERIFY_URL', url)
	const transport = vi.fn()
	vi.stubGlobal('fetch', transport)
	expect(await verifyEarlyAccessChallenge('token')).toBe('unavailable')
	expect(transport).not.toHaveBeenCalled()
})

it.each([
	'http://registry.example.test',
	'https://user:pass@registry.example.test',
	'https://registry.example.test/path',
	'https://registry.example.test/#fragment',
])('rejects invalid production registry %s', async (url) => {
	production()
	vi.stubEnv('EARLY_ACCESS_REGISTRY_URL', url)
	const transport = vi.fn()
	vi.stubGlobal('fetch', transport)
	expect((await postEarlyAccess(request())).status).toBe(503)
	expect(transport).not.toHaveBeenCalled()
})

it('rejects weak registry authentication', async () => {
	vi.stubEnv('EARLY_ACCESS_REGISTRY_KEY', 'short')
	const transport = vi.fn()
	vi.stubGlobal('fetch', transport)
	expect((await postEarlyAccess(request())).status).toBe(503)
	expect(transport).not.toHaveBeenCalled()
})
