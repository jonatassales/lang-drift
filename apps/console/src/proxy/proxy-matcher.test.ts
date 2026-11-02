import assert from 'node:assert/strict'
import test from 'node:test'

import { proxyMatcher } from './proxy-matcher.data.ts'

test('proxy matcher matches application routes', () => {
	const regex = new RegExp(`^${proxyMatcher[0]}$`, 'i')

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
	const regex = new RegExp(`^${proxyMatcher[0]}$`, 'i')

	const exclusions = [
		'/_next/static/immutable/chunks/39hocmh4qdt17.css',
		'/_next/static/chunks/main.js',
		'/_next/image',
		'/api/console/home',
		'/api/extub/products',
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
