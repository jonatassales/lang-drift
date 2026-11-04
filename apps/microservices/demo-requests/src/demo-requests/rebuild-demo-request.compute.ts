import type { DemoRequestEvent } from '@langdrift/events'
import { demoRequestEventNames } from '@langdrift/events'

import type {
	DemoRequestEmail,
	DemoRequestState,
	DemoRequestView,
} from './demo-request-view.types.js'

// Folds a request's events, in stream order, into its current state. The
// confirmation decides the state; the lead notification is only recorded.
export function rebuildDemoRequest(
	events: DemoRequestEvent[],
): DemoRequestView | null {
	const received = events.find(
		(event) => event.type === demoRequestEventNames.received,
	)
	if (!received || received.type !== demoRequestEventNames.received) {
		return null
	}

	let state: DemoRequestState = 'received'
	const emails: DemoRequestEmail[] = []
	for (const event of events) {
		if (event.type === demoRequestEventNames.emailSent) {
			emails.push({
				kind: event.payload.kind,
				outcome: 'sent',
				provider: event.payload.provider,
				messageId: event.payload.messageId,
				at: event.payload.sentAt,
			})
			if (event.payload.kind === 'confirmation') state = 'confirmed'
		}

		if (event.type === demoRequestEventNames.emailFailed) {
			emails.push({
				kind: event.payload.kind,
				outcome: 'failed',
				provider: event.payload.provider,
				reason: event.payload.reason,
				at: event.payload.failedAt,
			})
			const confirmationFailed =
				event.payload.kind === 'confirmation' && state !== 'confirmed'
			if (confirmationFailed) state = 'failed'
		}
	}

	return {
		id: received.aggregate.id,
		state,
		submission: received.payload,
		emails,
	}
}
