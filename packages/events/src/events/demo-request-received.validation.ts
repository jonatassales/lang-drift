import { z } from 'zod'

import { demoRequestAggregateSchema } from './demo-request-aggregate.validation.js'
import { demoRequestEventNames } from './demo-request-event-names.data.js'
import { demoRequestSubmissionSchema } from './demo-request-submission.validation.js'
import { eventEnvelopeSchema } from './event-envelope.validation.js'

// The first event of a demo request: what the website forwarded, and when.
export const demoRequestReceivedSchema = eventEnvelopeSchema.extend({
	type: z.literal(demoRequestEventNames.received),
	version: z.literal(1),
	aggregate: demoRequestAggregateSchema,
	payload: demoRequestSubmissionSchema
		.extend({
			submittedAt: z.iso.datetime(),
		})
		.strict(),
})
