import 'server-only'

import {
	BodyError,
	readBoundedBody,
} from '@components/early-access/request.body'
import {
	declaresOversizedBody,
	isPlainJSON,
	parseSubmissionJSON,
} from '@components/early-access/request.validation'

import {
	type BookDemoErrorCode,
	type BookDemoField,
	bookDemoSubmissionSchema,
	fieldErrorsOf,
} from './form.schema'

const bodyLimit = 4096

function reply(
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
			status,
			headers,
		},
	)
}

// The origin a request was addressed to. Next names every loopback host
// localhost in the request URL, so the Host header keeps 127.0.0.1 apart, as
// a browser does.
function addressedOrigin(request: Request): string {
	const addressed = new URL(request.url)
	addressed.host = request.headers.get('host') ?? addressed.host

	return addressed.origin
}

// A browser names the page behind a request in Origin and Sec-Fetch-Site.
// Book a demo has no configured public origin, so a page that names one must
// share the origin the request was addressed to. Browsers send Origin with
// every cross-origin POST, so a request without one is not a cross-origin
// browser request.
function isForeignRequest(request: Request): boolean {
	if (request.headers.get('sec-fetch-site') === 'cross-site') return true
	const origin = request.headers.get('origin')
	if (origin === null) return false

	return origin !== addressedOrigin(request)
}

// A demo request holds the early-access request contract before its body is
// read, and is validated as the form validates it, so that a destination slots
// in behind this same contract when one exists.
export async function postBookDemo(request: Request): Promise<Response> {
	if (request.method !== 'POST') return reply(405, 'INVALID_REQUEST')
	if (isForeignRequest(request)) return reply(403, 'INVALID_ORIGIN')
	if (!isPlainJSON(request)) return reply(415, 'INVALID_REQUEST')
	if (declaresOversizedBody(request, bodyLimit))
		return reply(413, 'INVALID_REQUEST')

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

	// Requests have no destination yet: refused as paused, and nothing is
	// stored, logged or forwarded.
	return reply(503, 'UNAVAILABLE')
}
