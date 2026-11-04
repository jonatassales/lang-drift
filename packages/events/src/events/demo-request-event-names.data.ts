// The event types of a demo request, in version 1. A handler subscribes to
// one of these names and a schema refines the envelope for it.
export const demoRequestEventNames = {
	received: 'demo-request.received',
	emailSent: 'demo-request.email-sent',
	emailFailed: 'demo-request.email-failed',
} as const
