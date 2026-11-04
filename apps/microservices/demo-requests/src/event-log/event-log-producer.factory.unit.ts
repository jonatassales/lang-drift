import { RedisStreamsClient } from '@langdrift/events/transport'
import { afterEach, describe, expect, it, vi } from 'vitest'

import { createEventLogProducer } from './event-log-producer.factory.js'

describe('demo-requests/event-log/event-log-producer.factory', () => {
	afterEach(() => {
		vi.unstubAllEnvs()
	})

	it('throws when REDIS_URL is not set', () => {
		vi.stubEnv('REDIS_URL', '')

		expect(() => createEventLogProducer()).toThrow(/REDIS_URL/)
	})

	it('builds a Redis Streams client for the demo-requests stream without connecting yet', () => {
		vi.stubEnv('REDIS_URL', 'redis://127.0.0.1:6379')

		expect(createEventLogProducer()).toBeInstanceOf(RedisStreamsClient)
	})
})
