import assert from 'node:assert/strict'
import test from 'node:test'

import { consolePathSafe, isEntityId } from './entity-id.ts'

test('the ids extub seeds and generates are entity ids', () => {
	for (const id of [
		'seed-microruntimes',
		'seed-product-2',
		'seed-runtime-4',
		'seed-secret-1',
		'seed-runtime-1-execution-1',
		'ticket-0042',
		'snake_case',
		'x'.repeat(128),
	])
		assert.equal(isEntityId(id), true, id)
})

test('dots, slashes, escapes, empty and over-long values are not', () => {
	for (const id of [
		'.',
		'..',
		'a.b',
		'/',
		'a/b',
		'%',
		'%2F',
		'..%2F',
		'',
		'x'.repeat(129),
		'seed runtime',
		'seed-runtime-1\n',
		'a?b',
		'a#b',
		'ação',
	])
		assert.equal(isEntityId(id), false, JSON.stringify(id))
})

test('a console path is safe when each segment is a name or an entity id', () => {
	assert.equal(consolePathSafe('home'), true)
	assert.equal(consolePathSafe('runtimes/seed-runtime-4/load-test'), true)
	assert.equal(
		consolePathSafe(
			'products/seed-microruntimes/tickets?reached=plan&cursor=a%2Fb',
		),
		true,
	)
	assert.equal(consolePathSafe('products/..%2F/secrets'), false)
	assert.equal(consolePathSafe('products/a.b'), false)
	assert.equal(consolePathSafe('products//secrets'), false)
	assert.equal(consolePathSafe('tickets/ticket-0042%3Freached=plan'), false)
})
