import 'server-only'

import {
	BodyError,
	readBoundedBody,
} from '@components/early-access/request.body'
import { parseSubmissionJSON } from '@components/early-access/request.validation'

import {
	type BookDemoErrorCode,
	type BookDemoField,
	bookDemoSubmissionSchema,
	fieldErrorsOf,
} from './form.schema'

const bodyLimit = 4096

function reply(
	status: number,
	code: 'VALIDATION_ERROR' | 'UNAVAILABLE',
	fieldErrors?: Partial<Record<BookDemoField, BookDemoErrorCode>>,
) {
	const headers: Record<string, string> = {
		'Cache-Control': 'no-store',
		'X-Content-Type-Options': 'nosniff',
	}
	if (status === 503) headers['Retry-After'] = '30'
	const named = fieldErrors && Object.keys(fieldErrors).length > 0

	return Response.json(
		{
			ok: false,
			code,
			...(named ? { fieldErrors } : {}),
		},
		{
			status,
			headers,
		},
	)
}

// A demo request is validated as the form validates it, then refused as
// paused: requests have no destination yet, and nothing is stored, logged or
// forwarded. The destination, when it exists, answers after the validation
// behind this same contract.
export async function postBookDemo(request: Request): Promise<Response> {
	let body: unknown
	try {
		const text = await readBoundedBody(request.body, bodyLimit, request.signal)
		body = parseSubmissionJSON(text)
	} catch (error) {
		const status = error instanceof BodyError ? error.status : 400

		return reply(status, 'VALIDATION_ERROR')
	}
	const result = bookDemoSubmissionSchema.safeParse(body)
	if (!result.success) {
		return reply(400, 'VALIDATION_ERROR', fieldErrorsOf(result.error))
	}

	return reply(503, 'UNAVAILABLE')
}
