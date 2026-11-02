import type { BookDemoErrorCode, BookDemoField } from './form.types'

export function reply(
	status: number,
	code:
		| 'INVALID_ORIGIN'
		| 'INVALID_REQUEST'
		| 'VALIDATION_ERROR'
		| 'UNAVAILABLE',
	fieldErrors?: Partial<Record<BookDemoField, BookDemoErrorCode>>,
) {
	const headers: Record<string, string> = {
		'Cache-Control': 'no-store',
		'X-Content-Type-Options': 'nosniff',
	}
	if (status === 503) headers['Retry-After'] = '30'
	const hasFieldErrors = fieldErrors && Object.keys(fieldErrors).length > 0

	return Response.json(
		{
			ok: false,
			code,
			...(hasFieldErrors ? { fieldErrors } : {}),
		},
		{
			headers,
			status,
		},
	)
}
