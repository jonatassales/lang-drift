import { bookDemoErrorCodes } from './book-demo-error-codes.data'
import type { BookDemoErrorCode } from './form.types'

export function isBookDemoErrorCode(
	value: unknown,
): value is BookDemoErrorCode {
	return bookDemoErrorCodes.some((code) => code === value)
}
