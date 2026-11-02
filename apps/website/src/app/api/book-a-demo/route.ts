import 'server-only'

import { reply } from '@components/book-a-demo/book-demo-reply.response'
import { fieldErrorsOf } from '@components/book-a-demo/field-errors.compute'
import { bookDemoSubmissionSchema } from '@components/book-a-demo/form.validation'
import { isForeignRequest } from '@components/book-a-demo/is-foreign-request.validation'
import { declaresOversizedBody } from '@components/early-access/declares-oversized-body.validation'
import { isPlainJSON } from '@components/early-access/is-plain-json.validation'
import { parseSubmissionJSON } from '@components/early-access/parse-submission-json.compute'
import {
	BodyError,
	readBoundedBody,
} from '@components/early-access/request.body'

export const runtime = 'nodejs'

const bodyLimit = 4096

// A demo request holds the early-access request contract before its body is
// read, and is validated as the form validates it, so that a destination slots
// in behind this same contract when one exists.
export async function POST(request: Request): Promise<Response> {
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
