import { describe, expect, it } from 'vitest'

import { parseStreamEntry } from './parse-stream-entry.compute.js'

const envelope = {
	id: '0d8ce74c-4756-4bfd-8d35-c820e302898e',
	type: 'demo-request.received',
	version: 1,
	aggregate: {
		type: 'demo-request',
		id: '701706be-7e70-4681-9608-826d3ad548b9',
	},
	occurredAt: '2026-11-03T14:05:00.000Z',
	correlationId: '701706be-7e70-4681-9608-826d3ad548b9',
	payload: {},
}

describe('infra/streaming/transport/parse-stream-entry.compute', () => {
	it('reads the type, the aggregate id and the envelope out of the entry fields', () => {
		const entry = parseStreamEntry('1762178700000-0', [
			'type',
			'demo-request.received',
			'aggregateId',
			'701706be-7e70-4681-9608-826d3ad548b9',
			'envelope',
			JSON.stringify(envelope),
		])

		expect(entry).toEqual({
			id: '1762178700000-0',
			type: 'demo-request.received',
			aggregateId: '701706be-7e70-4681-9608-826d3ad548b9',
			envelope,
		})
	})

	it('refuses an entry whose envelope is not JSON', () => {
		expect(() =>
			parseStreamEntry('1762178700000-0', [
				'type',
				'demo-request.received',
				'aggregateId',
				'701706be-7e70-4681-9608-826d3ad548b9',
				'envelope',
				'{not json',
			]),
		).toThrow(/envelope/)
	})

	it('refuses an entry without a type field', () => {
		expect(() =>
			parseStreamEntry('1762178700000-0', [
				'aggregateId',
				'701706be-7e70-4681-9608-826d3ad548b9',
				'envelope',
				JSON.stringify(envelope),
			]),
		).toThrow(/type/)
	})

	it('refuses an entry whose fields are missing', () => {
		expect(() => parseStreamEntry('1762178700000-0', null)).toThrow(/fields/)
	})
})
