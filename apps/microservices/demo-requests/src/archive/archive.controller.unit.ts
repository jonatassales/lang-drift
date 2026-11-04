import { RedisStreamsContext } from '@langdrift/events/transport'
import { describe, expect, it } from 'vitest'

import { ArchiveController } from './archive.controller.js'
import type { ArchiveService } from './archive.service.js'

class FakeArchive {
	readonly stored: [unknown, string][] = []

	async store(envelope: unknown, streamId: string): Promise<void> {
		this.stored.push([envelope, streamId])
	}
}

describe('demo-requests/archive/archive.controller', () => {
	it('stores every entry it is handed with the entry id from the context', async () => {
		const archive = new FakeArchive()
		const controller = new ArchiveController(
			archive as unknown as ArchiveService,
		)
		const context = new RedisStreamsContext([
			'demo-requests:events',
			'1762178700000-0',
			'archive',
			1,
		])

		await controller.archive(
			{
				type: 'demo-request.received',
			},
			context,
		)

		expect(archive.stored).toEqual([
			[
				{
					type: 'demo-request.received',
				},
				'1762178700000-0',
			],
		])
	})
})
