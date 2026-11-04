import type { z } from 'zod'

import type { demoRequestEmailFailedSchema } from './demo-request-email-failed.validation.js'
import type { demoRequestEmailKindSchema } from './demo-request-email-kind.validation.js'
import type { demoRequestEmailSentSchema } from './demo-request-email-sent.validation.js'
import type { demoRequestEventSchema } from './demo-request-event.validation.js'
import type { demoRequestReceivedSchema } from './demo-request-received.validation.js'
import type { demoRequestSubmissionSchema } from './demo-request-submission.validation.js'

export type DemoRequestSubmission = z.output<typeof demoRequestSubmissionSchema>
export type DemoRequestEmailKind = z.output<typeof demoRequestEmailKindSchema>
export type DemoRequestReceived = z.output<typeof demoRequestReceivedSchema>
export type DemoRequestEmailSent = z.output<typeof demoRequestEmailSentSchema>
export type DemoRequestEmailFailed = z.output<
	typeof demoRequestEmailFailedSchema
>
export type DemoRequestEvent = z.output<typeof demoRequestEventSchema>
