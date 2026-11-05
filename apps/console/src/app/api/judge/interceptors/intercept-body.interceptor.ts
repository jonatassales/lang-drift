import { judgeBodyLimit } from '@app/api/judge/body-limit.data'
import { judgeReply } from '@app/api/judge/judge-reply.map'

export interface ReadQuestionBody {
	body: unknown
}

// The body is read and parsed as JSON. One over the limit keeps its own
// status; one that is not JSON names the body as the field at fault.
export async function interceptBody(
	request: Request,
): Promise<Response | ReadQuestionBody> {
	const text = await request.text()

	if (Buffer.byteLength(text) > judgeBodyLimit) {
		return judgeReply(413, {
			code: 'PAYLOAD_TOO_LARGE',
		})
	}

	try {
		return {
			body: JSON.parse(text),
		}
	} catch {
		return judgeReply(400, {
			code: 'INVALID_QUESTION',
			field: 'body',
		})
	}
}
