import { demoRequestStream } from '@langdrift/events'
import type { StreamConnection } from '@langdrift/events/transport'
import { RedisStreamsClient } from '@langdrift/events/transport'

// Builds the producer that appends to the demo-requests stream; nothing
// connects before the first event.
export function createEventLogProducer(
	open: () => StreamConnection,
): RedisStreamsClient {
	return new RedisStreamsClient({
		stream: demoRequestStream,
		connect: open,
	})
}
