import type { z } from 'zod'

import type { BookDemoFieldErrors } from '@routes/book-a-demo/form.types'

import { isBookDemoErrorCode } from './is-book-demo-error-code.validation'

// The first code of each field that failed.
export function fieldErrorsOf(error: z.ZodError): BookDemoFieldErrors {
	const errors: BookDemoFieldErrors = {}

	for (const issue of error.issues) {
		const field = issue.path[0]

		const knownField =
			field === 'email' || field === 'company' || field === 'position'

		if (!knownField) continue

		const firstCodeOfField =
			!errors[field] && isBookDemoErrorCode(issue.message)

		if (!firstCodeOfField) continue

		errors[field] = issue.message
	}

	return errors
}
