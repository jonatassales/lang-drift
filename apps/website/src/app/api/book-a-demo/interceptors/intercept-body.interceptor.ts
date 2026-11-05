import { BodyError } from '@api/book-a-demo/body/body-error.type'
import { parseSubmissionJSON } from '@api/book-a-demo/body/parse-submission-json.compute'
import { readBoundedBody } from '@api/book-a-demo/body/read-bounded-body.adapter'
import { bookDemoReply } from '@api/book-a-demo/book-demo-reply.map'

import { bodyLimit } from './body-limit.data'

export interface ReadBody {
	body: unknown
}

// The body is read within the limit and parsed as strict JSON; a read that
// overflows keeps its own status, anything else malformed is a 400.
export async function interceptBody(
	request: Request,
): Promise<Response | ReadBody> {
	try {
		const text = await readBoundedBody(request.body, bodyLimit, request.signal)

		return {
			body: parseSubmissionJSON(text),
		}
	} catch (error) {
		const status = error instanceof BodyError ? error.status : 400

		return bookDemoReply(status, 'VALIDATION_ERROR')
	}
}
