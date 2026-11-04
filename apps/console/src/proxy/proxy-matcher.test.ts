import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import test from 'node:test'

function readProxyMatcher(): string {
	const proxySource = readFileSync(
		new URL('../proxy.ts', import.meta.url),
		'utf-8',
	)
	const match = proxySource.match(/matcher:\s*\[\s*['"]([^'"]+)['"],?\s*\]/)
	assert.ok(match?.[1], 'proxy.ts must export config.matcher')
	return match[1].replace(/\\\\/g, '\\')
}

test('proxy matcher matches application routes', () => {
	const pattern = readProxyMatcher()
	const regex = new RegExp(`^${pattern}$`, 'i')

	const routes = [
		'/',
		'/en',
		'/pt-br',
		'/zh-hant',
		'/ja',
		'/en/home',
		'/en/finops',
		'/en/products',
		'/en/settings',
		'/en/products/seed-microruntimes',
		'/en/runtimes/seed-runtime-1',
		'/en/products/seed-microruntimes/tickets/ticket-0042',
		'/login',
		'/sign-in',
		'/signup',
		'/en/login',
		'/en/sign-in',
	]

	for (const route of routes) {
		assert.equal(
			regex.test(route),
			true,
			`expected ${route} to match proxy matcher`,
		)
	}
})

test('proxy matcher excludes API, _next, and static asset files', () => {
	const pattern = readProxyMatcher()
	const regex = new RegExp(`^${pattern}$`, 'i')

	const exclusions = [
		'/_next/static/immutable/chunks/39hocmh4qdt17.css',
		'/_next/static/chunks/main.js',
		'/_next/image',
		'/api/console/home',
		'/api/console-stubs/products',
		'/favicon.ico',
		'/icon.svg',
		'/apple-icon.png',
		'/robots.txt',
		'/styles.css',
		'/bundle.js',
	]

	for (const path of exclusions) {
		assert.equal(
			regex.test(path),
			false,
			`expected ${path} to be excluded by proxy matcher`,
		)
	}
})
