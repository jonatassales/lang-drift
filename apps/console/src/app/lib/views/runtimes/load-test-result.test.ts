import assert from 'node:assert/strict'
import test from 'node:test'

import { isLoadTestResult } from './load-test-result.ts'

const result = {
	target: 'console-stubs-loopback',
	requests: 50,
	concurrency: 50,
	succeeded: 49,
	failed: 1,
	durationMs: 85,
	p50Ms: 10,
	p95Ms: 25,
	maxMs: 80,
	startedAt: '2026-10-21T10:00:00Z',
	completedAt: '2026-10-21T10:00:00.085Z',
}

test('burst evidence accepts observed partial failures and rejects fabricated or malformed results', () => {
	assert.ok(isLoadTestResult(result))
	for (const override of [
		{ target: 'https://example.com' },
		{ requests: 5000 },
		{ concurrency: 5000 },
		{ succeeded: 50 },
		{ failed: -1 },
		{ p95Ms: NaN },
		{ p95Ms: 90 },
		{ durationMs: 10 },
		{ completedAt: 'yesterday' },
		{ completedAt: '2020-01-01T00:00:00Z' },
	])
		assert.equal(
			isLoadTestResult({
				...result,
				...override,
			}),
			false,
		)
})
