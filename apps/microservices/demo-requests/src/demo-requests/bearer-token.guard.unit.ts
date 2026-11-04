import type { ExecutionContext } from '@nestjs/common'
import { UnauthorizedException } from '@nestjs/common'
import { afterEach, describe, expect, it, vi } from 'vitest'

import { BearerTokenGuard } from './bearer-token.guard.js'

const token = 'a-demo-requests-token-with-thirty-two-plus-characters'

function contextWith(authorization?: string): ExecutionContext {
	const request = {
		headers: authorization ? { authorization } : {},
	}

	return {
		switchToHttp: () => ({
			getRequest: () => request,
		}),
	} as unknown as ExecutionContext
}

describe('demo-requests/demo-requests/bearer-token.guard', () => {
	afterEach(() => {
		vi.unstubAllEnvs()
	})

	it('lets a request with the configured bearer token through', () => {
		vi.stubEnv('DEMO_REQUESTS_TOKEN', token)
		const guard = new BearerTokenGuard()

		expect(guard.canActivate(contextWith(`Bearer ${token}`))).toBe(true)
	})

	it('refuses a request without an authorization header', () => {
		vi.stubEnv('DEMO_REQUESTS_TOKEN', token)
		const guard = new BearerTokenGuard()

		expect(() => guard.canActivate(contextWith())).toThrow(
			UnauthorizedException,
		)
	})

	it('refuses a wrong token, a token of another length and a scheme that is not Bearer', () => {
		vi.stubEnv('DEMO_REQUESTS_TOKEN', token)
		const guard = new BearerTokenGuard()

		expect(() =>
			guard.canActivate(contextWith(`Bearer ${token.replace('a', 'b')}`)),
		).toThrow(UnauthorizedException)
		expect(() => guard.canActivate(contextWith('Bearer short'))).toThrow(
			UnauthorizedException,
		)
		expect(() => guard.canActivate(contextWith(`Basic ${token}`))).toThrow(
			UnauthorizedException,
		)
	})

	it('fails closed when DEMO_REQUESTS_TOKEN is not set', () => {
		vi.stubEnv('DEMO_REQUESTS_TOKEN', '')
		const guard = new BearerTokenGuard()

		expect(() => guard.canActivate(contextWith(`Bearer ${token}`))).toThrow(
			/DEMO_REQUESTS_TOKEN/,
		)
	})
})
