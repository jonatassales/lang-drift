import type { OnModuleInit } from '@nestjs/common'
import { Inject, Injectable } from '@nestjs/common'

import { eventEnvelopeSchema } from 'streaming'

import { archiveBucketToken } from './archive-bucket.token.js'
import { objectKeyOf } from './object-key.compute.js'
import type { ObjectStore } from './object-store.port.js'
import { objectStoreToken } from './object-store.token.js'

// Copies every event of the stream into the bucket as one object, so a
// trimmed stream never loses a request's history. The bucket is ensured
// once, when the module starts.
@Injectable()
export class ArchiveService implements OnModuleInit {
	constructor(
		@Inject(objectStoreToken) private readonly objectStore: ObjectStore,
		@Inject(archiveBucketToken) private readonly bucket: string,
	) {}

	async onModuleInit(): Promise<void> {
		const exists = await this.objectStore.bucketExists(this.bucket)

		if (exists) return

		await this.objectStore.makeBucket(this.bucket)
	}

	async store(envelope: unknown, streamId: string): Promise<void> {
		const parsed = eventEnvelopeSchema.safeParse(envelope)
		if (!parsed.success) {
			throw new Error(`Entry ${streamId} does not carry an event envelope`)
		}

		const body = Buffer.from(JSON.stringify(envelope), 'utf8')

		await this.objectStore.putObject(
			this.bucket,
			objectKeyOf(parsed.data.aggregate.id, streamId),
			body,
			body.byteLength,
			{
				'Content-Type': 'application/json',
				'X-Amz-Meta-Event-Type': parsed.data.type,
			},
		)
	}
}
