import type { StreamConnection } from 'infra/streaming/transport'
import {
	createProducerConnection,
	createStreamConnection,
} from 'infra/streaming/transport'

// A consumer's connection waits for Redis to come back; a producer's fails
// at once, so an append either lands now or is reported as not stored.
export interface EventLogConnectionOpener {
	consumer(): StreamConnection
	producer(): StreamConnection
}

// The one place the log's address is read. Each call opens a connection
// that connects on its first command or, for a producer, when asked to.
export function createEventLogConnectionOpener(): EventLogConnectionOpener {
	const url = process.env.REDIS_URL
	if (!url) throw new Error('Missing REDIS_URL')

	return {
		consumer: () => createStreamConnection(url),
		producer: () => createProducerConnection(url),
	}
}
