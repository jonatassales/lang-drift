import type {
	BookDemoErrorCode,
	BookDemoField,
} from '@routes/book-a-demo/form.types'

export type BookDemoReplyCode =
	| 'OK'
	| 'INVALID_ORIGIN'
	| 'INVALID_REQUEST'
	| 'VALIDATION_ERROR'
	| 'UNAVAILABLE'

export function reply(
	status: number,
	code: BookDemoReplyCode,
	fieldErrors?: Partial<Record<BookDemoField, BookDemoErrorCode>>,
) {
	const headers: Record<string, string> = {
		'Cache-Control': 'no-store',
		'X-Content-Type-Options': 'nosniff',
	}
	if (status === 503) headers['Retry-After'] = '30'
	const hasFieldErrors = fieldErrors && Object.keys(fieldErrors).length > 0
	const ok = status >= 200 && status < 300

	return Response.json(
		{
			code,
			ok,
			...(hasFieldErrors ? { fieldErrors } : {}),
		},
		{
			headers,
			status,
		},
	)
}
