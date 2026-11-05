import { fieldErrorsOf } from '@api/book-a-demo/field-errors.compute'
import { bookDemoSchema } from '@routes/book-a-demo/form.validation'

export function errorsOf(input: Record<string, unknown>) {
	const result = bookDemoSchema.safeParse(input)

	if (result.success) return {}

	return fieldErrorsOf(result.error)
}
