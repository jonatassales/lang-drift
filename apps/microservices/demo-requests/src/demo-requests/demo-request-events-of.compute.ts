import type { DemoRequestEvent } from 'infra/streaming'
import { demoRequestEventSchema } from 'infra/streaming'
import type { StreamEntry } from 'infra/streaming/transport'

// The demo-request events among a request's entries, in stream order; an
// envelope of another shape is left out of the fold.
export function demoRequestEventsOf(
	entries: StreamEntry[],
): DemoRequestEvent[] {
	const events: DemoRequestEvent[] = []
	for (const entry of entries) {
		const parsed = demoRequestEventSchema.safeParse(entry.envelope)
		if (parsed.success) events.push(parsed.data)
	}

	return events
}
