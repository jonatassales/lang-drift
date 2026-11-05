import type { z } from 'zod'

import type { bookDemoErrorCodes } from './book-demo-error-codes.data'
import type {
	bookDemoSchema,
	bookDemoSubmissionSchema,
} from './form.validation'

export type BookDemoErrorCode = (typeof bookDemoErrorCodes)[number]
export type BookDemoInput = z.input<typeof bookDemoSchema>
export type BookDemoData = z.output<typeof bookDemoSchema>
export type BookDemoSubmission = z.output<typeof bookDemoSubmissionSchema>
export type BookDemoField = keyof BookDemoData
export type BookDemoFieldErrors = Partial<
	Record<BookDemoField, BookDemoErrorCode>
>
export type BookDemoOutcome = 'idle' | 'sent' | 'unavailable'
