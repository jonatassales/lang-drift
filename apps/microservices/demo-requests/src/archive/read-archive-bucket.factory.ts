// The bucket the archive writes into, read where it is needed.
export function readArchiveBucket(): string {
	const bucket = process.env.MINIO_BUCKET
	if (!bucket) throw new Error('Missing MINIO_BUCKET')

	return bucket
}
