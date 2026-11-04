import { Redis } from 'ioredis'

import type { StreamConnection } from './stream-connection.types.js'

// Opens the connection a producer appends through. Nothing is queued while
// Redis is away: an append either lands now or fails now, so a caller can
// tell its own caller that nothing was stored. The producer connects
// explicitly before its first append.
export function createProducerConnection(url: string): StreamConnection {
	return new Redis(url, {
		connectTimeout: 2000,
		enableOfflineQueue: false,
		lazyConnect: true,
		maxRetriesPerRequest: 1,
	})
}
