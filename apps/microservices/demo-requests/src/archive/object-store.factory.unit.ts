import { afterEach, describe, expect, it, vi } from 'vitest'

import { createObjectStore } from './object-store.factory.js'

function stubMinio() {
	vi.stubEnv('MINIO_ENDPOINT', '127.0.0.1')
	vi.stubEnv('MINIO_PORT', '9000')
	vi.stubEnv('MINIO_USE_SSL', 'false')
	vi.stubEnv('MINIO_ACCESS_KEY', 'archive')
	vi.stubEnv('MINIO_SECRET_KEY', 'archive-secret-of-enough-length')
}

describe('demo-requests/archive/object-store.factory', () => {
	afterEach(() => {
		vi.unstubAllEnvs()
	})

	it('builds a client from the MinIO keys without connecting yet', () => {
		stubMinio()

		const store = createObjectStore()

		expect(typeof store.bucketExists).toBe('function')
		expect(typeof store.putObject).toBe('function')
	})

	it('throws naming the missing key', () => {
		stubMinio()
		vi.stubEnv('MINIO_SECRET_KEY', '')

		expect(() => createObjectStore()).toThrow(/MINIO_SECRET_KEY/)
	})

	it('refuses a port that is not a number and a TLS flag that is not true or false', () => {
		stubMinio()
		vi.stubEnv('MINIO_PORT', 'nine thousand')

		expect(() => createObjectStore()).toThrow(/MINIO_PORT/)

		stubMinio()
		vi.stubEnv('MINIO_USE_SSL', 'maybe')

		expect(() => createObjectStore()).toThrow(/MINIO_USE_SSL/)
	})
})
