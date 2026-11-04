import type { ObjectStore } from '@app/archive/object-store.port.js'

export interface StoredObject {
	bucket: string
	key: string
	body: string
	size: number
	metadata: Record<string, string>
}

// Records what the archive asks of the object store.
export class FakeObjectStore implements ObjectStore {
	readonly buckets = new Set<string>()
	readonly objects: StoredObject[] = []
	readonly calls: string[] = []

	constructor(existingBuckets: string[] = []) {
		for (const bucket of existingBuckets) this.buckets.add(bucket)
	}

	async bucketExists(bucket: string): Promise<boolean> {
		this.calls.push(`bucketExists ${bucket}`)

		return this.buckets.has(bucket)
	}

	async makeBucket(bucket: string): Promise<void> {
		this.calls.push(`makeBucket ${bucket}`)
		this.buckets.add(bucket)
	}

	async putObject(
		bucket: string,
		key: string,
		body: Buffer,
		size: number,
		metadata: Record<string, string>,
	): Promise<void> {
		this.calls.push(`putObject ${bucket} ${key}`)
		this.objects.push({
			bucket,
			key,
			body: body.toString('utf8'),
			size,
			metadata,
		})
	}
}
