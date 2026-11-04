import assert from 'node:assert/strict'
import test from 'node:test'

import {
	consoleHostAllowed,
	consoleStubsBase,
	crossSite,
	serverRead,
} from './console-stubs-link.ts'

function browserCall(headers: Record<string, string>) {
	return new Request('http://127.0.0.1:3130/api/console/products', {
		method: 'POST',
		headers: {
			host: '127.0.0.1:3130',
			...headers,
		},
	})
}

test('without a token the console serves loopback hosts only', () => {
	assert.equal(consoleHostAllowed('localhost', undefined), true)
	assert.equal(consoleHostAllowed('127.0.0.1', ''), true)
	assert.equal(
		consoleHostAllowed('langdrift-console.vercel.app', undefined),
		false,
	)
})

test('a token of only whitespace keeps the local contract', () => {
	assert.equal(
		consoleHostAllowed('langdrift-console.vercel.app', '  \n'),
		false,
	)
	assert.equal(
		consoleStubsBase('https://langdrift-stubs.vercel.app', ' '),
		null,
	)
})

test('with a token the console serves its hosted host', () => {
	assert.equal(
		consoleHostAllowed('langdrift-console.vercel.app', 'shared-token'),
		true,
	)
})

test('without a token console-stubs must be plain http on loopback', () => {
	assert.equal(
		consoleStubsBase('http://127.0.0.1:3401', undefined)?.href,
		'http://127.0.0.1:3401/',
	)
	assert.equal(
		consoleStubsBase('http://localhost:3402', undefined)?.href,
		'http://localhost:3402/',
	)
	assert.equal(
		consoleStubsBase('https://langdrift-stubs.vercel.app', undefined),
		null,
	)
	assert.equal(consoleStubsBase('http://example.com', undefined), null)
})

test('with a token console-stubs may be https on any host, or http on loopback', () => {
	assert.equal(
		consoleStubsBase('https://langdrift-stubs.vercel.app', 'shared-token')
			?.href,
		'https://langdrift-stubs.vercel.app/',
	)
	assert.equal(
		consoleStubsBase('http://127.0.0.1:3489', 'shared-token')?.href,
		'http://127.0.0.1:3489/',
	)
	assert.equal(consoleStubsBase('http://example.com', 'shared-token'), null)
})

test('console-stubs never carries credentials in its URL', () => {
	assert.equal(
		consoleStubsBase('https://user:secret@example.com', 'shared-token'),
		null,
	)
	assert.equal(consoleStubsBase('http://user@127.0.0.1:3401', undefined), null)
	assert.throws(() => consoleStubsBase('not a url', 'shared-token'))
})

test('a browser call from another site is cross-site', () => {
	assert.equal(crossSite(browserCall({ origin: 'https://example.com' })), true)
	assert.equal(crossSite(browserCall({ 'sec-fetch-site': 'cross-site' })), true)
	assert.equal(
		crossSite(
			browserCall({
				origin: 'http://127.0.0.1:3130',
				'sec-fetch-site': 'same-origin',
			}),
		),
		false,
	)
})

test('a server-rendered read carries the incoming host alone', () => {
	const read = serverRead('127.0.0.1:3130', 'products/seed-microruntimes')

	assert.equal(
		read.url,
		'http://127.0.0.1:3130/api/console/products/seed-microruntimes',
	)
	assert.equal(read.method, 'GET')
	assert.deepEqual([...read.headers], [])
	assert.equal(crossSite(read), false)
	assert.equal(
		crossSite(serverRead('langdrift-console.vercel.app', 'home')),
		false,
	)
})
