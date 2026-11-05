import { describe, expect, it } from 'vitest'

import { FakeObjectStore } from '@app/archive/testdata/object-store.mock.js'

import { ArchiveService } from './archive.service.js'

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
	payload: {
		company: 'Acme',
	},
}

describe('demo-requests/archive/archive.service', () => {
	it('creates the bucket on init when it is absent', async () => {
		const store = new FakeObjectStore()
		const archive = new ArchiveService(store, 'demo-requests')

		await archive.onModuleInit()

		expect(store.calls).toEqual([
			'bucketExists demo-requests',
			'makeBucket demo-requests',
		])
	})

	it('leaves the bucket alone on init when it exists', async () => {
		const store = new FakeObjectStore(['demo-requests'])
		const archive = new ArchiveService(store, 'demo-requests')

		await archive.onModuleInit()

		expect(store.calls).toEqual(['bucketExists demo-requests'])
	})

	it('stores an envelope under its aggregate and entry as JSON with the event type as metadata', async () => {
		const store = new FakeObjectStore(['demo-requests'])
		const archive = new ArchiveService(store, 'demo-requests')

		await archive.store(envelope, '1762178700000-0')

		expect(store.objects).toEqual([
			{
				bucket: 'demo-requests',
				key: 'demo-request/701706be-7e70-4681-9608-826d3ad548b9/1762178700000-0.json',
				body: JSON.stringify(envelope),
				size: Buffer.byteLength(JSON.stringify(envelope)),
				metadata: {
					'Content-Type': 'application/json',
					'X-Amz-Meta-Event-Type': 'demo-request.received',
				},
			},
		])
	})

	it('refuses an envelope that is not one, so the entry stays pending', async () => {
		const store = new FakeObjectStore(['demo-requests'])
		const archive = new ArchiveService(store, 'demo-requests')

		await expect(
			archive.store(
				{
					type: 'demo-request.received',
				},
				'1762178700000-0',
			),
		).rejects.toThrow(/envelope/)

		expect(store.objects).toEqual([])
	})
})
