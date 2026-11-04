import type {
	DemoRequestEmailKind,
	DemoRequestReceived,
} from '@langdrift/events'

export type DemoRequestState = 'received' | 'confirmed' | 'failed'

export interface DemoRequestEmailSentView {
	kind: DemoRequestEmailKind
	outcome: 'sent'
	provider: 'resend'
	messageId: string
	at: string
}

export interface DemoRequestEmailFailedView {
	kind: DemoRequestEmailKind
	outcome: 'failed'
	provider: 'resend'
	reason: string
	at: string
}

export type DemoRequestEmail =
	| DemoRequestEmailSentView
	| DemoRequestEmailFailedView

// One demo request as its events tell it: what was submitted, whether the
// visitor was confirmed, and every email the mailer sent or failed to send.
export interface DemoRequestView {
	id: string
	state: DemoRequestState
	submission: DemoRequestReceived['payload']
	emails: DemoRequestEmail[]
}
