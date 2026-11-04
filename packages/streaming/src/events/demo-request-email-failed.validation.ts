import { z } from 'zod'

import { demoRequestAggregateSchema } from './demo-request-aggregate.validation.js'
import { demoRequestEmailKindSchema } from './demo-request-email-kind.validation.js'
import { demoRequestEventNames } from './demo-request-event-names.data.js'
import { eventEnvelopeSchema } from './event-envelope.validation.js'

// The mailer could not deliver one of the request's emails; the reason is the
// provider's answer, with no secret in it.
export const demoRequestEmailFailedSchema = eventEnvelopeSchema.extend({
	type: z.literal(demoRequestEventNames.emailFailed),
	version: z.literal(1),
	aggregate: demoRequestAggregateSchema,
	payload: z
		.object({
			kind: demoRequestEmailKindSchema,
			provider: z.literal('resend'),
			reason: z.string().min(1).max(1000),
			failedAt: z.iso.datetime(),
		})
		.strict(),
})
