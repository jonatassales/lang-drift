import { describe, expect, it } from 'vitest'

import { toSseEvent } from './to-sse-event.map.js'

describe('demo-requests/demo-requests/to-sse-event.map', () => {
	it('sends an entry with its stream id as the event id, its type as the event name and the envelope as data', () => {
		const envelope = {
			type: 'demo-request.received',
			aggregate: {
				type: 'demo-request',
				id: '701706be-7e70-4681-9608-826d3ad548b9',
			},
		}

		expect(
			toSseEvent({
				kind: 'entry',
				entry: {
					id: '1762178700000-0',
					type: 'demo-request.received',
					aggregateId: '701706be-7e70-4681-9608-826d3ad548b9',
					envelope,
				},
			}),
		).toEqual({
			id: '1762178700000-0',
			type: 'demo-request.received',
			data: envelope,
		})
	})

	it('sends a ping as an event of its own with no data', () => {
		expect(
			toSseEvent({
				kind: 'ping',
			}),
		).toEqual({
			type: 'ping',
			data: '',
		})
	})
})
