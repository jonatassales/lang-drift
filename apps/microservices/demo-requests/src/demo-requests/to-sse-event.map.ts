import type { MessageEvent } from '@nestjs/common'

import type { FollowEvent } from '@app/event-log/event-log-reader.service.js'

// One follow event as the client receives it: an entry keeps its stream id
// as the event id, so a reconnecting client resumes with Last-Event-ID.
export function toSseEvent(event: FollowEvent): MessageEvent {
	if (event.kind === 'ping') {
		return {
			type: 'ping',
			data: '',
		}
	}

	return {
		id: event.entry.id,
		type: event.entry.type,
		data: event.entry.envelope as object,
	}
}
