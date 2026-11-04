import type { z } from 'zod'

import type { demoRequestEmailFailedSchema } from './demo-request-email-failed.validation'
import type { demoRequestEmailKindSchema } from './demo-request-email-kind.validation'
import type { demoRequestEmailSentSchema } from './demo-request-email-sent.validation'
import type { demoRequestEventSchema } from './demo-request-event.validation'
import type { demoRequestReceivedSchema } from './demo-request-received.validation'
import type { demoRequestSubmissionSchema } from './demo-request-submission.validation'

export type DemoRequestSubmission = z.output<typeof demoRequestSubmissionSchema>
export type DemoRequestEmailKind = z.output<typeof demoRequestEmailKindSchema>
export type DemoRequestReceived = z.output<typeof demoRequestReceivedSchema>
export type DemoRequestEmailSent = z.output<typeof demoRequestEmailSentSchema>
export type DemoRequestEmailFailed = z.output<
	typeof demoRequestEmailFailedSchema
>
export type DemoRequestEvent = z.output<typeof demoRequestEventSchema>
