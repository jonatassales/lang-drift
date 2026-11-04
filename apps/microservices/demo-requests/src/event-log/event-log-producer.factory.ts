import { demoRequestStream } from '@langdrift/events'
import {
	createStreamConnection,
	RedisStreamsClient,
} from '@langdrift/events/transport'

// Builds the producer that appends to the demo-requests stream. The URL is
// read here, where it is needed, and nothing connects before the first event.
export function createEventLogProducer(): RedisStreamsClient {
	const url = process.env.REDIS_URL
	if (!url) throw new Error('Missing REDIS_URL')

	return new RedisStreamsClient({
		stream: demoRequestStream,
		connect: () => createStreamConnection(url),
	})
}
