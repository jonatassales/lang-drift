import { z } from 'zod'

import { demoRequestAggregateSchema } from './demo-request-aggregate.validation'
import { demoRequestEmailKindSchema } from './demo-request-email-kind.validation'
import { demoRequestEventNames } from './demo-request-event-names.data'
import { eventEnvelopeSchema } from './event-envelope.validation'

// The mailer delivered one of the request's emails: the record of what was
// sent, to whom and when, with the provider's message id and never a key.
export const demoRequestEmailSentSchema = eventEnvelopeSchema.extend({
	type: z.literal(demoRequestEventNames.emailSent),
	version: z.literal(1),
	aggregate: demoRequestAggregateSchema,
	payload: z
		.object({
			kind: demoRequestEmailKindSchema,
			provider: z.literal('resend'),
			messageId: z.string().min(1),
			to: z.email(),
			sentAt: z.iso.datetime(),
		})
		.strict(),
})
