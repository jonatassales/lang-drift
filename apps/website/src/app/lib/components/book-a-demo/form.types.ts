import type { z } from 'zod'

import type {
	bookDemoSchema,
	bookDemoSubmissionSchema,
} from './form.validation'

export const bookDemoErrorCodes = [
	'emailRequired',
	'emailInvalid',
	'emailTooLong',
	'emailFreeMail',
	'companyRequired',
	'companyTooLong',
	'positionRequired',
] as const

export type BookDemoErrorCode = (typeof bookDemoErrorCodes)[number]
export type BookDemoInput = z.input<typeof bookDemoSchema>
export type BookDemoData = z.output<typeof bookDemoSchema>
export type BookDemoSubmission = z.output<typeof bookDemoSubmissionSchema>
export type BookDemoField = keyof BookDemoData
