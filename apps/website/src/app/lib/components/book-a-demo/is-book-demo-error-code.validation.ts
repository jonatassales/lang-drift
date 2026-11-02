import { type BookDemoErrorCode, bookDemoErrorCodes } from './form.types'

export function isBookDemoErrorCode(
	value: unknown,
): value is BookDemoErrorCode {
	return bookDemoErrorCodes.some((code) => code === value)
}
