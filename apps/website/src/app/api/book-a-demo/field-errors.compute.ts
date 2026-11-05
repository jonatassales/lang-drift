import type { z } from 'zod'

import type {
	BookDemoErrorCode,
	BookDemoField,
} from '@routes/book-a-demo/form.types'

import { isBookDemoErrorCode } from './is-book-demo-error-code.validation'

// The first code of each field that failed.
export function fieldErrorsOf(
	error: z.ZodError,
): Partial<Record<BookDemoField, BookDemoErrorCode>> {
	const errors: Partial<Record<BookDemoField, BookDemoErrorCode>> = {}
	for (const issue of error.issues) {
		const field = issue.path[0]
		if (field !== 'email' && field !== 'company' && field !== 'position') {
			continue
		}
		if (errors[field] || !isBookDemoErrorCode(issue.message)) continue
		errors[field] = issue.message
	}

	return errors
}
