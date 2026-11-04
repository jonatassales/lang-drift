import { Client } from 'minio'

import type { ObjectStore } from './object-store.port.js'

// Builds the MinIO client from the environment, read here where it is
// needed. Nothing connects before the first call.
export function createObjectStore(): ObjectStore {
	const endPoint = process.env.MINIO_ENDPOINT
	const port = process.env.MINIO_PORT
	const useSSL = process.env.MINIO_USE_SSL
	const accessKey = process.env.MINIO_ACCESS_KEY
	const secretKey = process.env.MINIO_SECRET_KEY

	if (!endPoint) throw new Error('Missing MINIO_ENDPOINT')
	if (!port || !/^\d+$/.test(port)) throw new Error('Invalid MINIO_PORT')
	if (useSSL !== 'true' && useSSL !== 'false') {
		throw new Error('Invalid MINIO_USE_SSL: expected true or false')
	}
	if (!accessKey) throw new Error('Missing MINIO_ACCESS_KEY')
	if (!secretKey) throw new Error('Missing MINIO_SECRET_KEY')

	return new Client({
		endPoint,
		port: Number(port),
		useSSL: useSSL === 'true',
		accessKey,
		secretKey,
	})
}
