import { demoRequestStream } from 'streaming'
import { RedisStreamsClient } from 'streaming/transport'

import type { EventLogConnectionOpener } from './event-log-connection-opener.factory.js'

// Builds the producer that appends to the demo-requests stream; nothing
// connects before the first event.
export function createEventLogProducer(
	open: EventLogConnectionOpener,
): RedisStreamsClient {
	return new RedisStreamsClient({
		stream: demoRequestStream,
		connect: () => open.producer(),
	})
}
