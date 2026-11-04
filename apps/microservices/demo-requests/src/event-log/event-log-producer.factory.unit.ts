import { RedisStreamsClient } from '@langdrift/events/transport'
import { describe, expect, it } from 'vitest'

import { FakeReaderConnection } from '@app/event-log/testdata/stream-connection.mock.js'

import { createEventLogProducer } from './event-log-producer.factory.js'

describe('demo-requests/event-log/event-log-producer.factory', () => {
	it('builds a Redis Streams client for the demo-requests stream without opening a connection yet', () => {
		let opened = 0
		const producer = createEventLogProducer({
			consumer: () => new FakeReaderConnection(),
			producer: () => {
				opened += 1

				return new FakeReaderConnection()
			},
		})

		expect(producer).toBeInstanceOf(RedisStreamsClient)
		expect(opened).toBe(0)
	})
})
