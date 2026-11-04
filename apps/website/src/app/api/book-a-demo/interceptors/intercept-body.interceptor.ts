import { reply } from '@components/book-a-demo/book-demo-reply.response'
import { parseSubmissionJSON } from '@components/early-access/parse-submission-json.compute'
import {
	BodyError,
	readBoundedBody,
} from '@components/early-access/request.body'

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

		return reply(status, 'VALIDATION_ERROR')
	}
}
