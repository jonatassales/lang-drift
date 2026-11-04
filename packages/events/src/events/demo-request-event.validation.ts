import { z } from 'zod'

import { demoRequestEmailFailedSchema } from './demo-request-email-failed.validation'
import { demoRequestEmailSentSchema } from './demo-request-email-sent.validation'
import { demoRequestReceivedSchema } from './demo-request-received.validation'

// Any event of a demo request, told apart by its type. A new version of a
// type is a new schema and a new member here.
export const demoRequestEventSchema = z.discriminatedUnion('type', [
	demoRequestReceivedSchema,
	demoRequestEmailSentSchema,
	demoRequestEmailFailedSchema,
])
