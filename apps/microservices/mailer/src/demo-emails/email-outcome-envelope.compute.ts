import type {
	DemoRequestEmailKind,
	DemoRequestEvent,
	DemoRequestReceived,
} from '@langdrift/events'
import {
	createEnvelope,
	demoRequestEventNames,
	demoRequestEventSchema,
} from '@langdrift/events'

export interface SentOutcome {
	kind: DemoRequestEmailKind
	outcome: 'sent'
	messageId: string
	to: string
	now: Date
}

export interface FailedOutcome {
	kind: DemoRequestEmailKind
	outcome: 'failed'
	reason: string
	now: Date
}

export type EmailOutcome = SentOutcome | FailedOutcome

// The record of one email's fate, caused by the received event: what was
// sent, to whom and when with the provider's id, or why it failed.
export function emailOutcomeEnvelope(
	received: DemoRequestReceived,
	outcome: EmailOutcome,
): DemoRequestEvent {
	const common = {
		version: 1,
		aggregate: received.aggregate,
		causationId: received.id,
		now: outcome.now,
	}

	if (outcome.outcome === 'sent') {
		return demoRequestEventSchema.parse(
			createEnvelope({
				...common,
				type: demoRequestEventNames.emailSent,
				payload: {
					kind: outcome.kind,
					provider: 'resend',
					messageId: outcome.messageId,
					to: outcome.to,
					sentAt: outcome.now.toISOString(),
				},
			}),
		)
	}

	return demoRequestEventSchema.parse(
		createEnvelope({
			...common,
			type: demoRequestEventNames.emailFailed,
			payload: {
				kind: outcome.kind,
				provider: 'resend',
				reason: outcome.reason,
				failedAt: outcome.now.toISOString(),
			},
		}),
	)
}
