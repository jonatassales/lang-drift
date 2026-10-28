import assert from 'node:assert/strict'
import test from 'node:test'

import { areaHome, areaOf } from './area.ts'

test('routes under /finops belong to FinOps', () => {
	assert.equal(areaOf('/finops'), 'finops')
	assert.equal(areaOf('/finops/anything'), 'finops')
})

test('every other console route belongs to Microruntimes', () => {
	for (const pathname of [
		'/home',
		'/products',
		'/products/seed-microruntimes/executions',
		'/runtimes/seed-runtime-1',
		'/settings',
		'/finopsish',
	]) {
		assert.equal(areaOf(pathname), 'microruntimes', pathname)
	}
})

test('each area opens on a route of its own', () => {
	assert.equal(areaOf(areaHome.microruntimes), 'microruntimes')
	assert.equal(areaOf(areaHome.finops), 'finops')
})
