import {
	type BookDemoErrorCode,
	bookDemoErrorCodes,
} from '@routes/book-a-demo/form.types'

export function isBookDemoErrorCode(
	value: unknown,
): value is BookDemoErrorCode {
	return bookDemoErrorCodes.some((code) => code === value)
}
