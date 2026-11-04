import { afterEach, describe, expect, it, vi } from 'vitest'

import { createEventLogConnectionOpener } from './event-log-connection-opener.factory.js'

describe('demo-requests/event-log/event-log-connection-opener.factory', () => {
	afterEach(() => {
		vi.unstubAllEnvs()
	})

	it('throws when REDIS_URL is not set', () => {
		vi.stubEnv('REDIS_URL', '')

		expect(() => createEventLogConnectionOpener()).toThrow(/REDIS_URL/)
	})

	it('opens a connection that can read the stream without connecting yet', () => {
		vi.stubEnv('REDIS_URL', 'redis://127.0.0.1:6379')

		const connection = createEventLogConnectionOpener()()

		expect(typeof connection.xreadgroup).toBe('function')
		expect(typeof connection.xrange).toBe('function')
		connection.disconnect()
	})
})
