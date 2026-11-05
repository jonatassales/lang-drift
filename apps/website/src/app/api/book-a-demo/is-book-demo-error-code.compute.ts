import { bookDemoErrorCodes } from '@routes/book-a-demo/book-demo-error-codes.data'
import type { BookDemoErrorCode } from '@routes/book-a-demo/form.types'

export function isBookDemoErrorCode(
	value: unknown,
): value is BookDemoErrorCode {
	return bookDemoErrorCodes.some((code) => code === value)
}
