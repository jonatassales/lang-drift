import { judgeBodyLimit } from '@app/api/judge/body-limit.data'
import { judgeReply } from '@app/api/judge/judge-reply.map'

// A declared length over the limit is refused before the body is read.
export function interceptBodySize(request: Request): Response | undefined {
	const declared = Number(request.headers.get('content-length') ?? '0')
	const overLimit = Number.isFinite(declared) && declared > judgeBodyLimit

	if (!overLimit) return undefined

	return judgeReply(413, {
		code: 'PAYLOAD_TOO_LARGE',
	})
}
