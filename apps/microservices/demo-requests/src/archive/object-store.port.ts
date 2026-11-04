// What the archive needs from an object store, narrowed to the three calls
// it makes. The MinIO client satisfies it; tests use a fake.
export interface ObjectStore {
	bucketExists(bucket: string): Promise<boolean>
	makeBucket(bucket: string): Promise<void>
	putObject(
		bucket: string,
		key: string,
		body: Buffer,
		size: number,
		metadata: Record<string, string>,
	): Promise<unknown>
}
