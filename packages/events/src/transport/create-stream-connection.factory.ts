import { Redis } from 'ioredis'

import type { StreamConnection } from './stream-connection.types.js'

// Opens one Redis connection for the transport. Commands wait for the
// connection instead of failing fast, and nothing connects before the first
// command, so a consumer can be constructed without Redis up yet.
export function createStreamConnection(url: string): StreamConnection {
	return new Redis(url, {
		lazyConnect: true,
		maxRetriesPerRequest: null,
	})
}
