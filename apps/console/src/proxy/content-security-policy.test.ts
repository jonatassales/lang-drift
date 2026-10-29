import assert from 'node:assert/strict'
import test from 'node:test'

import { createNonce, documentPolicy } from './content-security-policy.ts'

function sources(policy: string, directive: string) {
	const entry = policy
		.split('; ')
		.find((candidate) => candidate.startsWith(`${directive} `))

	return entry?.split(' ').slice(1)
}

test('every response gets a fresh nonce of 16 random bytes', () => {
	const nonces = new Set(Array.from({ length: 64 }, createNonce))

	assert.equal(nonces.size, 64)
	for (const nonce of nonces)
		assert.equal(Buffer.from(nonce, 'base64').length, 16)
})

test('scripts need this response nonce, never an inline grant', () => {
	const nonce = createNonce()
	const policy = documentPolicy(nonce, 'production')

	assert.deepEqual(sources(policy, 'script-src'), [
		"'self'",
		`'nonce-${nonce}'`,
		"'strict-dynamic'",
	])
	assert.notEqual(
		documentPolicy(createNonce(), 'production'),
		documentPolicy(createNonce(), 'production'),
	)
})

test('production allows the console origin and nothing else', () => {
	assert.equal(
		documentPolicy('bm9uY2U=', 'production'),
		[
			"default-src 'self'",
			"script-src 'self' 'nonce-bm9uY2U=' 'strict-dynamic'",
			"style-src 'self' 'unsafe-inline'",
			"img-src 'self' data: blob:",
			"font-src 'self'",
			"connect-src 'self'",
			"worker-src 'self' blob:",
			"manifest-src 'self'",
			"frame-src 'none'",
			"object-src 'none'",
			"base-uri 'self'",
			"form-action 'self'",
			"frame-ancestors 'none'",
		].join('; '),
	)
})

test('development adds eval for React and the development server socket', () => {
	const policy = documentPolicy('bm9uY2U=', 'development')

	assert.deepEqual(sources(policy, 'script-src'), [
		"'self'",
		"'nonce-bm9uY2U='",
		"'strict-dynamic'",
		"'unsafe-eval'",
	])
	assert.deepEqual(sources(policy, 'connect-src'), ["'self'", 'ws:'])
	assert.doesNotMatch(policy, /vercel\.live/u)
})

test('preview adds the Vercel Toolbar, which production never allows', () => {
	const policy = documentPolicy('bm9uY2U=', 'preview')

	for (const directive of [
		'script-src',
		'style-src',
		'img-src',
		'font-src',
		'connect-src',
	])
		assert.ok(sources(policy, directive)?.includes('https://vercel.live'))
	assert.deepEqual(sources(policy, 'frame-src'), ['https://vercel.live'])
	assert.doesNotMatch(policy, /unsafe-eval|ws:/u)
	assert.doesNotMatch(documentPolicy('bm9uY2U=', 'production'), /vercel\.live/u)
})

test('no environment allows inline scripts, plugins or framing', () => {
	for (const environment of ['development', 'preview', 'production'] as const) {
		const policy = documentPolicy('bm9uY2U=', environment)

		assert.ok(!sources(policy, 'script-src')?.includes("'unsafe-inline'"))
		assert.deepEqual(sources(policy, 'object-src'), ["'none'"])
		assert.deepEqual(sources(policy, 'frame-ancestors'), ["'none'"])
		assert.deepEqual(sources(policy, 'base-uri'), ["'self'"])
		assert.deepEqual(sources(policy, 'form-action'), ["'self'"])
	}
})
