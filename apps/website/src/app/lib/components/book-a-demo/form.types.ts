import type { z } from 'zod'

import type { bookDemoSchema } from './form.validation'

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
export type BookDemoField = keyof BookDemoData
