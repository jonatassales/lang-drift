import type { StreamConnection } from '@langdrift/events/transport'
import { createStreamConnection } from '@langdrift/events/transport'

// The one place the log's address is read. Each call opens a connection
// that connects on its first command.
export function createEventLogConnectionOpener(): () => StreamConnection {
	const url = process.env.REDIS_URL
	if (!url) throw new Error('Missing REDIS_URL')

	return () => createStreamConnection(url)
}
