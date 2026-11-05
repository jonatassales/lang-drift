import { isPlainJSON } from '@api/book-a-demo/body/is-plain-json.validation'
import { reply } from '@api/book-a-demo/book-demo-reply.response'

// A demo request is plain JSON: no encoding, no other media type.
export function interceptContentType(request: Request): Response | undefined {
	if (isPlainJSON(request)) return undefined

	return reply(415, 'INVALID_REQUEST')
}
