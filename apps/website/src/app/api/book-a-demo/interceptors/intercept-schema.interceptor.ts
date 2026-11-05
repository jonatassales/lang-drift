import { bookDemoReply } from '@api/book-a-demo/book-demo-reply.response'
import { fieldErrorsOf } from '@api/book-a-demo/field-errors-of.compute'
import type { BookDemoSubmission } from '@routes/book-a-demo/form.types'
import { bookDemoSubmissionSchema } from '@routes/book-a-demo/form.validation'

export interface AcceptedSubmission {
	submission: BookDemoSubmission
}

// The form's own schema decides, so the edge refuses exactly what the form
// refuses and names the field that failed.
export function interceptSchema(body: unknown): Response | AcceptedSubmission {
	const result = bookDemoSubmissionSchema.safeParse(body)

	if (!result.success) {
		return bookDemoReply(400, 'VALIDATION_ERROR', fieldErrorsOf(result.error))
	}

	return {
		submission: result.data,
	}
}
