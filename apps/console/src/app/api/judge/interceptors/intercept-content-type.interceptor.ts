import { judgeReply } from '@app/api/judge/judge-reply.map'

// A question is plain JSON: no content encoding, no other media type.
export function interceptContentType(request: Request): Response | undefined {
	const mediaType = request.headers
		.get('content-type')
		?.split(';')[0]
		.trim()
		.toLowerCase()
	const plainJSON =
		mediaType === 'application/json' && !request.headers.has('content-encoding')

	if (plainJSON) return undefined

	return judgeReply(415, {
		code: 'UNSUPPORTED_MEDIA_TYPE',
	})
}
